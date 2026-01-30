// Modal d’édition d’une expression CRON existante
import { useState } from "react";
import CronSelector from "../../components/organisms/CronSelector";
import { updateCronExpression } from "../../api/cronApi";
import { toast } from "react-toastify";

const EditCronModal = ({ task, onClose, onUpdated }) => {
  const [cronValue, setCronValue] = useState(task.cron_expression);
  const [loading, setLoading] = useState(false);

  const handleSave = async () => {
    setLoading(true);

    const res = await updateCronExpression(task.id, cronValue);

    setLoading(false);

    if (!res.success) {
      toast.error("Erreur lors de la mise à jour du CRON");
      return;
    }

    onUpdated?.();
  };

  return (
    <div className="modal" style={styles.overlay}>
      <div className="modal-content" style={styles.modal}>
        <h3 style={{ color: "black" }}>Modifier la planification CRON</h3>

        <p style={{ marginBottom: "10px", fontSize: "14px" }}>
          Tâche: <strong>{task.name || `#${task.id}`}</strong>
        </p>

        {/* --- Cron Selector --- */}
        <CronSelector onChange={(expr) => setCronValue(expr)} />

        <div style={{ marginTop: "15px", fontSize: "14px" }}>
          Nouveau CRON: <strong>{cronValue}</strong>
        </div>

        <div style={styles.buttons}>
          <button onClick={onClose} disabled={loading} style={styles.cancel}>
            Annuler
          </button>

          <button onClick={handleSave} disabled={loading} style={styles.save}>
            {loading ? "Mise à jour..." : "Enregistrer"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default EditCronModal;

const styles = {
  overlay: {
    background: "rgba(0,0,0,0.5)",
    position: "fixed",
    left: 0,
    right: 0,
    top: 0,
    bottom: 0,
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    zIndex: 1000,
  },
  modal: {
    background: "white",
    padding: "20px",
    borderRadius: "8px",
    width: "400px",
    maxWidth: "95%",
  },
  buttons: {
    marginTop: "20px",
    display: "flex",
    justifyContent: "flex-end",
    gap: "10px",
  },
  cancel: {
    background: "gray",
    color: "white",
    padding: "8px 12px",
    borderRadius: "5px",
  },
  save: {
    background: "green",
    color: "white",
    padding: "8px 12px",
    borderRadius: "5px",
  },
};
