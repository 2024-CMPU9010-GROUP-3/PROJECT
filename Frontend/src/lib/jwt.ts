'use server'

import { jwtDecode } from "jwt-decode";

export async function getDecodedAccessToken(token: string) {
  try {
    return jwtDecode(token);
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch(error) {
    return null;
  }
}