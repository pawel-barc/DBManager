import AddDatabase from "./Databases";
import "../../styles/organisms/Dashboard.css";
const Home = () => {
  return (
    <div className="dashboard">
      <h1>Welcome Home</h1>
      <AddDatabase />
    </div>
  );
};
export default Home;
