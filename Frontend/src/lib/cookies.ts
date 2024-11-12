'use server'

import {cookies} from "next/headers";

const cookiesAcceptedName = "magpie_cookies_accepted"

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
  await cookieStore.set(cookiesAcceptedName, "true");
}

export async function unsetCookiesAccepted() {
  const cookieStore = cookies();
  await cookieStore.delete(cookiesAcceptedName);
}