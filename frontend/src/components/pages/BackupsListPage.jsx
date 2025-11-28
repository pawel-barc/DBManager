import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getAllBackups, deleteBackup } from "../../api/backupApi";
import { toast } from "react-toastify";

const BackupsListPage = () => {
  const { databaseId } = useParams();
  const [backups, setBackups] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchBackups = async () => {
    setLoading(true);
    const response = await getAllBackups(databaseId);
    setLoading(false);
    if (response.success) {
      setBackups(response.data);
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
  }, [databaseId]);

  if (loading) return <p>Chargement...</p>;

  return (
    <div>
      <h2>Backups de la base {databaseId}</h2>
      {backups.length === 0 ? (
        <p>Aucun backup</p>
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
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default BackupsListPage;
