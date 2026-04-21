import type { UserRole } from "../enums/user-role.enum";

export interface IUser {
  id: string;
  /**
   * Creation timestamp
   */
  createdAt: string;
  /**
   * IUser email
   */
  email: string;
  role: UserRole;
}
