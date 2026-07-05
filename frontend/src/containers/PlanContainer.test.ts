import { act, renderHook } from "@testing-library/react";

import type { Note } from "@/domain/notes/types";
import type { Pin } from "@/domain/pins/types";
import type { Plan } from "@/domain/plans/types";
import { usePlanContainer } from "./PlanContainer";

const mockNote: Note = {
  id: "note-1",
  pin_id: "pin-1",
  content: "Visit the observation deck",
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-01T00:00:00Z",
};

const mockPin: Pin = {
  id: "pin-1",
  plan_id: "plan-1",
  name: "Tokyo Tower",
  latitude: 35.6586,
  longitude: 139.7454,
  category: "sightseeing",
  colour: "#FF5733",
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-01T00:00:00Z",
};

const mockPlan: Plan = {
  id: "plan-1",
  share_token: "abc123",
  title: "Tokyo Trip",
  pins: [{ ...mockPin, notes: [mockNote] }],
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-01T00:00:00Z",
};

function mockFetchOnce(body: unknown, status = 200) {
  global.fetch = jest.fn().mockResolvedValueOnce({
    ok: true,
    status,
    json: async () => body,
  } as Response);
}

beforeEach(() => {
  jest.resetAllMocks();
});

describe("usePlanContainer", () => {
  describe("initial load", () => {
    it("fetches the plan on mount and converts to PlanViewModel", async () => {
      mockFetchOnce(mockPlan);

      const { result } = renderHook(() => usePlanContainer("abc123"));

      await act(async () => {
        await Promise.resolve();
      });

      expect(result.current.planVM).not.toBeNull();
      expect(result.current.planVM?.shareToken).toBe("abc123");
      expect(result.current.planVM?.title).toBe("Tokyo Trip");
      expect(result.current.planVM?.pins).toHaveLength(1);
      expect(result.current.planVM?.pins[0].name).toBe("Tokyo Tower");
      expect(result.current.planVM?.pins[0].notes).toHaveLength(1);
      expect(result.current.planVM?.pins[0].notes[0].content).toBe("Visit the observation deck");
    });

    it("sets error state when plan fetch fails", async () => {
      global.fetch = jest.fn().mockResolvedValueOnce({
        ok: false,
        status: 404,
        statusText: "Not Found",
        json: async () => ({ error: "plan not found" }),
      } as Response);

      const { result } = renderHook(() => usePlanContainer("invalid-token"));

      await act(async () => {
        await Promise.resolve();
      });

      expect(result.current.planVM).toBeNull();
      expect(result.current.error).toBe("plan not found");
    });
  });

  describe("onCreatePin", () => {
    it("appends new PinViewModel to planVM.pins", async () => {
      global.fetch = jest
        .fn()
        .mockResolvedValueOnce({ ok: true, status: 200, json: async () => mockPlan } as Response)
        .mockResolvedValueOnce({ ok: true, status: 201, json: async () => mockPin } as Response);

      const { result } = renderHook(() => usePlanContainer("abc123"));

      await act(async () => {
        await Promise.resolve();
      });

      const initialCount = result.current.planVM?.pins.length ?? 0;

      const newPin: Pin = {
        ...mockPin,
        id: "pin-2",
        name: "Skytree",
      };
      global.fetch = jest.fn().mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => newPin,
      } as Response);

      await act(async () => {
        await result.current.onCreatePin({
          name: "Skytree",
          latitude: 35.7101,
          longitude: 139.8107,
          category: "sightseeing",
          colour: "#3399FF",
        });
      });

      expect(result.current.planVM?.pins.length).toBe(initialCount + 1);
      expect(result.current.planVM?.pins.at(-1)?.name).toBe("Skytree");
    });
  });

  describe("onUpdatePin", () => {
    it("replaces the updated pin in planVM.pins while keeping its notes", async () => {
      const updatedPin: Pin = { ...mockPin, colour: "#00FF00" };

      global.fetch = jest
        .fn()
        .mockResolvedValueOnce({ ok: true, status: 200, json: async () => mockPlan } as Response)
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          json: async () => updatedPin,
        } as Response);

      const { result } = renderHook(() => usePlanContainer("abc123"));

      await act(async () => {
        await Promise.resolve();
      });

      await act(async () => {
        await result.current.onUpdatePin("pin-1", { colour: "#00FF00" });
      });

      expect(result.current.planVM?.pins[0].colour).toBe("#00FF00");
      expect(result.current.planVM?.pins[0].notes).toHaveLength(1);
    });
  });

  describe("onDeletePin", () => {
    it("removes the deleted pin from planVM.pins", async () => {
      global.fetch = jest
        .fn()
        .mockResolvedValueOnce({ ok: true, status: 200, json: async () => mockPlan } as Response)
        .mockResolvedValueOnce({ ok: true, status: 204 } as Response);

      const { result } = renderHook(() => usePlanContainer("abc123"));

      await act(async () => {
        await Promise.resolve();
      });

      expect(result.current.planVM?.pins).toHaveLength(1);

      await act(async () => {
        await result.current.onDeletePin("pin-1");
      });

      expect(result.current.planVM?.pins).toHaveLength(0);
    });
  });

  describe("onCreateNote", () => {
    it("appends new NoteViewModel to the correct pin's notes", async () => {
      const newNote: Note = {
        id: "note-2",
        pin_id: "pin-1",
        content: "Bring camera",
        created_at: "2024-01-02T00:00:00Z",
        updated_at: "2024-01-02T00:00:00Z",
      };

      global.fetch = jest
        .fn()
        .mockResolvedValueOnce({ ok: true, status: 200, json: async () => mockPlan } as Response)
        .mockResolvedValueOnce({
          ok: true,
          status: 201,
          json: async () => newNote,
        } as Response);

      const { result } = renderHook(() => usePlanContainer("abc123"));

      await act(async () => {
        await Promise.resolve();
      });

      await act(async () => {
        await result.current.onCreateNote("pin-1", "Bring camera");
      });

      const pin = result.current.planVM?.pins[0];
      expect(pin?.notes).toHaveLength(2);
      expect(pin?.notes.at(-1)?.content).toBe("Bring camera");
    });
  });

  describe("onUpdateNote", () => {
    it("replaces the updated note in the correct pin's notes", async () => {
      const updatedNote: Note = { ...mockNote, content: "Updated content" };

      global.fetch = jest
        .fn()
        .mockResolvedValueOnce({ ok: true, status: 200, json: async () => mockPlan } as Response)
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          json: async () => updatedNote,
        } as Response);

      const { result } = renderHook(() => usePlanContainer("abc123"));

      await act(async () => {
        await Promise.resolve();
      });

      await act(async () => {
        await result.current.onUpdateNote("pin-1", "note-1", "Updated content");
      });

      expect(result.current.planVM?.pins[0].notes[0].content).toBe("Updated content");
    });
  });

  describe("onDeleteNote", () => {
    it("removes the deleted note from the correct pin's notes", async () => {
      global.fetch = jest
        .fn()
        .mockResolvedValueOnce({ ok: true, status: 200, json: async () => mockPlan } as Response)
        .mockResolvedValueOnce({ ok: true, status: 204 } as Response);

      const { result } = renderHook(() => usePlanContainer("abc123"));

      await act(async () => {
        await Promise.resolve();
      });

      expect(result.current.planVM?.pins[0].notes).toHaveLength(1);

      await act(async () => {
        await result.current.onDeleteNote("pin-1", "note-1");
      });

      expect(result.current.planVM?.pins[0].notes).toHaveLength(0);
    });
  });
});
