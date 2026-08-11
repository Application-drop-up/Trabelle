import { act, renderHook } from "@testing-library/react";

import type { User } from "@/domain/user/types";
import { useRegisterContainer } from "./RegisterContainer";

const mockUser: User = {
  id: "user-1",
  email: "taro@example.com",
  name: "Taro",
  created_at: "2024-01-01T00:00:00Z",
};

beforeEach(() => {
  jest.resetAllMocks();
});

describe("useRegisterContainer", () => {
  it("initializes with empty fields and no error", () => {
    const { result } = renderHook(() => useRegisterContainer());

    expect(result.current.email).toBe("");
    expect(result.current.password).toBe("");
    expect(result.current.name).toBe("");
    expect(result.current.error).toBeNull();
    expect(result.current.loading).toBe(false);
  });

  describe("onChangeEmail / onChangePassword / onChangeName", () => {
    it("updates each field independently", () => {
      const { result } = renderHook(() => useRegisterContainer());

      act(() => {
        result.current.onChangeEmail("taro@example.com");
        result.current.onChangePassword("password123");
        result.current.onChangeName("Taro");
      });

      expect(result.current.email).toBe("taro@example.com");
      expect(result.current.password).toBe("password123");
      expect(result.current.name).toBe("Taro");
    });
  });

  describe("onSubmit", () => {
    it("returns null when a field is blank", async () => {
      global.fetch = jest.fn();
      const { result } = renderHook(() => useRegisterContainer());

      act(() => {
        result.current.onChangeEmail("taro@example.com");
        result.current.onChangePassword("password123");
      });

      let returned: User | null = undefined as unknown as User | null;
      await act(async () => {
        returned = await result.current.onSubmit();
      });

      expect(returned).toBeNull();
      expect(global.fetch).not.toHaveBeenCalled();
    });

    it("returns null when a field is whitespace only", async () => {
      global.fetch = jest.fn();
      const { result } = renderHook(() => useRegisterContainer());

      act(() => {
        result.current.onChangeEmail("taro@example.com");
        result.current.onChangePassword("password123");
        result.current.onChangeName("   ");
      });

      let returned: User | null = undefined as unknown as User | null;
      await act(async () => {
        returned = await result.current.onSubmit();
      });

      expect(returned).toBeNull();
      expect(global.fetch).not.toHaveBeenCalled();
    });

    it("calls registerUser with trimmed fields and returns user on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 201,
        json: async () => mockUser,
      } as Response);

      const { result } = renderHook(() => useRegisterContainer());

      act(() => {
        result.current.onChangeEmail("  taro@example.com  ");
        result.current.onChangePassword("  password123  ");
        result.current.onChangeName("  Taro  ");
      });

      let returned: User | null = null;
      await act(async () => {
        returned = await result.current.onSubmit();
      });

      expect(returned).toEqual(mockUser);
      const fetchCall = (fetch as jest.Mock).mock.calls[0];
      expect(JSON.parse(fetchCall[1].body)).toEqual({
        email: "taro@example.com",
        password: "password123",
        name: "Taro",
      });
    });

    it("sets error state when registerUser fails", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 409,
        statusText: "Conflict",
        json: async () => ({ error: "email already taken" }),
      } as Response);

      const { result } = renderHook(() => useRegisterContainer());

      act(() => {
        result.current.onChangeEmail("taro@example.com");
        result.current.onChangePassword("password123");
        result.current.onChangeName("Taro");
      });

      await act(async () => {
        await result.current.onSubmit();
      });

      expect(result.current.error).toBe("email already taken");
    });
  });
});
