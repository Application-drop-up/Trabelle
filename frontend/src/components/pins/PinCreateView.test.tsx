import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import type { PinViewModel } from "@/mappers/planMapper";
import { PinCreateView } from "./PinCreateView";

const mockOnCreate = jest.fn();

jest.mock("@/containers/PinCreateContainer", () => ({
  usePinCreateContainer: () => ({
    onCreate: mockOnCreate,
    loading: false,
    error: null,
  }),
}));

const defaultProps = {
  planId: "plan-1",
  name: "Tokyo Tower",
  latitude: 35.6586,
  longitude: 139.7454,
  onCreated: jest.fn(),
  onCancel: jest.fn(),
};

beforeEach(() => {
  jest.clearAllMocks();
});

describe("PinCreateView", () => {
  it("renders pin name, category select, colour buttons, and action buttons", () => {
    render(<PinCreateView {...defaultProps} />);

    expect(screen.getByText("Tokyo Tower")).toBeInTheDocument();
    expect(screen.getByRole("combobox")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /add pin/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /cancel/i })).toBeInTheDocument();
  });

  it("calls onCancel when cancel button is clicked", () => {
    const onCancel = jest.fn();
    render(<PinCreateView {...defaultProps} onCancel={onCancel} />);

    fireEvent.click(screen.getByRole("button", { name: /cancel/i }));

    expect(onCancel).toHaveBeenCalledTimes(1);
  });

  it("calls onCreate with selected values and invokes onCreated on success", async () => {
    const mockVM: PinViewModel = {
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
    mockOnCreate.mockResolvedValue(mockVM);
    const onCreated = jest.fn();

    render(<PinCreateView {...defaultProps} onCreated={onCreated} />);

    fireEvent.change(screen.getByRole("combobox"), { target: { value: "sightseeing" } });
    fireEvent.click(screen.getByRole("button", { name: /add pin/i }));

    await waitFor(() => {
      expect(mockOnCreate).toHaveBeenCalledWith(
        expect.objectContaining({
          name: "Tokyo Tower",
          latitude: 35.6586,
          longitude: 139.7454,
          category: "sightseeing",
        }),
      );
      expect(onCreated).toHaveBeenCalledWith(mockVM);
    });
  });

  it("does not call onCreated when onCreate returns null", async () => {
    mockOnCreate.mockResolvedValue(null);
    const onCreated = jest.fn();

    render(<PinCreateView {...defaultProps} onCreated={onCreated} />);

    fireEvent.click(screen.getByRole("button", { name: /add pin/i }));

    await waitFor(() => {
      expect(mockOnCreate).toHaveBeenCalled();
      expect(onCreated).not.toHaveBeenCalled();
    });
  });
});
