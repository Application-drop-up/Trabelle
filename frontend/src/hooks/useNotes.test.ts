import { act, renderHook } from "@testing-library/react";

import type { Note } from "@/domain/notes/types";
import { useNotes } from "./useNotes";

const mockNote: Note = {
  id: "note-1",
  pin_id: "pin-1",
  content: "Test note content",
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-02T00:00:00Z",
};

beforeEach(() => {
  jest.resetAllMocks();
});

describe("useNotes", () => {
  describe("createNote", () => {
    it("returns created note on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 201,
        json: async () => mockNote,
      } as Response);

      const { result } = renderHook(() => useNotes());

      let returned: Note | null = null;
      await act(async () => {
        returned = await result.current.createNote("plan-1", "pin-1", "Test note content");
      });

      expect(returned).toEqual(mockNote);
      expect(result.current.error).toBeNull();
    });

    it("sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 400,
        statusText: "Bad Request",
        json: async () => ({ error: "invalid note content" }),
      } as Response);

      const { result } = renderHook(() => useNotes());

      let returned: Note | null = null;
      await act(async () => {
        returned = await result.current.createNote("plan-1", "pin-1", "");
      });

      expect(returned).toBeNull();
      expect(result.current.error).toBe("invalid note content");
    });

    it("calls fetch with correct URL and body", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 201,
        json: async () => mockNote,
      } as Response);

      const { result } = renderHook(() => useNotes());

      await act(async () => {
        await result.current.createNote("plan-1", "pin-1", "Test note content");
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/plans/plan-1/pins/pin-1/notes"),
        expect.objectContaining({ method: "POST" }),
      );
    });
  });

  describe("updateNote", () => {
    it("returns updated note on success", async () => {
      const updatedNote: Note = { ...mockNote, content: "Updated content" };
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => updatedNote,
      } as Response);

      const { result } = renderHook(() => useNotes());

      let returned: Note | null = null;
      await act(async () => {
        returned = await result.current.updateNote("plan-1", "pin-1", "note-1", "Updated content");
      });

      expect(returned).toEqual(updatedNote);
      expect(result.current.error).toBeNull();
    });

    it("sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 404,
        statusText: "Not Found",
        json: async () => ({ error: "note not found" }),
      } as Response);

      const { result } = renderHook(() => useNotes());

      await act(async () => {
        await result.current.updateNote("plan-1", "pin-1", "note-1", "Updated content");
      });

      expect(result.current.error).toBe("note not found");
    });
  });

  describe("deleteNote", () => {
    it("returns true on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 204,
      } as Response);

      const { result } = renderHook(() => useNotes());

      let ok = false;
      await act(async () => {
        ok = await result.current.deleteNote("plan-1", "pin-1", "note-1");
      });

      expect(ok).toBe(true);
      expect(result.current.error).toBeNull();
    });

    it("returns false and sets error on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 404,
        statusText: "Not Found",
        json: async () => ({ error: "note not found" }),
      } as Response);

      const { result } = renderHook(() => useNotes());

      let ok = true;
      await act(async () => {
        ok = await result.current.deleteNote("plan-1", "pin-1", "note-1");
      });

      expect(ok).toBe(false);
      expect(result.current.error).toBe("note not found");
    });

    it("calls fetch with correct URL", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 204,
      } as Response);

      const { result } = renderHook(() => useNotes());

      await act(async () => {
        await result.current.deleteNote("plan-1", "pin-1", "note-1");
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/plans/plan-1/pins/pin-1/notes/note-1"),
        expect.objectContaining({ method: "DELETE" }),
      );
    });
  });
});
