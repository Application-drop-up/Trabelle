"use client";

import { useRegisterContainer } from "@/containers/RegisterContainer";

export function RegisterView() {
  const {
    email,
    password,
    name,
    loading,
    error,
    onChangeEmail,
    onChangePassword,
    onChangeName,
    onSubmit,
  } = useRegisterContainer();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onSubmit();
  };

  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-6">
      <h1 className="text-2xl font-semibold">新規登録</h1>
      <form onSubmit={handleSubmit} className="flex w-full max-w-sm flex-col gap-4">
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
        <label className="flex flex-col gap-1 text-sm">
          パスワード
          <input
            type="password"
            value={password}
            onChange={(e) => onChangePassword(e.target.value)}
            disabled={loading}
            className="rounded border px-3 py-2 outline-none focus:ring-2"
          />
          <span className="text-xs text-gray-500">8文字以上の半角英数字</span>
        </label>
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
        {error && <p className="text-sm text-red-600">{error}</p>}
        <button
          type="submit"
          disabled={loading || !email.trim() || !password.trim() || !name.trim()}
          className="rounded bg-black px-4 py-2 text-white disabled:opacity-50"
        >
          {loading ? "登録中…" : "登録する"}
        </button>
        <p className="text-center text-xs text-gray-500">
          登録することで利用規約に同意したものとします
        </p>
      </form>
    </div>
  );
}
