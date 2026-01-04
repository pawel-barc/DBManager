import { Outlet } from "react-router-dom";
import HeaderUnlogged from "../components/organisms/HeaderUnlogged";
import "../styles/layouts/PublicLayout.css";
//Ce composant sert de structure de base pour les pages nécessitant une authentification
const PublicLayout = () => {
  return (
    <div className="layout-container">
      <HeaderUnlogged />
      <main className="main-content">
        <h1 className="app-name">SAFEBASE</h1>
        <Outlet />
      </main>
    </div>
  );
};

export default PublicLayout;
