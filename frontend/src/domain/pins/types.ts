import { z } from "zod";

import { noteSchema } from "@/domain/notes/types";

const pinCategorySchema = z.enum(["restaurant", "hotel", "sightseeing", "transport", "other"]);

export type PinCategory = z.infer<typeof pinCategorySchema>;

export const pinSchema = z.object({
  id: z.string(),
  plan_id: z.string(),
  name: z.string(),
  latitude: z.number(),
  longitude: z.number(),
  category: pinCategorySchema,
  colour: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type Pin = z.infer<typeof pinSchema>;

export const pinWithNotesSchema = pinSchema.extend({
  notes: z.array(noteSchema),
});

export type PinWithNotes = z.infer<typeof pinWithNotesSchema>;

export interface CreatePinInput {
  name: string;
  latitude: number;
  longitude: number;
  category: PinCategory;
  colour: string;
  // Set when the Pin is created from a Spot search result, so the backend
  // can cache that Spot for future searches.
  place_id?: string;
  address?: string;
}

export interface UpdatePinInput {
  category?: PinCategory;
  colour?: string;
}
