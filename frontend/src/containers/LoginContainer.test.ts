import { act, renderHook } from "@testing-library/react";

import { useLoginContainer } from "./LoginContainer";

beforeEach(() => {
  jest.resetAllMocks();
});

describe("useLoginContainer", () => {
  it("initializes on the credentials step with empty fields and no error", () => {
    const { result } = renderHook(() => useLoginContainer());

    expect(result.current.step).toBe("credentials");
    expect(result.current.email).toBe("");
    expect(result.current.password).toBe("");
    expect(result.current.code).toBe("");
    expect(result.current.error).toBeNull();
    expect(result.current.loading).toBe(false);
  });

  describe("onChangeEmail / onChangePassword / onChangeCode", () => {
    it("updates each field independently", () => {
      const { result } = renderHook(() => useLoginContainer());

      act(() => {
        result.current.onChangeEmail("taro@example.com");
        result.current.onChangePassword("password123");
        result.current.onChangeCode("123456");
      });

      expect(result.current.email).toBe("taro@example.com");
      expect(result.current.password).toBe("password123");
      expect(result.current.code).toBe("123456");
    });
  });

  describe("onSubmitCredentials", () => {
    it("returns false and stays on the credentials step when a field is blank", async () => {
      global.fetch = jest.fn();
      const { result } = renderHook(() => useLoginContainer());

      act(() => {
        result.current.onChangeEmail("taro@example.com");
      });

      let returned = true;
      await act(async () => {
        returned = await result.current.onSubmitCredentials();
      });

      expect(returned).toBe(false);
      expect(result.current.step).toBe("credentials");
      expect(global.fetch).not.toHaveBeenCalled();
    });

    it("returns false when a field is whitespace only", async () => {
      global.fetch = jest.fn();
      const { result } = renderHook(() => useLoginContainer());

      act(() => {
        result.current.onChangeEmail("taro@example.com");
        result.current.onChangePassword("   ");
      });

      let returned = true;
      await act(async () => {
        returned = await result.current.onSubmitCredentials();
      });

      expect(returned).toBe(false);
      expect(global.fetch).not.toHaveBeenCalled();
    });

    it("calls loginStart with trimmed fields and advances to the otp step on success", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({ message: "verification code sent" }),
      } as Response);

      const { result } = renderHook(() => useLoginContainer());

      act(() => {
        result.current.onChangeEmail("  taro@example.com  ");
        result.current.onChangePassword("  password123  ");
      });

      let returned = false;
      await act(async () => {
        returned = await result.current.onSubmitCredentials();
      });

      expect(returned).toBe(true);
      expect(result.current.step).toBe("otp");
      const fetchCall = (fetch as jest.Mock).mock.calls[0];
      expect(JSON.parse(fetchCall[1].body)).toEqual({
        email: "taro@example.com",
        password: "password123",
      });
    });

    it("sets error state and stays on the credentials step when loginStart fails", async () => {
      global.fetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 401,
        statusText: "Unauthorized",
        json: async () => ({ message: "invalid email or password" }),
      } as Response);

      const { result } = renderHook(() => useLoginContainer());

      act(() => {
        result.current.onChangeEmail("taro@example.com");
        result.current.onChangePassword("wrong-password");
      });

      await act(async () => {
        await result.current.onSubmitCredentials();
      });

      expect(result.current.error).toBe("invalid email or password");
      expect(result.current.step).toBe("credentials");
    });
  });

  describe("onSubmitCode", () => {
    it("returns false when the code is blank", async () => {
      global.fetch = jest.fn();
      const { result } = renderHook(() => useLoginContainer());

      let returned = true;
      await act(async () => {
        returned = await result.current.onSubmitCode();
      });

      expect(returned).toBe(false);
      expect(global.fetch).not.toHaveBeenCalled();
    });

    it("calls loginVerify with the email and trimmed code and returns true on success", async () => {
      global.fetch = jest
        .fn()
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          json: async () => ({ message: "verification code sent" }),
        } as Response)
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          json: async () => ({ message: "login successful" }),
        } as Response);

      const { result } = renderHook(() => useLoginContainer());

      act(() => {
        result.current.onChangeEmail("taro@example.com");
        result.current.onChangePassword("password123");
      });

      await act(async () => {
        await result.current.onSubmitCredentials();
      });

      act(() => {
        result.current.onChangeCode("  123456  ");
      });

      let returned = false;
      await act(async () => {
        returned = await result.current.onSubmitCode();
      });

      expect(returned).toBe(true);
      const fetchCall = (fetch as jest.Mock).mock.calls[1];
      expect(JSON.parse(fetchCall[1].body)).toEqual({
        email: "taro@example.com",
        code: "123456",
      });
    });

    it("sets error state when loginVerify fails", async () => {
      global.fetch = jest
        .fn()
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          json: async () => ({ message: "verification code sent" }),
        } as Response)
        .mockResolvedValueOnce({
          ok: false,
          status: 401,
          statusText: "Unauthorized",
          json: async () => ({ message: "invalid email or code" }),
        } as Response);

      const { result } = renderHook(() => useLoginContainer());

      act(() => {
        result.current.onChangeEmail("taro@example.com");
        result.current.onChangePassword("password123");
      });

      await act(async () => {
        await result.current.onSubmitCredentials();
      });

      act(() => {
        result.current.onChangeCode("000000");
      });

      await act(async () => {
        await result.current.onSubmitCode();
      });

      expect(result.current.error).toBe("invalid email or code");
    });
  });
});
