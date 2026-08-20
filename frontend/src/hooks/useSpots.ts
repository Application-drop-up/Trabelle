"use client";

import { useCallback, useState } from "react";
import { z } from "zod";

import { apiClient } from "@/lib/apiClient";
import { errorMessages } from "@/lib/messages";
import { spotSchema, type Spot } from "@/domain/spots/types";

export function useSpots() {
  const [spots, setSpots] = useState<Spot[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const searchSpots = useCallback(async (query: string): Promise<Spot[]> => {
    if (!query.trim()) {
      setSpots([]);
      return [];
    }
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.get(
        z.array(spotSchema),
        `/spots/search?query=${encodeURIComponent(query)}`,
      );
      setSpots(result);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.spots.search);
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  const clearSpots = useCallback(() => setSpots([]), []);

  return { spots, loading, error, searchSpots, clearSpots };
}
