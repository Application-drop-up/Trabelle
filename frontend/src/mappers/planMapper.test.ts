import type { Note } from "@/domain/notes/types";
import type { Pin, PinWithNotes } from "@/domain/pins/types";
import type { Plan } from "@/domain/plans/types";
import type { Spot } from "@/domain/spots/types";
import {
  toNoteViewModel,
  toPinViewModel,
  toPinWithNotesViewModel,
  toPlanViewModel,
  toSpotViewModel,
} from "./planMapper";

const mockNote: Note = {
  id: "note-1",
  pin_id: "pin-1",
  content: "Test note content",
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-02T00:00:00Z",
};

const mockPin: Pin = {
  id: "pin-1",
  plan_id: "plan-1",
  name: "Tokyo Tower",
  latitude: 35.6586,
  longitude: 139.7454,
  category: "sightseeing",
  colour: "#FF5733",
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-02T00:00:00Z",
};

const mockPinWithNotes: PinWithNotes = {
  ...mockPin,
  notes: [mockNote],
};

const mockPlan: Plan = {
  id: "plan-1",
  share_token: "abc123",
  title: "Tokyo Trip",
  pins: [mockPinWithNotes],
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-02T00:00:00Z",
};

const mockSpot: Spot = {
  place_id: "spot-1",
  name: "Shibuya Crossing",
  address: "Shibuya, Tokyo",
  latitude: 35.6595,
  longitude: 139.7004,
};

describe("toNoteViewModel", () => {
  it("converts snake_case Note to camelCase NoteViewModel", () => {
    const vm = toNoteViewModel(mockNote);

    expect(vm.id).toBe("note-1");
    expect(vm.pinId).toBe("pin-1");
    expect(vm.content).toBe("Test note content");
    expect(vm.createdAt).toEqual(new Date("2024-01-01T00:00:00Z"));
    expect(vm.updatedAt).toEqual(new Date("2024-01-02T00:00:00Z"));
  });

  it("parses created_at and updated_at as Date objects", () => {
    const vm = toNoteViewModel(mockNote);

    expect(vm.createdAt).toBeInstanceOf(Date);
    expect(vm.updatedAt).toBeInstanceOf(Date);
  });
});

describe("toPinViewModel", () => {
  it("converts snake_case Pin to camelCase PinViewModel with empty notes by default", () => {
    const vm = toPinViewModel(mockPin);

    expect(vm.id).toBe("pin-1");
    expect(vm.planId).toBe("plan-1");
    expect(vm.name).toBe("Tokyo Tower");
    expect(vm.latitude).toBe(35.6586);
    expect(vm.longitude).toBe(139.7454);
    expect(vm.category).toBe("sightseeing");
    expect(vm.colour).toBe("#FF5733");
    expect(vm.notes).toEqual([]);
  });

  it("includes provided notes", () => {
    const noteVm = toNoteViewModel(mockNote);
    const vm = toPinViewModel(mockPin, [noteVm]);

    expect(vm.notes).toHaveLength(1);
    expect(vm.notes[0].id).toBe("note-1");
  });

  it("parses created_at and updated_at as Date objects", () => {
    const vm = toPinViewModel(mockPin);

    expect(vm.createdAt).toBeInstanceOf(Date);
    expect(vm.updatedAt).toBeInstanceOf(Date);
  });
});

describe("toPinWithNotesViewModel", () => {
  it("converts PinWithNotes including nested notes", () => {
    const vm = toPinWithNotesViewModel(mockPinWithNotes);

    expect(vm.id).toBe("pin-1");
    expect(vm.notes).toHaveLength(1);
    expect(vm.notes[0].id).toBe("note-1");
    expect(vm.notes[0].pinId).toBe("pin-1");
    expect(vm.notes[0].content).toBe("Test note content");
  });

  it("handles a pin with no notes", () => {
    const pinWithNoNotes: PinWithNotes = { ...mockPin, notes: [] };
    const vm = toPinWithNotesViewModel(pinWithNoNotes);

    expect(vm.notes).toEqual([]);
  });
});

describe("toPlanViewModel", () => {
  it("converts snake_case Plan to camelCase PlanViewModel", () => {
    const vm = toPlanViewModel(mockPlan);

    expect(vm.id).toBe("plan-1");
    expect(vm.shareToken).toBe("abc123");
    expect(vm.title).toBe("Tokyo Trip");
    expect(vm.createdAt).toBeInstanceOf(Date);
    expect(vm.updatedAt).toBeInstanceOf(Date);
  });

  it("converts nested pins with their notes", () => {
    const vm = toPlanViewModel(mockPlan);

    expect(vm.pins).toHaveLength(1);
    expect(vm.pins[0].id).toBe("pin-1");
    expect(vm.pins[0].notes).toHaveLength(1);
    expect(vm.pins[0].notes[0].id).toBe("note-1");
  });

  it("handles a plan with no pins", () => {
    const emptyPlan: Plan = { ...mockPlan, pins: [] };
    const vm = toPlanViewModel(emptyPlan);

    expect(vm.pins).toEqual([]);
  });
});

describe("toSpotViewModel", () => {
  it("converts snake_case Spot to camelCase SpotViewModel", () => {
    const vm = toSpotViewModel(mockSpot);

    expect(vm.placeId).toBe("spot-1");
    expect(vm.name).toBe("Shibuya Crossing");
    expect(vm.address).toBe("Shibuya, Tokyo");
    expect(vm.latitude).toBe(35.6595);
    expect(vm.longitude).toBe(139.7004);
  });
});
