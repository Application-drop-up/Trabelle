"use client";

import { useCallback } from "react";

import { usePins } from "@/hooks/usePins";

type UsePinDeleteContainerReturn = {
  loading: boolean;
  error: string | null;
  onDelete: (planId: string, pinId: string) => Promise<boolean>;
};

export function usePinDeleteContainer(): UsePinDeleteContainerReturn {
  const { deletePin, loading, error } = usePins();

  const onDelete = useCallback(
    async (planId: string, pinId: string): Promise<boolean> => {
      return deletePin(planId, pinId);
    },
    [deletePin],
  );

  return { loading, error, onDelete };
}
