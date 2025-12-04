// Cette page affiche la liste des bases et permet d'ouvrir ou fermer le formulaire pour en akouter une nouvelle.
import DatabasesList from "../organisms/DatabasesList";
import AddDatabase from "../organisms/AddDatabase";
import { useState } from "react";
const Databases = () => {
  const [showForm, setShowForm] = useState(false);
  return (
    <>
      <DatabasesList />
      {!showForm && (
        <button onClick={() => setShowForm(true)}>+ Add Database</button>
      )}

      {showForm && <AddDatabase onClose={() => setShowForm(false)} />}
    </>
  );
};
export default Databases;
