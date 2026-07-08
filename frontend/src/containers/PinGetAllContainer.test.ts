import { act, renderHook } from "@testing-library/react";

import type { Pin } from "@/domain/pins/types";
import { usePinGetAllContainer } from "./PinGetAllContainer";

const mockPin1: Pin = {
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

const mockPin2: Pin = {
  id: "pin-2",
  plan_id: "plan-1",
  name: "Shibuya Crossing",
  latitude: 35.6595,
  longitude: 139.7004,
  category: "other",
  colour: "#0000FF",
  created_at: "2024-01-03T00:00:00Z",
  updated_at: "2024-01-04T00:00:00Z",
};

beforeEach(() => {
  jest.resetAllMocks();
});

describe("usePinGetAllContainer", () => {
  it("initializes with empty pinVMs, loading false, and no error", () => {
    const { result } = renderHook(() => usePinGetAllContainer());

    expect(result.current.pinVMs).toEqual([]);
    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
  });

  describe("getAllPins", () => {
    it("returns all PinViewModels and updates state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [mockPin1, mockPin2],
      } as Response);

      const { result } = renderHook(() => usePinGetAllContainer());

      let vms: unknown[] = [];
      await act(async () => {
        vms = await result.current.getAllPins("plan-1");
      });

      expect(vms).toHaveLength(2);
      expect(vms[0]).toMatchObject({ id: "pin-1", name: "Tokyo Tower" });
      expect(vms[1]).toMatchObject({ id: "pin-2", name: "Shibuya Crossing" });
      expect(result.current.pinVMs).toHaveLength(2);
    });

    it("returns empty array when no pins exist", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [],
      } as Response);

      const { result } = renderHook(() => usePinGetAllContainer());

      let vms: unknown[] = [];
      await act(async () => {
        vms = await result.current.getAllPins("plan-1");
      });

      expect(vms).toHaveLength(0);
      expect(result.current.pinVMs).toHaveLength(0);
    });

    it("returns empty array and sets error on API failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 500,
        json: async () => ({ error: "internal server error" }),
      } as Response);

      const { result } = renderHook(() => usePinGetAllContainer());

      let vms: unknown[] = [];
      await act(async () => {
        vms = await result.current.getAllPins("plan-1");
      });

      expect(vms).toHaveLength(0);
      expect(result.current.error).toBe("internal server error");
    });

    it("calls fetch with correct URL", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [mockPin1],
      } as Response);

      const { result } = renderHook(() => usePinGetAllContainer());

      await act(async () => {
        await result.current.getAllPins("plan-1");
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/plans/plan-1/pins"),
        expect.anything(),
      );
    });
  });
});
