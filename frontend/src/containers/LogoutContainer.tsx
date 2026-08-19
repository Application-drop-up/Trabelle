"use client";

import { useCallback } from "react";

import { useUserContext } from "@/components/UserProvider";

type UseLogoutContainerReturn = {
  loading: boolean;
  error: string | null;
  onLogout: () => Promise<boolean>;
};

export function useLogoutContainer(): UseLogoutContainerReturn {
  const { logoutUser, loading, error } = useUserContext();

  const onLogout = useCallback(async (): Promise<boolean> => {
    return logoutUser();
  }, [logoutUser]);

  return { loading, error, onLogout };
}
