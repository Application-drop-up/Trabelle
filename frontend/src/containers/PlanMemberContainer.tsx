"use client";

import { useCallback, useEffect } from "react";

import { usePlanMembers } from "@/hooks/usePlanMembers";
import { toPlanMemberViewModel, type PlanMemberViewModel } from "@/mappers/planMemberMapper";

type UsePlanMemberContainerReturn = {
  memberVMs: PlanMemberViewModel[];
  loading: boolean;
  error: string | null;
  onAddMember: (userId: string) => Promise<boolean>;
  onRemoveMember: (userId: string) => Promise<boolean>;
};

export function usePlanMemberContainer(planId: string): UsePlanMemberContainerReturn {
  const { members, loading, error, listMembers, addMember, removeMember } = usePlanMembers();

  useEffect(() => {
    listMembers(planId);
  }, [planId, listMembers]);

  const onAddMember = useCallback(
    async (userId: string): Promise<boolean> => {
      const result = await addMember(planId, userId);
      return result !== null;
    },
    [planId, addMember],
  );

  const onRemoveMember = useCallback(
    (userId: string): Promise<boolean> => removeMember(planId, userId),
    [planId, removeMember],
  );

  const memberVMs = members.map(toPlanMemberViewModel);

  return { memberVMs, loading, error, onAddMember, onRemoveMember };
}

export type { PlanMemberViewModel };
