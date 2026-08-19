import { act, renderHook } from "@testing-library/react";

import type {
  LoginStartInput,
  LoginVerifyInput,
  RegisterUserInput,
  UpdateUserInput,
  User,
} from "@/domain/user/types";
import { useUser } from "./useUser";

const mockInput: RegisterUserInput = {
  email: "taro@example.com",
  password: "password123",
  name: "Taro",
};

const mockUpdateInput: UpdateUserInput = {
  name: "Jiro",
  email: "jiro@example.com",
};

const mockUser: User = {
  id: "user-1",
  email: "taro@example.com",
  name: "Taro",
  created_at: "2024-01-01T00:00:00Z",
};

const mockUpdatedUser: User = {
  ...mockUser,
  name: "Jiro",
  email: "jiro@example.com",
};

const mockLoginStartInput: LoginStartInput = {
  email: "taro@example.com",
  password: "password123",
};

const mockLoginVerifyInput: LoginVerifyInput = {
  email: "taro@example.com",
  code: "123456",
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
        json: async () => ({ message: "email already taken" }),
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

  describe("updateUser", () => {
    it("returns updated user and updates state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => mockUpdatedUser,
      } as Response);

      const { result } = renderHook(() => useUser());

      let returned: User | null = null;
      await act(async () => {
        returned = await result.current.updateUser(mockUser.id, mockUpdateInput);
      });

      expect(returned).toEqual(mockUpdatedUser);
      expect(result.current.user).toEqual(mockUpdatedUser);
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it("sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 404,
        statusText: "Not Found",
        json: async () => ({ message: "user not found" }),
      } as Response);

      const { result } = renderHook(() => useUser());

      await act(async () => {
        await result.current.updateUser(mockUser.id, mockUpdateInput);
      });

      expect(result.current.user).toBeNull();
      expect(result.current.error).toBe("user not found");
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
        result.current.updateUser(mockUser.id, mockUpdateInput);
      });

      expect(result.current.loading).toBe(true);

      await act(async () => {
        resolveFetch({
          ok: true,
          status: 200,
          json: async () => mockUpdatedUser,
        });
      });

      expect(result.current.loading).toBe(false);
    });

    it("calls fetch with the correct URL and body", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => mockUpdatedUser,
      } as Response);

      const { result } = renderHook(() => useUser());

      await act(async () => {
        await result.current.updateUser(mockUser.id, mockUpdateInput);
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining(`/api/v1/user/${mockUser.id}`),
        expect.objectContaining({
          method: "PATCH",
          body: JSON.stringify(mockUpdateInput),
        }),
      );
    });
  });

  describe("deleteUser", () => {
    it("returns true and clears user state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 204,
      } as Response);

      const { result } = renderHook(() => useUser());

      let returned = false;
      await act(async () => {
        returned = await result.current.deleteUser(mockUser.id);
      });

      expect(returned).toBe(true);
      expect(result.current.user).toBeNull();
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it("returns false and sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 404,
        statusText: "Not Found",
        json: async () => ({ message: "user not found" }),
      } as Response);

      const { result } = renderHook(() => useUser());

      let returned = true;
      await act(async () => {
        returned = await result.current.deleteUser(mockUser.id);
      });

      expect(returned).toBe(false);
      expect(result.current.error).toBe("user not found");
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
        result.current.deleteUser(mockUser.id);
      });

      expect(result.current.loading).toBe(true);

      await act(async () => {
        resolveFetch({ ok: true, status: 204 });
      });

      expect(result.current.loading).toBe(false);
    });

    it("calls fetch with the correct URL and method", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 204,
      } as Response);

      const { result } = renderHook(() => useUser());

      await act(async () => {
        await result.current.deleteUser(mockUser.id);
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining(`/api/v1/user/${mockUser.id}`),
        expect.objectContaining({ method: "DELETE" }),
      );
    });
  });

  describe("loginStart", () => {
    it("returns true on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({ message: "verification code sent" }),
      } as Response);

      const { result } = renderHook(() => useUser());

      let returned = false;
      await act(async () => {
        returned = await result.current.loginStart(mockLoginStartInput);
      });

      expect(returned).toBe(true);
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it("returns false and sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 401,
        statusText: "Unauthorized",
        json: async () => ({ message: "invalid email or password" }),
      } as Response);

      const { result } = renderHook(() => useUser());

      let returned = true;
      await act(async () => {
        returned = await result.current.loginStart(mockLoginStartInput);
      });

      expect(returned).toBe(false);
      expect(result.current.error).toBe("invalid email or password");
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
        result.current.loginStart(mockLoginStartInput);
      });

      expect(result.current.loading).toBe(true);

      await act(async () => {
        resolveFetch({
          ok: true,
          status: 200,
          json: async () => ({ message: "verification code sent" }),
        });
      });

      expect(result.current.loading).toBe(false);
    });

    it("calls fetch with the correct URL and body", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({ message: "verification code sent" }),
      } as Response);

      const { result } = renderHook(() => useUser());

      await act(async () => {
        await result.current.loginStart(mockLoginStartInput);
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/api/v1/login"),
        expect.objectContaining({
          method: "POST",
          body: JSON.stringify(mockLoginStartInput),
        }),
      );
    });
  });

  describe("loginVerify", () => {
    it("returns true on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({ message: "login successful" }),
      } as Response);

      const { result } = renderHook(() => useUser());

      let returned = false;
      await act(async () => {
        returned = await result.current.loginVerify(mockLoginVerifyInput);
      });

      expect(returned).toBe(true);
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it("returns false and sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 401,
        statusText: "Unauthorized",
        json: async () => ({ message: "invalid email or code" }),
      } as Response);

      const { result } = renderHook(() => useUser());

      let returned = true;
      await act(async () => {
        returned = await result.current.loginVerify(mockLoginVerifyInput);
      });

      expect(returned).toBe(false);
      expect(result.current.error).toBe("invalid email or code");
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
        result.current.loginVerify(mockLoginVerifyInput);
      });

      expect(result.current.loading).toBe(true);

      await act(async () => {
        resolveFetch({
          ok: true,
          status: 200,
          json: async () => ({ message: "login successful" }),
        });
      });

      expect(result.current.loading).toBe(false);
    });

    it("calls fetch with the correct URL and body", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({ message: "login successful" }),
      } as Response);

      const { result } = renderHook(() => useUser());

      await act(async () => {
        await result.current.loginVerify(mockLoginVerifyInput);
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/api/v1/login/verify"),
        expect.objectContaining({
          method: "POST",
          body: JSON.stringify(mockLoginVerifyInput),
        }),
      );
    });
  });

  describe("fetchCurrentUser", () => {
    it("returns the user and updates state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => mockUser,
      } as Response);

      const { result } = renderHook(() => useUser());

      let returned: User | null = null;
      await act(async () => {
        returned = await result.current.fetchCurrentUser();
      });

      expect(returned).toEqual(mockUser);
      expect(result.current.user).toEqual(mockUser);
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it("returns null and sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 401,
        statusText: "Unauthorized",
        json: async () => ({ message: "not authenticated" }),
      } as Response);

      const { result } = renderHook(() => useUser());

      let returned: User | null = mockUser;
      await act(async () => {
        returned = await result.current.fetchCurrentUser();
      });

      expect(returned).toBeNull();
      expect(result.current.user).toBeNull();
      expect(result.current.error).toBe("not authenticated");
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
        result.current.fetchCurrentUser();
      });

      expect(result.current.loading).toBe(true);

      await act(async () => {
        resolveFetch({ ok: true, status: 200, json: async () => mockUser });
      });

      expect(result.current.loading).toBe(false);
    });

    it("calls fetch with the correct URL", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => mockUser,
      } as Response);

      const { result } = renderHook(() => useUser());

      await act(async () => {
        await result.current.fetchCurrentUser();
      });

      const fetchCall = (fetch as jest.Mock).mock.calls[0];
      expect(fetchCall[0]).toEqual(expect.stringContaining("/api/v1/user/me"));
      expect(fetchCall[1]).not.toHaveProperty("method");
      expect(fetchCall[1]).not.toHaveProperty("body");
    });
  });

  describe("logoutUser", () => {
    it("returns true and clears user state on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 204,
      } as Response);

      const { result } = renderHook(() => useUser());

      let returned = false;
      await act(async () => {
        returned = await result.current.logoutUser();
      });

      expect(returned).toBe(true);
      expect(result.current.user).toBeNull();
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it("returns false and sets error state on failure", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 500,
        statusText: "Internal Server Error",
        json: async () => ({ message: "internal server error" }),
      } as Response);

      const { result } = renderHook(() => useUser());

      let returned = true;
      await act(async () => {
        returned = await result.current.logoutUser();
      });

      expect(returned).toBe(false);
      expect(result.current.error).toBe("internal server error");
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
        result.current.logoutUser();
      });

      expect(result.current.loading).toBe(true);

      await act(async () => {
        resolveFetch({ ok: true, status: 204 });
      });

      expect(result.current.loading).toBe(false);
    });

    it("calls fetch with the correct URL", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 204,
      } as Response);

      const { result } = renderHook(() => useUser());

      await act(async () => {
        await result.current.logoutUser();
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining("/api/v1/logout"),
        expect.objectContaining({ method: "POST" }),
      );
    });
  });
});
