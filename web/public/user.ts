export type UserRole = 'user' | 'moder' | 'admin';

export interface User {
  id: string;
  name: string;
  email: string;
  avatar?: string;
  role?: UserRole;
}