import fetchWithRefresh from "./fetchWithRefresh";

const addScheduledTask = async (values) => {
  const request = await fetchWithRefresh("http://localhost:8080/cron/create", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(values),
    credentials: "include",
  });
  const response = await request.json();
  return response;
};

export default addScheduledTask;
