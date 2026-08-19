export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
}

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

export interface MessageResponse {
  message: string;
}
