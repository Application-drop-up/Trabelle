import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { PinDeleteButton } from "./PinDeleteButton";

const mockOnDelete = jest.fn();
const mockUsePinDeleteContainer = jest.fn();

jest.mock("@/containers/PinDeleteContainer", () => ({
  usePinDeleteContainer: (...args: unknown[]) => mockUsePinDeleteContainer(...args),
}));

beforeEach(() => {
  jest.clearAllMocks();
  mockUsePinDeleteContainer.mockReturnValue({
    onDelete: mockOnDelete,
    loading: false,
    error: null,
  });
});

describe("PinDeleteButton", () => {
  it("renders a Delete button", () => {
    render(<PinDeleteButton planId="plan-1" pinId="pin-1" onDeleted={jest.fn()} />);

    expect(screen.getByRole("button", { name: /delete/i })).toBeInTheDocument();
  });

  it("calls onDelete with planId and pinId on click, then calls onDeleted on success", async () => {
    mockOnDelete.mockResolvedValue(true);
    const onDeleted = jest.fn();

    render(<PinDeleteButton planId="plan-1" pinId="pin-1" onDeleted={onDeleted} />);

    fireEvent.click(screen.getByRole("button", { name: /delete/i }));

    await waitFor(() => {
      expect(mockOnDelete).toHaveBeenCalledWith("plan-1", "pin-1");
      expect(onDeleted).toHaveBeenCalledTimes(1);
    });
  });

  it("does not call onDeleted when onDelete returns false", async () => {
    mockOnDelete.mockResolvedValue(false);
    const onDeleted = jest.fn();

    render(<PinDeleteButton planId="plan-1" pinId="pin-1" onDeleted={onDeleted} />);

    fireEvent.click(screen.getByRole("button", { name: /delete/i }));

    await waitFor(() => {
      expect(mockOnDelete).toHaveBeenCalled();
      expect(onDeleted).not.toHaveBeenCalled();
    });
  });

  it("shows 'Deleting…' and disables button when loading", () => {
    mockUsePinDeleteContainer.mockReturnValue({
      onDelete: mockOnDelete,
      loading: true,
      error: null,
    });

    render(<PinDeleteButton planId="plan-1" pinId="pin-1" onDeleted={jest.fn()} />);

    const button = screen.getByRole("button", { name: /deleting/i });
    expect(button).toBeDisabled();
  });
});
