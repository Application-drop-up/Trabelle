import type { PinWithNotes } from "@/domain/pins/types";

export interface Plan {
  id: string;
  share_token: string;
  title: string;
  pins: PinWithNotes[];
  created_at: string;
  updated_at: string;
}
