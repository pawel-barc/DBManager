// Cette page récupère l'ID depuis l'URL, cherche la base correspondante dans la liste des bases, et affiche ses détails.
import { useParams, useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import { getDatabases } from "../../api/databaseApi";
import DatabaseDetails from "../organisms/DatabaseDetails";

const DatabaseDetailsPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [db, setDb] = useState(null);

  useEffect(() => {
    const load = async () => {
      const response = await getDatabases();
      if (response.success) {
        const found = response.data.find((item) => item.id === Number(id));
        if (found) setDb(found);
        else navigate("/databases");
      }
    };
    load();
  }, [id]);

  if (!db) return <p>Chargement...</p>;

  return (
    <DatabaseDetails
      db={db}
      onClose={() => navigate("/databases")}
      onDeleted={() => navigate("/databases")}
    />
  );
};

export default DatabaseDetailsPage;
