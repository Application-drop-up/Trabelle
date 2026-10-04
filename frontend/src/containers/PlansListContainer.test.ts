import { renderHook, waitFor } from "@testing-library/react";

import { UserProvider } from "@/components/UserProvider";
import type { PlanSummary } from "@/domain/plans/types";
import type { User } from "@/domain/user/types";
import { usePlansListContainer } from "./PlansListContainer";

const mockUser: User = {
  id: "user-1",
  email: "taro@example.com",
  name: "Taro",
  created_at: "2024-01-01T00:00:00Z",
};

const mockPlan: PlanSummary = {
  id: "plan-1",
  share_token: "abc123",
  title: "Tokyo Trip",
  is_public: false,
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-02T00:00:00Z",
};

beforeEach(() => {
  jest.resetAllMocks();
});

describe("usePlansListContainer", () => {
  it("fetches the current user, then lists their plans", async () => {
    global.fetch = jest.fn().mockImplementation((url: string) => {
      if (url.includes("/user/me")) {
        return Promise.resolve({ ok: true, status: 200, json: async () => mockUser } as Response);
      }
      return Promise.resolve({ ok: true, status: 200, json: async () => [mockPlan] } as Response);
    });

    const { result } = renderHook(() => usePlansListContainer(), { wrapper: UserProvider });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.plans).toEqual([mockPlan]);
    expect(result.current.error).toBeNull();
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining("/api/v1/user/user-1/plans"),
      expect.anything(),
    );
  });

  it("exposes an error when fetching the current user fails", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      status: 401,
      statusText: "Unauthorized",
      json: async () => ({ message: "not authenticated" }),
    } as Response);

    const { result } = renderHook(() => usePlansListContainer(), { wrapper: UserProvider });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.plans).toEqual([]);
    expect(result.current.error).toBe("not authenticated");
  });
});
