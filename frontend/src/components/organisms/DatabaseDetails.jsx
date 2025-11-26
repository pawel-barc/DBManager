import DeleteDatabase from "./DeleteDatabase";
import { useState } from "react";

const DatabaseDetails = ({ db, onClose, onDeleted }) => {
  const [showDeleteModal, setShowDeleteModal] = useState(false);

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

        <div style={{ display: "flex", gap: "10px", marginTop: "20px" }}>
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
          {/* Backup button pour après */}
        </div>

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
      </div>
    </div>
  );
};

export default DatabaseDetails;
