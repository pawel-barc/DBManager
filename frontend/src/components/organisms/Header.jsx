// src/components/organisms/Header.jsx
import { Link, useNavigate } from "react-router-dom";
import useAuthStore from "../../store/AuthStore";
import Logout from "./Logout";

const Header = () => {
  const { isAuthenticated, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <header style={styles.header}>
      <h1 style={styles.title}>SafeBase</h1>
      {isAuthenticated && (
        <nav>
          <Link to="/profile" style={styles.link}>
            Profile
          </Link>
          <button onClick={handleLogout} style={styles.button}>
            <Logout />
          </button>
        </nav>
      )}
    </header>
  );
};

const styles = {
  header: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    padding: "10px 20px",
    backgroundColor: "#282c34",
    color: "white",
  },
  title: {
    margin: 0,
  },
  link: {
    color: "white",
    marginRight: "15px",
    textDecoration: "none",
  },
  button: {
    padding: "5px 10px",
    cursor: "pointer",
  },
};

export default Header;
