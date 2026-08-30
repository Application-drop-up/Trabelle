import { act, renderHook } from "@testing-library/react";

import type { PlanMember } from "@/domain/planMembers/types";
import { usePlanMembers } from "./usePlanMembers";

const mockMember: PlanMember = {
  id: "member-1",
  plan_id: "plan-1",
  user_id: "user-1",
  created_at: "2024-01-01T00:00:00Z",
};

beforeEach(() => {
  jest.resetAllMocks();
});

describe("usePlanMembers", () => {
  describe("listMembers", () => {
    it("returns members and updates state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => [mockMember],
      } as Response);

      const { result } = renderHook(() => usePlanMembers());

      let returned: PlanMember[] = [];
      await act(async () => {
        returned = await result.current.listMembers("plan-1");
      });

      expect(returned).toHaveLength(1);
      expect(returned[0]).toEqual(mockMember);
      expect(result.current.members).toHaveLength(1);
      expect(result.current.error).toBeNull();
    });

    it("sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 500,
        statusText: "Internal Server Error",
        json: async () => ({ message: "internal server error" }),
      } as Response);

      const { result } = renderHook(() => usePlanMembers());

      let returned: PlanMember[] = [];
      await act(async () => {
        returned = await result.current.listMembers("plan-1");
      });

      expect(returned).toHaveLength(0);
      expect(result.current.members).toHaveLength(0);
      expect(result.current.error).toBe("internal server error");
    });
  });

  describe("addMember", () => {
    it("returns the member and appends it to state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 201,
        json: async () => mockMember,
      } as Response);

      const { result } = renderHook(() => usePlanMembers());

      let returned: PlanMember | null = null;
      await act(async () => {
        returned = await result.current.addMember("plan-1", "user-1");
      });

      expect(returned).toEqual(mockMember);
      expect(result.current.members).toHaveLength(1);
      expect(result.current.error).toBeNull();
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/api/v1/plans/plan-1/members"),
        expect.objectContaining({ method: "POST" }),
      );
    });

    it("sets error state and returns null when the user is already a member", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 409,
        statusText: "Conflict",
        json: async () => ({ message: "user is already a member of this plan" }),
      } as Response);

      const { result } = renderHook(() => usePlanMembers());

      let returned: PlanMember | null = mockMember;
      await act(async () => {
        returned = await result.current.addMember("plan-1", "user-1");
      });

      expect(returned).toBeNull();
      expect(result.current.error).toBe("user is already a member of this plan");
    });
  });

  describe("removeMember", () => {
    it("returns true and removes the member from state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 204,
        json: async () => null,
      } as Response);

      const { result } = renderHook(() => usePlanMembers());

      let removed = false;
      await act(async () => {
        removed = await result.current.removeMember("plan-1", "user-1");
      });

      expect(removed).toBe(true);
      expect(result.current.error).toBeNull();
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/api/v1/plans/plan-1/members/user-1"),
        expect.objectContaining({ method: "DELETE" }),
      );
    });

    it("returns false and sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 404,
        statusText: "Not Found",
        json: async () => ({ message: "plan member not found" }),
      } as Response);

      const { result } = renderHook(() => usePlanMembers());

      let removed = true;
      await act(async () => {
        removed = await result.current.removeMember("plan-1", "user-1");
      });

      expect(removed).toBe(false);
      expect(result.current.error).toBe("plan member not found");
    });
  });
});
