import type { LatLng, MapAdapter, MapMarker } from "./adapter";

type TtMarkerLike = { remove(): void };
type TtMapLike = {
  remove(): void;
  setCenter(center: { lng: number; lat: number }): void;
  setZoom(zoom: number): void;
  fitBounds(bounds: [[number, number], [number, number]], options?: { padding?: number }): void;
};

type TtRuntime = {
  map(options: {
    key: string;
    container: HTMLElement;
    center: { lng: number; lat: number };
    zoom: number;
  }): TtMapLike;
  Marker: new (options?: { color?: string }) => {
    setLngLat(pos: { lng: number; lat: number }): { addTo(map: TtMapLike): TtMarkerLike };
  };
};

export class TomTomMapAdapter implements MapAdapter {
  private map: TtMapLike | null = null;
  private markerMap = new Map<string, TtMarkerLike>();
  private tt: TtRuntime | null = null;

  async initialize(container: HTMLElement, center: LatLng, zoom: number): Promise<void> {
    const mod = await import("@tomtom-international/web-sdk-maps");
    await import("@tomtom-international/web-sdk-maps/dist/maps.css");
    // mod is the tt namespace (esModuleInterop wraps export= as default)
    this.tt = (mod.default ?? mod) as unknown as TtRuntime;
    this.map = this.tt.map({
      key: process.env.NEXT_PUBLIC_TOMTOM_API_KEY ?? "",
      container,
      center: { lng: center.lng, lat: center.lat },
      zoom,
    });
  }

  addMarker({ id, position, color }: MapMarker): void {
    if (!this.tt || !this.map) return;
    const m = new this.tt.Marker({ color })
      .setLngLat({ lng: position.lng, lat: position.lat })
      .addTo(this.map);
    this.markerMap.set(id, m);
  }

  removeMarker(id: string): void {
    const m = this.markerMap.get(id);
    if (!m) return;
    m.remove();
    this.markerMap.delete(id);
  }

  fitToMarkers(markers: MapMarker[]): void {
    if (!this.map || markers.length === 0) return;
    if (markers.length === 1) {
      const { lat, lng } = markers[0].position;
      this.map.setCenter({ lng, lat });
      this.map.setZoom(14);
      return;
    }
    const lats = markers.map((m) => m.position.lat);
    const lngs = markers.map((m) => m.position.lng);
    this.map.fitBounds(
      [
        [Math.min(...lngs), Math.min(...lats)],
        [Math.max(...lngs), Math.max(...lats)],
      ],
      { padding: 60 },
    );
  }

  destroy(): void {
    this.markerMap.forEach((m) => m.remove());
    this.markerMap.clear();
    this.map?.remove();
    this.map = null;
    this.tt = null;
  }
}
