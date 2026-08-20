import { z } from "zod";

export const spotSchema = z.object({
  place_id: z.string(),
  name: z.string(),
  address: z.string(),
  latitude: z.number(),
  longitude: z.number(),
});

export type Spot = z.infer<typeof spotSchema>;
