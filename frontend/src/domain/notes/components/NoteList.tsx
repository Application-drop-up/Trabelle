import type { Note } from "@/domain/notes/types";

export function NoteList({ notes }: { notes: Note[] }) {
  return (
    <ul className="ml-4 list-disc">
      {notes.map((note) => (
        <li key={note.id}>{note.content}</li>
      ))}
    </ul>
  );
}
