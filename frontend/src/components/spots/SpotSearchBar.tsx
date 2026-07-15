"use client";

import { useEffect, useRef } from "react";

import type { SpotViewModel } from "@/mappers/planMapper";

type Props = {
  query: string;
  spotVMs: SpotViewModel[];
  loading: boolean;
  error: string | null;
  onChangeQuery: (value: string) => void;
  onSearch: () => Promise<void>;
  onClear: () => void;
  onSelect: (spot: SpotViewModel) => void;
};

const DEBOUNCE_MS = 300;

export function SpotSearchBar({
  query,
  spotVMs,
  loading,
  error,
  onChangeQuery,
  onSearch,
  onClear,
  onSelect,
}: Props) {
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    if (timerRef.current) clearTimeout(timerRef.current);

    if (!query.trim()) {
      onClear();
      return;
    }

    timerRef.current = setTimeout(() => {
      onSearch();
    }, DEBOUNCE_MS);

    return () => {
      if (timerRef.current) clearTimeout(timerRef.current);
    };
    // onSearch/onClear identity isn't relevant to when the debounce should fire
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query]);

  return (
    <div className="relative">
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
        <input
          type="text"
          value={query}
          onChange={(e) => onChangeQuery(e.target.value)}
          placeholder="スポットを検索..."
          className="flex-1 bg-transparent text-sm text-zinc-900 outline-none placeholder:text-zinc-400"
        />
      </div>

      {loading && <p className="mt-1 px-1 text-xs text-zinc-400">検索中...</p>}
      {error && <p className="mt-1 px-1 text-xs text-red-600">{error}</p>}

      {spotVMs.length > 0 && (
        <ul className="absolute z-10 mt-1 w-full rounded-lg border border-zinc-200 bg-white shadow-lg">
          {spotVMs.map((spot) => (
            <li key={spot.placeId}>
              <button
                type="button"
                onClick={() => onSelect(spot)}
                className="flex w-full flex-col items-start px-3 py-2 text-left hover:bg-zinc-50"
              >
                <span className="text-sm font-medium text-zinc-900">{spot.name}</span>
                <span className="text-xs text-zinc-400">{spot.address}</span>
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
