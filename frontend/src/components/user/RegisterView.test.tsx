import { fireEvent, render, screen } from "@testing-library/react";

import type { User } from "@/domain/user/types";
import { RegisterView } from "./RegisterView";

const mockOnSubmit = jest.fn();
const mockOnChangeEmail = jest.fn();
const mockOnChangePassword = jest.fn();
const mockOnChangeName = jest.fn();
const mockUseRegisterContainer = jest.fn();

jest.mock("@/containers/RegisterContainer", () => ({
  useRegisterContainer: (...args: unknown[]) => mockUseRegisterContainer(...args),
}));

const mockUser: User = {
  id: "user-1",
  email: "taro@example.com",
  name: "Taro",
  created_at: "2024-01-01T00:00:00Z",
};

beforeEach(() => {
  jest.clearAllMocks();
  mockUseRegisterContainer.mockReturnValue({
    email: "",
    password: "",
    name: "",
    loading: false,
    error: null,
    onChangeEmail: mockOnChangeEmail,
    onChangePassword: mockOnChangePassword,
    onChangeName: mockOnChangeName,
    onSubmit: mockOnSubmit,
  });
});

describe("RegisterView", () => {
  it("renders the email, password, and name fields", () => {
    render(<RegisterView />);

    expect(screen.getByText("メールアドレス")).toBeInTheDocument();
    expect(screen.getByText("パスワード")).toBeInTheDocument();
    expect(screen.getByText("ユーザー名")).toBeInTheDocument();
  });

  it("calls onChangeEmail/onChangePassword/onChangeName when typing", () => {
    render(<RegisterView />);

    fireEvent.change(screen.getByLabelText("メールアドレス"), {
      target: { value: "taro@example.com" },
    });
    fireEvent.change(screen.getByLabelText("パスワード", { exact: false }), {
      target: { value: "password123" },
    });
    fireEvent.change(screen.getByLabelText("ユーザー名"), {
      target: { value: "Taro" },
    });

    expect(mockOnChangeEmail).toHaveBeenCalledWith("taro@example.com");
    expect(mockOnChangePassword).toHaveBeenCalledWith("password123");
    expect(mockOnChangeName).toHaveBeenCalledWith("Taro");
  });

  it("disables the submit button when a field is empty", () => {
    render(<RegisterView />);

    expect(screen.getByRole("button", { name: "登録する" })).toBeDisabled();
  });

  it("enables the submit button when all fields are filled", () => {
    mockUseRegisterContainer.mockReturnValue({
      email: "taro@example.com",
      password: "password123",
      name: "Taro",
      loading: false,
      error: null,
      onChangeEmail: mockOnChangeEmail,
      onChangePassword: mockOnChangePassword,
      onChangeName: mockOnChangeName,
      onSubmit: mockOnSubmit,
    });

    render(<RegisterView />);

    expect(screen.getByRole("button", { name: "登録する" })).toBeEnabled();
  });

  it("calls onSubmit when the form is submitted", () => {
    mockOnSubmit.mockResolvedValue(mockUser);
    mockUseRegisterContainer.mockReturnValue({
      email: "taro@example.com",
      password: "password123",
      name: "Taro",
      loading: false,
      error: null,
      onChangeEmail: mockOnChangeEmail,
      onChangePassword: mockOnChangePassword,
      onChangeName: mockOnChangeName,
      onSubmit: mockOnSubmit,
    });

    render(<RegisterView />);
    fireEvent.click(screen.getByRole("button", { name: "登録する" }));

    expect(mockOnSubmit).toHaveBeenCalledTimes(1);
  });

  it("shows loading state on the submit button", () => {
    mockUseRegisterContainer.mockReturnValue({
      email: "taro@example.com",
      password: "password123",
      name: "Taro",
      loading: true,
      error: null,
      onChangeEmail: mockOnChangeEmail,
      onChangePassword: mockOnChangePassword,
      onChangeName: mockOnChangeName,
      onSubmit: mockOnSubmit,
    });

    render(<RegisterView />);

    expect(screen.getByRole("button", { name: "登録中…" })).toBeDisabled();
  });

  it("shows an error message on failure", () => {
    mockUseRegisterContainer.mockReturnValue({
      email: "taro@example.com",
      password: "password123",
      name: "Taro",
      loading: false,
      error: "email already taken",
      onChangeEmail: mockOnChangeEmail,
      onChangePassword: mockOnChangePassword,
      onChangeName: mockOnChangeName,
      onSubmit: mockOnSubmit,
    });

    render(<RegisterView />);

    expect(screen.getByText("email already taken")).toBeInTheDocument();
  });
});
