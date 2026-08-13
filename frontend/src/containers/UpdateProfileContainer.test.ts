import { act, renderHook } from "@testing-library/react";

import type { User } from "@/domain/user/types";
import { useUpdateProfileContainer } from "./UpdateProfileContainer";

const mockUser: User = {
  id: "user-1",
  email: "taro@example.com",
  name: "Taro",
  created_at: "2024-01-01T00:00:00Z",
};

const initial = { name: "Taro", email: "taro@example.com" };

beforeEach(() => {
  jest.resetAllMocks();
});

describe("useUpdateProfileContainer", () => {
  it("initializes with the given name and email and no error", () => {
    const { result } = renderHook(() => useUpdateProfileContainer(mockUser.id, initial));

    expect(result.current.name).toBe("Taro");
    expect(result.current.email).toBe("taro@example.com");
    expect(result.current.error).toBeNull();
    expect(result.current.loading).toBe(false);
  });

  describe("onChangeName / onChangeEmail", () => {
    it("updates each field independently", () => {
      const { result } = renderHook(() => useUpdateProfileContainer(mockUser.id, initial));

      act(() => {
        result.current.onChangeName("Jiro");
        result.current.onChangeEmail("jiro@example.com");
      });

      expect(result.current.name).toBe("Jiro");
      expect(result.current.email).toBe("jiro@example.com");
    });
  });

  describe("onSubmitUpdate", () => {
    it("returns null when a field is blank", async () => {
      global.fetch = jest.fn();
      const { result } = renderHook(() => useUpdateProfileContainer(mockUser.id, initial));

      act(() => {
        result.current.onChangeName("");
      });

      let returned: User | null = undefined as unknown as User | null;
      await act(async () => {
        returned = await result.current.onSubmitUpdate();
      });

      expect(returned).toBeNull();
      expect(global.fetch).not.toHaveBeenCalled();
    });

    it("returns null when a field is whitespace only", async () => {
      global.fetch = jest.fn();
      const { result } = renderHook(() => useUpdateProfileContainer(mockUser.id, initial));

      act(() => {
        result.current.onChangeEmail("   ");
      });

      let returned: User | null = undefined as unknown as User | null;
      await act(async () => {
        returned = await result.current.onSubmitUpdate();
      });

      expect(returned).toBeNull();
      expect(global.fetch).not.toHaveBeenCalled();
    });

    it("calls updateUser with trimmed fields and returns user on success", async () => {
      const updatedUser: User = { ...mockUser, name: "Jiro", email: "jiro@example.com" };
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => updatedUser,
      } as Response);

      const { result } = renderHook(() => useUpdateProfileContainer(mockUser.id, initial));

      act(() => {
        result.current.onChangeName("  Jiro  ");
        result.current.onChangeEmail("  jiro@example.com  ");
      });

      let returned: User | null = null;
      await act(async () => {
        returned = await result.current.onSubmitUpdate();
      });

      expect(returned).toEqual(updatedUser);
      const fetchCall = (fetch as jest.Mock).mock.calls[0];
      expect(fetchCall[0]).toContain(`/api/v1/user/${mockUser.id}`);
      expect(JSON.parse(fetchCall[1].body)).toEqual({ name: "Jiro", email: "jiro@example.com" });
    });

    it("sets error state when updateUser fails", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 409,
        statusText: "Conflict",
        json: async () => ({ error: "email already taken" }),
      } as Response);

      const { result } = renderHook(() => useUpdateProfileContainer(mockUser.id, initial));

      await act(async () => {
        await result.current.onSubmitUpdate();
      });

      expect(result.current.error).toBe("email already taken");
    });
  });
});
