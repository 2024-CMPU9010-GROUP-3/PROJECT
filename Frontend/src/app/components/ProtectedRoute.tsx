"use client"; 

import { useEffect } from "react";
import { useRouter } from 'next/compat/router';

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const router = useRouter();

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            if (router) {
                router.push('/login'); // if no token, redirect to login page
            }
        }
    }, [router]);

    return <>{children}</>;
};

export default ProtectedRoute;
