import { useEffect, useState } from "react";
import { getDatabases } from "../../api/databaseApi";

const DatabasesList = () => {
  const [databases, setDatabases] = useState([]);

  useEffect(() => {
    const load = async () => {
      const response = await getDatabases();
      if (response.success) {
        setDatabases(response.data);
      }
    };
    load();
  }, []);

  return (
    <div>
      <h2>Mes bases de données</h2>
      {databases.length === 0 ? (
        <p>Aucun base enregistrée</p>
      ) : (
        <ul>
          {databases.map((db) => (
            <li key={db.id}>
              <strong>{db.name}</strong> - {db.type}
              <br />
              {db.host} : {db.port}
              <br />
              User: {db.db_username}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};
export default DatabasesList;
