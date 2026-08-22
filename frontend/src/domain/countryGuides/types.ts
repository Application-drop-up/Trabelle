import { z } from "zod";

const countryGuideItemCategorySchema = z.enum(["entry_card", "sim_recommendation", "packing_tip"]);

export type CountryGuideItemCategory = z.infer<typeof countryGuideItemCategorySchema>;

export const countryGuideItemSchema = z.object({
  id: z.string(),
  category: countryGuideItemCategorySchema,
  title: z.string(),
  description: z.string(),
  url: z.string(),
  is_mandatory: z.boolean(),
});

export type CountryGuideItem = z.infer<typeof countryGuideItemSchema>;

export const countryGuideSchema = z.object({
  id: z.string(),
  country_code: z.string(),
  country_name: z.string(),
  items: z.array(countryGuideItemSchema),
});

export type CountryGuide = z.infer<typeof countryGuideSchema>;
