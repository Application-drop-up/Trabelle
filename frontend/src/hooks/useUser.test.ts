import { act, renderHook } from "@testing-library/react";

import type { RegisterUserInput, User } from "@/domain/user/types";
import { useUser } from "./useUser";

const mockInput: RegisterUserInput = {
  email: "taro@example.com",
  password: "password123",
  name: "Taro",
};

const mockUser: User = {
  id: "user-1",
  email: "taro@example.com",
  name: "Taro",
  created_at: "2024-01-01T00:00:00Z",
};

beforeEach(() => {
  jest.resetAllMocks();
});

describe("useUser", () => {
  describe("registerUser", () => {
    it("returns registered user and updates state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 201,
        json: async () => mockUser,
      } as Response);

      const { result } = renderHook(() => useUser());

      let returned: User | null = null;
      await act(async () => {
        returned = await result.current.registerUser(mockInput);
      });

      expect(returned).toEqual(mockUser);
      expect(result.current.user).toEqual(mockUser);
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it("sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 409,
        statusText: "Conflict",
        json: async () => ({ error: "email already taken" }),
      } as Response);

      const { result } = renderHook(() => useUser());

      await act(async () => {
        await result.current.registerUser(mockInput);
      });

      expect(result.current.user).toBeNull();
      expect(result.current.error).toBe("email already taken");
      expect(result.current.loading).toBe(false);
    });

    it("sets loading true while fetching", async () => {
      let resolveFetch!: (value: unknown) => void;
      global.fetch = jest.fn().mockReturnValue(
        new Promise((resolve) => {
          resolveFetch = resolve;
        }),
      );

      const { result } = renderHook(() => useUser());

      act(() => {
        result.current.registerUser(mockInput);
      });

      expect(result.current.loading).toBe(true);

      await act(async () => {
        resolveFetch({
          ok: true,
          status: 201,
          json: async () => mockUser,
        });
      });

      expect(result.current.loading).toBe(false);
    });

    it("calls fetch with the correct URL and body", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 201,
        json: async () => mockUser,
      } as Response);

      const { result } = renderHook(() => useUser());

      await act(async () => {
        await result.current.registerUser(mockInput);
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/api/v1/user/register"),
        expect.objectContaining({
          method: "POST",
          body: JSON.stringify(mockInput),
        }),
      );
    });
  });
});
