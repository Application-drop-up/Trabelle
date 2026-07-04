import { act, renderHook } from "@testing-library/react";

import type { Pin } from "@/domain/pins/types";
import { usePins } from "./usePins";

const mockPin: Pin = {
  id: "pin-1",
  plan_id: "plan-1",
  name: "Tokyo Tower",
  latitude: 35.6586,
  longitude: 139.7454,
  category: "sightseeing",
  colour: "#FF5733",
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-02T00:00:00Z",
};

beforeEach(() => {
  jest.resetAllMocks();
});

describe("usePins", () => {
  describe("createPin", () => {
    it("returns created pin and appends to state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 201,
        json: async () => mockPin,
      } as Response);

      const { result } = renderHook(() => usePins());

      let returned: Pin | null = null;
      await act(async () => {
        returned = await result.current.createPin("plan-1", {
          name: "Tokyo Tower",
          latitude: 35.6586,
          longitude: 139.7454,
          category: "sightseeing",
          colour: "#FF5733",
        });
      });

      expect(returned).toEqual(mockPin);
      expect(result.current.pins).toHaveLength(1);
      expect(result.current.pins[0]).toEqual(mockPin);
      expect(result.current.error).toBeNull();
    });

    it("sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 400,
        statusText: "Bad Request",
        json: async () => ({ error: "invalid pin data" }),
      } as Response);

      const { result } = renderHook(() => usePins());

      let returned: Pin | null = null;
      await act(async () => {
        returned = await result.current.createPin("plan-1", {
          name: "Tokyo Tower",
          latitude: 35.6586,
          longitude: 139.7454,
          category: "sightseeing",
          colour: "#FF5733",
        });
      });

      expect(returned).toBeNull();
      expect(result.current.pins).toHaveLength(0);
      expect(result.current.error).toBe("invalid pin data");
    });

    it("calls fetch with correct URL", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 201,
        json: async () => mockPin,
      } as Response);

      const { result } = renderHook(() => usePins());

      await act(async () => {
        await result.current.createPin("plan-1", {
          name: "Tokyo Tower",
          latitude: 35.6586,
          longitude: 139.7454,
          category: "sightseeing",
          colour: "#FF5733",
        });
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/plans/plan-1/pins"),
        expect.objectContaining({ method: "POST" }),
      );
    });
  });

  describe("updatePin", () => {
    it("returns updated pin and replaces in state on success", async () => {
      const updatedPin: Pin = { ...mockPin, colour: "#00FF00" };
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => updatedPin,
      } as Response);

      const { result } = renderHook(() => usePins());

      // Pre-populate pins
      act(() => {
        result.current.setPins([mockPin]);
      });

      let returned: Pin | null = null;
      await act(async () => {
        returned = await result.current.updatePin("plan-1", "pin-1", { colour: "#00FF00" });
      });

      expect(returned).toEqual(updatedPin);
      expect(result.current.pins[0].colour).toBe("#00FF00");
    });

    it("sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 404,
        statusText: "Not Found",
        json: async () => ({ error: "pin not found" }),
      } as Response);

      const { result } = renderHook(() => usePins());

      await act(async () => {
        await result.current.updatePin("plan-1", "pin-1", { colour: "#00FF00" });
      });

      expect(result.current.error).toBe("pin not found");
    });
  });

  describe("deletePin", () => {
    it("returns true and removes pin from state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 204,
      } as Response);

      const { result } = renderHook(() => usePins());

      act(() => {
        result.current.setPins([mockPin]);
      });

      let ok = false;
      await act(async () => {
        ok = await result.current.deletePin("plan-1", "pin-1");
      });

      expect(ok).toBe(true);
      expect(result.current.pins).toHaveLength(0);
    });

    it("returns false and sets error on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 404,
        statusText: "Not Found",
        json: async () => ({ error: "pin not found" }),
      } as Response);

      const { result } = renderHook(() => usePins());

      let ok = true;
      await act(async () => {
        ok = await result.current.deletePin("plan-1", "pin-1");
      });

      expect(ok).toBe(false);
      expect(result.current.error).toBe("pin not found");
    });
  });
});
