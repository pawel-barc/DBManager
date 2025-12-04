import fetchWithRefresh from "./fetchWithRefresh";
// L'ajout de la base des données
const addDatabase = async (values) => {
  const request = await fetchWithRefresh(
    "http://localhost:8080/databases/add",
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(values),
      credentials: "include",
    }
  );
  const response = await request.json();
  return response;
};

// Récupérer les bases des données d'un utilisateur
const getDatabases = async () => {
  const request = await fetchWithRefresh(
    "http://localhost:8080/databases/list",
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    }
  );
  const response = await request.json();
  return response;
};
// Suppression d'une base de données
const deleteDatabase = async (id) => {
  const request = await fetchWithRefresh(
    `http://localhost:8080/databases/${id}`,
    {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    }
  );
  const response = await request.json();
  return response;
};

export { addDatabase, getDatabases, deleteDatabase };
