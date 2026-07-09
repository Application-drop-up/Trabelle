"use client";

import type { PlanViewModel } from "@/containers/PlanContainer";
import type { MapMarker } from "@/lib/map/adapter";
import { MapView } from "@/components/map/MapView";
import { PinCreateOverlay } from "@/components/pins/PinCreateOverlay";

import { BottomSheet } from "./BottomSheet";
import { PlanPanel } from "./PlanPanel";
import { Sidebar } from "./Sidebar";

// TODO: remove — temporary preview spot
const PREVIEW_SPOT = {
  placeId: "preview-1",
  name: "Tokyo Tower",
  address: "4 Chome-2-8 Shibakoen, Minato City, Tokyo",
  latitude: 35.6586,
  longitude: 139.7454,
};

const DEFAULT_CENTER = { lat: 35.6762, lng: 139.6503 } as const;

type Props = {
  planVM: PlanViewModel;
};

export function PlanLayout({ planVM }: Props) {
  const markers: MapMarker[] = planVM.pins.map((pin) => ({
    id: pin.id,
    position: { lat: pin.latitude, lng: pin.longitude },
    color: pin.colour,
  }));

  return (
    <div className="relative flex flex-1 overflow-hidden">
      <Sidebar>
        <PlanPanel planVM={planVM} />
      </Sidebar>

      <div className="relative flex-1">
        <MapView center={DEFAULT_CENTER} zoom={12} markers={markers} />
        <PinCreateOverlay
          planId={planVM.id}
          spot={PREVIEW_SPOT}
          onCreated={(_pin) => {}}
          onCancel={() => {}}
        />
      </div>

      <BottomSheet>
        <PlanPanel planVM={planVM} />
      </BottomSheet>
    </div>
  );
}
