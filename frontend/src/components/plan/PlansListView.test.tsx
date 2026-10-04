import { render, screen } from "@testing-library/react";

import { PlansListView } from "./PlansListView";

const mockUsePlansListContainer = jest.fn();

jest.mock("@/containers/PlansListContainer", () => ({
  usePlansListContainer: () => mockUsePlansListContainer(),
}));

const plans = [
  {
    id: "plan-1",
    share_token: "abc123",
    title: "Tokyo Trip",
    is_public: false,
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-02T00:00:00Z",
  },
];

beforeEach(() => {
  jest.clearAllMocks();
  mockUsePlansListContainer.mockReturnValue({ plans, loading: false, error: null });
});

describe("PlansListView", () => {
  it("renders the plans with links to each one", () => {
    render(<PlansListView />);

    const link = screen.getByRole("link", { name: "Tokyo Trip" });
    expect(link).toBeInTheDocument();
    expect(link).toHaveAttribute("href", "/plans/abc123");
  });

  it("links to the create-plan page", () => {
    render(<PlansListView />);

    expect(screen.getByRole("link", { name: "新規作成" })).toHaveAttribute("href", "/plans/new");
  });

  it("shows a loading indicator while loading", () => {
    mockUsePlansListContainer.mockReturnValue({ plans: [], loading: true, error: null });

    render(<PlansListView />);

    expect(screen.getByLabelText("読み込み中")).toBeInTheDocument();
  });

  it("shows an empty state when there are no plans", () => {
    mockUsePlansListContainer.mockReturnValue({ plans: [], loading: false, error: null });

    render(<PlansListView />);

    expect(screen.getByText("参加しているプランはまだありません")).toBeInTheDocument();
  });

  it("shows an error message on failure", () => {
    mockUsePlansListContainer.mockReturnValue({
      plans: [],
      loading: false,
      error: "Failed to fetch your plans",
    });

    render(<PlansListView />);

    expect(screen.getByText("Failed to fetch your plans")).toBeInTheDocument();
  });
});
