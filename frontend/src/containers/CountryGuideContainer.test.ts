import { act, renderHook, waitFor } from "@testing-library/react";

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
  it("lists guides on mount", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => [mockGuide],
    } as Response);

    const { result } = renderHook(() => useCountryGuideContainer());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.guideVMs).toHaveLength(1);
    expect(result.current.guideVMs[0]).toEqual({
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
    expect(result.current.selectedGuideVM).toBeNull();
  });

  it("fetches the guide detail when a code is selected", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => mockGuide,
    } as Response);

    const { result } = renderHook(() => useCountryGuideContainer());

    act(() => {
      result.current.onSelectCode("TH");
    });

    await waitFor(() => {
      expect(result.current.selectedGuideVM).not.toBeNull();
    });

    expect(result.current.selectedGuideVM?.countryCode).toBe("TH");
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining("/api/v1/country-guides/TH"),
      expect.anything(),
    );
  });

  it("exposes an error when the selected guide fetch fails", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      status: 404,
      statusText: "Not Found",
      json: async () => ({ message: "country guide not found" }),
    } as Response);

    const { result } = renderHook(() => useCountryGuideContainer());

    act(() => {
      result.current.onSelectCode("ZZ");
    });

    await waitFor(() => {
      expect(result.current.error).toBe("country guide not found");
    });

    expect(result.current.selectedGuideVM).toBeNull();
  });
});
