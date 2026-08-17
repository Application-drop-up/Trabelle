"use client";

import { useEffect } from "react";

import type { User } from "@/domain/user/types";
import { useUser } from "@/hooks/useUser";

type UseUserProfileContainerReturn = {
  user: User | null;
  loading: boolean;
  error: string | null;
};

export function useUserProfileContainer(): UseUserProfileContainerReturn {
  const { user, loading, error, fetchCurrentUser } = useUser();

  useEffect(() => {
    fetchCurrentUser();
  }, [fetchCurrentUser]);

  return { user, loading, error };
}
