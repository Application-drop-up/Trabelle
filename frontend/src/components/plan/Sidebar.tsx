import type { ReactNode } from "react";

type Props = {
  children: ReactNode;
};

export function Sidebar({ children }: Props) {
  return (
    <aside className="hidden h-full w-80 flex-shrink-0 flex-col border-r border-zinc-200 bg-white md:flex">
      {children}
    </aside>
  );
}
