type CardLayoutProps = {
  children: JSX.Element;
};

export default function CardLayout({ children }: CardLayoutProps) {
  return (
    <div className="flex justify-center p-20">
      <main className="flex-auto rounded-md bg-zinc-50 p-16">
        {children}
      </main>
    </div>
  );
}
