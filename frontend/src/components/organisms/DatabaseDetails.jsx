import { Link } from "react-router-dom";
import { useState } from "react";
import DeleteDatabase from "./DeleteDatabase";

import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faDatabase,
  faServer,
  faUser,
  faNetworkWired,
  faTrash,
  faXmark,
  faRotateLeft,
  faFolderOpen,
} from "@fortawesome/free-solid-svg-icons";

import "../../styles/organisms/DatabaseDetails.css";

const DatabaseDetails = ({ db, onClose, onDeleted }) => {
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="db-details-modal" onClick={(e) => e.stopPropagation()}>
        {/* HEADER */}
        {/* HEADER */}
        <header className="db-details-header">
          <h2>
            <FontAwesomeIcon icon={faDatabase} /> {db.name}
          </h2>
          <button className="close-btn" onClick={onClose} aria-label="Close">
            <FontAwesomeIcon icon={faXmark} />
          </button>
        </header>

        {/* INFO */}
        <section className="db-details-info">
          <div>
            <FontAwesomeIcon icon={faServer} />
            <span>
              <strong>Type:</strong> {db.type}
            </span>
          </div>

          <div>
            <FontAwesomeIcon icon={faNetworkWired} />
            <span>
              <strong>Host:</strong> {db.host}
            </span>
          </div>

          <div>
            <FontAwesomeIcon icon={faNetworkWired} />
            <span>
              <strong>Port:</strong> {db.port}
            </span>
          </div>

          <div>
            <FontAwesomeIcon icon={faUser} />
            <span>
              <strong>Utilisateur:</strong> {db.db_username}
            </span>
          </div>
        </section>

        {/* ACTIONS */}
        <section className="db-details-actions">
          <Link
            to={`/databases/${db.id}/backups`}
            className="action-btn primary"
          >
            <FontAwesomeIcon icon={faFolderOpen} />
            Sauvegardes
          </Link>

          <Link
            to={`/databases/${db.id}/restauration`}
            className="action-btn warning"
          >
            <FontAwesomeIcon icon={faRotateLeft} />
            Restaurer
          </Link>

          <button
            className="action-btn danger"
            onClick={() => setShowDeleteModal(true)}
          >
            <FontAwesomeIcon icon={faTrash} />
            Supprimer
          </button>
        </section>

        {/* DELETE MODAL */}
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
