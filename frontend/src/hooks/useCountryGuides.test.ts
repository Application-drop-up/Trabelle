import { act, renderHook } from "@testing-library/react";

import type { CountryGuide } from "@/domain/countryGuides/types";
import { useCountryGuides } from "./useCountryGuides";

const mockGuide: CountryGuide = {
  id: "guide-1",
  country_code: "TH",
  country_name: "Thailand",
  items: [
    {
      id: "item-1",
      category: "entry_card",
      title: "TDAC",
      description: "Apply online within 72h before arrival",
      url: "https://tdac.immigration.go.th",
      is_mandatory: true,
    },
  ],
};

beforeEach(() => {
  jest.resetAllMocks();
});

describe("useCountryGuides", () => {
  describe("listCountryGuides", () => {
    it("returns guides and updates state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [mockGuide],
      } as Response);

      const { result } = renderHook(() => useCountryGuides());

      let returned: CountryGuide[] = [];
      await act(async () => {
        returned = await result.current.listCountryGuides();
      });

      expect(returned).toHaveLength(1);
      expect(returned[0]).toEqual(mockGuide);
      expect(result.current.guides).toHaveLength(1);
      expect(result.current.error).toBeNull();
    });

    it("sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 500,
        statusText: "Internal Server Error",
        json: async () => ({ message: "internal server error" }),
      } as Response);

      const { result } = renderHook(() => useCountryGuides());

      let returned: CountryGuide[] = [];
      await act(async () => {
        returned = await result.current.listCountryGuides();
      });

      expect(returned).toHaveLength(0);
      expect(result.current.guides).toHaveLength(0);
      expect(result.current.error).toBe("internal server error");
    });
  });

  describe("getCountryGuide", () => {
    it("returns the guide on success without touching guides state", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => mockGuide,
      } as Response);

      const { result } = renderHook(() => useCountryGuides());

      let returned: CountryGuide | null = null;
      await act(async () => {
        returned = await result.current.getCountryGuide("TH");
      });

      expect(returned).toEqual(mockGuide);
      expect(result.current.guides).toHaveLength(0);
      expect(result.current.error).toBeNull();
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/api/v1/country-guides/TH"),
        expect.any(Object),
      );
    });

    it("sets error state and returns null for an unknown code", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 404,
        statusText: "Not Found",
        json: async () => ({ message: "country guide not found" }),
      } as Response);

      const { result } = renderHook(() => useCountryGuides());

      let returned: CountryGuide | null = mockGuide;
      await act(async () => {
        returned = await result.current.getCountryGuide("ZZ");
      });

      expect(returned).toBeNull();
      expect(result.current.error).toBe("country guide not found");
    });
  });
});
