"use client";

import { useEffect, useRef, useState } from "react";

import type { LatLng, MapAdapter, MapMarker } from "@/lib/map/adapter";

type Props = {
  adapter: MapAdapter;
  center: LatLng;
  zoom: number;
  markers: MapMarker[];
};

export default function MapCanvas({ adapter, center, zoom, markers }: Props) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [ready, setReady] = useState(false);
  const prevMarkerIds = useRef<Set<string>>(new Set());

  useEffect(() => {
    if (!containerRef.current) return;
    let alive = true;
    adapter.initialize(containerRef.current, center, zoom).then(() => {
      if (alive) setReady(true);
    });
    return () => {
      alive = false;
      adapter.destroy();
      setReady(false);
      prevMarkerIds.current = new Set();
    };
    // center and zoom are intentionally excluded: repositioning the map
    // after initialization is a separate concern not needed for the skeleton.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [adapter]);

  useEffect(() => {
    if (!ready) return;
    const currentIds = new Set(markers.map((m) => m.id));
    for (const id of prevMarkerIds.current) {
      if (!currentIds.has(id)) adapter.removeMarker(id);
    }
    for (const marker of markers) {
      if (!prevMarkerIds.current.has(marker.id)) adapter.addMarker(marker);
    }
    prevMarkerIds.current = currentIds;
  }, [ready, markers, adapter]);

  return <div ref={containerRef} className="h-full w-full" />;
}
