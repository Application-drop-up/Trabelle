import { render, screen } from "@testing-library/react";

import type { PinViewModel } from "@/mappers/planMapper";
import { PinListView } from "./PinListView";

const mockGetAllPins = jest.fn().mockResolvedValue([]);
const mockUsePinGetAllContainer = jest.fn();

jest.mock("@/containers/PinGetAllContainer", () => ({
  usePinGetAllContainer: (...args: unknown[]) => mockUsePinGetAllContainer(...args),
}));

const mockPinVM: PinViewModel = {
  id: "pin-1",
  planId: "plan-1",
  name: "Tokyo Tower",
  latitude: 35.6586,
  longitude: 139.7454,
  category: "sightseeing",
  colour: "#FF5733",
  notes: [],
  createdAt: new Date("2024-01-01"),
  updatedAt: new Date("2024-01-01"),
};

beforeEach(() => {
  jest.clearAllMocks();
  mockUsePinGetAllContainer.mockReturnValue({
    pinVMs: [],
    loading: false,
    error: null,
    getAllPins: mockGetAllPins,
  });
});

describe("PinListView", () => {
  it("calls getAllPins with planId on mount", () => {
    render(<PinListView planId="plan-1" />);

    expect(mockGetAllPins).toHaveBeenCalledWith("plan-1");
  });

  it("shows loading state", () => {
    mockUsePinGetAllContainer.mockReturnValue({
      pinVMs: [],
      loading: true,
      error: null,
      getAllPins: mockGetAllPins,
    });

    render(<PinListView planId="plan-1" />);

    expect(screen.getByText("Loading pins…")).toBeInTheDocument();
  });

  it("shows error message on failure", () => {
    mockUsePinGetAllContainer.mockReturnValue({
      pinVMs: [],
      loading: false,
      error: "Failed to fetch pins",
      getAllPins: mockGetAllPins,
    });

    render(<PinListView planId="plan-1" />);

    expect(screen.getByText("Failed to fetch pins")).toBeInTheDocument();
  });

  it("renders PinList with fetched pins", () => {
    mockUsePinGetAllContainer.mockReturnValue({
      pinVMs: [mockPinVM],
      loading: false,
      error: null,
      getAllPins: mockGetAllPins,
    });

    render(<PinListView planId="plan-1" />);

    expect(screen.getByText("Tokyo Tower")).toBeInTheDocument();
  });
});
