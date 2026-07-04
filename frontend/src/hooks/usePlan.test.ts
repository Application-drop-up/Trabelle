import { act, renderHook } from "@testing-library/react";

import type { Plan } from "@/domain/plans/types";
import { usePlan } from "./usePlan";

const mockPlan: Plan = {
  id: "plan-1",
  share_token: "abc123",
  title: "Tokyo Trip",
  pins: [],
  created_at: "2024-01-01T00:00:00Z",
  updated_at: "2024-01-02T00:00:00Z",
};

beforeEach(() => {
  jest.resetAllMocks();
});

describe("usePlan", () => {
  describe("createPlan", () => {
    it("returns created plan and updates state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 201,
        json: async () => mockPlan,
      } as Response);

      const { result } = renderHook(() => usePlan());

      let returned: Plan | null = null;
      await act(async () => {
        returned = await result.current.createPlan("Tokyo Trip");
      });

      expect(returned).toEqual(mockPlan);
      expect(result.current.plan).toEqual(mockPlan);
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it("sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 400,
        statusText: "Bad Request",
        json: async () => ({ error: "invalid request body" }),
      } as Response);

      const { result } = renderHook(() => usePlan());

      await act(async () => {
        await result.current.createPlan("Tokyo Trip");
      });

      expect(result.current.plan).toBeNull();
      expect(result.current.error).toBe("invalid request body");
      expect(result.current.loading).toBe(false);
    });

    it("sets loading true while fetching", async () => {
      let resolveFetch!: (value: unknown) => void;
      global.fetch = jest.fn().mockReturnValue(
        new Promise((resolve) => {
          resolveFetch = resolve;
        }),
      );

      const { result } = renderHook(() => usePlan());

      act(() => {
        result.current.createPlan("Tokyo Trip");
      });

      expect(result.current.loading).toBe(true);

      await act(async () => {
        resolveFetch({
          ok: true,
          status: 201,
          json: async () => mockPlan,
        });
      });

      expect(result.current.loading).toBe(false);
    });
  });

  describe("getPlan", () => {
    it("returns plan and updates state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => mockPlan,
      } as Response);

      const { result } = renderHook(() => usePlan());

      let returned: Plan | null = null;
      await act(async () => {
        returned = await result.current.getPlan("abc123");
      });

      expect(returned).toEqual(mockPlan);
      expect(result.current.plan).toEqual(mockPlan);
      expect(result.current.error).toBeNull();
    });

    it("sets error state when plan is not found", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 404,
        statusText: "Not Found",
        json: async () => ({ error: "plan not found" }),
      } as Response);

      const { result } = renderHook(() => usePlan());

      await act(async () => {
        await result.current.getPlan("invalid-token");
      });

      expect(result.current.plan).toBeNull();
      expect(result.current.error).toBe("plan not found");
    });

    it("calls fetch with the correct URL", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => mockPlan,
      } as Response);

      const { result } = renderHook(() => usePlan());

      await act(async () => {
        await result.current.getPlan("abc123");
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/plans/abc123"),
        expect.any(Object),
      );
    });
  });
});
