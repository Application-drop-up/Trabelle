import { fireEvent, render, screen } from "@testing-library/react";

import { CountryGuidePanel } from "./CountryGuidePanel";

const mockOnSelectCode = jest.fn();
const mockUseCountryGuideContainer = jest.fn();

jest.mock("@/containers/CountryGuideContainer", () => ({
  useCountryGuideContainer: () => mockUseCountryGuideContainer(),
}));

const guideVMs = [
  { id: "guide-1", countryCode: "TH", countryName: "タイ", items: [] },
  { id: "guide-2", countryCode: "TW", countryName: "台湾", items: [] },
];

const selectedGuideVM = {
  id: "guide-1",
  countryCode: "TH",
  countryName: "タイ",
  items: [
    {
      id: "item-1",
      category: "entry_card" as const,
      title: "TDAC",
      description: "出発72時間前からオンラインで申請可能",
      url: "https://tdac.immigration.go.th",
      isMandatory: true,
    },
    {
      id: "item-2",
      category: "packing_tip" as const,
      title: "電源プラグはBFタイプ",
      description: "",
      url: "",
      isMandatory: false,
    },
  ],
};

beforeEach(() => {
  jest.clearAllMocks();
  mockUseCountryGuideContainer.mockReturnValue({
    guideVMs,
    selectedCode: "",
    selectedGuideVM: null,
    loading: false,
    error: null,
    onSelectCode: mockOnSelectCode,
  });
});

describe("CountryGuidePanel", () => {
  it("renders the country options from guideVMs", () => {
    render(<CountryGuidePanel />);

    fireEvent.mouseDown(screen.getByRole("combobox", { name: "渡航先" }));

    expect(screen.getByRole("option", { name: "タイ" })).toBeInTheDocument();
    expect(screen.getByRole("option", { name: "台湾" })).toBeInTheDocument();
  });

  it("calls onSelectCode when a country is selected", () => {
    render(<CountryGuidePanel />);

    fireEvent.mouseDown(screen.getByRole("combobox", { name: "渡航先" }));
    fireEvent.click(screen.getByRole("option", { name: "タイ" }));

    expect(mockOnSelectCode).toHaveBeenCalledWith("TH");
  });

  it("shows a loading indicator while loading", () => {
    mockUseCountryGuideContainer.mockReturnValue({
      guideVMs: [],
      selectedCode: "",
      selectedGuideVM: null,
      loading: true,
      error: null,
      onSelectCode: mockOnSelectCode,
    });

    render(<CountryGuidePanel />);

    expect(screen.getByLabelText("読み込み中")).toBeInTheDocument();
  });

  it("shows an error message on failure", () => {
    mockUseCountryGuideContainer.mockReturnValue({
      guideVMs,
      selectedCode: "ZZ",
      selectedGuideVM: null,
      loading: false,
      error: "country guide not found",
      onSelectCode: mockOnSelectCode,
    });

    render(<CountryGuidePanel />);

    expect(screen.getByText("country guide not found")).toBeInTheDocument();
  });

  it("renders the selected guide's items, including a mandatory badge", () => {
    mockUseCountryGuideContainer.mockReturnValue({
      guideVMs,
      selectedCode: "TH",
      selectedGuideVM,
      loading: false,
      error: null,
      onSelectCode: mockOnSelectCode,
    });

    render(<CountryGuidePanel />);

    expect(screen.getByText("TDAC")).toBeInTheDocument();
    expect(screen.getByText("電源プラグはBFタイプ")).toBeInTheDocument();
    expect(screen.getByText("必須")).toBeInTheDocument();
    expect(screen.getByText("詳細を見る")).toBeInTheDocument();
  });
});
