import { NavLink } from "react-router-dom";
import Logout from "./Logout";

import logo from "../../assets/img/logo.png";
import dbs from "../../assets/img/dbs.png";
import alert from "../../assets/img/alerts.png";
import profile from "../../assets/img/profile.png";
import logoutIcon from "../../assets/img/logout.png";

import "../../styles/organisms/HeaderLogged.css";

const HeaderLogged = () => {
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
      <div className="sidebar-logged__icon sidebar-logged__logout">
        <Logout />
      </div>
    </aside>
  );
};

export default HeaderLogged;
