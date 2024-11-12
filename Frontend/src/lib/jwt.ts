'use server'

import { jwtDecode } from "jwt-decode";

export async function getDecodedAccessToken(token: string) {
  try {
    return jwtDecode(token);
  } catch(Error) {
    return null;
  }
}