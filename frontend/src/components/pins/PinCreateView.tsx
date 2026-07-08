"use client";

import { useState } from "react";

import { usePinCreateContainer } from "@/containers/PinCreateContainer";
import type { PinCategory } from "@/domain/pins/types";
import type { PinViewModel } from "@/mappers/planMapper";

type Props = {
  planId: string;
  name: string;
  latitude: number;
  longitude: number;
  onCreated: (pin: PinViewModel) => void;
  onCancel: () => void;
};

const CATEGORIES: PinCategory[] = ["restaurant", "hotel", "sightseeing", "transport", "other"];
const COLOURS = ["#FF5733", "#33FF57", "#3357FF", "#FF33A8", "#FFD700"];

export function PinCreateView({ planId, name, latitude, longitude, onCreated, onCancel }: Props) {
  const [category, setCategory] = useState<PinCategory>("other");
  const [colour, setColour] = useState(COLOURS[0]);
  const { onCreate, loading, error } = usePinCreateContainer(planId);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const vm = await onCreate({ name, latitude, longitude, category, colour });
    if (vm) onCreated(vm);
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4 p-4">
      <p className="text-sm font-medium text-zinc-900">{name}</p>
      <div className="flex flex-col gap-1">
        <label className="text-xs text-zinc-500">Category</label>
        <select
          value={category}
          onChange={(e) => setCategory(e.target.value as PinCategory)}
          disabled={loading}
          className="rounded border px-2 py-1 text-sm"
        >
          {CATEGORIES.map((c) => (
            <option key={c} value={c}>
              {c}
            </option>
          ))}
        </select>
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-xs text-zinc-500">Colour</label>
        <div className="flex gap-2">
          {COLOURS.map((c) => (
            <button
              key={c}
              type="button"
              onClick={() => setColour(c)}
              className={`h-6 w-6 rounded-full border-2 ${colour === c ? "border-zinc-900" : "border-transparent"}`}
              style={{ backgroundColor: c }}
            />
          ))}
        </div>
      </div>
      {error && <p className="text-sm text-red-600">{error}</p>}
      <div className="flex gap-2">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 rounded bg-black px-3 py-2 text-sm text-white disabled:opacity-50"
        >
          {loading ? "Adding…" : "Add Pin"}
        </button>
        <button
          type="button"
          onClick={onCancel}
          disabled={loading}
          className="rounded border px-3 py-2 text-sm text-zinc-700 disabled:opacity-50"
        >
          Cancel
        </button>
      </div>
    </form>
  );
}
