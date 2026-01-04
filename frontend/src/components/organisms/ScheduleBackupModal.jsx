// Composant affichant une fenêtre modale permettant à l'utilisateur de programmer un backup récurrent pour une base de données.
import { useState } from "react";
import CronSelector from "./CronSelector";
import { addScheduledTask } from "../../api/cronApi";
import { toast } from "react-toastify";
import "../../styles/organisms/ScheduleBackupModal.css";

const ScheduleBackupModal = ({ databaseId, onClose }) => {
  const [cronExpr, setCronExpr] = useState("* * * * *");
  const [loading, setLoading] = useState(false);

  const handleSave = async () => {
    setLoading(true);
    try {
      const response = await addScheduledTask({
        database_id: databaseId,
        cron_expression: cronExpr,
      });
      if (response.success) {
        toast.success("Tâche CRON enregistrée avec succès !");
        onClose();
      } else {
        toast.error(
          response.message || "Erreur lors de la création de la tâche CRON"
        );
      }
    } catch (err) {
      console.error(err);
      toast.error("Erreur lors de la création de la tâche CRON");
    }
    setLoading(false);
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-card" onClick={(e) => e.stopPropagation()}>
        <h3>Planifier un backup automatique</h3>

        <CronSelector onChange={(value) => setCronExpr(value)} />

        <div className="modal-actions">
          <button
            className="modal-btn primary"
            onClick={handleSave}
            disabled={loading}
          >
            {loading ? "Enregistrement..." : "Enregistrer"}
          </button>

          <button className="modal-btn secondary" onClick={onClose}>
            Annuler
          </button>
        </div>
      </div>
    </div>
  );
};

export default ScheduleBackupModal;
