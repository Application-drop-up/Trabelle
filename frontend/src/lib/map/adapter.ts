export type LatLng = { lat: number; lng: number };

export type MapMarker = {
  id: string;
  position: LatLng;
  color: string;
};

export interface MapAdapter {
  initialize(container: HTMLElement, center: LatLng, zoom: number): Promise<void>;
  addMarker(marker: MapMarker): void;
  removeMarker(id: string): void;
  destroy(): void;
}
