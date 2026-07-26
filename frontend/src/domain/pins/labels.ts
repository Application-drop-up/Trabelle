import type { PinCategory } from "./types";

export const PIN_CATEGORIES: PinCategory[] = [
  "restaurant",
  "hotel",
  "sightseeing",
  "transport",
  "other",
];

export const PIN_CATEGORY_LABELS: Record<PinCategory, string> = {
  restaurant: "レストラン",
  hotel: "ホテル",
  sightseeing: "観光",
  transport: "交通",
  other: "その他",
};
