"use client"
import {getCookiesAccepted} from '@/lib/cookies';
// SessionContext.tsx

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import {deleteSessionFromCookies, saveSessionToCookies} from '@/lib/cookies';

// Define the shape of the session context
interface SessionContextType {
  sessionToken: string;
  setSessionToken: React.Dispatch<React.SetStateAction<string>>;
  sessionUUID: string;
  setSessionUUID: React.Dispatch<React.SetStateAction<string>>;
  isUserLoggedIn: boolean;
  setIsUserLoggedIn: React.Dispatch<React.SetStateAction<boolean>>;
}

// Create a default value for the context
const defaultSessionContext: SessionContextType = {
  sessionToken: "",
  setSessionToken: () => {},
  sessionUUID: "",
  setSessionUUID: () => {},
  isUserLoggedIn: false,
  setIsUserLoggedIn: () => {},
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
  const [isUserLoggedIn, setIsUserLoggedIn] = useState<boolean>(false);
  const [isCookiesAccepted, setIsCookiesAccepted] = useState<boolean>(false);

  const initializeSession = async () => {
    const { token, uuid } = await fetchSessionData();
    setSessionToken(token);
    setSessionUUID(uuid);
  };

  const myfirstfunction = async () => {
    setIsCookiesAccepted(await getCookiesAccepted());
  }

  useEffect(() => {
    initializeSession();
    myfirstfunction();
  }, []);

  useEffect(() => {
    console.log("USE EFFECT FIRING");
    console.log("START OF USE EFFECT >>>", sessionToken, sessionUUID, isUserLoggedIn)
    try{
      const update = async () => {
        console.log("UPDATE CALLED >>>", sessionToken, sessionUUID)
        console.log("TRY COMMIT COOKIES >>>", isCookiesAccepted, sessionToken, sessionUUID, isUserLoggedIn)
        
        if(!isCookiesAccepted || (!sessionToken && !sessionUUID)) {
          await deleteSessionFromCookies();
          return;
        }
        
        if(isCookiesAccepted && sessionToken && sessionUUID) {
          console.log("NOW TRY TO REALLY SET COOKIES")
          await saveSessionToCookies(sessionToken, sessionUUID);
          console.log("AFTER SET COOKIES")
          return;
        }
      }
      update()
      console.log("END OF USE EFFECT")

    } catch (error){
      console.log("USE EFFECT ERROR", error)
    }
    finally {
      console.log("USE EFFECT DID FIRE")
    }

    // tryCommitToCookies()
  }, [sessionToken, sessionUUID, isUserLoggedIn])

  return (
    <SessionContext.Provider value={{ sessionToken, setSessionToken, sessionUUID, setSessionUUID, isUserLoggedIn, setIsUserLoggedIn }}>
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
