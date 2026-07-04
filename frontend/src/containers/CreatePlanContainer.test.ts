import { act, renderHook } from "@testing-library/react";

import type { Plan } from "@/domain/plans/types";
import { useCreatePlanContainer } from "./CreatePlanContainer";

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

describe("useCreatePlanContainer", () => {
  it("initializes with empty title and no error", () => {
    const { result } = renderHook(() => useCreatePlanContainer());

    expect(result.current.title).toBe("");
    expect(result.current.error).toBeNull();
    expect(result.current.loading).toBe(false);
  });

  describe("onChangeTitle", () => {
    it("updates title state", () => {
      const { result } = renderHook(() => useCreatePlanContainer());

      act(() => {
        result.current.onChangeTitle("Tokyo Trip");
      });

      expect(result.current.title).toBe("Tokyo Trip");
    });
  });

  describe("onSubmit", () => {
    it("returns null when title is blank", async () => {
      global.fetch = jest.fn();
      const { result } = renderHook(() => useCreatePlanContainer());

      let returned: Plan | null = undefined as unknown as Plan | null;
      await act(async () => {
        returned = await result.current.onSubmit();
      });

      expect(returned).toBeNull();
      expect(global.fetch).not.toHaveBeenCalled();
    });

    it("returns null when title is whitespace only", async () => {
      global.fetch = jest.fn();
      const { result } = renderHook(() => useCreatePlanContainer());

      act(() => {
        result.current.onChangeTitle("   ");
      });

      let returned: Plan | null = undefined as unknown as Plan | null;
      await act(async () => {
        returned = await result.current.onSubmit();
      });

      expect(returned).toBeNull();
      expect(global.fetch).not.toHaveBeenCalled();
    });

    it("calls createPlan with trimmed title and returns plan on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 201,
        json: async () => mockPlan,
      } as Response);

      const { result } = renderHook(() => useCreatePlanContainer());

      act(() => {
        result.current.onChangeTitle("  Tokyo Trip  ");
      });

      let returned: Plan | null = null;
      await act(async () => {
        returned = await result.current.onSubmit();
      });

      expect(returned).toEqual(mockPlan);
      const fetchCall = (fetch as jest.Mock).mock.calls[0];
      expect(JSON.parse(fetchCall[1].body)).toEqual({ title: "Tokyo Trip" });
    });

    it("sets error state when createPlan fails", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 400,
        statusText: "Bad Request",
        json: async () => ({ error: "title is required" }),
      } as Response);

      const { result } = renderHook(() => useCreatePlanContainer());

      act(() => {
        result.current.onChangeTitle("Tokyo Trip");
      });

      await act(async () => {
        await result.current.onSubmit();
      });

      expect(result.current.error).toBe("title is required");
    });
  });
});
