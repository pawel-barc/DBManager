// Ce composant affiche un formulaire permettant de tester une connexion à une base de données puis de l'enregistrer si le test réussit.
import { useState } from "react";
import testConnectionApi from "../../api/testConnection";
import { toast } from "react-toastify";
import "../../styles/organisms/AddDatabase.css";
import { addDatabase } from "../../api/databaseApi";

const AddDatabase = ({ onClose }) => {
  const [form, setForm] = useState({
    name: "",
    type: "postgres",
    host: "",
    port: 5432,
    db_username: "",
    db_password: "",
  });

  const [connectionOk, setConnectionOk] = useState(false);
  const [loadingTest, setLoadingTest] = useState(false);
  const [loadingAdd, setLoadingAdd] = useState(false);

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
    setConnectionOk(false);
  };

  const handleTest = async () => {
    setLoadingTest(true);
    const response = await testConnectionApi(form);
    setLoadingTest(false);

    if (response.success) {
      toast.success("Connexion testée avec succès !");
      setConnectionOk(true);
    } else {
      toast.error(response.message || "Test échoué");
      setConnectionOk(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!connectionOk) {
      toast.error("Veuillez tester la connexion avant d'ajouter !");
      return;
    }

    setLoadingAdd(true);
    const apiResponse = await addDatabase(form);
    setLoadingAdd(false);

    if (apiResponse.success) {
      toast.success("Base ajoutée avec succès");
      onClose?.();
    } else {
      toast.error(apiResponse.message || "Erreur");
    }
  };
  return (
    <div className="modal-overlay">
      <div className="modal">
        <button className="close-btn" onClick={onClose}>
          ×
        </button>
        <form onSubmit={handleSubmit}>
          <h2>Add Database</h2>

          <input
            name="name"
            placeholder="Nom de la base"
            value={form.name}
            onChange={handleChange}
            required
          />

          <select name="type" value={form.type} onChange={handleChange}>
            <option value="postgres">PostgreSQL</option>
            <option value="mysql">MySQL</option>
          </select>

          <input
            name="host"
            placeholder="Host"
            value={form.host}
            onChange={handleChange}
            required
          />

          <input
            type="number"
            name="port"
            placeholder="Port"
            value={form.port}
            onChange={handleChange}
            required
          />

          <input
            name="db_username"
            placeholder="Nom d'utilisateur"
            value={form.db_username}
            onChange={handleChange}
            required
          />

          <input
            name="db_password"
            placeholder="Mot de passe"
            value={form.db_password}
            onChange={handleChange}
            required
          />

          <button type="button" onClick={handleTest} disabled={loadingTest}>
            {loadingTest ? "Testing..." : "Test Connection"}
          </button>

          <button type="submit" disabled={!connectionOk || loadingAdd}>
            {loadingAdd ? "Adding..." : "Add Database"}
          </button>

          <button type="button" onClick={onClose}>
            Cancel
          </button>
        </form>
      </div>
    </div>
  );
};

export default AddDatabase;
