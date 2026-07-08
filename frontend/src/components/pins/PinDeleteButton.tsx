"use client";

import { usePinDeleteContainer } from "@/containers/PinDeleteContainer";

type Props = {
  planId: string;
  pinId: string;
  onDeleted: () => void;
};

export function PinDeleteButton({ planId, pinId, onDeleted }: Props) {
  const { onDelete, loading } = usePinDeleteContainer();

  const handleClick = async () => {
    const ok = await onDelete(planId, pinId);
    if (ok) onDeleted();
  };

  return (
    <button
      type="button"
      onClick={handleClick}
      disabled={loading}
      className="rounded border border-red-300 px-3 py-1.5 text-sm text-red-600 hover:bg-red-50 disabled:opacity-50"
    >
      {loading ? "Deleting…" : "Delete"}
    </button>
  );
}
