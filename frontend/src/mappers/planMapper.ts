import type { Pin } from "@/domain/pins/types";
import type { Spot } from "@/domain/spots/types";

export type SpotViewModel = {
  placeId: string;
  name: string;
  address: string;
  latitude: number;
  longitude: number;
};

export type PinViewModel = {
  id: string;
  planId: string;
  name: string;
  latitude: number;
  longitude: number;
  category: string;
  colour: string;
};

export function toSpotViewModel(spot: Spot): SpotViewModel {
  return {
    placeId: spot.place_id,
    name: spot.name,
    address: spot.address,
    latitude: spot.latitude,
    longitude: spot.longitude,
  };
}

export function toPinViewModel(pin: Pin): PinViewModel {
  return {
    id: pin.id,
    planId: pin.plan_id,
    name: pin.name,
    latitude: pin.latitude,
    longitude: pin.longitude,
    category: pin.category,
    colour: pin.colour,
  };
}
