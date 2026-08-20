"use client";

import { useCallback, useState } from "react";

import { apiClient } from "@/lib/apiClient";
import { errorMessages } from "@/lib/messages";
import type { CreatePinInput, Pin, UpdatePinInput } from "@/domain/pins/types";

export function usePins() {
  const [pins, setPins] = useState<Pin[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const listPins = useCallback(async (planId: string): Promise<Pin[]> => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.get<Pin[]>(`/plans/${planId}/pins`);
      setPins(result);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.pins.fetch);
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  const createPin = useCallback(
    async (planId: string, data: CreatePinInput): Promise<Pin | null> => {
      setLoading(true);
      setError(null);
      try {
        const result = await apiClient.post<Pin>(`/plans/${planId}/pins`, data);
        setPins((prev) => [...prev, result]);
        return result;
      } catch (err) {
        setError(err instanceof Error ? err.message : errorMessages.pins.create);
        return null;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const updatePin = useCallback(
    async (planId: string, pinId: string, data: UpdatePinInput): Promise<Pin | null> => {
      setLoading(true);
      setError(null);
      try {
        const result = await apiClient.patch<Pin>(`/plans/${planId}/pins/${pinId}`, data);
        setPins((prev) => prev.map((p) => (p.id === pinId ? result : p)));
        return result;
      } catch (err) {
        setError(err instanceof Error ? err.message : errorMessages.pins.update);
        return null;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const deletePin = useCallback(async (planId: string, pinId: string): Promise<boolean> => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.delete(`/plans/${planId}/pins/${pinId}`);
      setPins((prev) => prev.filter((p) => p.id !== pinId));
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.pins.delete);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  return { pins, setPins, loading, error, listPins, createPin, updatePin, deletePin };
}
