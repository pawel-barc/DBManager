import { useState } from "react";
import testConnectionApi from "../../api/testConnection";
import { toast } from "react-toastify";
import "../../styles/organisms/AddDatabase.css";

const AddDatabase = ({ onClose }) => {
  const [form, setForm] = useState({
    name: "",
    type: "postgres",
    host: "",
    port: 5432,
    db_username: "",
    db_password: "",
  });

  const [errorMessage, setErrorMessage] = useState(null);

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const apiResponse = await addDatabase(form);

      if (apiResponse.success) {
        toast.success("La base ajoutée avec succès");
        setForm({
          name: "",
          type: "postgres",
          host: "",
          port: 5432,
          db_username: "",
          db_password: "",
        });

        onClose?.();
      } else if (apiResponse.error) {
        setErrorMessage(apiResponse.error);
      }
    } catch (error) {
      console.error("Erreur capturée: ", error);
      setErrorMessage(error.message || "Une erreur est survenue");
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

          {errorMessage && <p style={{ color: "red" }}>{errorMessage}</p>}
          <button
            type="button"
            onClick={async () => {
              const res = await testConnectionApi(form);
              if (res.success) toast.success("Connection OK");
              else toast.error(res.error || "Error");
            }}
          >
            Test Connection
          </button>

          <button type="submit">Add</button>
          <button type="button" onClick={onClose}>
            Cancel
          </button>
        </form>
      </div>
    </div>
  );
};

export default AddDatabase;
