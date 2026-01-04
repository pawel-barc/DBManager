// Ce composant affiche le tableau de bord principal et gère l'ouverture du formulaire d'ajout de base de données
import { useState } from "react";
import AddDatabase from "../organisms/AddDatabase";
import DatabaseBox from "../organisms/DatabaseBox";
import BackupBox from "../organisms/BackupBox";
import AlertBox from "../organisms/AlertBox";
import useAuthStore from "../../store/AuthStore";
import "../../styles/pages/DashBoard.css";
import add from "../../assets/img/add.png";

const Home = () => {
  const [showForm, setShowForm] = useState(false);
  const currentUser = useAuthStore((state) => state.currentUser);
  const username = currentUser?.name;

  return (
    <>
      <h1 className="welcome-text">Bonjour, {username}</h1>
      <main className="dashboard-main">
        <DatabaseBox />
        <BackupBox />
        <AlertBox />
      </main>

      {!showForm && (
        <div className="add-db-div">
          <button
            className="add-db-btn"
            onClick={() => setShowForm(true)}
            aria-label="Add database"
          >
            <img src={add} alt="Ajout une base" className="add-db-btn__icon" />
          </button>
          <p>Ajouter une nouvelle base des donnés</p>
        </div>
      )}
      {showForm && <AddDatabase onClose={() => setShowForm(false)} />}
    </>
  );
};

export default Home;
