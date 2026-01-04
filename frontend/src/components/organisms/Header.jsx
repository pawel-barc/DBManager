import { Link, NavLink, useNavigate } from "react-router-dom";
import useAuthStore from "../../store/AuthStore";
import Logout from "./Logout";

import logo from "../../assets/img/logos.png";
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
      <div className="sidebar-logged__logo">
        <img src={logo} alt="SafeBase logo" />
        <h1 className="sidebar-logged__title">SafeBase</h1>
      </div>

      <nav className="sidebar-logged__nav">
        <NavLink
          to="/"
          end
          className={({ isActive }) =>
            `sidebar-logged__link ${isActive ? "active" : ""}`
          }
        >
          Dashboard
        </NavLink>

        <NavLink
          to="/profile"
          className={({ isActive }) =>
            `sidebar-logged__link ${isActive ? "active" : ""}`
          }
        >
          Profile
        </NavLink>

        <NavLink
          to="/databases"
          className={({ isActive }) =>
            `sidebar-logged__link ${isActive ? "active" : ""}`
          }
        >
          Databases
        </NavLink>

        <NavLink
          to="/alerts"
          className={({ isActive }) =>
            `sidebar-logged__link ${isActive ? "active" : ""}`
          }
        >
          Alerts
        </NavLink>

        <button className="sidebar-logged__logout" onClick={handleLogout}>
          <Logout />
        </button>
      </nav>
    </aside>
  );
};

export default HeaderLogged;
