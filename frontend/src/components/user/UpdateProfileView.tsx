"use client";

import { useUpdateProfileContainer } from "@/containers/UpdateProfileContainer";

type UpdateProfileViewProps = {
  userId: string;
  initial: { name: string; email: string };
};

export function UpdateProfileView({ userId, initial }: UpdateProfileViewProps) {
  const { name, email, loading, error, onChangeName, onChangeEmail, onSubmitUpdate } =
    useUpdateProfileContainer(userId, initial);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onSubmitUpdate();
  };

  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-6">
      <h1 className="text-2xl font-semibold">プロフィール編集</h1>
      <form onSubmit={handleSubmit} className="flex w-full max-w-sm flex-col gap-4">
        <label className="flex flex-col gap-1 text-sm">
          ユーザー名
          <input
            type="text"
            value={name}
            onChange={(e) => onChangeName(e.target.value)}
            disabled={loading}
            className="rounded border px-3 py-2 outline-none focus:ring-2"
          />
        </label>
        <label className="flex flex-col gap-1 text-sm">
          メールアドレス
          <input
            type="email"
            value={email}
            onChange={(e) => onChangeEmail(e.target.value)}
            disabled={loading}
            className="rounded border px-3 py-2 outline-none focus:ring-2"
          />
        </label>
        {error && <p className="text-sm text-red-600">{error}</p>}
        <button
          type="submit"
          disabled={loading || !name.trim() || !email.trim()}
          className="rounded bg-black px-4 py-2 text-white disabled:opacity-50"
        >
          {loading ? "更新中…" : "更新する"}
        </button>
      </form>
    </div>
  );
}
