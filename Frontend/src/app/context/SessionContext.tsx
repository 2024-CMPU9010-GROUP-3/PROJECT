"use client"
import {getCookiesAccepted} from '@/lib/cookies';
// SessionContext.tsx

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import {deleteSessionFromCookies, saveSessionToCookies} from '../components/serverActions/actions';

// Define the shape of the session context
interface SessionContextType {
  sessionToken: string;
  setSessionToken: React.Dispatch<React.SetStateAction<string>>;
  sessionUUID: string;
  setSessionUUID: React.Dispatch<React.SetStateAction<string>>;
}

// Create a default value for the context
const defaultSessionContext: SessionContextType = {
  sessionToken: "",
  setSessionToken: () => {},
  sessionUUID: "",
  setSessionUUID: () => {},
};

// Create context with the defined type
const SessionContext = createContext<SessionContextType>(defaultSessionContext);

// Define props for SessionProvider
interface SessionProviderProps {
  children: ReactNode;
}

// SessionProvider component
export const SessionProvider: React.FC<SessionProviderProps> = ({ children }) => {
  const [sessionToken, setSessionToken] = useState<string>("");
  const [sessionUUID, setSessionUUID] = useState<string>("");

  const initializeSession = async () => {
    const { token, uuid } = await fetchSessionData();
    setSessionToken(token);
    setSessionUUID(uuid);
  };

  useEffect(() => {
    initializeSession();
  }, []);

  useEffect(() => {
    const tryCommitToCookies = async () => {
      const cookiesAccepted = await getCookiesAccepted();

      if(!cookiesAccepted || (!sessionToken && !sessionUUID)) {
        await deleteSessionFromCookies();
        return;
      }

      if(cookiesAccepted && sessionToken && sessionUUID) {
        await saveSessionToCookies(sessionToken, sessionUUID);
        return;
      }
    }
    tryCommitToCookies()
  }, [sessionToken, sessionUUID])

  return (
    <SessionContext.Provider value={{ sessionToken, setSessionToken, sessionUUID, setSessionUUID }}>
      {children}
    </SessionContext.Provider>
  );
};

// Custom hook to access session context
export const useSession = () => useContext(SessionContext);

// Fetch session data from API
async function fetchSessionData() {
  const response = await fetch('/api/session');
  const data = await response.json();
  return { token: data.token, uuid: data.uuid };
}
