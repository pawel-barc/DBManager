import { useState } from "react";
import AddDatabase from "../organisms/AddDatabase";
import DatabaseBox from "../organisms/DatabaseBox";
import BackupBox from "../organisms/BackupBox";

const Home = () => {
  const [showForm, setShowForm] = useState(false);

  return (
    <>
      <h1>Welcome Home</h1>
      <DatabaseBox />
      <BackupBox />
      {!showForm && (
        <button onClick={() => setShowForm(true)}>+ Add Database</button>
      )}

      {showForm && <AddDatabase onClose={() => setShowForm(false)} />}
    </>
  );
};

export default Home;
