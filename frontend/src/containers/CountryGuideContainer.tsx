"use client";

import { useCallback, useEffect, useState } from "react";

import { useCountryGuides } from "@/hooks/useCountryGuides";
import { toCountryGuideViewModel, type CountryGuideViewModel } from "@/mappers/countryGuideMapper";

type UseCountryGuideContainerReturn = {
  guideVMs: CountryGuideViewModel[];
  selectedCode: string;
  selectedGuideVM: CountryGuideViewModel | null;
  loading: boolean;
  error: string | null;
  onSelectCode: (countryCode: string) => void;
};

export function useCountryGuideContainer(): UseCountryGuideContainerReturn {
  const { guides, loading, error, listCountryGuides, getCountryGuide } = useCountryGuides();
  const [selectedCode, setSelectedCode] = useState("");
  const [selectedGuideVM, setSelectedGuideVM] = useState<CountryGuideViewModel | null>(null);

  useEffect(() => {
    listCountryGuides();
  }, [listCountryGuides]);

  const onSelectCode = useCallback((countryCode: string) => {
    setSelectedCode(countryCode);
  }, []);

  useEffect(() => {
    if (!selectedCode) {
      setSelectedGuideVM(null);
      return;
    }
    getCountryGuide(selectedCode).then((guide) => {
      if (guide) setSelectedGuideVM(toCountryGuideViewModel(guide));
    });
  }, [selectedCode, getCountryGuide]);

  const guideVMs = guides.map(toCountryGuideViewModel);

  return { guideVMs, selectedCode, selectedGuideVM, loading, error, onSelectCode };
}

export type { CountryGuideViewModel };
