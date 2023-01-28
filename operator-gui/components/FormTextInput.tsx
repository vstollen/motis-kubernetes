interface FormInputProps {
  id: string;
  type: "text" | "url";
  value: string;
  onChange: (newValue: string) => void;
  label?: string;
  placeholder?: string;
}

export default function FormTextInput(props: FormInputProps) {
  const { id, label, value, onChange, type, placeholder } = props;

  return (
    <div>
      {label && (
        <>
          <label htmlFor={id}>{label}</label>
          <br />
        </>
      )}
      <input
        type={type}
        id={id}
        value={value}
        onChange={e => onChange(e.target.value)}
        placeholder={placeholder}
        className="mt-3 w-full rounded-md focus:border-rose-500 focus:ring-rose-500"
      />
    </div>
  );
}
