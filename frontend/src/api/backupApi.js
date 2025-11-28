import fetchWithRefresh from "./fetchWithRefresh";
// API pour exécuter une sauvegarde
const createBackup = async (databaseId, name, version = "v1") => {
  const request = await fetchWithRefresh(
    `http://localhost:8080/backups/${databaseId}/create`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, version }),
      credentials: "include",
    }
  );
  const response = await request.json();
  return response;
};
