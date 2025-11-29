// Ce composant affiche la liste des backups d'une base de données, et permet de les télécharger ou les supprimer
import { useEffect, useState } from "react";
import { getDatabaseBackups, deleteBackup } from "../../api/backupApi";
import { toast } from "react-toastify";

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

  if (loading) return <p>Chargement des backups...</p>;

  return (
    <div style={{ marginTop: "20px" }}>
      <h3>Backups</h3>
      {backups.length === 0 ? (
        <p>Aucun backup disponible</p>
      ) : (
        <ul>
          {backups.map((b) => (
            <li key={b.id} style={{ marginBottom: "10px" }}>
              {b.name} - {b.status} - {b.backup_date}
              <button
                style={{ marginLeft: "10px", color: "red" }}
                onClick={() => handleDelete(b.id)}
              >
                Supprimer
              </button>
              <button
                style={{ marginLeft: "5px", color: "green" }}
                onClick={() => handleDownload(b.id)}
              >
                Télécharger
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default DatabaseBackups;
