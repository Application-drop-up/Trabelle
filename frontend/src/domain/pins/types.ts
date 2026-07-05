export type PinCategory = "restaurant" | "hotel" | "sightseeing" | "transport" | "other";

export type Pin = {
  id: string;
  plan_id: string;
  name: string;
  latitude: number;
  longitude: number;
  category: PinCategory;
  colour: string;
  created_at: string;
  updated_at: string;
};

export type CreatePinInput = {
  name: string;
  latitude: number;
  longitude: number;
  category: PinCategory;
  colour: string;
};

export type UpdatePinInput = {
  category?: PinCategory;
  colour?: string;
};
