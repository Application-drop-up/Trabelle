import { z } from "zod";

export const noteSchema = z.object({
  id: z.string(),
  pin_id: z.string(),
  content: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type Note = z.infer<typeof noteSchema>;

export interface CreateNoteInput {
  content: string;
}

export interface UpdateNoteInput {
  content: string;
}
