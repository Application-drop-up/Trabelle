"use client";

import { usePlanContainer } from "@/containers/PlanContainer";
import type { MapMarker } from "@/lib/map/adapter";
import { MapView } from "@/components/map/MapView";

const DEFAULT_CENTER = { lat: 35.6762, lng: 139.6503 } as const;

type Props = { shareToken: string };

export function PlanView({ shareToken }: Props) {
  const { planVM, loading, error } = usePlanContainer(shareToken);

  if (loading) {
    return (
      <div className="flex flex-1 items-center justify-center">
        <p className="text-sm text-gray-500">Loading…</p>
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

  const markers: MapMarker[] = planVM.pins.map((pin) => ({
    id: pin.id,
    position: { lat: pin.latitude, lng: pin.longitude },
    color: pin.colour,
  }));

  return (
    <div className="flex flex-1 overflow-hidden">
      <aside className="flex w-72 flex-col gap-4 overflow-y-auto border-r p-4">
        <h1 className="text-xl font-semibold">{planVM.title}</h1>
        <ul className="flex flex-col gap-2">
          {planVM.pins.map((pin) => (
            <li key={pin.id} className="rounded border p-3 text-sm">
              <p className="font-medium">{pin.name}</p>
              <p className="text-gray-500">{pin.category}</p>
            </li>
          ))}
        </ul>
        {planVM.pins.length === 0 && <p className="text-sm text-gray-400">No pins yet.</p>}
      </aside>
      <main className="flex flex-1">
        <MapView center={DEFAULT_CENTER} zoom={12} markers={markers} />
      </main>
    </div>
  );
}
