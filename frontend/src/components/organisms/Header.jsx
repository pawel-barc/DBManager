import { NavLink, useNavigate } from "react-router-dom";
import useAuthStore from "../../store/AuthStore";
import Logout from "./Logout";
import "../../../public/styles/header.css";

const iconSrc = (base, isActive) =>
  `../../../public/icons/${base}_${isActive ? "active" : "inactive"}.svg`;

const navItems = [
  { to: "/", base: "logo", alt: "logo" },
  { to: "/backups", base: "manage_db", alt: "manage database" },
  { to: "/alerts", base: "alerts", alt: "alerts" },
  { to: "/profile", base: "profile", alt: "profile" },
];

const Header = () => {
  const { isAuthenticated, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <header style={styles.header}>
      {isAuthenticated && (
        <nav className="nav-header">
          {navItems.map(({ to, base, alt }) => (
            <NavLink key={to} to={to} style={styles.link} aria-label={alt}>
              {({ isActive }) => (
                <img src={iconSrc(base, isActive)} alt={alt} />
              )}
            </NavLink>
          ))}

          <button onClick={handleLogout} style={styles.button} aria-label="logout">
            <Logout />
          </button>
        </nav>
      )}
    </header>
  );
};

const styles = {
};

export default Header;
