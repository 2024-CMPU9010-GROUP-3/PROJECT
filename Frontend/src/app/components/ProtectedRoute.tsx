"use client";

import { useEffect } from "react";
import { useRouter, usePathname } from 'next/navigation';
import {getToken, getUUID} from "@/lib/session";

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const router = useRouter();
    const pathname = usePathname();

    useEffect(() => {
      const checkAuthentication = async () => {
        try {
            const token = await getToken();
            // Redirect to login if not authenticated and not on an unprotected path
            if (!token) {
                router.push('/login');
            }
        } catch (error) {
            console.error("Error during authentication check:", error);
        }
    };

    checkAuthentication();
    }, [router, pathname]);

    return <>{children}</>;
};

export default ProtectedRoute;
