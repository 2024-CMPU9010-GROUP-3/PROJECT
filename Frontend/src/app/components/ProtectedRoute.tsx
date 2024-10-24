'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { verifySession } from '@/lib/dal';

export default function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const router = useRouter();

  useEffect(() => {
    const checkAuth = async () => {
      const session = await verifySession();
      if (!session) {
        router.push('/login');
      }
    };

    checkAuth();
  }, [router]);

  return <>{children}</>;
}
