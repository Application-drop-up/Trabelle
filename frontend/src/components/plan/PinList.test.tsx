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
    render(<PinList pins={[]} onCreateNote={jest.fn()} onUpdateNote={jest.fn()} />);
    expect(screen.getByText("ピンがまだありません")).toBeInTheDocument();
  });

  it("shows existing notes and an add-note form when expanded", () => {
    render(<PinList pins={[mockPin]} onCreateNote={jest.fn()} onUpdateNote={jest.fn()} />);

    fireEvent.click(screen.getByText("Tokyo Tower"));

    expect(screen.getByText("Visit the observation deck")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("メモを追加...")).toBeInTheDocument();
  });

  it("disables the add button until text is entered", () => {
    render(<PinList pins={[mockPin]} onCreateNote={jest.fn()} onUpdateNote={jest.fn()} />);

    fireEvent.click(screen.getByText("Tokyo Tower"));

    expect(screen.getByRole("button", { name: "追加" })).toBeDisabled();

    fireEvent.change(screen.getByPlaceholderText("メモを追加..."), {
      target: { value: "Bring camera" },
    });

    expect(screen.getByRole("button", { name: "追加" })).not.toBeDisabled();
  });

  it("calls onCreateNote with the pin id and content, then clears the input", async () => {
    const onCreateNote = jest.fn().mockResolvedValue(undefined);
    render(<PinList pins={[mockPin]} onCreateNote={onCreateNote} onUpdateNote={jest.fn()} />);

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

  it("shows an edit button for each note and enters edit mode when clicked", () => {
    render(<PinList pins={[mockPin]} onCreateNote={jest.fn()} onUpdateNote={jest.fn()} />);

    fireEvent.click(screen.getByText("Tokyo Tower"));
    fireEvent.click(screen.getByRole("button", { name: "編集" }));

    expect(screen.getByDisplayValue("Visit the observation deck")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "保存" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "キャンセル" })).toBeInTheDocument();
  });

  it("calls onUpdateNote with the pin id, note id, and new content", async () => {
    const onUpdateNote = jest.fn().mockResolvedValue(undefined);
    render(<PinList pins={[mockPin]} onCreateNote={jest.fn()} onUpdateNote={onUpdateNote} />);

    fireEvent.click(screen.getByText("Tokyo Tower"));
    fireEvent.click(screen.getByRole("button", { name: "編集" }));

    fireEvent.change(screen.getByDisplayValue("Visit the observation deck"), {
      target: { value: "Great view of Mt. Fuji" },
    });
    fireEvent.click(screen.getByRole("button", { name: "保存" }));

    await waitFor(() => {
      expect(onUpdateNote).toHaveBeenCalledWith("pin-1", "note-1", "Great view of Mt. Fuji");
    });

    expect(screen.queryByRole("button", { name: "保存" })).not.toBeInTheDocument();
  });

  it("reverts to the original content on cancel without saving", () => {
    const onUpdateNote = jest.fn();
    render(<PinList pins={[mockPin]} onCreateNote={jest.fn()} onUpdateNote={onUpdateNote} />);

    fireEvent.click(screen.getByText("Tokyo Tower"));
    fireEvent.click(screen.getByRole("button", { name: "編集" }));
    fireEvent.change(screen.getByDisplayValue("Visit the observation deck"), {
      target: { value: "discard me" },
    });
    fireEvent.click(screen.getByRole("button", { name: "キャンセル" }));

    expect(onUpdateNote).not.toHaveBeenCalled();
    expect(screen.getByText("Visit the observation deck")).toBeInTheDocument();
  });
});
