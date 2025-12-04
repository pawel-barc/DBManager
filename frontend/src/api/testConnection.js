import fetchWithRefresh from "./fetchWithRefresh";
// Api pour tester la connexion à la base de données
const testConnectionApi = async (values) => {
  const request = await fetchWithRefresh(
    "http://localhost:8080/databases/test",
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify(values),
    }
  );
  const response = await request.json();
  return response;
};

export default testConnectionApi;
