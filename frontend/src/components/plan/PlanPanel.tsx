"use client";

import type { PlanViewModel } from "@/containers/PlanContainer";

import { PinList } from "./PinList";

type Props = {
  planVM: PlanViewModel;
};

export function PlanPanel({ planVM }: Props) {
  return (
    <div className="flex h-full flex-col">
      <div className="border-b border-zinc-200 px-4 py-3">
        <h1 className="font-semibold text-zinc-900">{planVM.title}</h1>
      </div>
      <div className="border-b border-zinc-200 px-4 py-3">
        <div className="flex items-center gap-2 rounded-lg bg-zinc-100 px-3 py-2">
          <svg
            className="h-4 w-4 flex-shrink-0 text-zinc-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth={2}
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M21 21l-4.35-4.35M17 11A6 6 0 1 1 5 11a6 6 0 0 1 12 0z"
            />
          </svg>
          <span className="text-sm text-zinc-400">スポットを検索...</span>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto p-4">
        <PinList pins={planVM.pins} />
      </div>
    </div>
  );
}
