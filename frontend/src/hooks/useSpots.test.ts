import { act, renderHook } from "@testing-library/react";

import type { Spot } from "@/domain/spots/types";
import { useSpots } from "./useSpots";

const mockSpot: Spot = {
  place_id: "spot-1",
  name: "Shibuya Crossing",
  address: "Shibuya, Tokyo",
  latitude: 35.6595,
  longitude: 139.7004,
};

beforeEach(() => {
  jest.resetAllMocks();
});

describe("useSpots", () => {
  describe("searchSpots", () => {
    it("returns spots and updates state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [mockSpot],
      } as Response);

      const { result } = renderHook(() => useSpots());

      let returned: Spot[] = [];
      await act(async () => {
        returned = await result.current.searchSpots("Shibuya");
      });

      expect(returned).toHaveLength(1);
      expect(returned[0]).toEqual(mockSpot);
      expect(result.current.spots).toHaveLength(1);
      expect(result.current.error).toBeNull();
    });

    it("sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 500,
        statusText: "Internal Server Error",
        json: async () => ({ error: "search failed" }),
      } as Response);

      const { result } = renderHook(() => useSpots());

      let returned: Spot[] = [];
      await act(async () => {
        returned = await result.current.searchSpots("Shibuya");
      });

      expect(returned).toHaveLength(0);
      expect(result.current.spots).toHaveLength(0);
      expect(result.current.error).toBe("search failed");
    });

    it("returns empty array and clears spots when query is blank", async () => {
      const { result } = renderHook(() => useSpots());

      // Pre-populate spots
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [mockSpot],
      } as Response);

      await act(async () => {
        await result.current.searchSpots("Shibuya");
      });

      expect(result.current.spots).toHaveLength(1);

      let returned: Spot[] = [];
      await act(async () => {
        returned = await result.current.searchSpots("   ");
      });

      expect(returned).toHaveLength(0);
      expect(result.current.spots).toHaveLength(0);
      expect(fetch).toHaveBeenCalledTimes(1); // second call did not hit fetch
    });

    it("calls fetch with URL-encoded query", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [],
      } as Response);

      const { result } = renderHook(() => useSpots());

      await act(async () => {
        await result.current.searchSpots("東京タワー");
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("query=%E6%9D%B1%E4%BA%AC%E3%82%BF%E3%83%AF%E3%83%BC"),
        expect.any(Object),
      );
    });
  });

  describe("clearSpots", () => {
    it("resets spots to empty array", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [mockSpot],
      } as Response);

      const { result } = renderHook(() => useSpots());

      await act(async () => {
        await result.current.searchSpots("Shibuya");
      });

      expect(result.current.spots).toHaveLength(1);

      act(() => {
        result.current.clearSpots();
      });

      expect(result.current.spots).toHaveLength(0);
    });
  });
});
