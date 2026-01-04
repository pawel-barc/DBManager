import Header from "../components/organisms/Header";
import { Outlet } from "react-router-dom";
import "../styles/organisms/layouts/AppLayout.css";

const AppLayout = () => {
  return (
    <div className="app-layout">
      <Header />
      <main className="app-layout__content">
        <Outlet />
      </main>
    </div>
  );
};

export default AppLayout;
