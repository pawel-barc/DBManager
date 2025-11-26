import { useState } from "react";
import DatabaseDetails from "./DatabaseDetails";

const DatabaseTile = ({ db, onDeleted }) => {
  const [showDetails, setShowDetails] = useState(false);

  return (
    <>
      <div
        onClick={() => setShowDetails(true)}
        style={{
          padding: "20px",
          margin: "10px",
          border: "1px solid #ccc",
          borderRadius: "10px",
          cursor: "pointer",
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
        <h3>{db.name}</h3>
        <small>{db.type}</small>
      </div>

      {showDetails && (
        <DatabaseDetails
          db={db}
          onClose={() => setShowDetails(false)}
          onDeleted={onDeleted}
        />
      )}
    </>
  );
};

export default DatabaseTile;
