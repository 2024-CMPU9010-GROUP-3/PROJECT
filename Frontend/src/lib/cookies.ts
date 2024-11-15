'use server'

import {cookies} from "next/headers";

const cookiesAcceptedName = "magpie_cookies_accepted"
const cookiesAcceptedMaxAge = 365 * 24 * 60 * 60; // one year

export async function getCookiesAccepted() {
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