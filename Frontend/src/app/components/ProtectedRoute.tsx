"use client";

import { useEffect } from "react";
import { useRouter, usePathname } from 'next/navigation';
import { useSession } from "@/app/context/SessionContext"


const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const router = useRouter();
  const pathname = usePathname();
  const { sessionToken } = useSession();

    useEffect(() => {
      const checkAuthentication = () => {
        try {
            // Redirect to login if not authenticated and not on an unprotected path
            if (!sessionToken) {
              setTimeout(() => {
                router.push('/login');
              })
            }
        } catch (error) {
            console.error("Error during authentication check:", error);
        }
    };

    checkAuthentication();
    }, [router, pathname, sessionToken]);

    return <>{children}</>;
};

export default ProtectedRoute;
