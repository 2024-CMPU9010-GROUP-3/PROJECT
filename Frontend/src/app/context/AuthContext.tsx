// src/app/context/AuthContext.tsx

"use client";

import React, { createContext, useContext, useState, useEffect } from "react";
import { getToken, getUUID } from "@/lib/session";

interface AuthContextType {
  isLoggedIn: boolean;
  loading: boolean;
  setIsLoggedIn: (value: boolean) => void; // 添加 setIsLoggedIn 属性
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  setSessionData: (value: any) => void; // add setSessionData property
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  sessionData: any; // add sessionData property
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [loading, setLoading] = useState(true);
  const [sessionData, setSessionData] = useState(null); // session data

  useEffect(() => {
    const checkSession = async () => {
      const token = await getToken();
      const uuid = await getUUID();
      setIsLoggedIn(!!token && !!uuid);
      setLoading(false);
    };

    checkSession();
  }, []);

  return (
    <AuthContext.Provider
      value={{
        isLoggedIn,
        loading,
        setIsLoggedIn,
        setSessionData,
        sessionData,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
