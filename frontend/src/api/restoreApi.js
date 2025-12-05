import fetchWithRefresh from "./fetchWithRefresh";

// API pour restaurer une base à partir d'un backup
const restoreBackup = async (backupId) => {
  const request = await fetchWithRefresh(
    `http://localhost:8080/backups/${backupId}/restore`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
    }
  );

  const response = await request.json();
  return response;
};

export default restoreBackup;
