"use client";

import { useEffect } from "react";

import { useUserContext } from "@/components/UserProvider";
import { usePlan } from "@/hooks/usePlan";
import type { PlanSummary } from "@/domain/plans/types";

type UsePlansListContainerReturn = {
  plans: PlanSummary[];
  loading: boolean;
  error: string | null;
};

export function usePlansListContainer(): UsePlansListContainerReturn {
  const { user, loading: userLoading, error: userError, fetchCurrentUser } = useUserContext();
  const { plans, loading: plansLoading, error: plansError, listPlansForUser } = usePlan();

  useEffect(() => {
    fetchCurrentUser();
  }, [fetchCurrentUser]);

  useEffect(() => {
    if (user) listPlansForUser(user.id);
  }, [user, listPlansForUser]);

  return { plans, loading: userLoading || plansLoading, error: userError ?? plansError };
}
