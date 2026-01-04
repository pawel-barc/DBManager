// src/router/Router.jsx
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import useAuthStore from "../store/AuthStore";

import Register from "../components/pages/Register";
import Login from "../components/pages/Login";
import Profile from "../components/pages/Profile";
import Dashboard from "../components/pages/DashBoard";
import Logout from "../components/organisms/Logout";
import Databases from "../components/pages/Databases";
import DBDetails from "../components/pages/DBDetails";
import BackupsList from "../components/pages/BackupsList";
import BackupsManagement from "../components/pages/BackupsManagement";
import DatabaseRestore from "../components/pages/DatabaseRestore";
import AlertListPage from "../components/pages/AlertsList";
import PublicLayout from "../layout/PublicLayout";
import PrivateLayout from "../layout/PrivateLayout"

const PrivateRoute = ({ children }) => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  return isAuthenticated ? children : <Navigate to="/login" replace />;
};

const PublicRoute = ({ children }) => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  return !isAuthenticated ? children : <Navigate to="/" replace />;
};

const Router = () => {
  return (
    <BrowserRouter>
      <Routes>
        {/* PUBLIC */}
        <Route
          element={
            <PublicRoute>
              <PublicLayout />
            </PublicRoute>
          }
        >
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
        </Route>

        {/* PRIVATE */}
        <Route
          element={
            <PrivateRoute>
              <PrivateLayout />
            </PrivateRoute>
          }
          >
          <Route path="/" element={<Dashboard />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/databases" element={<Databases />} />
          <Route path="/backups" element={<BackupsList />} />
          <Route path="/databases/:id" element={<DBDetails />} />
          <Route
            path="/databases/:id/backups"
            element={<BackupsManagement />}
          />
          <Route
            path="/databases/:id/restauration"
            element={<DatabaseRestore />}
          />
          <Route path="/alerts" element={<AlertListPage />} />
          <Route path="/logout" element={<Logout />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
