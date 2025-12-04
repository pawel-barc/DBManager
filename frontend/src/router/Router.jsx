// src/router/Router.jsx
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import useAuthStore from "../store/AuthStore";

import Register from "../components/pages/Register";
import Login from "../components/pages/Login";
import Profile from "../components/pages/Profile";
import Dashboard from "../components/pages/DashBoard";
import Logout from "../components/organisms/Logout";
import Header from "../components/organisms/Header";
import Databases from "../components/pages/Databases";
import DBDetails from "../components/pages/DBDetails";
import BackupsList from "../components/pages/BackupsList";
import BackupsManagement from "../components/pages/BackupsManagement";

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
      <Header />
      <Routes>
        {/* Public routes */}
        <Route
          path="/login"
          element={
            <PublicRoute>
              <Login />
            </PublicRoute>
          }
        />
        <Route
          path="/register"
          element={
            <PublicRoute>
              <Register />
            </PublicRoute>
          }
        />

        {/* Private routes */}
        <Route
          path="/"
          element={
            <PrivateRoute>
              <Dashboard />
            </PrivateRoute>
          }
        />
        <Route
          path="/profile"
          element={
            <PrivateRoute>
              <Profile />
            </PrivateRoute>
          }
        />
        <Route
          path="/databases"
          element={
            <PrivateRoute>
              <Databases />
            </PrivateRoute>
          }
        />
        <Route
          path="/backups"
          element={
            <PrivateRoute>
              <BackupsList />
            </PrivateRoute>
          }
        />
        <Route
          path="/databases/:id"
          element={
            <PrivateRoute>
              <DBDetails />
            </PrivateRoute>
          }
        />
        <Route
          path="/databases/:id/backups"
          element={
            <PrivateRoute>
              <BackupsManagement />
            </PrivateRoute>
          }
        />

        <Route
          path="/logout"
          element={
            <PrivateRoute>
              <Logout />
            </PrivateRoute>
          }
        />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
