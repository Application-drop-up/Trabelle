import { renderHook, act } from "@testing-library/react";
import { usePinUpdateContainer } from "./PinUpdateContainer";
import type { PinViewModel } from "@/mappers/planMapper";

const mockUpdatePin = jest.fn();
const mockUsePins = jest.fn();

jest.mock("@/hooks/usePins", () => ({
  usePins: (...args: unknown[]) => mockUsePins(...args),
}));

const mockPinVM: PinViewModel = {
  id: "pin-1",
  planId: "plan-1",
  name: "Tokyo Tower",
  latitude: 35.6586,
  longitude: 139.7454,
  category: "sightseeing",
  colour: "#3357FF",
  notes: [],
  createdAt: new Date("2024-01-01"),
  updatedAt: new Date("2024-01-01"),
};

beforeEach(() => {
  jest.clearAllMocks();
  mockUsePins.mockReturnValue({
    updatePin: mockUpdatePin,
    loading: false,
    error: null,
  });
});

describe("usePinUpdateContainer", () => {
  it("returns PinViewModel on successful update", async () => {
    mockUpdatePin.mockResolvedValue({
      id: "pin-1",
      plan_id: "plan-1",
      name: "Tokyo Tower",
      latitude: 35.6586,
      longitude: 139.7454,
      category: "sightseeing",
      colour: "#3357FF",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z",
    });

    const { result } = renderHook(() => usePinUpdateContainer());
    let vm: PinViewModel | null = null;

    await act(async () => {
      vm = await result.current.onUpdate("plan-1", "pin-1", {
        category: "sightseeing",
        colour: "#3357FF",
      });
    });

    expect(mockUpdatePin).toHaveBeenCalledWith("plan-1", "pin-1", {
      category: "sightseeing",
      colour: "#3357FF",
    });
    expect(vm).toMatchObject(mockPinVM);
  });

  it("returns null when updatePin fails", async () => {
    mockUpdatePin.mockResolvedValue(null);

    const { result } = renderHook(() => usePinUpdateContainer());
    let vm: PinViewModel | null = null;

    await act(async () => {
      vm = await result.current.onUpdate("plan-1", "pin-1", { category: "hotel" });
    });

    expect(vm).toBeNull();
  });

  it("exposes loading and error from usePins", () => {
    mockUsePins.mockReturnValue({
      updatePin: mockUpdatePin,
      loading: true,
      error: "Network error",
    });

    const { result } = renderHook(() => usePinUpdateContainer());

    expect(result.current.loading).toBe(true);
    expect(result.current.error).toBe("Network error");
  });
});
