"use client";

import DeleteIcon from "@mui/icons-material/Delete";

import { useDeleteAccountContainer } from "@/containers/DeleteAccountContainer";

type DeleteAccountViewProps = {
  userId: string;
};

export function DeleteAccountView({ userId }: DeleteAccountViewProps) {
  const { isConfirming, loading, error, onRequestDelete, onCancelDelete, onConfirmDelete } =
    useDeleteAccountContainer(userId);

  if (!isConfirming) {
    return (
      <button
        type="button"
        onClick={onRequestDelete}
        aria-label="アカウントを削除"
        className="rounded p-2 text-red-600 hover:bg-red-50"
      >
        <DeleteIcon />
      </button>
    );
  }

  return (
    <div className="flex flex-col gap-3 rounded border border-red-200 p-4">
      <p className="text-sm">本当に削除しますか？</p>
      {error && <p className="text-sm text-red-600">{error}</p>}
      <div className="flex gap-2">
        <button
          type="button"
          onClick={onCancelDelete}
          disabled={loading}
          className="rounded border px-4 py-2 disabled:opacity-50"
        >
          キャンセル
        </button>
        <button
          type="button"
          onClick={onConfirmDelete}
          disabled={loading}
          className="rounded bg-red-600 px-4 py-2 text-white disabled:opacity-50"
        >
          {loading ? "削除中…" : "削除する"}
        </button>
      </div>
    </div>
  );
}
