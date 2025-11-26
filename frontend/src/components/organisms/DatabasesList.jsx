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
    <div>
      <h2>Mes bases de données</h2>
      {databases.length === 0 ? (
        <p>Aucune base enregistrée</p>
      ) : (
        <div style={{ display: "flex", gap: "10px", flexWrap: "wrap" }}>
          {databases.map((db) => (
            <div
              key={db.id}
              onClick={() => goToDetails(db.id)}
              style={{
                border: "1px solid #ccc",
                padding: "15px",
                borderRadius: "10px",
                cursor: "pointer",
                width: "200px",
              }}
            >
              <strong>{db.name}</strong>
              <p>{db.type}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default DatabasesList;
