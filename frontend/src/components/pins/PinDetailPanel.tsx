"use client";

import type { NoteViewModel, PinViewModel } from "@/mappers/planMapper";

type Props = {
  pinVM: PinViewModel;
};

function NoteItem({ note }: { note: NoteViewModel }) {
  return <li className="text-sm text-zinc-700">{note.content}</li>;
}

export function PinDetailPanel({ pinVM }: Props) {
  return (
    <div className="flex flex-col gap-4 p-4">
      <div className="flex items-center gap-3">
        <span
          className="h-4 w-4 flex-shrink-0 rounded-full"
          style={{ backgroundColor: pinVM.colour }}
        />
        <h2 className="text-base font-semibold text-zinc-900">{pinVM.name}</h2>
      </div>
      <div className="flex flex-col gap-1">
        <span className="text-xs text-zinc-500">Category</span>
        <span className="text-sm text-zinc-800">{pinVM.category}</span>
      </div>
      <div className="flex flex-col gap-1">
        <span className="text-xs text-zinc-500">Notes</span>
        {pinVM.notes.length === 0 ? (
          <p className="text-sm text-zinc-400">No notes yet</p>
        ) : (
          <ul className="flex flex-col gap-1">
            {pinVM.notes.map((note) => (
              <NoteItem key={note.id} note={note} />
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
