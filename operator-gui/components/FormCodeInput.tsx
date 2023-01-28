interface FormInputProps {
  id: string;
  label: string;
  value: string;
  onChange: (newValue: string) => void;
  placeholder?: string;
}

export default function FormCodeInput(props: FormInputProps) {
  const { id, label, value, onChange, placeholder } = props;

  return(
    <div>
      <label htmlFor={id}>{label}</label>
      <br />
      <textarea
        id={id}
        placeholder={placeholder}
        value={value}
        onChange={e => onChange(e.target.value)}
        className="mt-3 w-full h-96 rounded-md font-mono focus:ring-rose-500 focus:border-rose-500" />
    </div>
  );
}
