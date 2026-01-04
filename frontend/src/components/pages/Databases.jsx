// Cette page affiche la liste des bases et permet d'ouvrir ou fermer le formulaire pour en akouter une nouvelle.
import DatabasesList from "../organisms/DatabasesList";
import AddDatabase from "../organisms/AddDatabase";
import add from "../../assets/img/add.png";
import "../../styles/pages/Databases.css";
import { useState } from "react";
const Databases = () => {
  const [showForm, setShowForm] = useState(false);
  return (
    <>
      {" "}
      <h2 className="intro-text">Mes bases de données</h2>
      {!showForm && (
        <div className="add-db-div2">
          <button
            className="add-db-btn"
            onClick={() => setShowForm(true)}
            aria-label="Add database"
          >
            <img src={add} alt="Ajout une base" className="add-db-btn__icon" />
          </button>
          <p>Ajouter une nouvelle base des données</p>
        </div>
      )}
      {showForm && <AddDatabase onClose={() => setShowForm(false)} />}
      <DatabasesList />
    </>
  );
};
export default Databases;
