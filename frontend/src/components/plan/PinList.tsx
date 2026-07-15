"use client";

import { useState } from "react";

import type { PinViewModel } from "@/containers/PlanContainer";
import { PIN_CATEGORY_LABELS } from "@/domain/pins/labels";

type Props = {
  pins: PinViewModel[];
};

export function PinList({ pins }: Props) {
  const [openPinId, setOpenPinId] = useState<string | null>(null);

  if (pins.length === 0) {
    return <p className="text-sm text-zinc-400">ピンがまだありません</p>;
  }

  return (
    <ul className="flex flex-col gap-2">
      {pins.map((pin) => {
        const isOpen = openPinId === pin.id;
        return (
          <li key={pin.id} className="rounded-lg border border-zinc-200 bg-white">
            <button
              type="button"
              className="flex w-full items-center gap-3 p-3 text-left"
              onClick={() => setOpenPinId(isOpen ? null : pin.id)}
            >
              <span
                className="h-3 w-3 flex-shrink-0 rounded-full"
                style={{ backgroundColor: pin.colour }}
              />
              <span className="flex-1 text-sm font-medium text-zinc-900">{pin.name}</span>
              <span className="text-xs text-zinc-400">{PIN_CATEGORY_LABELS[pin.category]}</span>
              <svg
                className={`h-4 w-4 flex-shrink-0 text-zinc-400 transition-transform duration-200 ${isOpen ? "rotate-180" : ""}`}
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                strokeWidth={2}
              >
                <path strokeLinecap="round" strokeLinejoin="round" d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            {isOpen && (
              <div className="border-t border-zinc-100 px-3 pb-3 pt-2">
                {pin.notes.length === 0 ? (
                  <p className="text-xs text-zinc-400">メモなし</p>
                ) : (
                  <ul className="flex flex-col gap-1">
                    {pin.notes.map((note) => (
                      <li key={note.id} className="text-sm text-zinc-700">
                        {note.content}
                      </li>
                    ))}
                  </ul>
                )}
              </div>
            )}
          </li>
        );
      })}
    </ul>
  );
}
