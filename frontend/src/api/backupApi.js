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

// API pour récupérer tout les backups
const getAllBackups = async () => {
  const request = await fetchWithRefresh(`http://localhost:8080/backups`, {
    method: "GET",
    headers: { "Content-Type": "application/json" },
  });
  const response = await request.json();
  return response;
};

// API pour récupérer les sauvegardes existantes
const getDatabaseBackups = async (databaseId) => {
  const request = await fetchWithRefresh(
    `http://localhost:8080/backups/${databaseId}`,
    {
      method: "GET",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
    }
  );
  const response = await request.json();
  return response;
};

// API pour supprimer un backup
const deleteBackup = async (backupId) => {
  const request = await fetchWithRefresh(
    `http://localhost:8080/backups/${backupId}`,
    {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
    }
  );
  const response = await request.json();
  return response;
};

export { createBackup, getAllBackups, getDatabaseBackups, deleteBackup };
