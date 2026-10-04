"use client";

import { useCallback, useState } from "react";

import { apiClient } from "@/lib/apiClient";
import { errorMessages } from "@/lib/messages";
import { planSchema, planSummarySchema, type Plan, type PlanSummary } from "@/domain/plans/types";

export function usePlan() {
  const [plan, setPlan] = useState<Plan | null>(null);
  const [plans, setPlans] = useState<PlanSummary[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createPlan = useCallback(async (title: string): Promise<Plan | null> => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.post(planSchema, "/plans", { title });
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
      const result = await apiClient.get(planSchema, `/plans/${shareToken}`);
      setPlan(result);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.plan.fetch);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const listPlansForUser = useCallback(async (userId: string): Promise<PlanSummary[]> => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.get(planSummarySchema.array(), `/api/v1/user/${userId}/plans`);
      setPlans(result);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.plan.listForUser);
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  return { plan, setPlan, plans, loading, error, createPlan, getPlan, listPlansForUser };
}
