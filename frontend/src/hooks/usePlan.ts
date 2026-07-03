"use client";

import { useCallback, useState } from "react";

import { apiClient } from "@/lib/apiClient";
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
      setError(err instanceof Error ? err.message : "Failed to create plan");
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
      setError(err instanceof Error ? err.message : "Failed to fetch plan");
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  return { plan, setPlan, loading, error, createPlan, getPlan };
}
