"use client";

import { useState } from "react";

import type { CreatePinInput, PinCategory } from "@/domain/pins/types";

type PinCreatePanelProps = {
  name: string;
  latitude: number;
  longitude: number;
  loading: boolean;
  error: string | null;
  onSubmit: (input: CreatePinInput) => void;
  onCancel: () => void;
};

const CATEGORIES: { value: PinCategory; label: string }[] = [
  { value: "restaurant", label: "Restaurant" },
  { value: "hotel", label: "Hotel" },
  { value: "sightseeing", label: "Sightseeing" },
  { value: "transport", label: "Transport" },
  { value: "other", label: "Other" },
];

const DEFAULT_COLOURS = ["#FF5733", "#33A1FF", "#33FF57", "#FF33A1", "#FFD700"];

export function PinCreatePanel({
  name,
  latitude,
  longitude,
  loading,
  error,
  onSubmit,
  onCancel,
}: PinCreatePanelProps) {
  const [category, setCategory] = useState<PinCategory>("other");
  const [colour, setColour] = useState(DEFAULT_COLOURS[0]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({ name, latitude, longitude, category, colour });
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-3 p-4">
      <p className="font-medium">{name}</p>

      <div className="flex flex-col gap-1">
        <label className="text-sm text-gray-600">Category</label>
        <select
          value={category}
          onChange={(e) => setCategory(e.target.value as PinCategory)}
          className="rounded border px-2 py-1 text-sm"
        >
          {CATEGORIES.map((cat) => (
            <option key={cat.value} value={cat.value}>
              {cat.label}
            </option>
          ))}
        </select>
      </div>

      <div className="flex flex-col gap-1">
        <label className="text-sm text-gray-600">Colour</label>
        <div className="flex gap-2">
          {DEFAULT_COLOURS.map((c) => (
            <button
              key={c}
              type="button"
              onClick={() => setColour(c)}
              className="h-6 w-6 rounded-full border-2"
              style={{
                backgroundColor: c,
                borderColor: colour === c ? "black" : "transparent",
              }}
            />
          ))}
        </div>
      </div>

      {error && <p className="text-sm text-red-500">{error}</p>}

      <div className="flex gap-2">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 rounded bg-blue-500 px-3 py-1.5 text-sm text-white disabled:opacity-50"
        >
          {loading ? "Saving..." : "Add Pin"}
        </button>
        <button type="button" onClick={onCancel} className="rounded border px-3 py-1.5 text-sm">
          Cancel
        </button>
      </div>
    </form>
  );
}
