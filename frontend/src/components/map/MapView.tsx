"use client";

import dynamic from "next/dynamic";
import { useMemo } from "react";

import type { LatLng, MapMarker } from "@/lib/map/adapter";
import { TomTomMapAdapter } from "@/lib/map/tomtomAdapter";

const MapCanvas = dynamic(() => import("./MapCanvas"), { ssr: false });

type Props = {
  center: LatLng;
  zoom?: number;
  markers?: MapMarker[];
  onMarkerClick?: (id: string) => void;
};

export function MapView({ center, zoom = 12, markers = [], onMarkerClick }: Props) {
  const adapter = useMemo(() => new TomTomMapAdapter(), []);

  return (
    <MapCanvas
      adapter={adapter}
      center={center}
      zoom={zoom}
      markers={markers}
      onMarkerClick={onMarkerClick}
    />
  );
}
