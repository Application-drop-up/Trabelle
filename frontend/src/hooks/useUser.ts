"use client";

import { useCallback, useState } from "react";

import { apiClient } from "@/lib/apiClient";
import type {
  LoginStartInput,
  LoginVerifyInput,
  MessageResponse,
  RegisterUserInput,
  UpdateUserInput,
  User,
} from "@/domain/user/types";

export function useUser() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const registerUser = useCallback(async (input: RegisterUserInput): Promise<User | null> => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.post<User>("/api/v1/user/register", input);
      setUser(result);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to register user");
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const updateUser = useCallback(
    async (id: string, input: UpdateUserInput): Promise<User | null> => {
      setLoading(true);
      setError(null);
      try {
        const result = await apiClient.patch<User>(`/api/v1/user/${id}`, input);
        setUser(result);
        return result;
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to update user");
        return null;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const deleteUser = useCallback(async (id: string): Promise<boolean> => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.delete(`/api/v1/user/${id}`);
      setUser(null);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete user");
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const loginStart = useCallback(async (input: LoginStartInput): Promise<boolean> => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.post<MessageResponse>("/api/v1/login", input);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to start login");
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const loginVerify = useCallback(async (input: LoginVerifyInput): Promise<boolean> => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.post<MessageResponse>("/api/v1/login/verify", input);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to verify login");
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchCurrentUser = useCallback(async (): Promise<User | null> => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.get<User>("/api/v1/user/me");
      setUser(result);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch current user");
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const logoutUser = useCallback(async (): Promise<boolean> => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.post<void>("/api/v1/logout", {});
      setUser(null);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to log out");
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    user,
    setUser,
    loading,
    error,
    registerUser,
    updateUser,
    deleteUser,
    loginStart,
    loginVerify,
    fetchCurrentUser,
    logoutUser,
  };
}
