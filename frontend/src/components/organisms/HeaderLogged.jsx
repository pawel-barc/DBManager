import { NavLink, useNavigate } from "react-router-dom";
import useAuthStore from "../../store/AuthStore";
import Logout from "./Logout";

import logo from "../../assets/img/logo.png";
import dbs from "../../assets/img/dbs.png";
import alert from "../../assets/img/alerts.png";
import profile from "../../assets/img/profile.png";
import logoutIcon from "../../assets/img/logout.png";

import "../../styles/organisms/HeaderLogged.css";

const HeaderLogged = () => {
  const logout = useAuthStore((state) => state.logout);
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <aside className="sidebar-logged">
      {/* LOGO */}
      <NavLink to="/" end className="sidebar-logged__logo">
        <img src={logo} alt="SafeBase logo" />
      </NavLink>

      {/* ICON NAV */}
      <nav className="sidebar-logged__icons">
        <NavLink
          to="/databases"
          className={({ isActive }) =>
            `sidebar-logged__icon ${isActive ? "active" : ""}`
          }
        >
          <img src={dbs} alt="Databases" />
        </NavLink>

        <NavLink
          to="/alerts"
          className={({ isActive }) =>
            `sidebar-logged__icon ${isActive ? "active" : ""}`
          }
        >
          <img src={alert} alt="Alerts" />
        </NavLink>

        <NavLink
          to="/profile"
          className={({ isActive }) =>
            `sidebar-logged__icon ${isActive ? "active" : ""}`
          }
        >
          <img src={profile} alt="Profile" />
        </NavLink>
      </nav>

      {/* LOGOUT */}
      <button
        className="sidebar-logged__icon sidebar-logged__logout"
        onClick={handleLogout}
        aria-label="Logout"
        title="Logout"
      >
        <img src={logoutIcon} alt="Logout" />
      </button>
    </aside>
  );
};

export default HeaderLogged;
