"use client"; // 确保这是一个客户端组件
import React, { createContext, useContext, useState, useEffect } from 'react';

// 更新 AuthContext 的默认值
const AuthContext = createContext<{ isLoggedIn: boolean; setIsLoggedIn: React.Dispatch<React.SetStateAction<boolean>> } | null>(null);

const handleLogin = async (credentials: { username: string; password: string }) => {
    const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);
    const [isLoggingIn, setIsLoggingIn] = useState<boolean>(false);
    setIsLoggingIn(true); // 添加状态以指示正在登录
    try {
        const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/auth/User/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(credentials),
        }); 

        const data = await response.json();
        if (data.success) {
            setIsLoggedIn(true);
            // 处理成功登录的逻辑
            localStorage.setItem("token", data.token);
        } else {
            // 处理登录失败的逻辑
            console.error(data.message);
        }
    } catch (error) {
        console.error("Login error:", error);
    } finally {
        setIsLoggingIn(false); // 无论成功与否，最终都要停止“logging in”
    }
};

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [isLoggedIn, setIsLoggedIn] = useState(false); // 管理登录状态
  
  useEffect(() => {
    const token = localStorage.getItem("token"); // 从 localStorage 获取 token
    console.log("useEffecttoken", token);
    if (token) {
      setIsLoggedIn(true); // 如果 token 存在，设置为已登录
    }
  }, []);

  return (
    <AuthContext.Provider value={{ isLoggedIn, setIsLoggedIn }}>
      {children}
    </AuthContext.Provider>
  );
};

// 导出 useAuth 钩子以便在其他组件中使用
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
