// Sélecteur interactif permettant de construire une expression CRON et d'en visualiser l’exécution programmée.
import { useState, useEffect } from "react";

const DAYS = [
  { label: "Lundi", value: "1" },
  { label: "Mardi", value: "2" },
  { label: "Mercredi", value: "3" },
  { label: "Jeudi", value: "4" },
  { label: "Vendredi", value: "5" },
  { label: "Samedi", value: "6" },
  { label: "Dimanche", value: "0" },
];

const CronSelector = ({ onChange }) => {
  const [mode, setMode] = useState("daily");
  const [hour, setHour] = useState("0");
  const [minute, setMinute] = useState("0");
  const [customCron, setCustomCron] = useState("* * * * *");
  const [weekDays, setWeekDays] = useState([]);

  useEffect(() => {
    let expr = "* * * * *";

    switch (mode) {
      case "hourly":
        expr = `${minute} * * * *`;
        break;

      case "daily":
        expr = `${minute} ${hour} * * *`;
        break;

      case "weekly": {
        const days = weekDays.length > 0 ? weekDays.join(",") : "*";
        expr = `${minute} ${hour} * * ${days}`;
        break;
      }

      case "custom":
        expr = customCron;
        break;
    }

    onChange(expr);
  }, [mode, hour, minute, weekDays, customCron]);

  const toggleDay = (value) => {
    setWeekDays((prev) =>
      prev.includes(value) ? prev.filter((d) => d !== value) : [...prev, value]
    );
  };

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
      <label>
        Type de planification:
        <select value={mode} onChange={(e) => setMode(e.target.value)}>
          <option value="hourly">Chaque heure</option>
          <option value="daily">Chaque jour</option>
          <option value="weekly">Jours spécifiques</option>
          <option value="custom">Avancé (CRON)</option>
        </select>
      </label>

      {(mode === "daily" || mode === "hourly" || mode === "weekly") && (
        <div style={{ display: "flex", gap: "10px" }}>
          <input
            type="number"
            min="0"
            max="59"
            value={minute}
            onChange={(e) => setMinute(e.target.value)}
            placeholder="Minute"
            style={{ width: "60px" }}
          />

          {(mode === "daily" || mode === "weekly") && (
            <input
              type="number"
              min="0"
              max="23"
              value={hour}
              onChange={(e) => setHour(e.target.value)}
              placeholder="Heure"
              style={{ width: "60px" }}
            />
          )}
        </div>
      )}

      {mode === "weekly" && (
        <div style={{ display: "flex", flexDirection: "column" }}>
          {DAYS.map((d) => (
            <label key={d.value}>
              <input
                type="checkbox"
                value={d.value}
                checked={weekDays.includes(d.value)}
                onChange={() => toggleDay(d.value)}
              />
              {d.label}
            </label>
          ))}
        </div>
      )}

      {mode === "custom" && (
        <input
          type="text"
          value={customCron}
          onChange={(e) => setCustomCron(e.target.value)}
          placeholder="* * * * *"
          style={{ width: "200px" }}
        />
      )}

      <div style={{ marginTop: "10px" }}>
        CRON généré : <strong>{customCron}</strong>
      </div>
    </div>
  );
};

export default CronSelector;
