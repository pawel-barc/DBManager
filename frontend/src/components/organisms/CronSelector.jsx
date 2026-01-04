// Sélecteur interactif permettant de construire ou modifier une expression CRON.
// Permet aussi de charger un CRON existant (mode édition).

import { useState, useEffect } from "react";

const DAYS = [
  { label: "Tous les jours", value: "*" },
  { label: "Lundi", value: "1" },
  { label: "Mardi", value: "2" },
  { label: "Mercredi", value: "3" },
  { label: "Jeudi", value: "4" },
  { label: "Vendredi", value: "5" },
  { label: "Samedi", value: "6" },
  { label: "Dimanche", value: "0" },
];

const CronSelector = ({ value, onChange }) => {
  const [mode, setMode] = useState("daily");
  const [hour, setHour] = useState("0");
  const [minute, setMinute] = useState("0");
  const [weekDay, setWeekDay] = useState("*");
  const [customCron, setCustomCron] = useState("* * * * *");
  const [generated, setGenerated] = useState("* * * * *");

  useEffect(() => {
    if (!value) return;

    const parts = value.split(" ");
    if (parts.length !== 5) return;

    const [m, h, , , d] = parts;

    setGenerated(value);
    setCustomCron(value);

    if (d !== "*" && d !== "?" && d !== "") {
      setMode("weekly");
      setMinute(m);
      setHour(h);
      setWeekDay(d); // ⚠️ un seul jour attendu
    } else if (h !== "*" && m !== "*") {
      setMode("daily");
      setMinute(m);
      setHour(h);
    } else if (h === "*" && m !== "*") {
      setMode("hourly");
      setMinute(m);
    } else {
      setMode("custom");
    }
  }, [value]);

  useEffect(() => {
    let expr = "* * * * *";

    switch (mode) {
      case "hourly":
        expr = `${minute} * * * *`;
        break;

      case "daily":
        expr = `${minute} ${hour} * * *`;
        break;

      case "weekly":
        expr = `${minute} ${hour} * * ${weekDay}`;
        break;

      case "custom":
        expr = customCron;
        break;
    }

    setGenerated(expr);
    onChange(expr);
  }, [mode, hour, minute, weekDay, customCron]);

  return (
    <div className="cron-selector">
      <div className="cron-field">
        <label>Type de planification</label>
        <select value={mode} onChange={(e) => setMode(e.target.value)}>
          <option value="hourly">Chaque heure</option>
          <option value="daily">Chaque jour</option>
          <option value="weekly">Jour spécifique</option>
          <option value="custom">Avancé (CRON)</option>
        </select>
      </div>

      {(mode === "hourly" || mode === "daily" || mode === "weekly") && (
        <div className="cron-time-row">
          <input
            type="number"
            min="0"
            max="59"
            value={minute}
            onChange={(e) => setMinute(e.target.value)}
            placeholder="Minute"
          />

          {(mode === "daily" || mode === "weekly") && (
            <input
              type="number"
              min="0"
              max="23"
              value={hour}
              onChange={(e) => setHour(e.target.value)}
              placeholder="Heure"
            />
          )}
        </div>
      )}

      {mode === "weekly" && (
        <div className="cron-field">
          <label>Jour de la semaine</label>
          <select value={weekDay} onChange={(e) => setWeekDay(e.target.value)}>
            {DAYS.map((d) => (
              <option key={d.value} value={d.value}>
                {d.label}
              </option>
            ))}
          </select>
          <span className="cron-hint">
            Une seule journée peut être sélectionnée
          </span>
        </div>
      )}

      {mode === "custom" && (
        <div className="cron-field">
          <label>Expression CRON</label>
          <input
            type="text"
            value={customCron}
            onChange={(e) => setCustomCron(e.target.value)}
            placeholder="* * * * *"
          />
        </div>
      )}

      <div className="cron-hint">
        CRON généré : <strong>{generated}</strong>
      </div>
    </div>
  );
};

export default CronSelector;
