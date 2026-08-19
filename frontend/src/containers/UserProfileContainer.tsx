"use client";

import { useEffect } from "react";

import type { User } from "@/domain/user/types";
import { useUserContext } from "@/components/UserProvider";

type UseUserProfileContainerReturn = {
  user: User | null;
  loading: boolean;
  error: string | null;
};

export function useUserProfileContainer(): UseUserProfileContainerReturn {
  const { user, loading, error, fetchCurrentUser } = useUserContext();

  useEffect(() => {
    fetchCurrentUser();
  }, [fetchCurrentUser]);

  return { user, loading, error };
}
