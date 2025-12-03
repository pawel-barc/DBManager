// Page de gestion des sauvegardes pour une base donnée
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import DatabaseBackups from "../organisms/DatabaseBackups";
import ScheduleBackupModal from "../organisms/ScheduleBackupModal";

import {
  getUserScheduledTasks,
  toggleTaskActive,
  updateCronExpression,
} from "../../api/cronApi";

import { createBackup } from "../../api/backupApi";
import { toast } from "react-toastify";
import { humanReadableCron } from "../../utils/helper";
import EditCronModal from "../../components/organisms/EditCronModal";

const BackupsManagement = () => {
  const { id } = useParams(); // databaseId
  const databaseId = Number(id);

  const [tasks, setTasks] = useState([]);
  const [loadingTasks, setLoadingTasks] = useState(true);

  const [showCronModal, setShowCronModal] = useState(false); // create new
  const [editTask, setEditTask] = useState(null); // edit existing cron

  const [manualName, setManualName] = useState("");
  const [runningBackup, setRunningBackup] = useState(false);

  // ---- FETCH CRON TASKS ----
  const fetchTasks = async () => {
    setLoadingTasks(true);
    const res = await getUserScheduledTasks();
    setLoadingTasks(false);

    if (!res.success) {
      toast.error("Erreur lors du chargement des tâches planifiées");
      return;
    }

    const filtered = res.data.filter((t) => t.database_id === databaseId);
    setTasks(filtered);
  };

  const handleToggleActive = async (taskId, active) => {
    const res = await toggleTaskActive(taskId, active);

    if (res.success) {
      toast.success("Tâche mise à jour");
      fetchTasks();
    } else {
      toast.error("Erreur lors de la mise à jour");
    }
  };

  useEffect(() => {
    fetchTasks();
  }, [databaseId]);

  // ---- BACKUP NOW ----
  const handleManualBackup = async () => {
    if (runningBackup) return;

    setRunningBackup(true);

    const res = await createBackup(databaseId, manualName || null);

    setRunningBackup(false);

    if (!res.success) {
      toast.error("Erreur lors du backup manuel");
      return;
    }

    toast.success("Backup créé !");
    setManualName("");
  };

  return (
    <div className="page-container">
      <h2>Gestion des sauvegardes</h2>

      {/* --- NEW CRON SCHEDULER --- */}
      <button
        style={{ marginBottom: "20px", background: "blue", color: "white" }}
        onClick={() => setShowCronModal(true)}
      >
        Planifier un backup automatique
      </button>

      {showCronModal && (
        <ScheduleBackupModal
          databaseId={databaseId}
          onClose={() => setShowCronModal(false)}
        />
      )}

      {/* ----- EDIT EXISTING CRON MODAL ----- */}
      {editTask && (
        <EditCronModal
          task={editTask}
          onClose={() => setEditTask(null)}
          onUpdated={() => {
            setEditTask(null);
            fetchTasks();
            toast.success("Expression CRON mise à jour !");
          }}
        />
      )}

      {/* ─────────────────────────── */}
      {/*       BACKUP NOW            */}
      {/* ─────────────────────────── */}
      <h3>Backup manuel immédiat</h3>

      <div style={{ marginBottom: "20px" }}>
        <input
          type="text"
          placeholder="Nom du backup (optionnel)"
          value={manualName}
          onChange={(e) => setManualName(e.target.value)}
          style={{
            padding: "8px",
            marginRight: "10px",
            borderRadius: "5px",
            border: "1px solid #ccc",
          }}
        />

        <button
          onClick={handleManualBackup}
          disabled={runningBackup}
          style={{
            background: runningBackup ? "gray" : "green",
            color: "white",
            padding: "8px 12px",
            borderRadius: "5px",
          }}
        >
          {runningBackup ? "En cours..." : "Exécuter un backup maintenant"}
        </button>
      </div>

      {/* --- BACKUPS LIST --- */}
      <DatabaseBackups databaseId={databaseId} />

      {/* --- CRON TASKS --- */}
      <h3 style={{ marginTop: "30px" }}>Backups Automatiques (CRON)</h3>

      {loadingTasks ? (
        <p>Chargement...</p>
      ) : tasks.length === 0 ? (
        <p>Aucune tâche planifiée.</p>
      ) : (
        <ul>
          {tasks.map((t) => (
            <li key={t.id} style={{ marginBottom: "10px" }}>
              {humanReadableCron(t.cron_expression)}— Dernier run:{" "}
              {t.last_run_at || "jamais"}
              <button
                style={{ marginLeft: "10px" }}
                onClick={() => handleToggleActive(t.id, !t.is_active)}
              >
                {t.is_active ? "Désactiver" : "Activer"}
              </button>
              {/* --- OPEN CRON EDIT MODAL --- */}
              <button
                style={{ marginLeft: "10px" }}
                onClick={() => setEditTask(t)}
              >
                Modifier
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default BackupsManagement;
