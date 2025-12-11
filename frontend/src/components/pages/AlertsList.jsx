// Cette page récupère et affiche la liste des notifications d'un utilisateur.
// Un clic sur une notification la marque immédiatement comme "lue" (is_read = true)
// et met à jour son style ainsi que son statut dans la base de données.
import { useEffect, useState } from "react";
import {
  getAllUserAlerts,
  markAlertAsRead,
  markAllAsRead,
  deleteAlert,
} from "../../api/alertApi";
import { toast } from "react-toastify";

const AlertListPage = () => {
  const [alerts, setAlerts] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchAlerts = async () => {
    setLoading(true);

    const response = await getAllUserAlerts();
    setLoading(false);

    // structure backend : { success: true, data:[...] }
    if (response.success) {
      setAlerts(response.data || []);
      console.log(response.data);
    } else {
      toast.error(
        response.message || "Erreur lors du chargement des notifications"
      );
    }
  };

  useEffect(() => {
    fetchAlerts();
  }, []);

  const handleClickAlert = async (alertId) => {
    try {
      const response = await markAlertAsRead(alertId);

      if (!response.success) {
        toast.error("Impossible de marquer comme lu.");
        return;
      }

      setAlerts((prev) =>
        prev.map((a) => (a.id === alertId ? { ...a, is_read: true } : a))
      );
    } catch (err) {
      console.error(err);
      toast.error("Erreur lors du traitement de la notification.");
    }
  };

  const handleMarkAllRead = async () => {
    const response = await markAllAsRead();
    if (response.success) {
      setAlerts((prev) => prev.map((a) => ({ ...a, is_read: true })));
    } else {
      toast.error("Impossible de marquer toutes les notifications commes lues");
    }
  };

  const handleDeleteAlert = async (alertId) => {
    const response = await deleteAlert(alertId);
    if (response.success) {
      setAlerts((prev) => prev.filter((a) => a.id != alertId));
      toast.success("Alerte supprimée avec succès");
    } else {
      toast.error("Erreur lors de la suppression");
    }
  };

  if (loading) return <p>Chargement...</p>;

  return (
    <div>
      <h2>Notifications</h2>
      <button onClick={handleMarkAllRead}>Tout marquer comme lu</button>

      {alerts.length === 0 ? (
        <p>Aucune notification.</p>
      ) : (
        <ul style={{ padding: 0, listStyle: "none" }}>
          {alerts.map((a) => (
            <li
              key={a.id}
              onClick={() => handleClickAlert(a.id)}
              style={{
                marginBottom: "10px",
                padding: "10px",
                borderRadius: "6px",
                cursor: "pointer",
                backgroundColor: a.is_read ? "#fff" : "#e9f3ff",
                border: "1px solid #ddd",
                position: "relative",
              }}
            >
              <strong>{a.alert_type}</strong> — {a.message}
              <br />
              <small style={{ color: "#666" }}>{a.created_at}</small>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  handleDeleteAlert(a.id);
                }}
                style={{
                  position: "absolute",
                  border: "2px solid black",
                  borderRadius: "4px ",
                  bottom: "2px",
                  right: "-80px",
                  background: "transparent",
                  cursor: "pointer",
                  fontSize: "26px",
                  color: "#900",
                  fontWeight: "bold",
                }}
              >
                x
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default AlertListPage;
