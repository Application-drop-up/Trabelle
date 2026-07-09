import { render, screen } from "@testing-library/react";

import type { PinViewModel } from "@/mappers/planMapper";
import { PinDetailPanel } from "./PinDetailPanel";

const mockPinVM: PinViewModel = {
  id: "pin-1",
  planId: "plan-1",
  name: "Tokyo Tower",
  latitude: 35.6586,
  longitude: 139.7454,
  category: "sightseeing",
  colour: "#FF5733",
  notes: [],
  createdAt: new Date("2024-01-01"),
  updatedAt: new Date("2024-01-01"),
};

describe("PinDetailPanel", () => {
  it("renders pin name, category, and colour indicator", () => {
    render(<PinDetailPanel pinVM={mockPinVM} />);

    expect(screen.getByText("Tokyo Tower")).toBeInTheDocument();
    expect(screen.getByText("sightseeing")).toBeInTheDocument();
    expect(screen.getByText("Category")).toBeInTheDocument();
  });

  it("shows 'No notes yet' when notes array is empty", () => {
    render(<PinDetailPanel pinVM={mockPinVM} />);

    expect(screen.getByText("No notes yet")).toBeInTheDocument();
  });

  it("renders note contents when notes exist", () => {
    const pinWithNotes: PinViewModel = {
      ...mockPinVM,
      notes: [
        {
          id: "note-1",
          pinId: "pin-1",
          content: "Great view from the top",
          createdAt: new Date("2024-01-01"),
          updatedAt: new Date("2024-01-01"),
        },
        {
          id: "note-2",
          pinId: "pin-1",
          content: "Recommended at sunset",
          createdAt: new Date("2024-01-01"),
          updatedAt: new Date("2024-01-01"),
        },
      ],
    };

    render(<PinDetailPanel pinVM={pinWithNotes} />);

    expect(screen.getByText("Great view from the top")).toBeInTheDocument();
    expect(screen.getByText("Recommended at sunset")).toBeInTheDocument();
    expect(screen.queryByText("No notes yet")).not.toBeInTheDocument();
  });
});
