"use client";

import { useState } from "react";

import type { NoteViewModel, PinViewModel } from "@/containers/PlanContainer";

type Props = {
  pins: PinViewModel[];
  onCreateNote: (pinId: string, content: string) => Promise<void>;
  onUpdateNote: (pinId: string, noteId: string, content: string) => Promise<void>;
};

export function PinList({ pins, onCreateNote, onUpdateNote }: Props) {
  const [openPinId, setOpenPinId] = useState<string | null>(null);

  if (pins.length === 0) {
    return <p className="text-sm text-zinc-400">ピンがまだありません</p>;
  }

  return (
    <ul className="flex flex-col gap-2">
      {pins.map((pin) => (
        <PinListItem
          key={pin.id}
          pin={pin}
          isOpen={openPinId === pin.id}
          onToggle={() => setOpenPinId(openPinId === pin.id ? null : pin.id)}
          onCreateNote={onCreateNote}
          onUpdateNote={onUpdateNote}
        />
      ))}
    </ul>
  );
}

type PinListItemProps = {
  pin: PinViewModel;
  isOpen: boolean;
  onToggle: () => void;
  onCreateNote: (pinId: string, content: string) => Promise<void>;
  onUpdateNote: (pinId: string, noteId: string, content: string) => Promise<void>;
};

function PinListItem({ pin, isOpen, onToggle, onCreateNote, onUpdateNote }: PinListItemProps) {
  const [newNote, setNewNote] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const handleAddNote = async () => {
    if (!newNote.trim()) return;
    setSubmitting(true);
    await onCreateNote(pin.id, newNote.trim());
    setSubmitting(false);
    setNewNote("");
  };

  return (
    <li className="rounded-lg border border-zinc-200 bg-white">
      <button
        type="button"
        className="flex w-full items-center gap-3 p-3 text-left"
        onClick={onToggle}
      >
        <span
          className="h-3 w-3 flex-shrink-0 rounded-full"
          style={{ backgroundColor: pin.colour }}
        />
        <span className="flex-1 text-sm font-medium text-zinc-900">{pin.name}</span>
        <span className="text-xs text-zinc-400">{pin.category}</span>
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
                <NoteRow key={note.id} pinId={pin.id} note={note} onUpdateNote={onUpdateNote} />
              ))}
            </ul>
          )}
          <div className="mt-2 flex gap-2">
            <input
              type="text"
              value={newNote}
              onChange={(e) => setNewNote(e.target.value)}
              disabled={submitting}
              placeholder="メモを追加..."
              className="flex-1 rounded border px-2 py-1.5 text-sm text-zinc-900 outline-none focus:ring-2 focus:ring-zinc-300"
            />
            <button
              type="button"
              onClick={handleAddNote}
              disabled={submitting || !newNote.trim()}
              className="rounded-lg bg-black px-3 py-1.5 text-sm font-medium text-white disabled:opacity-50"
            >
              {submitting ? "追加中…" : "追加"}
            </button>
          </div>
        </div>
      )}
    </li>
  );
}

type NoteRowProps = {
  pinId: string;
  note: NoteViewModel;
  onUpdateNote: (pinId: string, noteId: string, content: string) => Promise<void>;
};

function NoteRow({ pinId, note, onUpdateNote }: NoteRowProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [content, setContent] = useState(note.content);
  const [submitting, setSubmitting] = useState(false);

  const handleSave = async () => {
    if (!content.trim()) return;
    setSubmitting(true);
    await onUpdateNote(pinId, note.id, content.trim());
    setSubmitting(false);
    setIsEditing(false);
  };

  const handleCancel = () => {
    setContent(note.content);
    setIsEditing(false);
  };

  if (isEditing) {
    return (
      <li className="flex flex-col gap-1">
        <input
          type="text"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          disabled={submitting}
          className="rounded border px-2 py-1 text-sm text-zinc-900 outline-none focus:ring-2 focus:ring-zinc-300"
        />
        <div className="flex gap-2">
          <button
            type="button"
            onClick={handleSave}
            disabled={submitting || !content.trim()}
            className="rounded bg-black px-2 py-1 text-xs font-medium text-white disabled:opacity-50"
          >
            {submitting ? "保存中…" : "保存"}
          </button>
          <button
            type="button"
            onClick={handleCancel}
            disabled={submitting}
            className="rounded border px-2 py-1 text-xs text-zinc-700 disabled:opacity-50"
          >
            キャンセル
          </button>
        </div>
      </li>
    );
  }

  return (
    <li className="flex items-start justify-between gap-2 text-sm text-zinc-700">
      <span className="flex-1">{note.content}</span>
      <button
        type="button"
        onClick={() => setIsEditing(true)}
        className="flex-shrink-0 text-xs text-zinc-400 hover:text-zinc-700"
      >
        編集
      </button>
    </li>
  );
}
