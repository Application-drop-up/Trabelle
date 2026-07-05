"use client";

import { useCallback } from "react";

import type { CreatePinInput } from "@/domain/pins/types";
import { useCreatePin } from "@/hooks/usePins";
import { toPinViewModel, type PinViewModel } from "@/mappers/planMapper";

type UsePinCreateContainerReturn = {
  loading: boolean;
  error: string | null;
  onCreate: (input: CreatePinInput) => Promise<PinViewModel | null>;
};

export function usePinCreateContainer(planId: string): UsePinCreateContainerReturn {
  const { createPin, loading, error } = useCreatePin(planId);

  const onCreate = useCallback(
    async (input: CreatePinInput): Promise<PinViewModel | null> => {
      const pin = await createPin(input);
      if (!pin) return null;
      return toPinViewModel(pin);
    },
    [createPin],
  );

  return { loading, error, onCreate };
}
