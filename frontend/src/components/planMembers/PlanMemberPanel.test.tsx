import { fireEvent, render, screen } from "@testing-library/react";

import { PlanMemberPanel } from "./PlanMemberPanel";

const mockOnAddMember = jest.fn();
const mockOnRemoveMember = jest.fn();
const mockUsePlanMemberContainer = jest.fn();

jest.mock("@/containers/PlanMemberContainer", () => ({
  usePlanMemberContainer: (...args: unknown[]) => mockUsePlanMemberContainer(...args),
}));

const memberVMs = [
  { id: "member-1", planId: "plan-1", userId: "user-1", createdAt: "2024-01-01T00:00:00Z" },
];

beforeEach(() => {
  jest.clearAllMocks();
  mockOnAddMember.mockResolvedValue(true);
  mockUsePlanMemberContainer.mockReturnValue({
    memberVMs,
    loading: false,
    error: null,
    onAddMember: mockOnAddMember,
    onRemoveMember: mockOnRemoveMember,
  });
});

describe("PlanMemberPanel", () => {
  it("renders the current members", () => {
    render(<PlanMemberPanel planId="plan-1" />);

    expect(screen.getByText("user-1")).toBeInTheDocument();
  });

  it("calls onAddMember with the entered user id and clears the input on success", async () => {
    render(<PlanMemberPanel planId="plan-1" />);

    fireEvent.change(screen.getByLabelText("ユーザーID"), { target: { value: "user-2" } });
    fireEvent.click(screen.getByRole("button", { name: "追加" }));

    await screen.findByText("user-1");
    expect(mockOnAddMember).toHaveBeenCalledWith("user-2");
  });

  it("disables the add button when the input is empty", () => {
    render(<PlanMemberPanel planId="plan-1" />);

    expect(screen.getByRole("button", { name: "追加" })).toBeDisabled();
  });

  it("calls onRemoveMember when the delete button is clicked", () => {
    render(<PlanMemberPanel planId="plan-1" />);

    fireEvent.click(screen.getByRole("button", { name: "削除" }));

    expect(mockOnRemoveMember).toHaveBeenCalledWith("user-1");
  });

  it("shows a loading indicator while loading", () => {
    mockUsePlanMemberContainer.mockReturnValue({
      memberVMs: [],
      loading: true,
      error: null,
      onAddMember: mockOnAddMember,
      onRemoveMember: mockOnRemoveMember,
    });

    render(<PlanMemberPanel planId="plan-1" />);

    expect(screen.getByLabelText("読み込み中")).toBeInTheDocument();
  });

  it("shows an error message on failure", () => {
    mockUsePlanMemberContainer.mockReturnValue({
      memberVMs,
      loading: false,
      error: "user is already a member of this plan",
      onAddMember: mockOnAddMember,
      onRemoveMember: mockOnRemoveMember,
    });

    render(<PlanMemberPanel planId="plan-1" />);

    expect(screen.getByText("user is already a member of this plan")).toBeInTheDocument();
  });
});
