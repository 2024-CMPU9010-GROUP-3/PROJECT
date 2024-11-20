"use client";

import { useEffect } from "react";
import { useRouter, usePathname } from 'next/navigation';
import { useSession } from "@/app/context/SessionContext"


const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const router = useRouter();
  const pathname = usePathname();
  const { sessionToken } = useSession();

    useEffect(() => {
      if (!sessionToken) {
        router.push('/login');
      }
    }, [router, pathname, sessionToken]);

    return <>{children}</>;
};

export default ProtectedRoute;
