"use client";

import { useCallback, useState } from "react";

import { apiClient } from "@/lib/apiClient";
import { errorMessages } from "@/lib/messages";
import type { Plan } from "@/domain/plans/types";

export function usePlan() {
  const [plan, setPlan] = useState<Plan | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createPlan = useCallback(async (title: string): Promise<Plan | null> => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.post<Plan>("/plans", { title });
      setPlan(result);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.plan.create);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const getPlan = useCallback(async (shareToken: string): Promise<Plan | null> => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.get<Plan>(`/plans/${shareToken}`);
      setPlan(result);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.plan.fetch);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  return { plan, setPlan, loading, error, createPlan, getPlan };
}
