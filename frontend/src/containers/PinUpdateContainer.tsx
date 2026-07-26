"use client";

import { useCallback } from "react";

import type { UpdatePinInput } from "@/domain/pins/types";
import { usePins } from "@/hooks/usePins";
import { toPinViewModel, type PinViewModel } from "@/mappers/planMapper";

type UsePinUpdateContainerReturn = {
  loading: boolean;
  error: string | null;
  onUpdate: (planId: string, pinId: string, data: UpdatePinInput) => Promise<PinViewModel | null>;
};

export function usePinUpdateContainer(): UsePinUpdateContainerReturn {
  const { updatePin, loading, error } = usePins();

  const onUpdate = useCallback(
    async (planId: string, pinId: string, data: UpdatePinInput): Promise<PinViewModel | null> => {
      const pin = await updatePin(planId, pinId, data);
      if (!pin) return null;
      return toPinViewModel(pin);
    },
    [updatePin],
  );

  return { loading, error, onUpdate };
}
