"use client";

import type { PlanViewModel } from "@/containers/PlanContainer";
import type { useSpotSearchContainer, SpotViewModel } from "@/containers/SpotSearchContainer";
import { SpotSearchBar } from "@/components/spots/SpotSearchBar";

import { PinList } from "./PinList";

type Props = {
  planVM: PlanViewModel;
  spotSearch: ReturnType<typeof useSpotSearchContainer>;
  onSelectSpot: (spot: SpotViewModel) => void;
  onCreateNote: (pinId: string, content: string) => Promise<void>;
  onUpdateNote: (pinId: string, noteId: string, content: string) => Promise<void>;
  onDeleteNote: (pinId: string, noteId: string) => Promise<void>;
};

export function PlanPanel({
  planVM,
  spotSearch,
  onSelectSpot,
  onCreateNote,
  onUpdateNote,
  onDeleteNote,
}: Props) {
  return (
    <div className="flex h-full flex-col">
      <div className="border-b border-zinc-200 px-4 py-3">
        <h1 className="font-semibold text-zinc-900">{planVM.title}</h1>
      </div>
      <div className="border-b border-zinc-200 px-4 py-3">
        <SpotSearchBar
          query={spotSearch.query}
          spotVMs={spotSearch.spotVMs}
          loading={spotSearch.loading}
          error={spotSearch.error}
          onChangeQuery={spotSearch.onChangeQuery}
          onSearch={spotSearch.onSearch}
          onClear={spotSearch.onClear}
          onSelect={onSelectSpot}
        />
      </div>
      <div className="flex-1 overflow-y-auto p-4">
        <PinList
          pins={planVM.pins}
          onCreateNote={onCreateNote}
          onUpdateNote={onUpdateNote}
          onDeleteNote={onDeleteNote}
        />
      </div>
    </div>
  );
}
