'use server'

import {jwtDecode, JwtPayload} from "jwt-decode";
import {cookies} from "next/headers";

const cookiesAcceptedName = "magpie_cookies_accepted"
const tokenCookieName = "magpie_token";
const uuidCookieName = "magpie_uuid";
const cookiesAcceptedMaxAge = 365 * 24 * 60 * 60; // one year

export async function getCookiesAccepted() {
  console.log("GET COOKIES ACCEPTED CALLED")
  const cookieStore = cookies()
  const cookiesAccepted = cookieStore.get(cookiesAcceptedName);
  if (cookiesAccepted && cookiesAccepted.value === "true"){
    return true;
  } else {
    return false;
  }
}

export async function setCookiesAccepted(){
  const cookieStore = cookies();
  cookieStore.set({
    name: cookiesAcceptedName,
    value: "true",
    maxAge: cookiesAcceptedMaxAge,
  });
}

export async function unsetCookiesAccepted() {
  const cookieStore = cookies();
  cookieStore.delete(cookiesAcceptedName);
}

export async function saveSessionToCookies(sessionToken: string, sessionUUID: string) {
  console.log("SAVE SESSION TO COOKIES >>>", sessionToken, sessionUUID)
  const cookieStore = cookies();


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