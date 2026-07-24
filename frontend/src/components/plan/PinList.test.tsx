import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import type { PinViewModel } from "@/containers/PlanContainer";
import { PinList } from "./PinList";

const mockPin: PinViewModel = {
  id: "pin-1",
  planId: "plan-1",
  name: "Tokyo Tower",
  latitude: 35.6586,
  longitude: 139.7454,
  category: "sightseeing",
  colour: "#FF5733",
  notes: [
    {
      id: "note-1",
      pinId: "pin-1",
      content: "Visit the observation deck",
      createdAt: new Date("2024-01-01"),
      updatedAt: new Date("2024-01-01"),
    },
  ],
  createdAt: new Date("2024-01-01"),
  updatedAt: new Date("2024-01-01"),
};

describe("PinList", () => {
  it("shows a placeholder when there are no pins", () => {
    render(<PinList pins={[]} onCreateNote={jest.fn()} />);
    expect(screen.getByText("ピンがまだありません")).toBeInTheDocument();
  });

  it("shows existing notes and an add-note form when expanded", () => {
    render(<PinList pins={[mockPin]} onCreateNote={jest.fn()} />);

    fireEvent.click(screen.getByText("Tokyo Tower"));

    expect(screen.getByText("Visit the observation deck")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("メモを追加...")).toBeInTheDocument();
  });

  it("disables the add button until text is entered", () => {
    render(<PinList pins={[mockPin]} onCreateNote={jest.fn()} />);

    fireEvent.click(screen.getByText("Tokyo Tower"));

    expect(screen.getByRole("button", { name: "追加" })).toBeDisabled();

    fireEvent.change(screen.getByPlaceholderText("メモを追加..."), {
      target: { value: "Bring camera" },
    });

    expect(screen.getByRole("button", { name: "追加" })).not.toBeDisabled();
  });

  it("calls onCreateNote with the pin id and content, then clears the input", async () => {
    const onCreateNote = jest.fn().mockResolvedValue(undefined);
    render(<PinList pins={[mockPin]} onCreateNote={onCreateNote} />);

    fireEvent.click(screen.getByText("Tokyo Tower"));
    fireEvent.change(screen.getByPlaceholderText("メモを追加..."), {
      target: { value: "Bring camera" },
    });
    fireEvent.click(screen.getByRole("button", { name: "追加" }));

    await waitFor(() => {
      expect(onCreateNote).toHaveBeenCalledWith("pin-1", "Bring camera");
    });

    expect(screen.getByPlaceholderText("メモを追加...")).toHaveValue("");
  });
});
