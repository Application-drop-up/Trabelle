"use client";

import { useCallback, useState } from "react";

import { apiClient, voidSchema } from "@/lib/apiClient";
import { errorMessages } from "@/lib/messages";
import {
  messageResponseSchema,
  userSchema,
  type LoginStartInput,
  type LoginVerifyInput,
  type RegisterUserInput,
  type UpdateUserInput,
  type User,
} from "@/domain/user/types";

export function useUser() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const registerUser = useCallback(async (input: RegisterUserInput): Promise<User | null> => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.post(userSchema, "/api/v1/user/register", input);
      setUser(result);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.user.register);
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
        const result = await apiClient.patch(userSchema, `/api/v1/user/${id}`, input);
        setUser(result);
        return result;
      } catch (err) {
        setError(err instanceof Error ? err.message : errorMessages.user.update);
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
      setError(err instanceof Error ? err.message : errorMessages.user.delete);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const loginStart = useCallback(async (input: LoginStartInput): Promise<boolean> => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.post(messageResponseSchema, "/api/v1/login", input);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.user.loginStart);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const loginVerify = useCallback(async (input: LoginVerifyInput): Promise<boolean> => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.post(messageResponseSchema, "/api/v1/login/verify", input);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.user.loginVerify);
      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchCurrentUser = useCallback(async (): Promise<User | null> => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.get(userSchema, "/api/v1/user/me");
      setUser(result);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.user.fetchCurrentUser);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const logoutUser = useCallback(async (): Promise<boolean> => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.post(voidSchema, "/api/v1/logout", {});
      setUser(null);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : errorMessages.user.logout);
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
