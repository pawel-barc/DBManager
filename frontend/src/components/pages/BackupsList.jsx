// Cette page récupère et affiche la liste des backups pour une base de donnée et permet de supprimer un backup.
import { useEffect, useState } from "react";
import { getAllBackups, deleteBackup } from "../../api/backupApi";
import { toast } from "react-toastify";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTrash, faDatabase } from "@fortawesome/free-solid-svg-icons";
import { Tooltip } from "react-tooltip";
import "../../styles/pages/BackupsList.css";
const BackupsListPage = () => {
  const [backups, setBackups] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchBackups = async () => {
    setLoading(true);
    const response = await getAllBackups();
    setLoading(false);
    if (response.success) {
      setBackups(response.data || []);
    } else {
      toast.error(response.message || "Erreur lors du chargement des backups");
    }
  };

  const handleDelete = async (backupId) => {
    if (!window.confirm("Voulez-vous vraiment supprimer ce backup ?")) return;
    const response = await deleteBackup(backupId);
    if (response.success) {
      toast.success("Backup supprimé !");
      fetchBackups();
    } else {
      toast.error(response.message || "Erreur lors de la suppression");
    }
  };

  useEffect(() => {
    fetchBackups();
  }, []);

  if (loading) return <p>Chargement...</p>;

  return (
    <div className="backups-container">
      <h2>Toutes les sauvegardes</h2>

      {backups.length === 0 ? (
        <p>Aucune sauvegarde.</p>
      ) : (
        <ul className="backups-list">
          {backups.map((b) => (
            <li key={b.id} className="backups-item">
              <div
                className={`backups-card ${
                  b.status === "success" ? "success" : "error"
                }`}
              >
                <FontAwesomeIcon icon={faDatabase} className="backup-icon" />

                <span className="backup-name">{b.name}</span>

                <span className={`backup-status ${b.status}`}>{b.status}</span>

                <span className="backup-date">{b.backup_date}</span>
              </div>

              <button
                className="delete-backup-btn"
                onClick={() => handleDelete(b.id)}
                data-tooltip-id="delete-backup-tooltip"
                aria-label="Delete backup"
              >
                <FontAwesomeIcon icon={faTrash} />
              </button>
            </li>
          ))}
        </ul>
      )}

      <Tooltip id="delete-backup-tooltip" place="right">
        Supprimer la sauvegarde
      </Tooltip>
    </div>
  );
};

export default BackupsListPage;
