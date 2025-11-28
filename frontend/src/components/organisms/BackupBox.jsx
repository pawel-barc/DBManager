import { useEffect, useState } from "react";
import { getAllBackups } from "../../api/backupApi";
import { useNavigate } from "react-router-dom";
import { toast } from "react-toastify";

const BackupBox = ({ databaseId }) => {
  const [count, setCount] = useState(0);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchBackups = async () => {
      try {
        const response = await getAllBackups();
        if (response.success && Array.isArray(response.data)) {
          setCount(response.data.length);
        } else {
          setCount(0);
        }
      } catch (err) {
        console.error("Erreur lors du chargement des backups", err);
        toast.error("Erreur lors du chargement des backups");
        setCount(0);
      }
    };
    fetchBackups();
  }, []);

  const handleClick = () => {
    navigate(`/databases/${databaseId}/backups`);
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
      <h3>Backups</h3>
      <p style={{ fontSize: "2rem", fontWeight: "bold" }}>{count}</p>
      <small>Cliquer pour voir les backups</small>
    </div>
  );
};
export default BackupBox;
