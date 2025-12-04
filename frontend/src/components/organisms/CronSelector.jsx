// Sélecteur interactif permettant de construire ou modifier une expression CRON.
// Permet aussi de charger un CRON existant (mode édition).

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

const CronSelector = ({ value, onChange }) => {
  const [mode, setMode] = useState("daily");
  const [hour, setHour] = useState("0");
  const [minute, setMinute] = useState("0");
  const [weekDays, setWeekDays] = useState([]);
  const [customCron, setCustomCron] = useState("* * * * *");

  const [generated, setGenerated] = useState("* * * * *");

  //
  // --- 1) CHARGEMENT D’UNE EXPRESSION CRON EXISTANTE ---
  //
  useEffect(() => {
    if (!value) return;

    const parts = value.split(" ");
    if (parts.length !== 5) return;

    const [m, h, , , d] = parts;

    setCustomCron(value);
    setGenerated(value);

    if (d !== "*" && d !== "?" && d !== "") {
      setMode("weekly");
      setWeekDays(d.split(","));
    } else if (h !== "*") {
      setMode("daily");
      setHour(h);
      setMinute(m);
    } else if (h === "*" && m !== "*") {
      setMode("hourly");
      setMinute(m);
    } else {
      setMode("custom");
    }
  }, [value]);

  //
  // --- 2) GÉNÉRATION DE L’EXPRESSION CRON ---
  //
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
        expr = `${minute} ${hour} * * ${
          weekDays.length > 0 ? weekDays.join(",") : "*"
        }`;
        break;
      case "custom":
        expr = customCron;
        break;
    }

    setGenerated(expr);
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
        CRON généré : <strong>{generated}</strong>
      </div>
    </div>
  );
};

export default CronSelector;
