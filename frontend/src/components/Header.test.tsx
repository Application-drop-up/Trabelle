import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import type { User } from "@/domain/user/types";
import { Header } from "./Header";

const mockOnLogout = jest.fn();
const mockPush = jest.fn();
const mockUseUserContext = jest.fn();

jest.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
}));

jest.mock("@/components/UserProvider", () => ({
  useUserContext: (...args: unknown[]) => mockUseUserContext(...args),
}));

jest.mock("@/containers/LogoutContainer", () => ({
  useLogoutContainer: () => ({ loading: false, error: null, onLogout: mockOnLogout }),
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

describe("Header", () => {
  it("shows only the app name when logged out", () => {
    mockUseUserContext.mockReturnValue({ user: null });

    render(<Header />);

    expect(screen.getByText("Trabelle")).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "ログアウト" })).not.toBeInTheDocument();
  });

  it("shows the user name and logout button when logged in", () => {
    mockUseUserContext.mockReturnValue({ user: mockUser });

    render(<Header />);

    expect(screen.getByText("Taro さん")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "ログアウト" })).toBeInTheDocument();
  });

  it("redirects to /login when logout succeeds", async () => {
    mockUseUserContext.mockReturnValue({ user: mockUser });
    mockOnLogout.mockResolvedValue(true);

    render(<Header />);
    fireEvent.click(screen.getByRole("button", { name: "ログアウト" }));

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith("/login");
    });
  });

  it("does not redirect when logout fails", async () => {
    mockUseUserContext.mockReturnValue({ user: mockUser });
    mockOnLogout.mockResolvedValue(false);

    render(<Header />);
    fireEvent.click(screen.getByRole("button", { name: "ログアウト" }));

    await waitFor(() => {
      expect(mockOnLogout).toHaveBeenCalledTimes(1);
    });
    expect(mockPush).not.toHaveBeenCalled();
  });
});
