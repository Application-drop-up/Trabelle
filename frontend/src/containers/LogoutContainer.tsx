"use client";

import { useCallback } from "react";

import { useUser } from "@/hooks/useUser";

type UseLogoutContainerReturn = {
  loading: boolean;
  error: string | null;
  onLogout: () => Promise<boolean>;
};

export function useLogoutContainer(): UseLogoutContainerReturn {
  const { logoutUser, loading, error } = useUser();

  const onLogout = useCallback(async (): Promise<boolean> => {
    return logoutUser();
  }, [logoutUser]);

  return { loading, error, onLogout };
}
