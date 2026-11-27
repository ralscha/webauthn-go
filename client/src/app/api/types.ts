/* Do not change, this code is generated from Golang structs */

export interface PasskeyInfo {
  label: string;
  createdAt: string;
  lastUsed: string | null;
  current: boolean;
}
export interface SecretOutput {
  passkeys: PasskeyInfo[];
  message: string;
}
export interface UsernameInput {
  username: string;
}
export interface Errors {
  errors: { [key: string]: string[] };
}
