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
export default getAllUserAlerts;
