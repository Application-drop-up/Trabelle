"use client";

import { useState } from "react";

import { usePinCreateContainer } from "@/containers/PinCreateContainer";
import type { PinCategory } from "@/domain/pins/types";
import type { PinViewModel, SpotViewModel } from "@/mappers/planMapper";

type Props = {
  planId: string;
  spot: SpotViewModel;
  onCreated: (pin: PinViewModel) => void;
  onCancel: () => void;
};

const CATEGORIES: { value: PinCategory; label: string }[] = [
  { value: "restaurant", label: "Restaurant" },
  { value: "hotel", label: "Hotel" },
  { value: "sightseeing", label: "Sightseeing" },
  { value: "transport", label: "Transport" },
  { value: "other", label: "Other" },
];

const COLOURS = ["#FF5733", "#33FF57", "#3357FF", "#FF33A8", "#FFD700"];

export function PinCreateOverlay({ planId, spot, onCreated, onCancel }: Props) {
  const [name, setName] = useState(spot.name);
  const [category, setCategory] = useState<PinCategory>("other");
  const [colour, setColour] = useState(COLOURS[0]);
  const { onCreate, loading, error } = usePinCreateContainer(planId);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const vm = await onCreate({
      name: name.trim(),
      latitude: spot.latitude,
      longitude: spot.longitude,
      category,
      colour,
    });
    if (vm) onCreated(vm);
  };

  return (
    <div className="absolute bottom-4 left-4 z-10 w-72 rounded-xl bg-white p-4 shadow-lg">
      <form onSubmit={handleSubmit} className="flex flex-col gap-3">
        <div className="flex flex-col gap-1">
          <label className="text-xs text-zinc-500">Name</label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            disabled={loading}
            className="rounded border px-2 py-1.5 text-sm outline-none focus:ring-2 focus:ring-zinc-300"
          />
        </div>
        <div className="flex flex-col gap-1">
          <label className="text-xs text-zinc-500">Category</label>
          <select
            value={category}
            onChange={(e) => setCategory(e.target.value as PinCategory)}
            disabled={loading}
            className="rounded border px-2 py-1.5 text-sm"
          >
            {CATEGORIES.map((c) => (
              <option key={c.value} value={c.value}>
                {c.label}
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
                className={`h-6 w-6 rounded-full border-2 transition-transform ${
                  colour === c ? "scale-110 border-zinc-900" : "border-transparent"
                }`}
                style={{ backgroundColor: c }}
              />
            ))}
          </div>
        </div>
        {error && <p className="text-xs text-red-600">{error}</p>}
        <div className="flex gap-2 pt-1">
          <button
            type="submit"
            disabled={loading || !name.trim()}
            className="flex-1 rounded-lg bg-black px-3 py-2 text-sm font-medium text-white disabled:opacity-50"
          >
            {loading ? "Adding…" : "Add Pin"}
          </button>
          <button
            type="button"
            onClick={onCancel}
            disabled={loading}
            className="rounded-lg border px-3 py-2 text-sm text-zinc-700 disabled:opacity-50"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}
