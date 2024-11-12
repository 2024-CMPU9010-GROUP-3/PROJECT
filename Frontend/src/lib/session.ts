'use server'

import { cookies } from 'next/headers'
import {getCookiesAccepted} from './cookies';
import {getDecodedAccessToken} from './jwt';

const tokenCookieName = "magpie_token"
const uuidCookieName = "magpie_uuid"


let sessionToken = "";
let sessionUUID = "";

export async function getToken() {
  if (await getCookiesAccepted()) {
    console.log("cookies accepted")
    await loadSessionFromCookies();
  }
  return sessionToken;
}

export async function getUUID() {
  if (await getCookiesAccepted()) {
    await loadSessionFromCookies();
  }
  return sessionUUID;
}

export async function commitSessionToCookies(){
  const cookieStore = cookies()
  if (!(await getCookiesAccepted())){
    return;
  }

  // default to 1 day
  let expiry = 86400*1000;
  const tokenInfo = await getDecodedAccessToken(sessionToken);
  if(tokenInfo && tokenInfo.exp){
    expiry = tokenInfo.exp * 1000;
  }
  const expiryDate = new Date(expiry);
  cookieStore.set(tokenCookieName, sessionToken, {expires: expiryDate});
  cookieStore.set(uuidCookieName, sessionUUID, {expires: expiryDate})
}

export async function deleteSessionFromCookies(){
  const cookieStore = cookies()
  cookieStore.delete(tokenCookieName)
  cookieStore.delete(uuidCookieName)
}

export async function loadSessionFromCookies(){
  const cookieStore = cookies()
  if (!(await getCookiesAccepted())){
    deleteSessionFromCookies();
    return;
  }
  const tokenCookie = cookieStore.get(tokenCookieName);
  const uuidCookie = cookieStore.get(uuidCookieName);

  if(tokenCookie && uuidCookie){
    sessionToken = tokenCookie.value;
    sessionUUID = uuidCookie.value;
    return true;
  }
  return false;
}

export async function setToken(token: string) {
  sessionToken = token;
}

export async function setUUID(uuid: string) {
  sessionUUID = uuid;
}
