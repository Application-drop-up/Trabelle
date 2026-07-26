"use client";

import { useCallback } from "react";

import type { CreatePinInput } from "@/domain/pins/types";
import { usePins } from "@/hooks/usePins";
import { toPinViewModel, type PinViewModel } from "@/mappers/planMapper";

type UsePinCreateContainerReturn = {
  loading: boolean;
  error: string | null;
  onCreate: (data: CreatePinInput) => Promise<PinViewModel | null>;
};

export function usePinCreateContainer(planId: string): UsePinCreateContainerReturn {
  const { createPin, loading, error } = usePins();

  const onCreate = useCallback(
    async (data: CreatePinInput): Promise<PinViewModel | null> => {
      const pin = await createPin(planId, data);
      if (!pin) return null;
      return toPinViewModel(pin);
    },
    [planId, createPin],
  );

  return { loading, error, onCreate };
}
