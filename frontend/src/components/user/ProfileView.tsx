"use client";

import { Alert, CircularProgress } from "@mui/material";

import { useUserProfileContainer } from "@/containers/UserProfileContainer";

export function ProfileView() {
  const { user, loading, error } = useUserProfileContainer();

  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-6">
      <h1 className="text-2xl font-semibold">プロフィール</h1>

      {loading && <CircularProgress aria-label="読み込み中" />}

      {!loading && error && (
        <Alert severity="error" className="w-full max-w-sm">
          {error}
        </Alert>
      )}

      {!loading && !error && user && (
        <div className="flex w-full max-w-sm flex-col gap-4">
          <div className="flex flex-col gap-1 text-sm">
            <span className="text-gray-500">ユーザー名</span>
            <span>{user.name}</span>
          </div>
          <div className="flex flex-col gap-1 text-sm">
            <span className="text-gray-500">メールアドレス</span>
            <span>{user.email}</span>
          </div>
          <div className="flex flex-col gap-1 text-sm">
            <span className="text-gray-500">登録日</span>
            <span>{new Date(user.created_at).toLocaleDateString("ja-JP")}</span>
          </div>
        </div>
      )}
    </div>
  );
}
