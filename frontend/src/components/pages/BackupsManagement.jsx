// Page de gestion des sauvegardes pour une base donnée
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import DatabaseBackups from "../organisms/DatabaseBackups";
import ScheduleBackupModal from "../organisms/ScheduleBackupModal";
import "../../styles/pages/BackupsManagement.css";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTrash,
  faDownload,
  faPen,
  faPowerOff,
  faClock,
} from "@fortawesome/free-solid-svg-icons";
import { Tooltip } from "react-tooltip";
import {
  getUserScheduledTasks,
  toggleTaskActive,
  deleteScheduledTask,
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

  const [showCronModal, setShowCronModal] = useState(false);
  const [editTask, setEditTask] = useState(null);

  const [manualName, setManualName] = useState("");
  const [runningBackup, setRunningBackup] = useState(false);

  // ---- RÉCUPÉRER LES TÂCHES CRON ----
  const fetchTasks = async () => {
    setLoadingTasks(true);
    const response = await getUserScheduledTasks();
    setLoadingTasks(false);

    if (!response.success) {
      toast.error("Erreur lors du chargement des tâches planifiées");
      return;
    }

    const filtered = response.data.filter((t) => t.database_id === databaseId);
    setTasks(filtered);
  };

  const handleToggleActive = async (taskId, active) => {
    const response = await toggleTaskActive(taskId, active);

    if (response.success) {
      toast.success("Tâche mise à jour");
      fetchTasks();
    } else {
      toast.error("Erreur lors de la mise à jour");
    }
  };

  useEffect(() => {
    fetchTasks();
  }, [databaseId]);

  // ---- SAUVEGARDE IMMÉDIATE ----
  const handleManualBackup = async () => {
    if (runningBackup) return;

    setRunningBackup(true);

    const response = await createBackup(databaseId, manualName || null);

    setRunningBackup(false);

    if (!response.success) {
      toast.error("Erreur lors du backup manuel");
      return;
    }

    toast.success("Backup créé !");
    setManualName("");
  };

  const handleDeleteTask = async (taskId) => {
    if (!window.confirm("Supprimer cette tâche planifiée ?")) return;

    const response = await deleteScheduledTask(taskId);
    if (response.success) {
      toast.success("Tâche planifiée supprimée");
      fetchTasks();
    } else {
      toast.error("Erreur lors de la suppression");
    }
  };

  return (
    <div className="page-main-container">
      <h2>Gestion des sauvegardes</h2>
      {showCronModal && (
        <ScheduleBackupModal
          databaseId={databaseId}
          onClose={() => setShowCronModal(false)}
        />
      )}

      {/* ----- MODIFIER UNE TÂCHE CRON EXISTANTE ----- */}
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
      {/*    SAUVEGARDE MAINTENANT    */}
      {/* ─────────────────────────── */}

      <div className="manual-backup">
        <input
          type="text"
          placeholder="Nom du backup (optionnel)"
          value={manualName}
          onChange={(e) => setManualName(e.target.value)}
        />

        <button onClick={handleManualBackup} disabled={runningBackup}>
          {runningBackup ? "En cours..." : "Exécuter un backup maintenant"}
        </button>
      </div>
      {/* --- NOUVEAU PLANIFICATEUR CRON --- */}
      <button
        className="task-schedule-btn"
        onClick={() => setShowCronModal(true)}
      >
        Planifier un sauvegarde automatique
      </button>

      {/* --- LISTE DES SAUVEGARDES --- */}
      <DatabaseBackups databaseId={databaseId} />

      {/* --- TÂCHES CRON --- */}
      <h3 style={{ marginTop: "30px" }}>Sauvegardes Automatiques (CRON)</h3>

      {loadingTasks ? (
        <p>Chargement...</p>
      ) : tasks.length === 0 ? (
        <p>Aucune tâche planifiée.</p>
      ) : (
        <ul className="cron-list">
          {tasks.map((t) => (
            <li key={t.id} className="cron-item">
              {/* INFO */}
              <div className="cron-card">
                <FontAwesomeIcon icon={faClock} className="cron-icon" />

                <div className="cron-info-grid">
                  <strong>{humanReadableCron(t.cron_expression)}</strong>
                  <span className="cron-status">
                    {t.is_active ? "Actif" : "Inactif"}
                  </span>
                  <small>{t.last_run_at || "Jamais exécuté"}</small>
                </div>
              </div>

              {/* ACTIONS */}
              <div className="cron-actions">
                <button
                  className="cron-btn toggle"
                  onClick={() => handleToggleActive(t.id, !t.is_active)}
                  data-tooltip-id="cron-toggle"
                >
                  <FontAwesomeIcon icon={faPowerOff} />
                </button>

                <button
                  className="cron-btn edit"
                  onClick={() => setEditTask(t)}
                  data-tooltip-id="cron-edit"
                >
                  <FontAwesomeIcon icon={faPen} />
                </button>

                <button
                  className="cron-btn delete"
                  onClick={() => handleDeleteTask(t.id)}
                  data-tooltip-id="cron-delete"
                >
                  <FontAwesomeIcon icon={faTrash} />
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
      <Tooltip id="cron-toggle">Activer / Désactiver</Tooltip>
      <Tooltip id="cron-edit">Modifier</Tooltip>
      <Tooltip id="cron-delete">Supprimer</Tooltip>
    </div>
  );
};

export default BackupsManagement;
