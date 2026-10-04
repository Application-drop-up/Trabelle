"use client";

import Link from "next/link";
import { Alert, CircularProgress } from "@mui/material";

import { usePlansListContainer } from "@/containers/PlansListContainer";

export function PlansListView() {
  const { plans, loading, error } = usePlansListContainer();

  return (
    <div className="flex flex-1 flex-col items-center gap-6 p-6">
      <div className="flex w-full max-w-sm items-center justify-between">
        <h1 className="text-2xl font-semibold">マイプラン</h1>
        <Link href="/plans/new" className="text-sm text-blue-600 underline">
          新規作成
        </Link>
      </div>

      {loading && <CircularProgress aria-label="読み込み中" />}

      {!loading && error && (
        <Alert severity="error" className="w-full max-w-sm">
          {error}
        </Alert>
      )}

      {!loading && !error && plans.length === 0 && (
        <p className="text-sm text-gray-500">参加しているプランはまだありません</p>
      )}

      {!loading && !error && plans.length > 0 && (
        <ul className="flex w-full max-w-sm flex-col gap-2">
          {plans.map((plan) => (
            <li key={plan.id}>
              <Link
                href={`/plans/${plan.share_token}`}
                className="block rounded border px-4 py-3 text-sm hover:bg-gray-50"
              >
                {plan.title}
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
