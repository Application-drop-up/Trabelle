import type { Note } from "@/domain/notes/types";
import type { Pin, PinCategory, PinWithNotes } from "@/domain/pins/types";
import type { Plan } from "@/domain/plans/types";
import type { Spot } from "@/domain/spots/types";

export type NoteViewModel = {
  id: string;
  pinId: string;
  content: string;
  createdAt: Date;
  updatedAt: Date;
};

export type PinViewModel = {
  id: string;
  planId: string;
  name: string;
  latitude: number;
  longitude: number;
  category: PinCategory;
  colour: string;
  notes: NoteViewModel[];
  createdAt: Date;
  updatedAt: Date;
};

export type PlanViewModel = {
  id: string;
  shareToken: string;
  title: string;
  pins: PinViewModel[];
  createdAt: Date;
  updatedAt: Date;
};

export type SpotViewModel = {
  placeId: string;
  name: string;
  address: string;
  latitude: number;
  longitude: number;
};

export function toNoteViewModel(note: Note): NoteViewModel {
  return {
    id: note.id,
    pinId: note.pin_id,
    content: note.content,
    createdAt: new Date(note.created_at),
    updatedAt: new Date(note.updated_at),
  };
}

export function toPinViewModel(pin: Pin, notes: NoteViewModel[] = []): PinViewModel {
  return {
    id: pin.id,
    planId: pin.plan_id,
    name: pin.name,
    latitude: pin.latitude,
    longitude: pin.longitude,
    category: pin.category,
    colour: pin.colour,
    notes,
    createdAt: new Date(pin.created_at),
    updatedAt: new Date(pin.updated_at),
  };
}

export function toPinWithNotesViewModel(pin: PinWithNotes): PinViewModel {
  return toPinViewModel(pin, pin.notes.map(toNoteViewModel));
}

export function toPlanViewModel(plan: Plan): PlanViewModel {
  return {
    id: plan.id,
    shareToken: plan.share_token,
    title: plan.title,
    pins: plan.pins.map(toPinWithNotesViewModel),
    createdAt: new Date(plan.created_at),
    updatedAt: new Date(plan.updated_at),
  };
}

export function toSpotViewModel(spot: Spot): SpotViewModel {
  return {
    placeId: spot.place_id,
    name: spot.name,
    address: spot.address,
    latitude: spot.latitude,
    longitude: spot.longitude,
  };
}
