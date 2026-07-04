"use client";

import { useRef, useState } from "react";
import type { ReactNode } from "react";

type SnapPoint = "peek" | "half" | "full";

const SNAP_TRANSLATE: Record<SnapPoint, string> = {
  peek: "calc(100% - 5rem)",
  half: "45%",
  full: "0%",
};

const SNAP_ORDER: SnapPoint[] = ["peek", "half", "full"];

type Props = {
  children: ReactNode;
};

export function BottomSheet({ children }: Props) {
  const [snap, setSnap] = useState<SnapPoint>("peek");
  const dragStartY = useRef<number | null>(null);

  const handlePointerDown = (e: React.PointerEvent<HTMLDivElement>) => {
    dragStartY.current = e.clientY;
    e.currentTarget.setPointerCapture(e.pointerId);
  };

  const handlePointerUp = (e: React.PointerEvent<HTMLDivElement>) => {
    if (dragStartY.current === null) return;
    const delta = dragStartY.current - e.clientY;
    const idx = SNAP_ORDER.indexOf(snap);

    if (Math.abs(delta) < 8) {
      setSnap(SNAP_ORDER[(idx + 1) % SNAP_ORDER.length]);
    } else if (delta > 40 && idx < SNAP_ORDER.length - 1) {
      setSnap(SNAP_ORDER[idx + 1]);
    } else if (delta < -40 && idx > 0) {
      setSnap(SNAP_ORDER[idx - 1]);
    }

    dragStartY.current = null;
  };

  return (
    <div
      className="absolute inset-x-0 bottom-0 h-full transition-transform duration-300 md:hidden"
      style={{ transform: `translateY(${SNAP_TRANSLATE[snap]})` }}
    >
      <div className="flex h-full flex-col rounded-t-2xl bg-white shadow-[0_-4px_24px_rgba(0,0,0,0.10)]">
        <div
          className="flex cursor-grab touch-none justify-center pb-2 pt-3 active:cursor-grabbing"
          onPointerDown={handlePointerDown}
          onPointerUp={handlePointerUp}
        >
          <div className="h-1 w-10 rounded-full bg-zinc-300" />
        </div>
        <div className="flex-1 overflow-hidden">{children}</div>
      </div>
    </div>
  );
}
