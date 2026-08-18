import { render, screen } from "@testing-library/react";

import type { User } from "@/domain/user/types";
import { ProfileView } from "./ProfileView";

const mockUseUserProfileContainer = jest.fn();

jest.mock("@/containers/UserProfileContainer", () => ({
  useUserProfileContainer: (...args: unknown[]) => mockUseUserProfileContainer(...args),
}));

const mockUser: User = {
  id: "user-1",
  email: "taro@example.com",
  name: "Taro",
  created_at: "2024-01-01T00:00:00Z",
};

beforeEach(() => {
  jest.clearAllMocks();
});

describe("ProfileView", () => {
  it("shows a loading indicator while loading", () => {
    mockUseUserProfileContainer.mockReturnValue({ user: null, loading: true, error: null });

    render(<ProfileView />);

    expect(screen.getByLabelText("読み込み中")).toBeInTheDocument();
  });

  it("shows an error message on failure", () => {
    mockUseUserProfileContainer.mockReturnValue({
      user: null,
      loading: false,
      error: "not authenticated",
    });

    render(<ProfileView />);

    expect(screen.getByText("not authenticated")).toBeInTheDocument();
  });

  it("shows the user's profile on success", () => {
    mockUseUserProfileContainer.mockReturnValue({ user: mockUser, loading: false, error: null });

    render(<ProfileView />);

    expect(screen.getByText("Taro")).toBeInTheDocument();
    expect(screen.getByText("taro@example.com")).toBeInTheDocument();
  });
});
