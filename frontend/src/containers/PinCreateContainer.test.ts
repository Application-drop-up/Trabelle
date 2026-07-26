import { act, renderHook } from "@testing-library/react";

import type { Pin } from "@/domain/pins/types";
import { usePinCreateContainer } from "./PinCreateContainer";

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

describe("usePinCreateContainer", () => {
  it("initializes with loading false and no error", () => {
    const { result } = renderHook(() => usePinCreateContainer("plan-1"));

    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
  });

  describe("onCreate", () => {
    it("returns PinViewModel on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 201,
        json: async () => mockPin,
      } as Response);

      const { result } = renderHook(() => usePinCreateContainer("plan-1"));

      let vm = null;
      await act(async () => {
        vm = await result.current.onCreate({
          name: "Tokyo Tower",
          latitude: 35.6586,
          longitude: 139.7454,
          category: "sightseeing",
          colour: "#FF5733",
        });
      });

      expect(vm).not.toBeNull();
      expect(vm!.id).toBe("pin-1");
      expect(vm!.planId).toBe("plan-1");
      expect(vm!.name).toBe("Tokyo Tower");
      expect(vm!.category).toBe("sightseeing");
      expect(vm!.colour).toBe("#FF5733");
    });

    it("returns null and sets error on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 400,
        json: async () => ({ error: "invalid pin data" }),
      } as Response);

      const { result } = renderHook(() => usePinCreateContainer("plan-1"));

      let vm = null;
      await act(async () => {
        vm = await result.current.onCreate({
          name: "Tokyo Tower",
          latitude: 35.6586,
          longitude: 139.7454,
          category: "sightseeing",
          colour: "#FF5733",
        });
      });

      expect(vm).toBeNull();
      expect(result.current.error).toBe("invalid pin data");
    });

    it("calls fetch with correct URL and method", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 201,
        json: async () => mockPin,
      } as Response);

      const { result } = renderHook(() => usePinCreateContainer("plan-1"));

      await act(async () => {
        await result.current.onCreate({
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
});
