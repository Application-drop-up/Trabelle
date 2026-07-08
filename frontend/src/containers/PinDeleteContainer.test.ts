import { act, renderHook } from "@testing-library/react";

import { usePinDeleteContainer } from "./PinDeleteContainer";

beforeEach(() => {
  jest.resetAllMocks();
});

describe("usePinDeleteContainer", () => {
  it("initializes with loading false and no error", () => {
    const { result } = renderHook(() => usePinDeleteContainer());

    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
  });

  describe("onDelete", () => {
    it("returns true on successful deletion", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 204,
      } as Response);

      const { result } = renderHook(() => usePinDeleteContainer());

      let ok = false;
      await act(async () => {
        ok = await result.current.onDelete("plan-1", "pin-1");
      });

      expect(ok).toBe(true);
      expect(result.current.error).toBeNull();
    });

    it("returns false and sets error on API failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 404,
        json: async () => ({ error: "pin not found" }),
      } as Response);

      const { result } = renderHook(() => usePinDeleteContainer());

      let ok = true;
      await act(async () => {
        ok = await result.current.onDelete("plan-1", "pin-1");
      });

      expect(ok).toBe(false);
      expect(result.current.error).toBe("pin not found");
    });

    it("calls fetch with correct URL and method", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 204,
      } as Response);

      const { result } = renderHook(() => usePinDeleteContainer());

      await act(async () => {
        await result.current.onDelete("plan-1", "pin-1");
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/plans/plan-1/pins/pin-1"),
        expect.objectContaining({ method: "DELETE" }),
      );
    });
  });
});
