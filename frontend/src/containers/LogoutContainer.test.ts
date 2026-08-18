import { act, renderHook } from "@testing-library/react";

import { useLogoutContainer } from "./LogoutContainer";

beforeEach(() => {
  jest.resetAllMocks();
});

describe("useLogoutContainer", () => {
  it("initializes with no error and not loading", () => {
    const { result } = renderHook(() => useLogoutContainer());

    expect(result.current.error).toBeNull();
    expect(result.current.loading).toBe(false);
  });

  describe("onLogout", () => {
    it("calls the logout endpoint and returns true on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 204,
      } as Response);

      const { result } = renderHook(() => useLogoutContainer());

      let returned = false;
      await act(async () => {
        returned = await result.current.onLogout();
      });

      expect(returned).toBe(true);
      const fetchCall = (fetch as jest.Mock).mock.calls[0];
      expect(fetchCall[0]).toContain("/api/v1/logout");
      expect(fetchCall[1]).toEqual(expect.objectContaining({ method: "POST" }));
    });

    it("returns false and sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 500,
        statusText: "Internal Server Error",
        json: async () => ({ message: "internal server error" }),
      } as Response);

      const { result } = renderHook(() => useLogoutContainer());

      let returned = true;
      await act(async () => {
        returned = await result.current.onLogout();
      });

      expect(returned).toBe(false);
      expect(result.current.error).toBe("internal server error");
    });
  });
});
