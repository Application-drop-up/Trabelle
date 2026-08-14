import { act, renderHook } from "@testing-library/react";

import { useDeleteAccountContainer } from "./DeleteAccountContainer";

const userId = "user-1";

beforeEach(() => {
  jest.resetAllMocks();
});

describe("useDeleteAccountContainer", () => {
  it("initializes with isConfirming false and no error", () => {
    const { result } = renderHook(() => useDeleteAccountContainer(userId));

    expect(result.current.isConfirming).toBe(false);
    expect(result.current.error).toBeNull();
    expect(result.current.loading).toBe(false);
  });

  describe("onRequestDelete / onCancelDelete", () => {
    it("toggles isConfirming", () => {
      const { result } = renderHook(() => useDeleteAccountContainer(userId));

      act(() => {
        result.current.onRequestDelete();
      });
      expect(result.current.isConfirming).toBe(true);

      act(() => {
        result.current.onCancelDelete();
      });
      expect(result.current.isConfirming).toBe(false);
    });
  });

  describe("onConfirmDelete", () => {
    it("calls deleteUser with the userId and returns true on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 204,
      } as Response);

      const { result } = renderHook(() => useDeleteAccountContainer(userId));

      let returned = false;
      await act(async () => {
        returned = await result.current.onConfirmDelete();
      });

      expect(returned).toBe(true);
      const fetchCall = (fetch as jest.Mock).mock.calls[0];
      expect(fetchCall[0]).toContain(`/api/v1/user/${userId}`);
      expect(fetchCall[1]).toEqual(expect.objectContaining({ method: "DELETE" }));
    });

    it("returns false and sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 404,
        statusText: "Not Found",
        json: async () => ({ error: "user not found" }),
      } as Response);

      const { result } = renderHook(() => useDeleteAccountContainer(userId));

      let returned = true;
      await act(async () => {
        returned = await result.current.onConfirmDelete();
      });

      expect(returned).toBe(false);
      expect(result.current.error).toBe("user not found");
    });
  });
});
