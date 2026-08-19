import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { LoginView } from "./LoginView";

const mockOnSubmitCredentials = jest.fn();
const mockOnSubmitCode = jest.fn();
const mockOnChangeEmail = jest.fn();
const mockOnChangePassword = jest.fn();
const mockOnChangeCode = jest.fn();
const mockUseLoginContainer = jest.fn();
const mockPush = jest.fn();

jest.mock("@/containers/LoginContainer", () => ({
  useLoginContainer: (...args: unknown[]) => mockUseLoginContainer(...args),
}));

jest.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
}));

function credentialsState(overrides = {}) {
  return {
    step: "credentials",
    email: "",
    password: "",
    code: "",
    loading: false,
    error: null,
    onChangeEmail: mockOnChangeEmail,
    onChangePassword: mockOnChangePassword,
    onChangeCode: mockOnChangeCode,
    onSubmitCredentials: mockOnSubmitCredentials,
    onSubmitCode: mockOnSubmitCode,
    ...overrides,
  };
}

beforeEach(() => {
  jest.clearAllMocks();
  mockUseLoginContainer.mockReturnValue(credentialsState());
});

describe("LoginView", () => {
  describe("credentials step", () => {
    it("renders the email and password fields", () => {
      render(<LoginView />);

      expect(screen.getByLabelText("メールアドレス")).toBeInTheDocument();
      expect(screen.getByLabelText("パスワード")).toBeInTheDocument();
    });

    it("calls onChangeEmail/onChangePassword when typing", () => {
      render(<LoginView />);

      fireEvent.change(screen.getByLabelText("メールアドレス"), {
        target: { value: "taro@example.com" },
      });
      fireEvent.change(screen.getByLabelText("パスワード"), {
        target: { value: "password123" },
      });

      expect(mockOnChangeEmail).toHaveBeenCalledWith("taro@example.com");
      expect(mockOnChangePassword).toHaveBeenCalledWith("password123");
    });

    it("disables the submit button when a field is empty", () => {
      render(<LoginView />);

      expect(screen.getByRole("button", { name: "コードを送信" })).toBeDisabled();
    });

    it("enables the submit button when both fields are filled", () => {
      mockUseLoginContainer.mockReturnValue(
        credentialsState({ email: "taro@example.com", password: "password123" }),
      );

      render(<LoginView />);

      expect(screen.getByRole("button", { name: "コードを送信" })).toBeEnabled();
    });

    it("calls onSubmitCredentials when the form is submitted", () => {
      mockUseLoginContainer.mockReturnValue(
        credentialsState({ email: "taro@example.com", password: "password123" }),
      );

      render(<LoginView />);
      fireEvent.click(screen.getByRole("button", { name: "コードを送信" }));

      expect(mockOnSubmitCredentials).toHaveBeenCalledTimes(1);
    });

    it("shows loading state on the submit button", () => {
      mockUseLoginContainer.mockReturnValue(
        credentialsState({ email: "taro@example.com", password: "password123", loading: true }),
      );

      render(<LoginView />);

      expect(screen.getByRole("button", { name: "送信中…" })).toBeDisabled();
    });

    it("shows an error message on failure", () => {
      mockUseLoginContainer.mockReturnValue(
        credentialsState({ error: "invalid email or password" }),
      );

      render(<LoginView />);

      expect(screen.getByText("invalid email or password")).toBeInTheDocument();
    });
  });

  describe("otp step", () => {
    it("renders the code field instead of email/password", () => {
      mockUseLoginContainer.mockReturnValue(
        credentialsState({ step: "otp", email: "taro@example.com" }),
      );

      render(<LoginView />);

      expect(screen.getByLabelText("認証コード")).toBeInTheDocument();
      expect(screen.queryByLabelText("メールアドレス")).not.toBeInTheDocument();
      expect(screen.queryByLabelText("パスワード")).not.toBeInTheDocument();
      expect(screen.getByText(/taro@example.com/)).toBeInTheDocument();
    });

    it("calls onChangeCode when typing", () => {
      mockUseLoginContainer.mockReturnValue(
        credentialsState({ step: "otp", email: "taro@example.com" }),
      );

      render(<LoginView />);
      fireEvent.change(screen.getByLabelText("認証コード"), { target: { value: "123456" } });

      expect(mockOnChangeCode).toHaveBeenCalledWith("123456");
    });

    it("disables the submit button when the code is empty", () => {
      mockUseLoginContainer.mockReturnValue(
        credentialsState({ step: "otp", email: "taro@example.com" }),
      );

      render(<LoginView />);

      expect(screen.getByRole("button", { name: "ログイン" })).toBeDisabled();
    });

    it("calls onSubmitCode when the form is submitted", () => {
      mockUseLoginContainer.mockReturnValue(
        credentialsState({ step: "otp", email: "taro@example.com", code: "123456" }),
      );

      render(<LoginView />);
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      expect(mockOnSubmitCode).toHaveBeenCalledTimes(1);
    });

    it("redirects to /profile when onSubmitCode succeeds", async () => {
      mockOnSubmitCode.mockResolvedValue(true);
      mockUseLoginContainer.mockReturnValue(
        credentialsState({ step: "otp", email: "taro@example.com", code: "123456" }),
      );

      render(<LoginView />);
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith("/profile");
      });
    });

    it("does not redirect when onSubmitCode fails", async () => {
      mockOnSubmitCode.mockResolvedValue(false);
      mockUseLoginContainer.mockReturnValue(
        credentialsState({ step: "otp", email: "taro@example.com", code: "123456" }),
      );

      render(<LoginView />);
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(mockOnSubmitCode).toHaveBeenCalledTimes(1);
      });
      expect(mockPush).not.toHaveBeenCalled();
    });

    it("shows loading state on the submit button", () => {
      mockUseLoginContainer.mockReturnValue(
        credentialsState({ step: "otp", email: "taro@example.com", code: "123456", loading: true }),
      );

      render(<LoginView />);

      expect(screen.getByRole("button", { name: "ログイン中…" })).toBeDisabled();
    });

    it("shows an error message on failure", () => {
      mockUseLoginContainer.mockReturnValue(
        credentialsState({
          step: "otp",
          email: "taro@example.com",
          error: "invalid email or code",
        }),
      );

      render(<LoginView />);

      expect(screen.getByText("invalid email or code")).toBeInTheDocument();
    });
  });
});
