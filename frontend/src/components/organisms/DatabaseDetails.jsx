// Composant modal affichant les détails d'une base de données et offrant des actions rapides :
// gérer les sauvegardes, restaurer la base ou fermer le modal.
import { Link } from "react-router-dom";
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

        <div style={{ marginTop: "20px" }}>
          <Link to={`/databases/${db.id}/backups`} className="button-blue">
            Gérer les sauvegardes
          </Link>
          <br />
          <Link
            to={`/databases/${db.id}/restauration`}
            className="button-orange"
          >
            Restaurer la base
          </Link>
          <br />
          <button
            style={{ background: "red", color: "white", marginLeft: "10px" }}
            onClick={() => setShowDeleteModal(true)}
          >
            Supprimer la base
          </button>

          <button
            style={{ background: "gray", color: "white", marginLeft: "10px" }}
            onClick={onClose}
          >
            Fermer
          </button>
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
