"use client";

import { useCallback, useState } from "react";

import { useUserContext } from "@/components/UserProvider";

type LoginStep = "credentials" | "otp";

type UseLoginContainerReturn = {
  step: LoginStep;
  email: string;
  password: string;
  code: string;
  loading: boolean;
  error: string | null;
  onChangeEmail: (value: string) => void;
  onChangePassword: (value: string) => void;
  onChangeCode: (value: string) => void;
  onSubmitCredentials: () => Promise<boolean>;
  onSubmitCode: () => Promise<boolean>;
};

export function useLoginContainer(): UseLoginContainerReturn {
  const [step, setStep] = useState<LoginStep>("credentials");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [code, setCode] = useState("");
  const { loginStart, loginVerify, loading, error } = useUserContext();

  const onChangeEmail = useCallback((value: string) => {
    setEmail(value);
  }, []);

  const onChangePassword = useCallback((value: string) => {
    setPassword(value);
  }, []);

  const onChangeCode = useCallback((value: string) => {
    setCode(value);
  }, []);

  const onSubmitCredentials = useCallback(async (): Promise<boolean> => {
    if (!email.trim() || !password.trim()) return false;
    const success = await loginStart({ email: email.trim(), password: password.trim() });
    if (success) setStep("otp");
    return success;
  }, [email, password, loginStart]);

  const onSubmitCode = useCallback(async (): Promise<boolean> => {
    if (!code.trim()) return false;
    return loginVerify({ email: email.trim(), code: code.trim() });
  }, [email, code, loginVerify]);

  return {
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
  };
}
