"use client";

import { useCallback, useState } from "react";

import type { User } from "@/domain/user/types";
import { useUserContext } from "@/components/UserProvider";

type UseUpdateProfileContainerReturn = {
  name: string;
  email: string;
  loading: boolean;
  error: string | null;
  onChangeName: (value: string) => void;
  onChangeEmail: (value: string) => void;
  onSubmitUpdate: () => Promise<User | null>;
};

export function useUpdateProfileContainer(
  userId: string,
  initial: { name: string; email: string },
): UseUpdateProfileContainerReturn {
  const [name, setName] = useState(initial.name);
  const [email, setEmail] = useState(initial.email);
  const { updateUser, loading, error } = useUserContext();

  const onChangeName = useCallback((value: string) => {
    setName(value);
  }, []);

  const onChangeEmail = useCallback((value: string) => {
    setEmail(value);
  }, []);

  const onSubmitUpdate = useCallback(async (): Promise<User | null> => {
    if (!name.trim() || !email.trim()) return null;
    return updateUser(userId, { name: name.trim(), email: email.trim() });
  }, [userId, name, email, updateUser]);

  return { name, email, loading, error, onChangeName, onChangeEmail, onSubmitUpdate };
}
