type ButtonProps = {
  children: JSX.Element | string
  onClick?: () => void
}

export default function Button({ children, onClick }: ButtonProps) {
  return (
    <button
      className="float-right mt-8 rounded-md bg-rose-600 p-4 py-3 font-bold text-[hsl(347_100%_99%)] shadow shadow-rose-600/50 transition-all hover:bg-rose-500 hover:shadow-md hover:shadow-rose-500/50 active:shadow"
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
