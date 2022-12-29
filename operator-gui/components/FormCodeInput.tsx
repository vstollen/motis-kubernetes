interface FormInputProps {
  id: string;
  label: string;
  placeholder?: string;
}

export default function FormCodeInput(props: FormInputProps) {
  const { id, label, placeholder } = props;

  return(
    <div>
      <label htmlFor={id}>{label}</label>
      <br />
      <textarea
        id={id}
        placeholder={placeholder}
        className="mt-3 w-full h-96 rounded-md font-mono focus:ring-rose-500 focus:border-rose-500" />
    </div>
  );
}
