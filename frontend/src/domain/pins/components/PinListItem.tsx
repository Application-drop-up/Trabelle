import { NoteList } from "@/domain/notes/components/NoteList";
import type { PinWithNotes } from "@/domain/pins/types";

export function PinListItem({ pin }: { pin: PinWithNotes }) {
  return (
    <li>
      <p>
        {pin.name} ({pin.category})
      </p>
      <NoteList notes={pin.notes} />
    </li>
  );
}
