//Ce composant sert de structure de base pour les pages nécessitant une authentification
import { Outlet } from "react-router-dom";
import HeaderLogged from "../components/organisms/HeaderLogged";
import "../styles/layouts/PrivateLayout.css";

const PrivateLayout = () => {
  return (
    <div className="layout-container2">
      <HeaderLogged />
      <main className="main-content2">
        <h1 className="app-name">SAFEBASE</h1>
        <Outlet />
      </main>
    </div>
  );
};

export default PrivateLayout;
