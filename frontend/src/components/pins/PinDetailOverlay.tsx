"use client";

import { useState } from "react";

import { usePinUpdateContainer } from "@/containers/PinUpdateContainer";
import { PIN_CATEGORIES, PIN_CATEGORY_LABELS } from "@/domain/pins/labels";
import type { PinCategory } from "@/domain/pins/types";
import type { PinViewModel } from "@/mappers/planMapper";

type Props = {
  planId: string;
  pin: PinViewModel;
  onUpdated: (pin: PinViewModel) => void;
  onDelete: () => void;
  onClose: () => void;
};

const COLOURS = ["#FF5733", "#33FF57", "#3357FF", "#FF33A8", "#FFD700"];

export function PinDetailOverlay({ planId, pin, onUpdated, onDelete, onClose }: Props) {
  const [isEditing, setIsEditing] = useState(false);
  const [category, setCategory] = useState<PinCategory>(pin.category);
  const [colour, setColour] = useState(pin.colour);
  const { onUpdate, loading, error } = usePinUpdateContainer();

  const handleSave = async () => {
    const vm = await onUpdate(planId, pin.id, { category, colour });
    if (vm) {
      onUpdated(vm);
      setIsEditing(false);
    }
  };

  const handleCancelEdit = () => {
    setCategory(pin.category);
    setColour(pin.colour);
    setIsEditing(false);
  };

  return (
    <div className="absolute bottom-4 left-1/2 z-10 w-72 -translate-x-1/2 rounded-xl bg-white p-4 shadow-lg">
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-2">
          <span
            className="mt-0.5 h-3 w-3 flex-shrink-0 rounded-full"
            style={{ backgroundColor: isEditing ? colour : pin.colour }}
          />
          <span className="text-sm font-medium text-black">{pin.name}</span>
        </div>
        <button
          type="button"
          onClick={onClose}
          className="text-zinc-400 hover:text-zinc-700"
          aria-label="閉じる"
        >
          ✕
        </button>
      </div>

      {isEditing ? (
        <div className="mt-3 flex flex-col gap-3">
          <div className="flex flex-col gap-1">
            <label className="text-xs text-zinc-500">カテゴリ</label>
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value as PinCategory)}
              disabled={loading}
              className="rounded border px-2 py-1.5 text-sm text-black"
            >
              {PIN_CATEGORIES.map((value) => (
                <option key={value} value={value}>
                  {PIN_CATEGORY_LABELS[value]}
                </option>
              ))}
            </select>
          </div>
          <div className="flex flex-col gap-1">
            <label className="text-xs text-zinc-500">色</label>
            <div className="flex gap-2">
              {COLOURS.map((c) => (
                <button
                  key={c}
                  type="button"
                  onClick={() => setColour(c)}
                  className={`h-6 w-6 rounded-full border-2 transition-transform ${
                    colour === c ? "scale-110 border-zinc-900" : "border-transparent"
                  }`}
                  style={{ backgroundColor: c }}
                />
              ))}
            </div>
          </div>
          {error && <p className="text-xs text-red-600">{error}</p>}
          <div className="flex gap-2 pt-1">
            <button
              type="button"
              onClick={handleSave}
              disabled={loading}
              className="flex-1 rounded-lg bg-black px-3 py-2 text-sm font-medium text-white disabled:opacity-50"
            >
              {loading ? "保存中…" : "保存"}
            </button>
            <button
              type="button"
              onClick={handleCancelEdit}
              disabled={loading}
              className="rounded-lg border px-3 py-2 text-sm text-zinc-700 disabled:opacity-50"
            >
              キャンセル
            </button>
          </div>
        </div>
      ) : (
        <div className="mt-3 flex flex-col gap-3">
          <span className="text-xs text-zinc-500">{PIN_CATEGORY_LABELS[pin.category]}</span>
          <div className="flex gap-2">
            <button
              type="button"
              onClick={() => setIsEditing(true)}
              className="flex-1 rounded-lg border px-3 py-2 text-sm text-zinc-700"
            >
              編集
            </button>
            <button
              type="button"
              onClick={onDelete}
              className="flex-1 rounded-lg bg-red-500 px-3 py-2 text-sm font-medium text-white"
            >
              削除
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
