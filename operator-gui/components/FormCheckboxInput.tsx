interface FormInputProps {
  id: string;
  label: string;
  checked: boolean;
  onChange: () => void;
}

export default function FormCheckboxInput(props: FormInputProps) {
  const { id, label, checked, onChange } = props;

  return (
    <div className="flex items-center">
      <input
        type="checkbox"
        id={id}
        className="mr-3 h-4 w-4 rounded-md text-rose-600 focus:ring-rose-500"
        checked={checked}
        onChange={onChange}
      />
      <label htmlFor={id}>{label}</label>
    </div>
  );
}
