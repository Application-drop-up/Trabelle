"use client";

import { useCallback, useState } from "react";

import type { Plan } from "@/domain/plans/types";
import { usePlan } from "@/hooks/usePlan";

type UseCreatePlanContainerReturn = {
  title: string;
  loading: boolean;
  error: string | null;
  onChangeTitle: (value: string) => void;
  onSubmit: () => Promise<Plan | null>;
};

export function useCreatePlanContainer(): UseCreatePlanContainerReturn {
  const [title, setTitle] = useState("");
  const { createPlan, loading, error } = usePlan();

  const onChangeTitle = useCallback((value: string) => {
    setTitle(value);
  }, []);

  const onSubmit = useCallback(async (): Promise<Plan | null> => {
    if (!title.trim()) return null;
    return createPlan(title.trim());
  }, [title, createPlan]);

  return { title, loading, error, onChangeTitle, onSubmit };
}
