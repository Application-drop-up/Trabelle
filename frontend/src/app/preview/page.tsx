"use client";

import { PinCreateOverlay } from "@/components/pins/PinCreateOverlay";

const PREVIEW_SPOT = {
  placeId: "preview-1",
  name: "Tokyo Tower",
  address: "4 Chome-2-8 Shibakoen, Minato City, Tokyo",
  latitude: 35.6586,
  longitude: 139.7454,
};

export default function PreviewPage() {
  return (
    <div className="relative h-screen w-full bg-zinc-100">
      <div className="flex h-full items-center justify-center text-zinc-400">Map area</div>
      <PinCreateOverlay
        planId="preview-plan"
        spot={PREVIEW_SPOT}
        onCreated={(_pin) => {}}
        onCancel={() => {}}
      />
    </div>
  );
}
