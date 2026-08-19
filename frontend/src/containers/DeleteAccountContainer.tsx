"use client";

import { useCallback, useState } from "react";

import { useUserContext } from "@/components/UserProvider";

type UseDeleteAccountContainerReturn = {
  isConfirming: boolean;
  loading: boolean;
  error: string | null;
  onRequestDelete: () => void;
  onCancelDelete: () => void;
  onConfirmDelete: () => Promise<boolean>;
};

export function useDeleteAccountContainer(userId: string): UseDeleteAccountContainerReturn {
  const [isConfirming, setIsConfirming] = useState(false);
  const { deleteUser, loading, error } = useUserContext();

  const onRequestDelete = useCallback(() => {
    setIsConfirming(true);
  }, []);

  const onCancelDelete = useCallback(() => {
    setIsConfirming(false);
  }, []);

  const onConfirmDelete = useCallback(async (): Promise<boolean> => {
    return deleteUser(userId);
  }, [userId, deleteUser]);

  return { isConfirming, loading, error, onRequestDelete, onCancelDelete, onConfirmDelete };
}
