"use server";

import { cookies } from "next/headers";
import { getCookiesAccepted } from "./cookies";
import { getDecodedAccessToken } from "./jwt";

interface SessionData {
  token: string;
  uuid: string;
}

const tokenCookieName = "magpie_token";
const uuidCookieName = "magpie_uuid";

// export async function getToken() {
//   if (await getCookiesAccepted()) {
//     console.log("cookies accepted");
//     await loadSessionFromCookies();
//   }
//   return sessionToken;
// }

export async function getToken(): Promise<string> {
  const session = await getSession();
  return session?.token ?? "";
}

// export async function getUUID() {
//   if (await getCookiesAccepted()) {
//     await loadSessionFromCookies();
//   }
//   return sessionUUID;
// }

export async function getUUID(): Promise<string> {
  const session = await getSession();
  return session?.uuid ?? "";
}

export async function getSession(): Promise<SessionData | null> {
  // Only attempt to get cookies if they're accepted
  if (await getCookiesAccepted()) {
    const cookieStore = cookies();
    const tokenCookie = cookieStore.get(tokenCookieName);
    const uuidCookie = cookieStore.get(uuidCookieName);

    if (tokenCookie && uuidCookie) {
      return {
        token: tokenCookie.value,
        uuid: uuidCookie.value,
      };
    }
  } else {
    const data = localStorage.getItem("session");
    const parsedData = data ? JSON.parse(data) : null;
    return {
      token: parsedData?.token ?? "",
      uuid: parsedData?.uuid ?? "",
    };
  }
  return null;
}

export async function setSession(token: string, uuid: string): Promise<void> {
  // Only set cookies if they're accepted
  if (await getCookiesAccepted()) {
    await commitSessionToCookies(token, uuid);
  }
  // If cookies aren't accepted, we don't store the session anywhere on the server
}

export async function deleteSessionFromCookies(): Promise<void> {
  const cookieStore = cookies();
  cookieStore.delete(tokenCookieName);
  cookieStore.delete(uuidCookieName);
}

export async function commitSessionToCookies(
  token: string,
  uuid: string
): Promise<void> {
  const cookieStore = cookies();

  const tokenInfo = await getDecodedAccessToken(token);
  const expiry = tokenInfo?.exp
    ? new Date(tokenInfo.exp * 1000)
    : new Date(Date.now() + 86400 * 1000); // 1 day default

  const cookieOptions = {
    expires: expiry,
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax" as const,
    path: "/",
  };

  cookieStore.set(tokenCookieName, token, cookieOptions);
  cookieStore.set(uuidCookieName, uuid, cookieOptions);
}

// export async function loadSessionFromCookies() {
//   const cookieStore = cookies();
//   if (!(await getCookiesAccepted())) {
//     deleteSessionFromCookies();
//     return;
//   }
//   const tokenCookie = cookieStore.get(tokenCookieName);
//   const uuidCookie = cookieStore.get(uuidCookieName);

//   if (tokenCookie && uuidCookie) {
//     sessionToken = tokenCookie.value;
//     sessionUUID = uuidCookie.value;
//     return true;
//   }
//   return false;
// }

// export async function setToken(token: string) {
//   sessionToken = token;
// }

// export async function setUUID(uuid: string) {
//   sessionUUID = uuid;
// }
