"use client";

import { useEffect } from "react";

import { usePinGetAllContainer } from "@/containers/PinGetAllContainer";
import { PinList } from "@/components/plan/PinList";

type Props = {
  planId: string;
};

export function PinListView({ planId }: Props) {
  const { pinVMs, loading, error, getAllPins } = usePinGetAllContainer();

  useEffect(() => {
    getAllPins(planId);
  }, [planId, getAllPins]);

  if (loading) {
    return <p className="text-sm text-zinc-400">Loading pins…</p>;
  }

  if (error) {
    return <p className="text-sm text-red-600">{error}</p>;
  }

  return <PinList pins={pinVMs} />;
}
