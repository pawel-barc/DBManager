// Composant affichant une fenêtre modale permettant à l'utilisateur de programmer un backup récurrent pour une base de données.
import { useState } from "react";
import CronSelector from "./CronSelector";
import addScheduledTask from "../../api/cronApi";
import { toast } from "react-toastify";

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
    <div className="modal">
      <div className="modal-content">
        <h3>Planifier un backup automatique</h3>
        <CronSelector onChange={(value) => setCronExpr(value)} />
        <div style={{ marginTop: "20px", display: "flex", gap: "10px" }}>
          <button onClick={handleSave} disabled={loading}>
            {loading ? "Enregistrement..." : "Enregistrer"}
          </button>
          <button
            style={{ backgroundColor: "gray", color: "white" }}
            onClick={onClose}
          >
            Annuler
          </button>
        </div>
      </div>
    </div>
  );
};

export default ScheduleBackupModal;
