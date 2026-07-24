"use client";

import type { PlanViewModel } from "@/containers/PlanContainer";
import type { MapMarker } from "@/lib/map/adapter";
import { MapView } from "@/components/map/MapView";

import { BottomSheet } from "./BottomSheet";
import { PlanPanel } from "./PlanPanel";
import { Sidebar } from "./Sidebar";

const DEFAULT_CENTER = { lat: 35.6762, lng: 139.6503 } as const;

type Props = {
  planVM: PlanViewModel;
  onCreateNote: (pinId: string, content: string) => Promise<void>;
};

export function PlanLayout({ planVM, onCreateNote }: Props) {
  const markers: MapMarker[] = planVM.pins.map((pin) => ({
    id: pin.id,
    position: { lat: pin.latitude, lng: pin.longitude },
    color: pin.colour,
  }));

  return (
    <div className="relative flex flex-1 overflow-hidden">
      <Sidebar>
        <PlanPanel planVM={planVM} onCreateNote={onCreateNote} />
      </Sidebar>

      <div className="flex-1">
        <MapView center={DEFAULT_CENTER} zoom={12} markers={markers} />
      </div>

      <BottomSheet>
        <PlanPanel planVM={planVM} onCreateNote={onCreateNote} />
      </BottomSheet>
    </div>
  );
}
