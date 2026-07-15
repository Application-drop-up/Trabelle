import { fireEvent, render, screen } from "@testing-library/react";

import type { SpotViewModel } from "@/mappers/planMapper";
import { SpotSearchBar } from "./SpotSearchBar";

const mockSpots: SpotViewModel[] = [
  {
    placeId: "spot-1",
    name: "Tokyo Tower",
    address: "4 Chome-2-8 Shibakoen, Minato City, Tokyo",
    latitude: 35.6586,
    longitude: 139.7454,
  },
];

beforeEach(() => {
  jest.useFakeTimers();
});

afterEach(() => {
  jest.useRealTimers();
});

describe("SpotSearchBar", () => {
  it("calls onSearch after the debounce delay once the user stops typing", () => {
    const onSearch = jest.fn().mockResolvedValue(undefined);
    render(
      <SpotSearchBar
        query="Tokyo"
        spotVMs={[]}
        loading={false}
        error={null}
        onChangeQuery={jest.fn()}
        onSearch={onSearch}
        onClear={jest.fn()}
        onSelect={jest.fn()}
      />,
    );

    expect(onSearch).not.toHaveBeenCalled();

    jest.advanceTimersByTime(300);

    expect(onSearch).toHaveBeenCalledTimes(1);
  });

  it("calls onClear instead of onSearch when the query is empty", () => {
    const onSearch = jest.fn();
    const onClear = jest.fn();
    render(
      <SpotSearchBar
        query=""
        spotVMs={[]}
        loading={false}
        error={null}
        onChangeQuery={jest.fn()}
        onSearch={onSearch}
        onClear={onClear}
        onSelect={jest.fn()}
      />,
    );

    jest.advanceTimersByTime(300);

    expect(onClear).toHaveBeenCalled();
    expect(onSearch).not.toHaveBeenCalled();
  });

  it("renders results and calls onSelect when one is clicked", () => {
    const onSelect = jest.fn();
    render(
      <SpotSearchBar
        query="Tokyo"
        spotVMs={mockSpots}
        loading={false}
        error={null}
        onChangeQuery={jest.fn()}
        onSearch={jest.fn()}
        onClear={jest.fn()}
        onSelect={onSelect}
      />,
    );

    expect(screen.getByText("Tokyo Tower")).toBeInTheDocument();
    fireEvent.click(screen.getByText("Tokyo Tower"));

    expect(onSelect).toHaveBeenCalledWith(mockSpots[0]);
  });

  it("shows the error message when present", () => {
    render(
      <SpotSearchBar
        query="Tokyo"
        spotVMs={[]}
        loading={false}
        error="検索に失敗しました"
        onChangeQuery={jest.fn()}
        onSearch={jest.fn()}
        onClear={jest.fn()}
        onSelect={jest.fn()}
      />,
    );

    expect(screen.getByText("検索に失敗しました")).toBeInTheDocument();
  });
});
