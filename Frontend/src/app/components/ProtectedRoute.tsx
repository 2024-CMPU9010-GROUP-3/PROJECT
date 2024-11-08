"use client";

import { useEffect } from "react";
import { useRouter, usePathname } from 'next/navigation';

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const router = useRouter();
    const pathname = usePathname();

    useEffect(() => {
        const userId = localStorage.getItem('userId');
        const unprotectedPaths = ['/signup', '/forgot-password']; // add unprotected paths

        if (!userId && !unprotectedPaths.includes(pathname)) {
            router.push('/login'); // if no token, and not in unprotected paths, redirect to login page
        }
    }, [router, pathname]);

    return <>{children}</>;
};

export default ProtectedRoute;
