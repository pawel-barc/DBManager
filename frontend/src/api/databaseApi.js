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

export default addDatabase;
