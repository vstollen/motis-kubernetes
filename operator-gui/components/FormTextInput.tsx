interface FormInputProps {
  id: string;
  type: "text" | "url";
  label?: string;
  placeholder?: string;
}

export default function FormTextInput(props: FormInputProps) {
  const { id, label, type, placeholder } = props;

  return (
    <div>
      {label && (
        <>
          <label htmlFor={id}>{label}</label>
          <br />
        </>
      )}
      <input type={type} id={id} placeholder={placeholder} className="mt-3 w-full rounded-md focus:ring-rose-500 focus:border-rose-500" />
    </div>
  );
}
