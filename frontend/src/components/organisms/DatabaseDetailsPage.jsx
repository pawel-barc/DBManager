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
      const res = await getDatabases();
      if (res.success) {
        const found = res.data.find((item) => item.id === Number(id));
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
