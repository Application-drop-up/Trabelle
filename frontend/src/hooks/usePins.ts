"use client";

import { useCallback, useState } from "react";

import type { CreatePinInput, Pin } from "@/domain/pins/types";
import { apiClient } from "@/lib/apiClient";

export function useCreatePin(planId: string) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createPin = useCallback(
    async (input: CreatePinInput): Promise<Pin | null> => {
      setLoading(true);
      setError(null);
      try {
        const pin = await apiClient.post<Pin>(`/plans/${planId}/pins`, input);
        return pin;
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to create pin");
        return null;
      } finally {
        setLoading(false);
      }
    },
    [planId],
  );

  return { createPin, loading, error };
}

export function usePins(planId: string) {
  const [pins, setPins] = useState<Pin[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchPins = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.get<Pin[]>(`/plans/${planId}/pins`);
      setPins(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch pins");
    } finally {
      setLoading(false);
    }
  }, [planId]);

  return { pins, loading, error, fetchPins };
}
