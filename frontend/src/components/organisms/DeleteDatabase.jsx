// Ce composant affiche une fenêtre modale permettant de supprimer une base de données existante.
import { useState } from "react";
import { toast } from "react-toastify";
import { deleteDatabase } from "../../api/databaseApi";

const DeleteDatabase = ({ databaseId, onClose, onDeleted }) => {
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);

  const handleDelete = async () => {
    setLoading(true);
    setErrorMessage(null);
    try {
      const response = await deleteDatabase(databaseId);
      if (response.success) {
        toast.success("La base de données a été supprimée avec succès !");
        onDeleted?.();
        onClose?.();
      } else if (response.error) {
        setErrorMessage(response.error);
      }
    } catch (err) {
      setErrorMessage(err.message || "Une erreur est survenue");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal">
      <div className="modal-content">
        <h3>Supprimer la base de données</h3>
        <p>
          Êtes-vous sûr de vouloir supprimer cette base ? Cette action est
          irréversible.
        </p>

        {errorMessage && <p style={{ color: "red" }}>{errorMessage}</p>}

        <div className="modal-buttons">
          <button onClick={handleDelete} disabled={loading}>
            {loading ? "Suppression..." : "Oui, supprimer"}
          </button>
          <button onClick={onClose} disabled={loading}>
            Annuler
          </button>
        </div>
      </div>
    </div>
  );
};

export default DeleteDatabase;
