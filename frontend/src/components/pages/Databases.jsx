import { useState } from "react";
import addDatabase from "../../api/databaseApi";
import { toast } from "react-toastify";

const AddDatabase = () => {
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
      } else if (apiResponse.error) {
        setErrorMessage(apiResponse.error);
      }
    } catch (error) {
      console.error("Erreur capturée: ", error);
      setErrorMessage(error.message || "Une erreur est survenue");
    }
  };

  return (
    <form onSubmit={handleSubmit}>
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
      <button type="submit">Ajouter la base</button>
    </form>
  );
};
export default AddDatabase;
