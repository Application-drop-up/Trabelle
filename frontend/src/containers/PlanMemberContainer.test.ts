import { act, renderHook, waitFor } from "@testing-library/react";

import type { PlanMember } from "@/domain/planMembers/types";
import { usePlanMemberContainer } from "./PlanMemberContainer";

const mockMember: PlanMember = {
  id: "member-1",
  plan_id: "plan-1",
  user_id: "user-1",
  created_at: "2024-01-01T00:00:00Z",
};

beforeEach(() => {
  jest.resetAllMocks();
});

describe("usePlanMemberContainer", () => {
  it("lists members on mount", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => [mockMember],
    } as Response);

    const { result } = renderHook(() => usePlanMemberContainer("plan-1"));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.memberVMs).toEqual([
      { id: "member-1", planId: "plan-1", userId: "user-1", createdAt: "2024-01-01T00:00:00Z" },
    ]);
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining("/api/v1/plans/plan-1/members"),
      expect.anything(),
    );
  });

  it("adds a member and returns true on success", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      status: 201,
      json: async () => mockMember,
    } as Response);

    const { result } = renderHook(() => usePlanMemberContainer("plan-1"));

    let added = false;
    await act(async () => {
      added = await result.current.onAddMember("user-1");
    });

    expect(added).toBe(true);
    await waitFor(() => {
      expect(result.current.memberVMs.some((vm) => vm.userId === "user-1")).toBe(true);
    });
  });

  it("returns false when adding a member fails", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      status: 409,
      statusText: "Conflict",
      json: async () => ({ message: "user is already a member of this plan" }),
    } as Response);

    const { result } = renderHook(() => usePlanMemberContainer("plan-1"));

    let added = true;
    await act(async () => {
      added = await result.current.onAddMember("user-1");
    });

    expect(added).toBe(false);
    await waitFor(() => {
      expect(result.current.error).toBe("user is already a member of this plan");
    });
  });

  it("removes a member and returns true on success", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      status: 204,
      json: async () => null,
    } as Response);

    const { result } = renderHook(() => usePlanMemberContainer("plan-1"));

    let removed = false;
    await act(async () => {
      removed = await result.current.onRemoveMember("user-1");
    });

    expect(removed).toBe(true);
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining("/api/v1/plans/plan-1/members/user-1"),
      expect.objectContaining({ method: "DELETE" }),
    );
  });
});
