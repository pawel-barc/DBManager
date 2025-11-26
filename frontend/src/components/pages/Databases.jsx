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
