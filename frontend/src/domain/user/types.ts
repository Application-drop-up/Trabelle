import { z } from "zod";

export const userSchema = z.object({
  id: z.string(),
  email: z.string(),
  name: z.string(),
  created_at: z.string(),
});

export type User = z.infer<typeof userSchema>;

export const messageResponseSchema = z.object({
  message: z.string(),
  // Only present in non-production environments (see backend LoginStart handler).
  code: z.string().optional(),
});

export type MessageResponse = z.infer<typeof messageResponseSchema>;

export interface RegisterUserInput {
  email: string;
  password: string;
  name: string;
}

export interface UpdateUserInput {
  name: string;
  email: string;
}

export interface LoginStartInput {
  email: string;
  password: string;
}

export interface LoginVerifyInput {
  email: string;
  code: string;
}
