import { fireEvent, render, screen } from "@testing-library/react";

import type { User } from "@/domain/user/types";
import { UpdateProfileView } from "./UpdateProfileView";

const mockOnSubmitUpdate = jest.fn();
const mockOnChangeName = jest.fn();
const mockOnChangeEmail = jest.fn();
const mockUseUpdateProfileContainer = jest.fn();

jest.mock("@/containers/UpdateProfileContainer", () => ({
  useUpdateProfileContainer: (...args: unknown[]) => mockUseUpdateProfileContainer(...args),
}));

const mockUser: User = {
  id: "user-1",
  email: "taro@example.com",
  name: "Taro",
  created_at: "2024-01-01T00:00:00Z",
};

const initial = { name: "Taro", email: "taro@example.com" };

beforeEach(() => {
  jest.clearAllMocks();
  mockUseUpdateProfileContainer.mockReturnValue({
    name: "Taro",
    email: "taro@example.com",
    loading: false,
    error: null,
    onChangeName: mockOnChangeName,
    onChangeEmail: mockOnChangeEmail,
    onSubmitUpdate: mockOnSubmitUpdate,
  });
});

describe("UpdateProfileView", () => {
  it("passes userId and initial values to the container", () => {
    render(<UpdateProfileView userId={mockUser.id} initial={initial} />);

    expect(mockUseUpdateProfileContainer).toHaveBeenCalledWith(mockUser.id, initial);
  });

  it("renders the name and email fields pre-filled", () => {
    render(<UpdateProfileView userId={mockUser.id} initial={initial} />);

    expect(screen.getByDisplayValue("Taro")).toBeInTheDocument();
    expect(screen.getByDisplayValue("taro@example.com")).toBeInTheDocument();
  });

  it("calls onChangeName/onChangeEmail when typing", () => {
    render(<UpdateProfileView userId={mockUser.id} initial={initial} />);

    fireEvent.change(screen.getByLabelText("ユーザー名"), { target: { value: "Jiro" } });
    fireEvent.change(screen.getByLabelText("メールアドレス"), {
      target: { value: "jiro@example.com" },
    });

    expect(mockOnChangeName).toHaveBeenCalledWith("Jiro");
    expect(mockOnChangeEmail).toHaveBeenCalledWith("jiro@example.com");
  });

  it("disables the submit button when a field is empty", () => {
    mockUseUpdateProfileContainer.mockReturnValue({
      name: "",
      email: "taro@example.com",
      loading: false,
      error: null,
      onChangeName: mockOnChangeName,
      onChangeEmail: mockOnChangeEmail,
      onSubmitUpdate: mockOnSubmitUpdate,
    });

    render(<UpdateProfileView userId={mockUser.id} initial={initial} />);

    expect(screen.getByRole("button", { name: "更新する" })).toBeDisabled();
  });

  it("calls onSubmitUpdate when the form is submitted", () => {
    mockOnSubmitUpdate.mockResolvedValue(mockUser);

    render(<UpdateProfileView userId={mockUser.id} initial={initial} />);
    fireEvent.click(screen.getByRole("button", { name: "更新する" }));

    expect(mockOnSubmitUpdate).toHaveBeenCalledTimes(1);
  });

  it("shows loading state on the submit button", () => {
    mockUseUpdateProfileContainer.mockReturnValue({
      name: "Taro",
      email: "taro@example.com",
      loading: true,
      error: null,
      onChangeName: mockOnChangeName,
      onChangeEmail: mockOnChangeEmail,
      onSubmitUpdate: mockOnSubmitUpdate,
    });

    render(<UpdateProfileView userId={mockUser.id} initial={initial} />);

    expect(screen.getByRole("button", { name: "更新中…" })).toBeDisabled();
  });

  it("shows an error message on failure", () => {
    mockUseUpdateProfileContainer.mockReturnValue({
      name: "Taro",
      email: "taro@example.com",
      loading: false,
      error: "email already taken",
      onChangeName: mockOnChangeName,
      onChangeEmail: mockOnChangeEmail,
      onSubmitUpdate: mockOnSubmitUpdate,
    });

    render(<UpdateProfileView userId={mockUser.id} initial={initial} />);

    expect(screen.getByText("email already taken")).toBeInTheDocument();
  });
});
