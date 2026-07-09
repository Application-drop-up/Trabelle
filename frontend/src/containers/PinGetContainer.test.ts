import { act, renderHook } from "@testing-library/react";

import type { Pin } from "@/domain/pins/types";
import { usePinGetContainer } from "./PinGetContainer";

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

describe("usePinGetContainer", () => {
  it("initializes with pinVM null, loading false, and no error", () => {
    const { result } = renderHook(() => usePinGetContainer());

    expect(result.current.pinVM).toBeNull();
    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
  });

  describe("getPin", () => {
    it("returns PinViewModel when pin is found", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [mockPin],
      } as Response);

      const { result } = renderHook(() => usePinGetContainer());

      let vm = null;
      await act(async () => {
        vm = await result.current.getPin("plan-1", "pin-1");
      });

      expect(vm).not.toBeNull();
      expect(vm!.id).toBe("pin-1");
      expect(vm!.planId).toBe("plan-1");
      expect(vm!.name).toBe("Tokyo Tower");
      expect(result.current.pinVM).not.toBeNull();
      expect(result.current.pinVM!.id).toBe("pin-1");
    });

    it("returns null when pin is not found in list", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [mockPin],
      } as Response);

      const { result } = renderHook(() => usePinGetContainer());

      let vm = null;
      await act(async () => {
        vm = await result.current.getPin("plan-1", "pin-999");
      });

      expect(vm).toBeNull();
      expect(result.current.pinVM).toBeNull();
    });

    it("returns null and sets error on API failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 500,
        json: async () => ({ error: "internal server error" }),
      } as Response);

      const { result } = renderHook(() => usePinGetContainer());

      let vm = null;
      await act(async () => {
        vm = await result.current.getPin("plan-1", "pin-1");
      });

      expect(vm).toBeNull();
      expect(result.current.error).toBe("internal server error");
    });

    it("calls fetch with correct URL", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [mockPin],
      } as Response);

      const { result } = renderHook(() => usePinGetContainer());

      await act(async () => {
        await result.current.getPin("plan-1", "pin-1");
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/plans/plan-1/pins"),
        expect.anything(),
      );
    });
  });
});
