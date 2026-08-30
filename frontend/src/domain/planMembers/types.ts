import { z } from "zod";

export const planMemberSchema = z.object({
  id: z.string(),
  plan_id: z.string(),
  user_id: z.string(),
  created_at: z.string(),
});

export type PlanMember = z.infer<typeof planMemberSchema>;

export interface AddPlanMemberInput {
  user_id: string;
}
