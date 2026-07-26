import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import type { PinViewModel, SpotViewModel } from "@/mappers/planMapper";
import { PinCreateOverlay } from "./PinCreateOverlay";

const mockOnCreate = jest.fn();
const mockUsePinCreateContainer = jest.fn();

jest.mock("@/containers/PinCreateContainer", () => ({
  usePinCreateContainer: (...args: unknown[]) => mockUsePinCreateContainer(...args),
}));

const mockSpot: SpotViewModel = {
  placeId: "spot-1",
  name: "Tokyo Tower",
  address: "Tokyo",
  latitude: 35.6586,
  longitude: 139.7454,
};

const mockPinVM: PinViewModel = {
  id: "pin-1",
  planId: "plan-1",
  name: "Tokyo Tower",
  latitude: 35.6586,
  longitude: 139.7454,
  category: "other",
  colour: "#FF5733",
  notes: [],
  createdAt: new Date("2024-01-01"),
  updatedAt: new Date("2024-01-01"),
};

beforeEach(() => {
  jest.clearAllMocks();
  mockUsePinCreateContainer.mockReturnValue({
    onCreate: mockOnCreate,
    loading: false,
    error: null,
  });
});

describe("PinCreateOverlay", () => {
  it("pre-fills name from spot", () => {
    render(
      <PinCreateOverlay
        planId="plan-1"
        spot={mockSpot}
        onCreated={jest.fn()}
        onCancel={jest.fn()}
      />,
    );

    expect(screen.getByDisplayValue("Tokyo Tower")).toBeInTheDocument();
  });

  it("allows editing the name", () => {
    render(
      <PinCreateOverlay
        planId="plan-1"
        spot={mockSpot}
        onCreated={jest.fn()}
        onCancel={jest.fn()}
      />,
    );

    const input = screen.getByDisplayValue("Tokyo Tower");
    fireEvent.change(input, { target: { value: "TT" } });

    expect(screen.getByDisplayValue("TT")).toBeInTheDocument();
  });

  it("calls onCancel when Cancel is clicked", () => {
    const onCancel = jest.fn();
    render(
      <PinCreateOverlay
        planId="plan-1"
        spot={mockSpot}
        onCreated={jest.fn()}
        onCancel={onCancel}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "キャンセル" }));

    expect(onCancel).toHaveBeenCalledTimes(1);
  });

  it("submits with correct values and calls onCreated on success", async () => {
    mockOnCreate.mockResolvedValue(mockPinVM);
    const onCreated = jest.fn();

    render(
      <PinCreateOverlay
        planId="plan-1"
        spot={mockSpot}
        onCreated={onCreated}
        onCancel={jest.fn()}
      />,
    );

    fireEvent.change(screen.getByRole("combobox"), { target: { value: "sightseeing" } });
    fireEvent.click(screen.getByRole("button", { name: "追加" }));

    await waitFor(() => {
      expect(mockOnCreate).toHaveBeenCalledWith(
        expect.objectContaining({
          name: "Tokyo Tower",
          latitude: 35.6586,
          longitude: 139.7454,
          category: "sightseeing",
        }),
      );
      expect(onCreated).toHaveBeenCalledWith(mockPinVM);
    });
  });

  it("disables Add Pin button when name is empty", () => {
    render(
      <PinCreateOverlay
        planId="plan-1"
        spot={mockSpot}
        onCreated={jest.fn()}
        onCancel={jest.fn()}
      />,
    );

    fireEvent.change(screen.getByDisplayValue("Tokyo Tower"), { target: { value: "" } });

    expect(screen.getByRole("button", { name: "追加" })).toBeDisabled();
  });

  it("shows error message on failure", () => {
    mockUsePinCreateContainer.mockReturnValue({
      onCreate: mockOnCreate,
      loading: false,
      error: "Failed to create pin",
    });

    render(
      <PinCreateOverlay
        planId="plan-1"
        spot={mockSpot}
        onCreated={jest.fn()}
        onCancel={jest.fn()}
      />,
    );

    expect(screen.getByText("Failed to create pin")).toBeInTheDocument();
  });
});
