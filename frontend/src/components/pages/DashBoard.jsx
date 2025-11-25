import { useState } from "react";
import AddDatabase from "./Databases";

const Home = () => {
  const [showForm, setShowForm] = useState(false);

  return (
    <>
      <h1>Welcome Home</h1>

      {!showForm && (
        <button onClick={() => setShowForm(true)}>+ Add Database</button>
      )}

      {showForm && <AddDatabase onClose={() => setShowForm(false)} />}
    </>
  );
};

export default Home;
