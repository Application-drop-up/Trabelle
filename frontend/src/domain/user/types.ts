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
}
