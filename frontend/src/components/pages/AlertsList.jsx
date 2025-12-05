// Cette page récupère et affiche la liste des notifications d'un utilisateur.
import { useEffect, useState } from "react";
import getAllUserAlerts from "../../api/alertApi";
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
    } else {
      toast.error(
        response.message || "Erreur lors du chargement des notifications"
      );
    }
  };

  useEffect(() => {
    fetchAlerts();
  }, []);

  if (loading) return <p>Chargement...</p>;

  return (
    <div>
      <h2>Notifications</h2>

      {alerts.length === 0 ? (
        <p>Aucune notification.</p>
      ) : (
        <ul>
          {alerts.map((a) => (
            <li key={a.id} style={{ marginBottom: "10px" }}>
              <strong>{a.alert_type}</strong> — {a.message}
              <br />
              <small style={{ color: "#888" }}>{a.created_at}</small>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default AlertListPage;
