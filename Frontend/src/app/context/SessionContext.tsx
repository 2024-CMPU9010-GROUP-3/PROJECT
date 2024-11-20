"use client";

import { getCookiesAccepted, loadSessionFromCookies } from "@/lib/cookies";
import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import { deleteSessionFromCookies, saveSessionToCookies } from "@/lib/cookies";

// Define the shape of the session context
interface SessionContextType {
  sessionToken: string;
  setSessionToken: React.Dispatch<React.SetStateAction<string>>;
  sessionUUID: string;
  setSessionUUID: React.Dispatch<React.SetStateAction<string>>;
  isCookiesAccepted: boolean;
  setIsCookiesAccepted: React.Dispatch<React.SetStateAction<boolean>>;
}

// Create a default value for the context
const defaultSessionContext: SessionContextType = {
  sessionToken: "",
  setSessionToken: () => {},
  sessionUUID: "",
  setSessionUUID: () => {},
  isCookiesAccepted: false,
  setIsCookiesAccepted: () => {},
};

// Create context with the defined type
const SessionContext = createContext<SessionContextType>(defaultSessionContext);

// Define props for SessionProvider
interface SessionProviderProps {
  children: ReactNode;
}

// SessionProvider component
export const SessionProvider: React.FC<SessionProviderProps> = ({
  children,
}) => {
  const [sessionToken, setSessionToken] = useState<string>("");
  const [sessionUUID, setSessionUUID] = useState<string>("");
  const [isCookiesAccepted, setIsCookiesAccepted] = useState<boolean>(false);

  useEffect(() => {
    const { token, uuid } = loadSessionFromCookies()
    if (token && uuid){
      setSessionToken(token);
      setSessionUUID(uuid);
    }
    setIsCookiesAccepted(getCookiesAccepted());
  }, []);

  useEffect(() => {
    if (!isCookiesAccepted || (!sessionToken && !sessionUUID)) {
      deleteSessionFromCookies();
      return;
    }

    if (isCookiesAccepted && sessionToken && sessionUUID) {
      saveSessionToCookies(sessionToken, sessionUUID);
      return;
    }
  }, [sessionToken, sessionUUID, isCookiesAccepted]);

  return (
    <SessionContext.Provider
      value={{
        sessionToken,
        setSessionToken,
        sessionUUID,
        setSessionUUID,
        isCookiesAccepted,
        setIsCookiesAccepted
      }}
    >
      {children}
    </SessionContext.Provider>
  );
};

// Custom hook to access session context
export const useSession = () => useContext(SessionContext);
