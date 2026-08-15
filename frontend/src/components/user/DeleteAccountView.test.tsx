import { fireEvent, render, screen } from "@testing-library/react";

import { DeleteAccountView } from "./DeleteAccountView";

const mockOnRequestDelete = jest.fn();
const mockOnCancelDelete = jest.fn();
const mockOnConfirmDelete = jest.fn();
const mockUseDeleteAccountContainer = jest.fn();

jest.mock("@/containers/DeleteAccountContainer", () => ({
  useDeleteAccountContainer: (...args: unknown[]) => mockUseDeleteAccountContainer(...args),
}));

const userId = "user-1";

beforeEach(() => {
  jest.clearAllMocks();
  mockUseDeleteAccountContainer.mockReturnValue({
    isConfirming: false,
    loading: false,
    error: null,
    onRequestDelete: mockOnRequestDelete,
    onCancelDelete: mockOnCancelDelete,
    onConfirmDelete: mockOnConfirmDelete,
  });
});

describe("DeleteAccountView", () => {
  it("passes userId to the container", () => {
    render(<DeleteAccountView userId={userId} />);

    expect(mockUseDeleteAccountContainer).toHaveBeenCalledWith(userId);
  });

  it("renders only the trash icon button when not confirming", () => {
    render(<DeleteAccountView userId={userId} />);

    expect(screen.getByRole("button", { name: "アカウントを削除" })).toBeInTheDocument();
    expect(screen.queryByText("本当に削除しますか？")).not.toBeInTheDocument();
  });

  it("calls onRequestDelete when the trash icon button is clicked", () => {
    render(<DeleteAccountView userId={userId} />);

    fireEvent.click(screen.getByRole("button", { name: "アカウントを削除" }));

    expect(mockOnRequestDelete).toHaveBeenCalledTimes(1);
  });

  it("shows the confirmation prompt with cancel/delete buttons when confirming", () => {
    mockUseDeleteAccountContainer.mockReturnValue({
      isConfirming: true,
      loading: false,
      error: null,
      onRequestDelete: mockOnRequestDelete,
      onCancelDelete: mockOnCancelDelete,
      onConfirmDelete: mockOnConfirmDelete,
    });

    render(<DeleteAccountView userId={userId} />);

    expect(screen.getByText("本当に削除しますか？")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "キャンセル" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "削除する" })).toBeInTheDocument();
  });

  it("calls onCancelDelete when Cancel is clicked", () => {
    mockUseDeleteAccountContainer.mockReturnValue({
      isConfirming: true,
      loading: false,
      error: null,
      onRequestDelete: mockOnRequestDelete,
      onCancelDelete: mockOnCancelDelete,
      onConfirmDelete: mockOnConfirmDelete,
    });

    render(<DeleteAccountView userId={userId} />);
    fireEvent.click(screen.getByRole("button", { name: "キャンセル" }));

    expect(mockOnCancelDelete).toHaveBeenCalledTimes(1);
  });

  it("calls onConfirmDelete when the confirm button is clicked", () => {
    mockUseDeleteAccountContainer.mockReturnValue({
      isConfirming: true,
      loading: false,
      error: null,
      onRequestDelete: mockOnRequestDelete,
      onCancelDelete: mockOnCancelDelete,
      onConfirmDelete: mockOnConfirmDelete,
    });

    render(<DeleteAccountView userId={userId} />);
    fireEvent.click(screen.getByRole("button", { name: "削除する" }));

    expect(mockOnConfirmDelete).toHaveBeenCalledTimes(1);
  });

  it("shows loading state and disables buttons while deleting", () => {
    mockUseDeleteAccountContainer.mockReturnValue({
      isConfirming: true,
      loading: true,
      error: null,
      onRequestDelete: mockOnRequestDelete,
      onCancelDelete: mockOnCancelDelete,
      onConfirmDelete: mockOnConfirmDelete,
    });

    render(<DeleteAccountView userId={userId} />);

    expect(screen.getByRole("button", { name: "削除中…" })).toBeDisabled();
    expect(screen.getByRole("button", { name: "キャンセル" })).toBeDisabled();
  });

  it("shows an error message on failure", () => {
    mockUseDeleteAccountContainer.mockReturnValue({
      isConfirming: true,
      loading: false,
      error: "user not found",
      onRequestDelete: mockOnRequestDelete,
      onCancelDelete: mockOnCancelDelete,
      onConfirmDelete: mockOnConfirmDelete,
    });

    render(<DeleteAccountView userId={userId} />);

    expect(screen.getByText("user not found")).toBeInTheDocument();
  });
});
