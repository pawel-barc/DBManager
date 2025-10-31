import React, { useEffect, useMemo, useRef, useState } from "react";
import { toast } from "react-toastify";

// =====================
// API helpers (inchangés)
// =====================
async function apiListDatabases() {
  const r = await fetch("/api/databases", {
    headers: { Accept: "application/json" },
    credentials: "include", // utile si session cookie
  });
  if (r.status === 204) return []; // aucun contenu
  if (!r.ok) throw new Error("Impossible de charger les bases");
  const payload = await r.json().catch(() => {
    throw new Error("Réponse non-JSON de /api/databases");
  });
  // Accepte tableau direct OU {data:[...]}
  return Array.isArray(payload) ? payload : (Array.isArray(payload?.data) ? payload.data : []);
}

async function apiListBackups(database_id) {
  const url = new URL("/api/backups", window.location.origin);
  if (database_id) url.searchParams.set("database_id", database_id);
  const r = await fetch(url);
  if (!r.ok) throw new Error("Impossible de charger les backups");
  return r.json();
}
async function apiCreateBackup(body) {
  const r = await fetch("/api/backups", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!r.ok) {
    const err = await r.json().catch(() => ({}));
    throw new Error(err?.error || "Échec de la création du backup");
  }
  return r.json();
}
async function apiGetBackupStatus(id) {
  const r = await fetch(`/api/backups/${id}/status`);
  if (!r.ok) throw new Error("Statut du backup indisponible");
  return r.json();
}
async function apiGetBackupLog(id) {
  const r = await fetch(`/api/backups/${id}/log`);
  if (!r.ok) throw new Error("Log indisponible");
  return r.json();
}
function downloadHref(id) {
  return `/api/backups/${id}/download`;
}

// =====================
// Utils
// =====================
function humanSize(bytes) {
  if (bytes == null) return "—";
  const units = ["B", "KB", "MB", "GB", "TB"];
  let i = 0;
  let n = Number(bytes);
  while (n >= 1024 && i < units.length - 1) {
    n /= 1024;
    i++;
  }
  return `${n.toFixed(1)} ${units[i]}`;
}
function defaultBackupName({ project="app", env="prod", dbName="db", strategy="full", engine="pgdump", engineMajor="16" } = {}) {
  const dt = new Date();
  const pad = (x) => String(x).padStart(2, "0");
  const YYYY = dt.getUTCFullYear();
  const MM = pad(dt.getUTCMonth() + 1);
  const DD = pad(dt.getUTCDate());
  const HH = pad(dt.getUTCHours());
  const mm = pad(dt.getUTCMinutes());
  const ss = pad(dt.getUTCSeconds());
  return `${project}-${env}-${dbName}-${strategy}-${engine}-v${engineMajor}-${YYYY}${MM}${DD}T${HH}${mm}${ss}Z-001`;
}

function copy(text) {
  if (!text) return;
  navigator.clipboard?.writeText(text)
    .then(() => toast.success("Chemin copié"))
    .catch(() => {
      // fallback très simple
      const ta = document.createElement("textarea");
      ta.value = text;
      document.body.appendChild(ta);
      ta.select();
      try {
        document.execCommand("copy");
        toast.success("Chemin copié");
      } catch {
        toast.error("Impossible de copier");
      } finally {
        document.body.removeChild(ta);
      }
    });
}

// Présets d'outils/versions (à ajuster selon tes environnements)
const VERSION_PRESETS = {
  postgres: [
    "pg_dump 17",
    "pg_dump 16",
    "pg_dump 15",
  ],
  mysql: [
    "mysqldump 8.0",
    "mysqldump 5.7",
  ],
};
function inferFamily(dbType = "") {
  const t = String(dbType).toLowerCase();
  if (t.includes("postgre")) return "postgres";
  if (t.includes("mysql") || t.includes("maria")) return "mysql";
  return "other";
}

