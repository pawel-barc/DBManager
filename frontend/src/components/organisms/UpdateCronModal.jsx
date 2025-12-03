// Ce composant permet de modifier l'expression CRON d'une tâche existante.
// Il utilise CronSelector pour générer une nouvelle expression.

import { useState } from "react";
import { toast } from "react-toastify";
import CronSelector from "../molecules/CronSelector";
import { updateCronExpression } from "../../api/cronApi";

const UpdateCronModal = ({ task, onClose, onUpdated }) => {
  const [cron, setCron] = useState(task.cron_expression);
  const [saving, setSaving] = useState(false);

  const handleSave = async () => {
    setSaving(true);

    const res = await updateCronExpression(task.id, cron);

    setSaving(false);

    if (res.success) {
      toast.success("Expression CRON mise à jour !");
      onUpdated?.();
      onClose?.();
    } else {
      toast.error("Erreur lors de la mise à jour");
    }
  };

  return (
    <div className="modal">
      <div className="modal-content">
        <h3>Modifier la planification</h3>

        <CronSelector
          onChange={(expr) => setCron(expr)}
          initialCron={task.cron_expression}
        />

        <div style={{ marginTop: "20px" }}>
          <button
            onClick={handleSave}
            disabled={saving}
            style={{ marginRight: "10px", background: "green", color: "white" }}
          >
            {saving ? "Sauvegarde..." : "Enregistrer"}
          </button>

          <button
            onClick={onClose}
            style={{ background: "gray", color: "white" }}
          >
            Annuler
          </button>
        </div>
      </div>
    </div>
  );
};

export default UpdateCronModal;
