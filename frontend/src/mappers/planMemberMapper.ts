import type { PlanMember } from "@/domain/planMembers/types";

export type PlanMemberViewModel = {
  id: string;
  planId: string;
  userId: string;
  createdAt: string;
};

export function toPlanMemberViewModel(member: PlanMember): PlanMemberViewModel {
  return {
    id: member.id,
    planId: member.plan_id,
    userId: member.user_id,
    createdAt: member.created_at,
  };
}
