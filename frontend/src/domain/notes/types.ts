export interface Note {
  id: string;
  pin_id: string;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface CreateNoteInput {
  content: string;
}

export interface UpdateNoteInput {
  content: string;
}
