"use client";

import { Alert, Button, TextField } from "@mui/material";

import { useLoginContainer } from "@/containers/LoginContainer";

export function LoginView() {
  const {
    step,
    email,
    password,
    code,
    loading,
    error,
    onChangeEmail,
    onChangePassword,
    onChangeCode,
    onSubmitCredentials,
    onSubmitCode,
  } = useLoginContainer();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (step === "credentials") {
      await onSubmitCredentials();
    } else {
      await onSubmitCode();
    }
  };

  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-6">
      <h1 className="text-2xl font-semibold">ログイン</h1>
      <form onSubmit={handleSubmit} className="flex w-full max-w-sm flex-col gap-4">
        {step === "credentials" ? (
          <>
            <TextField
              label="メールアドレス"
              type="email"
              value={email}
              onChange={(e) => onChangeEmail(e.target.value)}
              disabled={loading}
              fullWidth
            />
            <TextField
              label="パスワード"
              type="password"
              value={password}
              onChange={(e) => onChangePassword(e.target.value)}
              disabled={loading}
              fullWidth
            />
          </>
        ) : (
          <>
            <p className="text-sm text-gray-600">
              {email} 宛に送信した認証コードを入力してください
            </p>
            <TextField
              label="認証コード"
              value={code}
              onChange={(e) => onChangeCode(e.target.value)}
              disabled={loading}
              fullWidth
            />
          </>
        )}
        {error && <Alert severity="error">{error}</Alert>}
        <Button
          type="submit"
          variant="contained"
          disabled={
            loading || (step === "credentials" ? !email.trim() || !password.trim() : !code.trim())
          }
        >
          {step === "credentials"
            ? loading
              ? "送信中…"
              : "コードを送信"
            : loading
              ? "ログイン中…"
              : "ログイン"}
        </Button>
      </form>
    </div>
  );
}
