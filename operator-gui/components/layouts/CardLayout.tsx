type CardLayoutProps = {
  children: JSX.Element;
};

export default function CardLayout({ children }: CardLayoutProps) {
  return (
    <>
      <main className="flex justify-center p-20">
        <div className="flex-auto max-w-3xl rounded-md bg-zinc-50 p-16">{children}</div>
      </main>
    </>
  );
}
