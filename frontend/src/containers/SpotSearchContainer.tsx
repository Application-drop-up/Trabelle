"use client";

import { useCallback, useState } from "react";

import type { CreatePinInput } from "@/domain/pins/types";
import { useSpots } from "@/hooks/useSpots";
import { toSpotViewModel, type SpotViewModel } from "@/mappers/planMapper";

type UseSpotSearchContainerReturn = {
  query: string;
  spotVMs: SpotViewModel[];
  loading: boolean;
  error: string | null;
  onChangeQuery: (value: string) => void;
  onSearch: () => Promise<void>;
  onSelectSpot: (spot: SpotViewModel) => CreatePinInput;
  onClear: () => void;
};

export function useSpotSearchContainer(): UseSpotSearchContainerReturn {
  const [query, setQuery] = useState("");
  const { spots, loading, error, searchSpots, clearSpots } = useSpots();

  const onChangeQuery = useCallback((value: string) => {
    setQuery(value);
  }, []);

  const onSearch = useCallback(async () => {
    await searchSpots(query);
  }, [query, searchSpots]);

  const onSelectSpot = useCallback((spot: SpotViewModel): CreatePinInput => {
    return {
      name: spot.name,
      latitude: spot.latitude,
      longitude: spot.longitude,
      category: "other",
      colour: "#FF5733",
    };
  }, []);

  const onClear = useCallback(() => {
    setQuery("");
    clearSpots();
  }, [clearSpots]);

  const spotVMs = spots.map(toSpotViewModel);

  return {
    query,
    spotVMs,
    loading,
    error,
    onChangeQuery,
    onSearch,
    onSelectSpot,
    onClear,
  };
}

export type { SpotViewModel };
