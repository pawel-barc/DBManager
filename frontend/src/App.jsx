import Router from "./router/Router";
import "./App.css";
import Header from "./components/organisms/Header";
import { ToastContainer } from "react-toastify";
function App() {
  return (
    <>
      <Router />
      <ToastContainer
        position="top-right"
        autoClose={3000}
        hideProgressBar={false}
        newestOnTop={false}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
      />
    </>
  );
}

export default App;
