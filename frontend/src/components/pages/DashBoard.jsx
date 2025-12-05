// Ce composant affiche le tableau de bord principal et gère l'ouverture du formulaire d'ajout de base de données
import { useState } from "react";
import AddDatabase from "../organisms/AddDatabase";
import DatabaseBox from "../organisms/DatabaseBox";
import BackupBox from "../organisms/BackupBox";
import AlertBox from "../organisms/AlertBox";

const Home = () => {
  const [showForm, setShowForm] = useState(false);

  return (
    <>
      <h1>Welcome Home</h1>
      <DatabaseBox />
      <BackupBox />
      <AlertBox />
      {!showForm && (
        <button onClick={() => setShowForm(true)}>+ Add Database</button>
      )}

      {showForm && <AddDatabase onClose={() => setShowForm(false)} />}
    </>
  );
};

export default Home;
