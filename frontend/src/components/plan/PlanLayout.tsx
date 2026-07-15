"use client";

import { useState } from "react";

import type { PinViewModel, PlanViewModel } from "@/containers/PlanContainer";
import type { MapMarker } from "@/lib/map/adapter";
import { MapView } from "@/components/map/MapView";
import { PinDetailOverlay } from "@/components/pins/PinDetailOverlay";

import { BottomSheet } from "./BottomSheet";
import { PlanPanel } from "./PlanPanel";
import { Sidebar } from "./Sidebar";

const DEFAULT_CENTER = { lat: 35.6762, lng: 139.6503 } as const;

type Props = {
  planVM: PlanViewModel;
  onDeletePin: (pinId: string) => Promise<void>;
  applyPinUpdate: (pin: PinViewModel) => void;
};

export function PlanLayout({ planVM, onDeletePin, applyPinUpdate }: Props) {
  const [selectedPinId, setSelectedPinId] = useState<string | null>(null);

  const markers: MapMarker[] = planVM.pins.map((pin) => ({
    id: pin.id,
    position: { lat: pin.latitude, lng: pin.longitude },
    color: pin.colour,
  }));

  const selectedPin = planVM.pins.find((pin) => pin.id === selectedPinId) ?? null;

  return (
    <div className="relative flex flex-1 overflow-hidden">
      <Sidebar>
        <PlanPanel planVM={planVM} />
      </Sidebar>

      <div className="flex-1">
        <MapView
          center={DEFAULT_CENTER}
          zoom={12}
          markers={markers}
          onMarkerClick={setSelectedPinId}
        />
      </div>

      <BottomSheet>
        <PlanPanel planVM={planVM} />
      </BottomSheet>

      {selectedPin && (
        <PinDetailOverlay
          planId={planVM.id}
          pin={selectedPin}
          onUpdated={(pin: PinViewModel) => applyPinUpdate(pin)}
          onDelete={async () => {
            await onDeletePin(selectedPin.id);
            setSelectedPinId(null);
          }}
          onClose={() => setSelectedPinId(null)}
        />
      )}
    </div>
  );
}
