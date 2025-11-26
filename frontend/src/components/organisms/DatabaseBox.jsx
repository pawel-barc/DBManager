import { useEffect, useState } from "react";
import { getDatabases } from "../../api/databaseApi";
import { useNavigate } from "react-router-dom";

const DatabaseBox = () => {
  const [count, setCount] = useState(0);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchDatabases = async () => {
      try {
        const res = await getDatabases();
        if (res.success) {
          setCount(res.data.length);
        }
      } catch (err) {
        console.error("Erreur lors du chargement des bases", err);
      }
    };
    fetchDatabases();
  }, []);

  const handleClick = () => {
    navigate("/databases");
  };

  return (
    <div
      className="dashboard-box"
      onClick={handleClick}
      style={{
        cursor: "pointer",
        padding: "20px",
        border: "1px solid #ccc",
        borderRadius: "10px",
        textAlign: "center",
        boxShadow: "0 2px 5px rgba(0,0,0,0.1)",
        transition: "all 0.2s",
      }}
      onMouseEnter={(e) =>
        (e.currentTarget.style.boxShadow = "0 4px 10px rgba(0,0,0,0.2)")
      }
      onMouseLeave={(e) =>
        (e.currentTarget.style.boxShadow = "0 2px 5px rgba(0,0,0,0.1)")
      }
    >
      <h3>Bases de données</h3>
      <p style={{ fontSize: "2rem", fontWeight: "bold" }}>{count}</p>
      <small>Cliquer pour voir la liste</small>
    </div>
  );
};

export default DatabaseBox;
