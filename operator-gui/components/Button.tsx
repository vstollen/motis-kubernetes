type ButtonProps = {
  children: JSX.Element | string;
  type: "primary" | "tertiary";
  onClick?: () => void;
};

const primaryButtonClasses = [
  "bg-rose-600",
  "text-[hsl(347_100%_99%)]",
  "hover:shadow-rose-500/50",
  "hover:bg-rose-500",
  "shadow",
  "shadow-rose-600/50",
  "hover:shadow-md",
  "active:shadow",
];

const secondaryButtonClasses = [
    "text-zinc-900",
    "hover:underline"
];

export default function Button({ children, type, onClick }: ButtonProps) {
  const buttonClasses = [
    "float-right",
    "mt-8",
    "ml-8",
    "rounded-md",
    "p-4",
    "py-3",
    "font-bold",
    "transition-all",
  ]

  switch (type) {
    case "primary":
      buttonClasses.push(...primaryButtonClasses);
      break;
    case "tertiary":
      buttonClasses.push(...secondaryButtonClasses);
      break;
  }

  return (
    <button
      className={buttonClasses.join(" ")}
      onClick={() => {
        if (onClick) {
          onClick();
        }
      }}
    >
      {children}
    </button>
  );
}
