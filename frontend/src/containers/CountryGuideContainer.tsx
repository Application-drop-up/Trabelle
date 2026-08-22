"use client";

import { useEffect, useState } from "react";

import { useCountryGuides } from "@/hooks/useCountryGuides";
import { toCountryGuideViewModel, type CountryGuideViewModel } from "@/mappers/countryGuideMapper";

type UseCountryGuideContainerReturn = {
  guideVM: CountryGuideViewModel | null;
  loading: boolean;
  error: string | null;
};

export function useCountryGuideContainer(countryCode: string): UseCountryGuideContainerReturn {
  const { loading, error, getCountryGuide } = useCountryGuides();
  const [guideVM, setGuideVM] = useState<CountryGuideViewModel | null>(null);

  useEffect(() => {
    if (!countryCode) {
      setGuideVM(null);
      return;
    }
    getCountryGuide(countryCode).then((guide) => {
      if (guide) setGuideVM(toCountryGuideViewModel(guide));
    });
  }, [countryCode, getCountryGuide]);

  return { guideVM, loading, error };
}

export type { CountryGuideViewModel };
