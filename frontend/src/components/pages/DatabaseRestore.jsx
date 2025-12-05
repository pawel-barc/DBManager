// Composant permettant à l'utilisateur de visualiser et de restaurer les backups d'une base de données.
// Affiche la liste des backups disponibles, gère l'état de chargement et indique quel backup est en cours de restauration.
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import restoreBackup from "../../api/restoreApi";
import { getDatabaseBackups } from "../../api/backupApi";
import { toast } from "react-toastify";

const DatabaseRestore = () => {
  const { id } = useParams();
  const databaseId = Number(id);

  const [backups, setBackups] = useState([]);
  const [loading, setLoading] = useState(false);
  const [restoring, setRestoring] = useState(null);

  const fetchBackups = async () => {
    setLoading(true);
    const response = await getDatabaseBackups(databaseId);
    setLoading(false);

    if (response.success) {
      setBackups(response.data || []);
    } else {
      toast.error(response.message || "Erreur lors du chargement des backups");
    }
  };

  useEffect(() => {
    fetchBackups();
  }, [databaseId]);

  const handleRestore = async (backupId) => {
    if (restoring) return;
    setRestoring(backupId);
    const response = await restoreBackup(backupId);
    setRestoring(null);

    if (response.success) {
      toast.success("Restauration en cours !");
      fetchBackups();
    } else {
      toast.error(response.message || "Erreur lors de la restauration");
    }
    setRestoring(null);
  };

  return (
    <div className="page-container">
      <h2>Restaurer la base</h2>

      {loading ? (
        <p>Chargement...</p>
      ) : backups.length === 0 ? (
        <p>Aucun backup trouvé.</p>
      ) : (
        <ul>
          {backups.map((b) => (
            <li key={b.id} style={{ marginBottom: "10px" }}>
              <strong>{b.name}</strong> —{" "}
              {new Date(b.backup_date).toLocaleString()} —{" "}
              {b.file_size?.toFixed(2)} MB
              <button
                onClick={() => handleRestore(b.id)}
                disabled={restoring === b.id}
                style={{
                  marginLeft: "10px",
                  background: restoring === b.id ? "gray" : "orange",
                  color: "white",
                }}
              >
                {restoring === b.id ? "Restoring..." : "Restaurer"}
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default DatabaseRestore;
