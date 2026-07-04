"use client";

import { useCallback, useState } from "react";

import { apiClient } from "@/lib/apiClient";
import type { Note } from "@/domain/notes/types";

export function useNotes() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createNote = useCallback(
    async (planId: string, pinId: string, content: string): Promise<Note | null> => {
      setLoading(true);
      setError(null);
      try {
        return await apiClient.post<Note>(`/plans/${planId}/pins/${pinId}/notes`, { content });
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to create note");
        return null;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const updateNote = useCallback(
    async (
      planId: string,
      pinId: string,
      noteId: string,
      content: string,
    ): Promise<Note | null> => {
      setLoading(true);
      setError(null);
      try {
        return await apiClient.patch<Note>(`/plans/${planId}/pins/${pinId}/notes/${noteId}`, {
          content,
        });
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to update note");
        return null;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const deleteNote = useCallback(
    async (planId: string, pinId: string, noteId: string): Promise<boolean> => {
      setLoading(true);
      setError(null);
      try {
        await apiClient.delete(`/plans/${planId}/pins/${pinId}/notes/${noteId}`);
        return true;
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to delete note");
        return false;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  return { loading, error, createNote, updateNote, deleteNote };
}
