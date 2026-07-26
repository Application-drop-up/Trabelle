"use client";

import { useCallback, useState } from "react";

import { usePins } from "@/hooks/usePins";
import { toPinViewModel, type PinViewModel } from "@/mappers/planMapper";

type UsePinGetContainerReturn = {
  pinVM: PinViewModel | null;
  loading: boolean;
  error: string | null;
  getPin: (planId: string, pinId: string) => Promise<PinViewModel | null>;
};

export function usePinGetContainer(): UsePinGetContainerReturn {
  const { listPins, loading, error } = usePins();
  const [pinVM, setPinVM] = useState<PinViewModel | null>(null);

  const getPin = useCallback(
    async (planId: string, pinId: string): Promise<PinViewModel | null> => {
      const pins = await listPins(planId);
      const pin = pins.find((p) => p.id === pinId) ?? null;
      if (!pin) return null;
      const vm = toPinViewModel(pin);
      setPinVM(vm);
      return vm;
    },
    [listPins],
  );

  return { pinVM, loading, error, getPin };
}
