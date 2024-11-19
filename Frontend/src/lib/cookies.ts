import { jwtDecode, JwtPayload } from "jwt-decode";
import Cookies from "js-cookie";

const cookiesAcceptedName = "magpie_cookies_accepted";
const tokenCookieName = "magpie_token";
const uuidCookieName = "magpie_uuid";
// const cookiesAcceptedMaxAge = 365 * 24 * 60 * 60; // one year

export function getCookiesAccepted() {
  try {
    const cookiesAccepted = Cookies.get(cookiesAcceptedName);
    return cookiesAccepted === "true";
  } catch (error) {
    console.error("Error getting cookie:", error);
    return false;
  }
}

export function setCookiesAccepted() {
  try {
    Cookies.set(cookiesAcceptedName, "true", {
      expires: 365, // 1 year
      path: "/",
      sameSite: "strict",
      secure: true, // for HTTPS
    });
    return true;
  } catch (error) {
    console.error("Error setting cookie:", error);
    return false;
  }
}

export async function unsetCookiesAccepted() {
  Cookies.remove(cookiesAcceptedName, { path: "/" });
  Cookies.remove(tokenCookieName, { path: "/" });
  Cookies.remove(uuidCookieName, { path: "/" });
}

export function saveSessionToCookies(
  sessionToken: string,
  sessionUUID: string
) {
  try {
    const decoded = jwtDecode<JwtPayload>(sessionToken);
    const expires = decoded.exp
      ? new Date(decoded.exp * 1000)
      : new Date(Date.now() + 86400000);

    const options = {
      expires,
      path: "/",
      sameSite: "strict" as const,
      secure: true,
    };

    Cookies.set(tokenCookieName, sessionToken, options);
    Cookies.set(uuidCookieName, sessionUUID, options);
    return true;
  } catch (error) {
    console.error("Error saving session:", error);
    return false;
  }
}

export function loadSessionFromCookies() {
  try {
    const token = Cookies.get(tokenCookieName);
    const uuid = Cookies.get(uuidCookieName);
    return {token: token, uuid: uuid};
  } catch (error) {
    console.log("Error loading session from cookies:", error);
    return {token: undefined, uuid: undefined}
  }
}

export function deleteSessionFromCookies() {
  try {
    Cookies.remove(tokenCookieName, { path: "/" });
    Cookies.remove(uuidCookieName, { path: "/" });
    return true;
  } catch (error) {
    console.error("Error deleting session:", error);
    return false;
  }
}
