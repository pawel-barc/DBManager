const Input = ({ type, name, id, value, onChange, onBlur }) => {
  return (
    <input
      type={type}
      name={name}
      id={id}
      value={value}
      onChange={onChange}
      onBlur={onBlur}
      className="input-field"
    />
  );
};

export default Input;
