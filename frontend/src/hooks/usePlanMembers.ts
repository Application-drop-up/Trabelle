"use client";

import { useCallback, useState } from "react";

import { apiClient } from "@/lib/apiClient";
import { errorMessages } from "@/lib/messages";
import { planMemberSchema, type PlanMember } from "@/domain/planMembers/types";

export function usePlanMembers() {
  const [members, setMembers] = useState<PlanMember[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const listMembers = useCallback(async (planId: string): Promise<PlanMember[]> => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.get(
        planMemberSchema.array(),
        `/api/v1/plans/${planId}/members`,
      );
      setMembers(result);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.planMembers.list);
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  const addMember = useCallback(
    async (planId: string, userId: string): Promise<PlanMember | null> => {
      setLoading(true);
      setError(null);
      try {
        const result = await apiClient.post(planMemberSchema, `/api/v1/plans/${planId}/members`, {
          user_id: userId,
        });
        setMembers((prev) => [...prev, result]);
        return result;
      } catch (err) {
        setError(err instanceof Error ? err.message : errorMessages.planMembers.add);
        return null;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const removeMember = useCallback(async (planId: string, userId: string): Promise<boolean> => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.delete(`/api/v1/plans/${planId}/members/${userId}`);
      setMembers((prev) => prev.filter((m) => m.user_id !== userId));
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.planMembers.remove);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  return { members, loading, error, listMembers, addMember, removeMember };
}
