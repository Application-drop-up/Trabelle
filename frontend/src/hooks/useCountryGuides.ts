"use client";

import { useCallback, useState } from "react";
import { z } from "zod";

import { apiClient } from "@/lib/apiClient";
import { errorMessages } from "@/lib/messages";
import { countryGuideSchema, type CountryGuide } from "@/domain/countryGuides/types";

export function useCountryGuides() {
  const [guides, setGuides] = useState<CountryGuide[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const listCountryGuides = useCallback(async (): Promise<CountryGuide[]> => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.get(z.array(countryGuideSchema), "/api/v1/country-guides");
      setGuides(result);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.countryGuides.list);
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  const getCountryGuide = useCallback(async (code: string): Promise<CountryGuide | null> => {
    setLoading(true);
    setError(null);
    try {
      return await apiClient.get(countryGuideSchema, `/api/v1/country-guides/${code}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.countryGuides.get);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  return { guides, loading, error, listCountryGuides, getCountryGuide };
}
