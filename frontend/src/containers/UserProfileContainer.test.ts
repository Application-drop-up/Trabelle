import { renderHook, waitFor } from "@testing-library/react";

import type { User } from "@/domain/user/types";
import { useUserProfileContainer } from "./UserProfileContainer";

const mockUser: User = {
  id: "user-1",
  email: "taro@example.com",
  name: "Taro",
  created_at: "2024-01-01T00:00:00Z",
};

beforeEach(() => {
  jest.resetAllMocks();
});

describe("useUserProfileContainer", () => {
  it("fetches the current user on mount and exposes it", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => mockUser,
    } as Response);

    const { result } = renderHook(() => useUserProfileContainer());

    expect(result.current.loading).toBe(true);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.user).toEqual(mockUser);
    expect(result.current.error).toBeNull();
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining("/api/v1/user/me"),
      expect.anything(),
    );
  });

  it("exposes an error when the fetch fails", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      status: 401,
      statusText: "Unauthorized",
      json: async () => ({ error: "not authenticated" }),
    } as Response);

    const { result } = renderHook(() => useUserProfileContainer());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.user).toBeNull();
    expect(result.current.error).toBe("not authenticated");
  });
});
