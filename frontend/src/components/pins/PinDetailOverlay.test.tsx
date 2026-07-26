import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import type { PinViewModel } from "@/mappers/planMapper";
import { PinDetailOverlay } from "./PinDetailOverlay";

const mockOnUpdate = jest.fn();
const mockUsePinUpdateContainer = jest.fn();

jest.mock("@/containers/PinUpdateContainer", () => ({
  usePinUpdateContainer: (...args: unknown[]) => mockUsePinUpdateContainer(...args),
}));

const mockPin: PinViewModel = {
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

const updatedPin: PinViewModel = { ...mockPin, category: "hotel", colour: "#33FF57" };

beforeEach(() => {
  jest.clearAllMocks();
  mockUsePinUpdateContainer.mockReturnValue({
    onUpdate: mockOnUpdate,
    loading: false,
    error: null,
  });
});

describe("PinDetailOverlay", () => {
  it("shows pin name and category in view mode", () => {
    render(
      <PinDetailOverlay
        planId="plan-1"
        pin={mockPin}
        onUpdated={jest.fn()}
        onDelete={jest.fn()}
        onClose={jest.fn()}
      />,
    );

    expect(screen.getByText("Tokyo Tower")).toBeInTheDocument();
    expect(screen.getByText("観光")).toBeInTheDocument();
  });

  it("calls onClose when close button is clicked", () => {
    const onClose = jest.fn();
    render(
      <PinDetailOverlay
        planId="plan-1"
        pin={mockPin}
        onUpdated={jest.fn()}
        onDelete={jest.fn()}
        onClose={onClose}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "閉じる" }));
    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it("calls onDelete when Delete is clicked", () => {
    const onDelete = jest.fn();
    render(
      <PinDetailOverlay
        planId="plan-1"
        pin={mockPin}
        onUpdated={jest.fn()}
        onDelete={onDelete}
        onClose={jest.fn()}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "削除" }));
    expect(onDelete).toHaveBeenCalledTimes(1);
  });

  it("opens edit mode when Edit is clicked", () => {
    render(
      <PinDetailOverlay
        planId="plan-1"
        pin={mockPin}
        onUpdated={jest.fn()}
        onDelete={jest.fn()}
        onClose={jest.fn()}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "編集" }));

    expect(screen.getByRole("combobox")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "保存" })).toBeInTheDocument();
  });

  it("saves updated values and calls onUpdated", async () => {
    mockOnUpdate.mockResolvedValue(updatedPin);
    const onUpdated = jest.fn();

    render(
      <PinDetailOverlay
        planId="plan-1"
        pin={mockPin}
        onUpdated={onUpdated}
        onDelete={jest.fn()}
        onClose={jest.fn()}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "編集" }));
    fireEvent.change(screen.getByRole("combobox"), { target: { value: "hotel" } });
    fireEvent.click(screen.getByRole("button", { name: "保存" }));

    await waitFor(() => {
      expect(mockOnUpdate).toHaveBeenCalledWith(
        "plan-1",
        "pin-1",
        expect.objectContaining({ category: "hotel" }),
      );
      expect(onUpdated).toHaveBeenCalledWith(updatedPin);
    });
  });

  it("reverts to view mode on Cancel without saving", () => {
    render(
      <PinDetailOverlay
        planId="plan-1"
        pin={mockPin}
        onUpdated={jest.fn()}
        onDelete={jest.fn()}
        onClose={jest.fn()}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "編集" }));
    fireEvent.click(screen.getByRole("button", { name: "キャンセル" }));

    expect(mockOnUpdate).not.toHaveBeenCalled();
    expect(screen.getByText("観光")).toBeInTheDocument();
  });

  it("shows error message when update fails", () => {
    mockUsePinUpdateContainer.mockReturnValue({
      onUpdate: mockOnUpdate,
      loading: false,
      error: "Failed to update pin",
    });

    render(
      <PinDetailOverlay
        planId="plan-1"
        pin={mockPin}
        onUpdated={jest.fn()}
        onDelete={jest.fn()}
        onClose={jest.fn()}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "編集" }));
    expect(screen.getByText("Failed to update pin")).toBeInTheDocument();
  });
});
