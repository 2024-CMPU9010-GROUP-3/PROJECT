'use client';

import { useUser } from '@/app/contexts/UserContext';
import { LogoutButton } from './LogoutButton';
import Link from 'next/link'; 

export function NavBar() {
  const { user } = useUser();

  return (
    <nav>
      {user ? (
        <div>
          <span>Welcome, {user.firstName}</span>
          <LogoutButton />
        </div>
      ) : (
        <Link href="/login">Login</Link>
      )}
    </nav>
  );
}
