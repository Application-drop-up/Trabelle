import type { CountryGuide, CountryGuideItem } from "@/domain/countryGuides/types";
import { toCountryGuideItemViewModel, toCountryGuideViewModel } from "./countryGuideMapper";

const mockItem: CountryGuideItem = {
  id: "item-1",
  category: "entry_card",
  title: "TDAC",
  description: "Apply online within 72h before arrival",
  url: "https://tdac.immigration.go.th",
  is_mandatory: true,
};

const mockGuide: CountryGuide = {
  id: "guide-1",
  country_code: "TH",
  country_name: "Thailand",
  items: [mockItem],
};

describe("toCountryGuideItemViewModel", () => {
  it("converts snake_case CountryGuideItem to camelCase ViewModel", () => {
    const vm = toCountryGuideItemViewModel(mockItem);

    expect(vm.id).toBe("item-1");
    expect(vm.category).toBe("entry_card");
    expect(vm.title).toBe("TDAC");
    expect(vm.description).toBe("Apply online within 72h before arrival");
    expect(vm.url).toBe("https://tdac.immigration.go.th");
    expect(vm.isMandatory).toBe(true);
  });
});

describe("toCountryGuideViewModel", () => {
  it("converts snake_case CountryGuide to camelCase ViewModel", () => {
    const vm = toCountryGuideViewModel(mockGuide);

    expect(vm.id).toBe("guide-1");
    expect(vm.countryCode).toBe("TH");
    expect(vm.countryName).toBe("Thailand");
    expect(vm.items).toHaveLength(1);
    expect(vm.items[0].id).toBe("item-1");
    expect(vm.items[0].isMandatory).toBe(true);
  });

  it("maps an empty items array", () => {
    const vm = toCountryGuideViewModel({ ...mockGuide, items: [] });

    expect(vm.items).toEqual([]);
  });
});
