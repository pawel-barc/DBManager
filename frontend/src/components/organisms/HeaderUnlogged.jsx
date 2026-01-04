import { Link } from "react-router-dom";
import "../../styles/organisms/HeaderUnlogged.css";
import logo from "../../assets/img/logo.png";

const HeaderUnlogged = () => {
  return (
    <aside className="sidebar-unlogged">
      <div className="sidebar-unlogged__logo">
        <img className="logo" src={logo} alt="Safebase logo" />
      </div>
    </aside>
  );
};

export default HeaderUnlogged;
