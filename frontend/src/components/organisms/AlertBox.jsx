// Ce composposant récupère le nombre total des alertes et le nombre de notifications non lues
// via l'API. Il affiche également un indicateur (pastille rouge) lorsqu'il existe des alertes
// non lues. Un clic redirige l'utilisateur vers la liste complète des notifications.
import { useEffect, useState } from "react";
import { getAllUserAlerts } from "../../api/alertApi";
import { useNavigate } from "react-router-dom";
import { toast } from "react-toastify";

const AlertBox = () => {
  const [count, setCount] = useState(0);
  const navigate = useNavigate();
  const [unreadCount, setUnreadCount] = useState(0);

  useEffect(() => {
    const fetchAlerts = async () => {
      try {
        const response = await getAllUserAlerts();
        if (response.success && Array.isArray(response.data)) {
          const alerts = response.data;
          setCount(alerts.length);
          const unread = alerts.filter((a) => !a.is_read).length;
          setUnreadCount(unread);
        } else {
          setCount(0);
          setUnreadCount(0);
        }
      } catch (err) {
        console.error("Erreur lors du chargement des alerts", err);
        toast.error("Erreur lors du chargement des alerts");
        setCount(0);
        setUnreadCount(0);
      }
    };
    fetchAlerts();
  }, []);

  const handleClick = () => {
    navigate("/alerts");
  };
  return (
    <div
      className="dashboard-box"
      onClick={handleClick}
      style={{
        cursor: "pointer",
        padding: "20px",
        border: "1px solid #ccc",
        borderRadius: "10px",
        textAlign: "center",
        boxShadow: "0 2px 5px rgba(0,0,0,0.1)",
        transition: "all 0.2s",
        position: "relative",
      }}
      onMouseEnter={(e) =>
        (e.currentTarget.style.boxShadow = "0 4px 10px rgba(0,0,0,0.2)")
      }
      onMouseLeave={(e) =>
        (e.currentTarget.style.boxShadow = "0 2px 5px rgba(0,0,0,0.1)")
      }
    >
      {unreadCount > 0 && (
        <div
          style={{
            position: "absolute",
            top: "10px",
            right: "10px",
            width: "14px",
            height: "14px",
            backgroundColor: "red",
            borderRadius: "50%",
            border: "2px solid white",
          }}
        ></div>
      )}
      <h3>Notifications</h3>
      <p style={{ fontSize: "2rem", fontWeight: "bold" }}>{count}</p>
      <small>Cliquer pour voir les Notifications</small>
    </div>
  );
};
export default AlertBox;
