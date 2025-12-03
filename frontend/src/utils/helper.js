// Fonction utilitaire qui convertit une expression CRON en texte lisible pour l'utilisateur,
// afin de faciliter la compréhension des horaires de backup automatique.
const DAYS_LABELS = {
  0: "Dimanche",
  1: "Lundi",
  2: "Mardi",
  3: "Mercredi",
  4: "Jeudi",
  5: "Vendredi",
  6: "Samedi",
};

export function humanReadableCron(cronExpr) {
  const parts = cronExpr.trim().split(" ");
  if (parts.length !== 5) return cronExpr;

  const [minute, hour, dayOfMonth, month, dayOfWeek] = parts;

  // semaine
  if (dayOfWeek !== "*" && dayOfMonth === "*" && month === "*") {
    const days = dayOfWeek.split(",").map((d) => DAYS_LABELS[d] || d);
    return `Chaque ${days.join(", ")} à ${hour.padStart(
      2,
      "0"
    )}:${minute.padStart(2, "0")}`;
  }

  // jour
  if (dayOfWeek === "*" && dayOfMonth === "*" && month === "*") {
    return `Chaque jour à ${hour.padStart(2, "0")}:${minute.padStart(2, "0")}`;
  }

  // heure
  if (
    dayOfWeek === "*" &&
    dayOfMonth === "*" &&
    month === "*" &&
    hour === "*"
  ) {
    return `Chaque heure à ${minute.padStart(2, "0")} min`;
  }

  return cronExpr;
}
