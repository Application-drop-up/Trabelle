import type {
  CountryGuide,
  CountryGuideItem,
  CountryGuideItemCategory,
} from "@/domain/countryGuides/types";

export type CountryGuideItemViewModel = {
  id: string;
  category: CountryGuideItemCategory;
  title: string;
  description: string;
  url: string;
  isMandatory: boolean;
};

export type CountryGuideViewModel = {
  id: string;
  countryCode: string;
  countryName: string;
  items: CountryGuideItemViewModel[];
};

export function toCountryGuideItemViewModel(item: CountryGuideItem): CountryGuideItemViewModel {
  return {
    id: item.id,
    category: item.category,
    title: item.title,
    description: item.description,
    url: item.url,
    isMandatory: item.is_mandatory,
  };
}

export function toCountryGuideViewModel(guide: CountryGuide): CountryGuideViewModel {
  return {
    id: guide.id,
    countryCode: guide.country_code,
    countryName: guide.country_name,
    items: guide.items.map(toCountryGuideItemViewModel),
  };
}