// =====================
// Composant
// =====================
export default function BackupManager() {
  const [databases, setDatabases] = useState([]);
  const [loadingDBs, setLoadingDBs] = useState(true);
  const [selectedDatabaseId, setSelectedDatabaseId] = useState(
    () => localStorage.getItem("bm.selectedDatabaseId") || ""
  );

  const [backups, setBackups] = useState([]);
  const [loadingBackups, setLoadingBackups] = useState(false);

  // Formulaire
  const [form, setForm] = useState({
    name: "",
    version: "",
  });
  const [versionMode, setVersionMode] = useState(
    () => localStorage.getItem("bm.versionMode") || "preset" // 'preset' | 'custom'
  );
  const [selectedPreset, setSelectedPreset] = useState(
    () => localStorage.getItem("bm.selectedPreset") || ""
  );

  // UX: filtres & auto-refresh
  const [statusFilter, setStatusFilter] = useState("all"); // all|done|running|queued|error
  const [search, setSearch] = useState("");
  const [autoRefresh, setAutoRefresh] = useState(
    () => localStorage.getItem("bm.autoRefresh") === "true"
  );
  const autoRefreshRef = useRef(null);

  const [selectedLog, setSelectedLog] = useState(null); // {id, log}
  const pollingRef = useRef(null);

  const selDb = useMemo(
    () => databases.find((d) => String(d.id) === String(selectedDatabaseId)),
    [databases, selectedDatabaseId]
  );
  const dbFamily = inferFamily(selDb?.type);

  // Charger bases
  useEffect(() => {
    (async () => {
      try {
        const data = await apiListDatabases();
        setDatabases(data || []);
        if (!selectedDatabaseId && data?.length) {
          setSelectedDatabaseId(String(data[0].id));
        }
      } catch (e) {
        toast.error(e.message);
      } finally {
        setLoadingDBs(false);
      }
    })();
  }, []); // only once

  // Persist selection DB
  useEffect(() => {
    if (selectedDatabaseId) localStorage.setItem("bm.selectedDatabaseId", selectedDatabaseId);
  }, [selectedDatabaseId]);

  // Charger backups quand la base change
  useEffect(() => {
    if (!selectedDatabaseId) return;
    (async () => {
      setLoadingBackups(true);
      try {
        const list = await apiListBackups(selectedDatabaseId);
        setBackups(list || []);
      } catch (e) {
        toast.error(e.message);
      } finally {
        setLoadingBackups(false);
      }
    })();
  }, [selectedDatabaseId]);

  // Pré-remplir nom par défaut quand la base ou son nom changent
  useEffect(() => {
    setForm((f) => ({
      ...f,
      name: defaultBackupName(selDb?.name),
    }));
  }, [selDb?.name]);

  // Préremplir version selon la famille + preset sélectionné
  useEffect(() => {
    if (versionMode === "preset") {
      const presets = VERSION_PRESETS[dbFamily] || [];
      // si le preset courant n'est pas compatible, on prend le 1er disponible
      if (!presets.includes(selectedPreset)) {
        const next = presets[0] || "";
        setSelectedPreset(next);
        setForm((f) => ({ ...f, version: next || "" }));
        localStorage.setItem("bm.selectedPreset", next);
      } else {
        setForm((f) => ({ ...f, version: selectedPreset }));
      }
    }
  }, [dbFamily, versionMode]); // re-évalue quand la DB change ou le mode change

  // Persistance des choix UX
  useEffect(() => {
    localStorage.setItem("bm.versionMode", versionMode);
  }, [versionMode]);
  useEffect(() => {
    localStorage.setItem("bm.selectedPreset", selectedPreset);
  }, [selectedPreset]);
  useEffect(() => {
    localStorage.setItem("bm.autoRefresh", String(autoRefresh));
  }, [autoRefresh]);

  function onFormChange(e) {
    const { name, value } = e.target;
    setForm((p) => ({ ...p, [name]: value }));
  }

  async function refreshBackups() {
    if (!selectedDatabaseId) return;
    setLoadingBackups(true);
    try {
      setBackups(await apiListBackups(selectedDatabaseId));
    } catch (e) {
      toast.error(e.message);
    } finally {
      setLoadingBackups(false);
    }
  }

  function stopPolling() {
    if (pollingRef.current) {
      clearInterval(pollingRef.current);
      pollingRef.current = null;
    }
  }
  async function startPolling(backupId) {
    stopPolling();
    pollingRef.current = setInterval(async () => {
      try {
        const st = await apiGetBackupStatus(backupId);
        if (st?.status === "running" || st?.status === "queued") return;
        stopPolling();
        if (st?.status === "done") {
          toast.success("Backup terminé");
          refreshBackups();
        } else if (st?.status === "error") {
          toast.error(st?.error || "Erreur du backup");
          refreshBackups();
        }
      } catch (e) {
        stopPolling();
        toast.error(e.message);
      }
    }, 1500);
  }

  // Auto-refresh de la liste (hors polling)
  useEffect(() => {
    if (!autoRefresh) {
      if (autoRefreshRef.current) {
        clearInterval(autoRefreshRef.current);
        autoRefreshRef.current = null;
      }
      return;
    }
    autoRefreshRef.current = setInterval(() => {
      // si on poll déjà un backup en cours, on évite un refresh "brutal"
      if (!pollingRef.current) refreshBackups();
    }, 10000);
    return () => {
      if (autoRefreshRef.current) clearInterval(autoRefreshRef.current);
    };
  }, [autoRefresh, selectedDatabaseId]);

  async function handleCreateBackup(e) {
    e?.preventDefault?.();
    if (!selectedDatabaseId) return toast.warn("Choisis une base");
    if (!form.name?.trim()) return toast.warn("Donne un nom de backup");
    // Version: si mode preset et pas de preset, c'est acceptable (backend optionnel), sinon trim
    const version = form.version?.trim() || undefined;

    try {
      const res = await apiCreateBackup({
        database_id: Number(selectedDatabaseId),
        name: form.name.trim(),
        version,
      });
      toast.info("Backup lancé…");
      if (res?.backup_id) startPolling(res.backup_id);
      else refreshBackups();
    } catch (e) {
      toast.error(e.message);
    }
  }

  async function openLog(id) {
    try {
      const r = await apiGetBackupLog(id);
      setSelectedLog({ id, log: r?.log ?? "" });
    } catch (e) {
      toast.error(e.message);
    }
  }

  // Filtrage local
  const filteredBackups = useMemo(() => {
    const q = search.trim().toLowerCase();
    return (backups || []).filter((b) => {
      const statusOk = statusFilter === "all" ? true : (b.status || "") === statusFilter;
      const text = `${b?.name || ""} ${b?.file_path || ""}`.toLowerCase();
      const searchOk = q ? text.includes(q) : true;
      return statusOk && searchOk;
    });
  }, [backups, statusFilter, search]);

  // UI helpers
  const presets = VERSION_PRESETS[dbFamily] || [];
  const showPresetSelect = versionMode === "preset";
  const showCustomInput = versionMode === "custom";

  return (
    <div className="min-h-screen w-full p-6 bg-gray-50">
      <div className="max-w-6xl mx-auto space-y-6">
        <header className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <h1 className="text-2xl font-semibold">Backups – DB Manager</h1>
          <div className="text-sm text-gray-500">
            Colonnes:&nbsp;
            <code>id</code>, <code>database_id</code>, <code>name</code>, <code>file_path</code>,{" "}
            <code>file_size</code>, <code>backup_date</code>, <code>status</code>, <code>version</code>,{" "}
            <code>log</code>
          </div>
        </header>

        {/* Sélecteur de base */}
        <section className="bg-white rounded-2xl shadow p-5">
          <div className="flex items-center justify-between mb-3">
            <h2 className="text-lg font-medium">Choisir une base</h2>
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                className="rounded border-gray-300"
                checked={autoRefresh}
                onChange={(e) => setAutoRefresh(e.target.checked)}
              />
              Auto-rafraîchir (10 s)
            </label>
          </div>

          {loadingDBs ? (
            <div className="text-gray-500">Chargement des bases…</div>
          ) : databases.length === 0 ? (
            <div className="text-red-600">Aucune base enregistrée. Ajoute d&apos;abord une base.</div>
          ) : (
            <div className="grid gap-3 sm:grid-cols-3">
              <label className="block">
                <span className="text-sm text-gray-600">Base</span>
                <select
                  className="mt-1 w-full rounded-xl border p-2"
                  value={selectedDatabaseId}
                  onChange={(e) => setSelectedDatabaseId(e.target.value)}
                >
                  {databases.map((d) => (
                    <option key={d.id} value={d.id}>
                      {d.name} ({d.type}) — {d.host}:{d.port}
                    </option>
                  ))}
                </select>
              </label>

              <div className="text-sm text-gray-500 flex items-end">
                Utilisateur:&nbsp;
                <span className="font-mono" title="db_username">
                  {selDb?.db_username || "—"}
                </span>
              </div>

              <div className="text-sm text-gray-500 flex items-end">
                Outil détecté:&nbsp;
                <span className="font-mono" title="Famille d’outil selon type de base">
                  {dbFamily === "postgres" ? "pg_dump" : dbFamily === "mysql" ? "mysqldump" : "—"}
                </span>
              </div>
            </div>
          )}
        </section>

        {/* Formulaire de création */}
        <form onSubmit={handleCreateBackup} className="bg-white rounded-2xl shadow p-5">
          <h2 className="text-lg font-medium mb-3">Créer un backup</h2>

          <div className="grid gap-4 md:grid-cols-3">
            <label className="block md:col-span-1">
              <span className="text-sm text-gray-600">Nom</span>
              <input
                name="name"
                value={form.name}
                onChange={onFormChange}
                className="mt-1 w-full rounded-xl border p-2"
                required
              />
            </label>

            <div className="md:col-span-2 grid gap-2">
              <div className="flex items-center gap-3">
                <span className="text-sm text-gray-600">Version / Outil</span>
                <div className="flex items-center gap-3 text-sm">
                  <label className="inline-flex items-center gap-1">
                    <input
                      type="radio"
                      name="versionMode"
                      value="preset"
                      checked={versionMode === "preset"}
                      onChange={() => setVersionMode("preset")}
                    />
                    Préset
                  </label>
                  <label className="inline-flex items-center gap-1">
                    <input
                      type="radio"
                      name="versionMode"
                      value="custom"
                      checked={versionMode === "custom"}
                      onChange={() => setVersionMode("custom")}
                    />
                    Autre…
                  </label>
                </div>
              </div>

              {showPresetSelect && (
                <div className="grid gap-2 md:grid-cols-2">
                  <label className="block">
                    <span className="sr-only">Préset d’outil</span>
                    <select
                      className="mt-1 w-full rounded-xl border p-2"
                      value={selectedPreset}
                      onChange={(e) => {
                        const val = e.target.value;
                        setSelectedPreset(val);
                        setForm((f) => ({ ...f, version: val }));
                      }}
                      disabled={presets.length === 0}
                      title={
                        presets.length
                          ? "Choisis un préset adapté au type de base"
                          : "Aucun préset disponible pour ce type de base"
                      }
                    >
                      {presets.length === 0 ? (
                        <option value="">— Aucun préset pour cette base —</option>
                      ) : (
                        presets.map((p) => (
                          <option key={p} value={p}>
                            {p}
                          </option>
                        ))
                      )}
                    </select>
                  </label>

                  <div className="text-xs text-gray-500 flex items-center">
                    Astuce: bascule sur “Autre…” pour saisir une version libre (ex: <code>pg_dump 16.3</code>).
                  </div>
                </div>
              )}

              {showCustomInput && (
                <label className="block">
                  <span className="sr-only">Version personnalisée</span>
                  <input
                    name="version"
                    value={form.version}
                    onChange={onFormChange}
                    placeholder={
                      dbFamily === "postgres"
                        ? "ex: pg_dump 16.3"
                        : dbFamily === "mysql"
                        ? "ex: mysqldump 8.0"
                        : "ex: pg_dump 16.3 / mysqldump 8.0"
                    }
                    className="mt-1 w-full rounded-xl border p-2"
                  />
                </label>
              )}
            </div>
          </div>

          <div className="mt-4 flex flex-wrap gap-3">
            <button type="submit" className="px-4 py-2 rounded-xl bg-black text-white hover:opacity-90">
              Lancer le backup
            </button>
            <button
              type="button"
              onClick={() => setForm((f) => ({ ...f, name: defaultBackupName(selDb?.name) }))}
              className="px-4 py-2 rounded-xl border"
              title="Génère un nom basé sur la date/heure et le nom de la base"
            >
              Nom par défaut
            </button>
            <button type="button" onClick={refreshBackups} className="px-4 py-2 rounded-xl border">
              Rafraîchir la liste
            </button>
          </div>
        </form>

        {/* Outils de liste (filtres, recherche) */}
        <section className="bg-white rounded-2xl shadow p-5">
          <div className="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
            <h2 className="text-lg font-medium">Backups existants</h2>
            <div className="flex flex-col gap-3 md:flex-row">
              <label className="block">
                <span className="text-sm text-gray-600">Filtrer par statut</span>
                <select
                  className="mt-1 w-full rounded-xl border p-2"
                  value={statusFilter}
                  onChange={(e) => setStatusFilter(e.target.value)}
                >
                  <option value="all">Tous</option>
                  <option value="done">done</option>
                  <option value="running">running</option>
                  <option value="queued">queued</option>
                  <option value="error">error</option>
                </select>
              </label>

              <label className="block">
                <span className="text-sm text-gray-600">Recherche</span>
                <input
                  className="mt-1 w-full rounded-xl border p-2"
                  placeholder="Rechercher (nom, chemin)…"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
              </label>
            </div>
          </div>

          <div className="overflow-x-auto mt-4">
            {loadingBackups ? (
              <div className="text-gray-500">Chargement…</div>
            ) : filteredBackups?.length ? (
              <table className="min-w-full text-sm">
                <thead>
                  <tr className="text-left border-b">
                    <th className="py-2 pr-4">id</th>
                    <th className="py-2 pr-4">database_id</th>
                    <th className="py-2 pr-4">name</th>
                    <th className="py-2 pr-4">file_path</th>
                    <th className="py-2 pr-4">file_size</th>
                    <th className="py-2 pr-4">backup_date</th>
                    <th className="py-2 pr-4">status</th>
                    <th className="py-2 pr-4">version</th>
                    <th className="py-2 pr-4">log</th>
                    <th className="py-2 pr-4">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredBackups.map((b) => (
                    <tr key={b.id} className="border-b last:border-0 align-top">
                      <td className="py-2 pr-4">{b.id}</td>
                      <td className="py-2 pr-4">{b.database_id}</td>
                      <td className="py-2 pr-4 font-mono">{b.name}</td>
                      <td className="py-2 pr-4 max-w-[22rem]">
                        <div className="flex items-center gap-2">
                          <div className="truncate" title={b.file_path || ""}>
                            {b.file_path || "—"}
                          </div>
                          {b.file_path && (
                            <button
                              onClick={() => copy(b.file_path)}
                              className="px-2 py-0.5 rounded-lg border text-xs"
                              title="Copier le chemin"
                            >
                              Copier
                            </button>
                          )}
                        </div>
                      </td>
                      <td className="py-2 pr-4">{humanSize(b.file_size)}</td>
                      <td className="py-2 pr-4" title={b.backup_date || ""}>
                        {b.backup_date ? new Date(b.backup_date).toLocaleString() : "—"}
                      </td>
                      <td className="py-2 pr-4">
                        {b.status === "done" && (
                          <span className="px-2 py-1 rounded-full bg-green-100 text-green-700">done</span>
                        )}
                        {b.status === "running" && (
                          <span className="px-2 py-1 rounded-full bg-blue-100 text-blue-700">running</span>
                        )}
                        {b.status === "queued" && (
                          <span className="px-2 py-1 rounded-full bg-gray-100 text-gray-700">queued</span>
                        )}
                        {b.status === "error" && (
                          <span className="px-2 py-1 rounded-full bg-red-100 text-red-700">error</span>
                        )}
                        {!b.status && <span className="text-gray-400">—</span>}
                      </td>
                      <td className="py-2 pr-4">{b.version || "—"}</td>
                      <td className="py-2 pr-4">
                        <button onClick={() => openLog(b.id)} className="px-2 py-1 rounded-lg border">
                          Voir
                        </button>
                      </td>
                      <td className="py-2 pr-4 space-x-2">
                        {b.status === "done" ? (
                          <a
                            href={downloadHref(b.id)}
                            className="px-3 py-1 rounded-lg bg-black text-white hover:opacity-90"
                          >
                            Télécharger
                          </a>
                        ) : (
                          <span className="text-gray-400">—</span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (
              <div className="text-gray-500">Aucun backup pour cette base.</div>
            )}
          </div>
        </section>

        {/* Modal log */}
        {selectedLog && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
            <div className="bg-white w-full max-w-3xl rounded-2xl shadow p-5">
              <div className="flex items-center justify-between mb-3">
                <h3 className="text-lg font-medium">Log du backup #{selectedLog.id}</h3>
                <button onClick={() => setSelectedLog(null)} className="px-2 py-1 rounded-lg border">
                  Fermer
                </button>
              </div>
              <pre className="whitespace-pre-wrap text-xs bg-gray-50 p-3 rounded-xl max-h-[60vh] overflow-auto">
                {selectedLog.log || "(vide)"}
              </pre>
            </div>
          </div>
        )}

        <footer className="text-xs text-gray-500">
          Le backend renseigne <code>file_path</code>, <code>file_size</code>, <code>backup_date</code>,{" "}
          <code>status</code> et <code>log</code> au fil de l&apos;exécution. Le téléchargement utilise l&apos;ID du
          backup.
        </footer>
      </div>
    </div>
  );
}
