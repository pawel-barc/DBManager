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

// API pour récupérer les tâches planifiées de l'utilisateur
const getUserScheduledTasks = async () => {
  const request = await fetchWithRefresh(
    "http://localhost:8080/scheduled-tasks",
    {
      method: "GET",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
    }
  );
  const response = await request.json();
  return response;
};

// API pour activer/désactiver une tâche planifiée
const toggleTaskActive = async (taskId, active) => {
  const request = await fetchWithRefresh(
    "http://localhost:8080/scheduled-tasks/toggle",
    {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ task_id: taskId, active }),
      credentials: "include",
    }
  );
  const response = await request.json();
  return response;
};

// API pour mettre à jour l'expression CRON d'une tâche
const updateCronExpression = async (taskId, cronExpression) => {
  const request = await fetchWithRefresh(
    "http://localhost:8080/scheduled-tasks/update-cron",
    {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id: taskId, cron_expression: cronExpression }),
      credentials: "include",
    }
  );
  const response = await request.json();
  return response;
};

export {
  addScheduledTask,
  getUserScheduledTasks,
  toggleTaskActive,
  updateCronExpression,
};
