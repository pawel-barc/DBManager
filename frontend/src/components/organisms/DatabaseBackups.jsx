// Ce composant affiche la liste des backups d'une base de données, et permet de les télécharger ou les supprimer
import { useEffect, useState } from "react";
import { getDatabaseBackups, deleteBackup } from "../../api/backupApi";
import { toast } from "react-toastify";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTrash,
  faDownload,
  faDatabase,
} from "@fortawesome/free-solid-svg-icons";
import "../../styles/organisms/DatabaseBackups.css";

const DatabaseBackups = ({ databaseId }) => {
  const [backups, setBackups] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchBackups = async () => {
    setLoading(true);
    const response = await getDatabaseBackups(databaseId);
    setLoading(false);

    if (response.success) {
      setBackups(Array.isArray(response.data) ? response.data : []);
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

  const handleDownload = (backupId) => {
    window.open(`http://localhost:8080/backups/${backupId}/download`, "_blank");
  };

  useEffect(() => {
    fetchBackups();
  }, [databaseId]);

  if (loading) return <p>Chargement des Sauvegardes...</p>;

  return (
    <div className="db-backups-container">
      <h3>Sauvegardes manuels </h3>

      {backups.length === 0 ? (
        <p>Aucune sauvegarde disponible</p>
      ) : (
        <ul className="db-backups-list">
          {backups.map((b) => (
            <li key={b.id} className="db-backup-item">
              {/* CARD */}
              <div
                className={`db-backup-card ${
                  b.status === "success" ? "success" : "error"
                }`}
              >
                <FontAwesomeIcon icon={faDatabase} className="db-backup-icon" />

                <div className="db-backup-info">
                  <strong>{b.name}</strong>
                  <span className="status">Status: {b.status}</span>
                  <small>{new Date(b.backup_date).toLocaleString()}</small>
                </div>
              </div>

              <div className="db-backup-actions">
                <button
                  className="db-backup-btn download"
                  onClick={() => handleDownload(b.id)}
                  aria-label="Télécharger"
                >
                  <FontAwesomeIcon icon={faDownload} />
                </button>

                <button
                  className="db-backup-btn delete"
                  onClick={() => handleDelete(b.id)}
                  aria-label="Supprimer"
                >
                  <FontAwesomeIcon icon={faTrash} />
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default DatabaseBackups;
