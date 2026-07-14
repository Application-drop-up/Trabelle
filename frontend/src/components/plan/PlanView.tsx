"use client";

import { usePlanContainer } from "@/containers/PlanContainer";

import { PlanLayout } from "./PlanLayout";

type Props = { shareToken: string };

export function PlanView({ shareToken }: Props) {
  const { planVM, loading, error, onDeletePin, applyPinUpdate } = usePlanContainer(shareToken);

  if (loading) {
    return (
      <div className="flex flex-1 items-center justify-center">
        <p className="text-sm text-zinc-500">読み込み中...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-1 items-center justify-center">
        <p className="text-sm text-red-600">{error}</p>
      </div>
    );
  }

  if (!planVM) return null;

  return <PlanLayout planVM={planVM} onDeletePin={onDeletePin} applyPinUpdate={applyPinUpdate} />;
}
