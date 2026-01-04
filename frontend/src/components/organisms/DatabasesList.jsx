// Ce composant récupère toutes les base de données depuis l'API et affiche une liste cliquable permettant d'ouvrir les détails de chaque base.
import { useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";
import { getDatabases } from "../../api/databaseApi";

const DatabasesList = () => {
  const [databases, setDatabases] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    const load = async () => {
      try {
        const response = await getDatabases();
        if (response.success) {
          setDatabases(Array.isArray(response.data) ? response.data : []);
        }
      } catch (err) {
        console.error("Erreur lors du chargement des bases", err);
        setDatabases([]);
      }
    };
    load();
  }, []);
  const goToDetails = (id) => {
    navigate(`/databases/${id}`);
  };

  return (
    <div className="db-list">
      {databases.length === 0 ? (
        <p>Aucune base enregistrée</p>
      ) : (
        <div style={{ display: "flex", gap: "20px", flexWrap: "wrap" }}>
          {databases.map((db) => (
            <div
              key={db.id}
              onClick={() => goToDetails(db.id)}
              style={{
                border: "1px solid #ccc",
                padding: "15px",
                borderRadius: "10px",
                cursor: "pointer",
                minWidth: "180px",
                color: "#0353a4",
                backgroundColor: "white",
                textAlign: "center",
                fontSize:"20px"
              }}
            >
              <strong>{db.name}</strong>
              <p style={{fontSize:"16px"}}>{db.type}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default DatabasesList;
