import fetchWithRefresh from "./fetchWithRefresh";

// Récupération des alertes d'un utilisateur
const getAllUserAlerts = async () => {
  const request = await fetchWithRefresh("http://localhost:8080/alerts", {
    method: "GET",
    credentials: "include",
  });
  const response = await request.json();
  return response;
};
// Marque une alerte comme lu
const markAlertAsRead = async (alertId) => {
  const request = await fetchWithRefresh(
    `http://localhost:8080/alerts/${alertId}/read`,
    {
      method: "PUT",
      credentials: "include",
    }
  );
  const response = await request.json();
  return response;
};

// Marque toutes les alertes comme lues
const markAllAsRead = async () => {
  const request = await fetchWithRefresh(
    `http://localhost:8080/alerts/read-all`,
    {
      method: "PUT",
      credentials: "include",
    }
  );
  const response = await request.json();
  return response;
};

// Supprime une alerte
const deleteAlert = async (alertId) => {
  const request = await fetchWithRefresh(
    `http://localhost:8080/alerts/${alertId}`,
    {
      method: "DELETE",
      credentials: "include",
    }
  );
  const response = await request.json();
  return response;
};
export { getAllUserAlerts, markAlertAsRead, markAllAsRead, deleteAlert };
