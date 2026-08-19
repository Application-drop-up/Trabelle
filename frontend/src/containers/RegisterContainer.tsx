"use client";

import { useCallback, useState } from "react";

import type { User } from "@/domain/user/types";
import { useUserContext } from "@/components/UserProvider";

type UseRegisterContainerReturn = {
  email: string;
  password: string;
  name: string;
  loading: boolean;
  error: string | null;
  onChangeEmail: (value: string) => void;
  onChangePassword: (value: string) => void;
  onChangeName: (value: string) => void;
  onSubmit: () => Promise<User | null>;
};

export function useRegisterContainer(): UseRegisterContainerReturn {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const { registerUser, loading, error } = useUserContext();

  const onChangeEmail = useCallback((value: string) => {
    setEmail(value);
  }, []);

  const onChangePassword = useCallback((value: string) => {
    setPassword(value);
  }, []);

  const onChangeName = useCallback((value: string) => {
    setName(value);
  }, []);

  const onSubmit = useCallback(async (): Promise<User | null> => {
    if (!email.trim() || !password.trim() || !name.trim()) return null;
    return registerUser({
      email: email.trim(),
      password: password.trim(),
      name: name.trim(),
    });
  }, [email, password, name, registerUser]);

  return {
    email,
    password,
    name,
    loading,
    error,
    onChangeEmail,
    onChangePassword,
    onChangeName,
    onSubmit,
  };
}
