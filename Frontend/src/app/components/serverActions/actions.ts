'use server';

import {cookies} from 'next/headers';
import { jwtDecode, JwtPayload } from "jwt-decode";

const tokenCookieName = "magpie_token";
const uuidCookieName = "magpie_uuid";

export async function saveSessionToCookies(sessionToken: string, sessionUUID: string) {
  const cookieStore = cookies();

  console.log("SAVE SESSION TO COOKIES >>>", sessionToken, sessionUUID)

  // Set cookies to store sessionToken and sessionUUID
  let expiryDate = new Date(Date.now() + 86400 * 1000); // 1 day expiry default
  const decoded = jwtDecode<JwtPayload>(sessionToken)
  if(decoded.exp){
    expiryDate = new Date(decoded.exp * 1000);
  }
  cookieStore.set(tokenCookieName, sessionToken, { expires: expiryDate });
  cookieStore.set(uuidCookieName, sessionUUID, { expires: expiryDate });
}

export async function deleteSessionFromCookies() {
  const cookieStore = cookies();
  cookieStore.delete(tokenCookieName);
  cookieStore.delete(uuidCookieName);
}
