import { z } from "zod";

import { pinWithNotesSchema } from "@/domain/pins/types";

export const planSchema = z.object({
  id: z.string(),
  share_token: z.string(),
  title: z.string(),
  pins: z.array(pinWithNotesSchema),
  created_at: z.string(),
  updated_at: z.string(),
});

export type Plan = z.infer<typeof planSchema>;

export const planSummarySchema = z.object({
  id: z.string(),
  share_token: z.string(),
  title: z.string(),
  is_public: z.boolean(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type PlanSummary = z.infer<typeof planSummarySchema>;
