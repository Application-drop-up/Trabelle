import type { Note } from "@/domain/notes/types";

export type PinCategory = "restaurant" | "hotel" | "sightseeing" | "transport" | "other";

export interface Pin {
  id: string;
  plan_id: string;
  name: string;
  latitude: number;
  longitude: number;
  category: PinCategory;
  colour: string;
  created_at: string;
  updated_at: string;
}

export interface PinWithNotes extends Pin {
  notes: Note[];
}

export interface CreatePinInput {
  name: string;
  latitude: number;
  longitude: number;
  category: PinCategory;
  colour: string;
}

export interface UpdatePinInput {
  category?: PinCategory;
  colour?: string;
}
