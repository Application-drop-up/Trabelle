import { renderHook, waitFor } from "@testing-library/react";

import type { CountryGuide } from "@/domain/countryGuides/types";
import { useCountryGuideContainer } from "./CountryGuideContainer";

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

describe("useCountryGuideContainer", () => {
  it("fetches the guide for the given country code on mount", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => mockGuide,
    } as Response);

    const { result } = renderHook(() => useCountryGuideContainer("TH"));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.guideVM).toEqual({
      id: "guide-1",
      countryCode: "TH",
      countryName: "Thailand",
      items: [
        {
          id: "item-1",
          category: "entry_card",
          title: "TDAC",
          description: "Apply online within 72h before arrival",
          url: "https://tdac.immigration.go.th",
          isMandatory: true,
        },
      ],
    });
    expect(result.current.error).toBeNull();
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining("/api/v1/country-guides/TH"),
      expect.anything(),
    );
  });

  it("exposes an error when the fetch fails", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      status: 404,
      statusText: "Not Found",
      json: async () => ({ message: "country guide not found" }),
    } as Response);

    const { result } = renderHook(() => useCountryGuideContainer("ZZ"));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.guideVM).toBeNull();
    expect(result.current.error).toBe("country guide not found");
  });

  it("does not fetch when countryCode is empty", () => {
    global.fetch = jest.fn();

    const { result } = renderHook(() => useCountryGuideContainer(""));

    expect(result.current.guideVM).toBeNull();
    expect(fetch).not.toHaveBeenCalled();
  });
});
