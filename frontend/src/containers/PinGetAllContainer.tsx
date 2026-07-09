"use client";

import { useCallback, useState } from "react";

import { usePins } from "@/hooks/usePins";
import { toPinViewModel, type PinViewModel } from "@/mappers/planMapper";

type UsePinGetAllContainerReturn = {
  pinVMs: PinViewModel[];
  loading: boolean;
  error: string | null;
  getAllPins: (planId: string) => Promise<PinViewModel[]>;
};

export function usePinGetAllContainer(): UsePinGetAllContainerReturn {
  const { listPins, loading, error } = usePins();
  const [pinVMs, setPinVMs] = useState<PinViewModel[]>([]);

  const getAllPins = useCallback(
    async (planId: string): Promise<PinViewModel[]> => {
      const pins = await listPins(planId);
      const vms = pins.map((p) => toPinViewModel(p));
      setPinVMs(vms);
      return vms;
    },
    [listPins],
  );

  return { pinVMs, loading, error, getAllPins };
}
