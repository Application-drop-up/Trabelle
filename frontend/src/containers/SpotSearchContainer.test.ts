import { act, renderHook } from "@testing-library/react";

import type { Spot } from "@/domain/spots/types";
import type { SpotViewModel } from "./SpotSearchContainer";
import { useSpotSearchContainer } from "./SpotSearchContainer";

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

describe("useSpotSearchContainer", () => {
  it("initializes with empty state", () => {
    const { result } = renderHook(() => useSpotSearchContainer());

    expect(result.current.query).toBe("");
    expect(result.current.spotVMs).toHaveLength(0);
    expect(result.current.error).toBeNull();
    expect(result.current.loading).toBe(false);
  });

  describe("onChangeQuery", () => {
    it("updates query state", () => {
      const { result } = renderHook(() => useSpotSearchContainer());

      act(() => {
        result.current.onChangeQuery("Shibuya");
      });

      expect(result.current.query).toBe("Shibuya");
    });
  });

  describe("onSearch", () => {
    it("searches and populates spotVMs as camelCase ViewModels", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [mockSpot],
      } as Response);

      const { result } = renderHook(() => useSpotSearchContainer());

      act(() => {
        result.current.onChangeQuery("Shibuya");
      });

      await act(async () => {
        await result.current.onSearch();
      });

      expect(result.current.spotVMs).toHaveLength(1);
      expect(result.current.spotVMs[0].placeId).toBe("spot-1");
      expect(result.current.spotVMs[0].name).toBe("Shibuya Crossing");
      expect(result.current.spotVMs[0].address).toBe("Shibuya, Tokyo");
    });

    it("sets error state on search failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 500,
        statusText: "Internal Server Error",
        json: async () => ({ error: "search failed" }),
      } as Response);

      const { result } = renderHook(() => useSpotSearchContainer());

      act(() => {
        result.current.onChangeQuery("Shibuya");
      });

      await act(async () => {
        await result.current.onSearch();
      });

      expect(result.current.error).toBe("search failed");
      expect(result.current.spotVMs).toHaveLength(0);
    });
  });

  describe("onSelectSpot", () => {
    it("converts SpotViewModel to CreatePinInput with default category and colour", () => {
      const { result } = renderHook(() => useSpotSearchContainer());

      const spotVM: SpotViewModel = {
        placeId: "spot-1",
        name: "Shibuya Crossing",
        address: "Shibuya, Tokyo",
        latitude: 35.6595,
        longitude: 139.7004,
      };

      const pinInput = result.current.onSelectSpot(spotVM);

      expect(pinInput.name).toBe("Shibuya Crossing");
      expect(pinInput.latitude).toBe(35.6595);
      expect(pinInput.longitude).toBe(139.7004);
      expect(pinInput.category).toBe("other");
      expect(pinInput.colour).toBe("#FF5733");
    });
  });

  describe("onClear", () => {
    it("resets query and spots", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [mockSpot],
      } as Response);

      const { result } = renderHook(() => useSpotSearchContainer());

      act(() => {
        result.current.onChangeQuery("Shibuya");
      });

      await act(async () => {
        await result.current.onSearch();
      });

      expect(result.current.spotVMs).toHaveLength(1);
      expect(result.current.query).toBe("Shibuya");

      act(() => {
        result.current.onClear();
      });

      expect(result.current.query).toBe("");
      expect(result.current.spotVMs).toHaveLength(0);
    });
  });
});
