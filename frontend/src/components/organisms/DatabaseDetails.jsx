import DeleteDatabase from "./DeleteDatabase";
import { useState } from "react";
import { toast } from "react-toastify";
import { createBackup } from "../../api/backupApi";
import DatabaseBackups from "./DatabaseBackups";
import ScheduleBackupModal from "./ScheduleBackupModal";

const DatabaseDetails = ({ db, onClose, onDeleted }) => {
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [backupName, setBackupName] = useState("");
  const [loading, setLoading] = useState(false);
  const [showCronModal, setShowCronModal] = useState(false);

  const handleCreateBackup = async () => {
    if (!backupName.trim()) {
      toast.error("Veuillez entrer un nom pour le backup!");
      return;
    }
    setLoading(true);
    const response = await createBackup(db.id, backupName, "v1");
    setLoading(false);

    if (response.success) {
      toast.success("Backup crée avec succès !");
      setBackupName("");
    } else {
      toast.error(response.message || "Erreur lors de la création du backup");
    }
  };

  return (
    <div className="modal">
      <div className="modal-content">
        <h2>Détails de la base: {db.name}</h2>
        <p>
          <strong>Type:</strong> {db.type}
        </p>
        <p>
          <strong>Host:</strong> {db.host}
        </p>
        <p>
          <strong>Port:</strong> {db.port}
        </p>
        <p>
          <strong>Utilisateur:</strong> {db.db_username}
        </p>

        <h3>Créer un backup</h3>
        <input
          type="text"
          placeholder="Nom du backup (ex: daily, manual...)"
          value={backupName}
          onChange={(e) => setBackupName(e.target.value)}
          style={{ padding: "8px", width: "250px", marginRight: "10px" }}
        />
        <button onClick={handleCreateBackup} disabled={loading}>
          {loading ? "Création..." : "Créer le backup"}
        </button>

        <div style={{ display: "flex", gap: "10px", marginTop: "20px" }}>
          <button
            style={{ backgroundColor: "blue", color: "white" }}
            onClick={() => setShowCronModal(true)}
          >
            Planifier un backup
          </button>
          <button
            style={{ backgroundColor: "red", color: "white" }}
            onClick={() => setShowDeleteModal(true)}
          >
            Supprimer
          </button>
          <button
            style={{ backgroundColor: "gray", color: "white" }}
            onClick={onClose}
          >
            Fermer
          </button>
        </div>

        {showCronModal && (
          <ScheduleBackupModal
            databaseId={db.id}
            onClose={() => setShowCronModal(false)}
          />
        )}

        {showDeleteModal && (
          <DeleteDatabase
            databaseId={db.id}
            onClose={() => setShowDeleteModal(false)}
            onDeleted={() => {
              onDeleted?.();
              onClose?.();
            }}
          />
        )}

        <DatabaseBackups databaseId={db.id} />
      </div>
    </div>
  );
};

export default DatabaseDetails;
