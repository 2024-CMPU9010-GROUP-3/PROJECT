'use client';

import { Button } from "@/components/ui/registry/button";
import { logout } from "@/app/actions";
import { useRouter } from "next/navigation";

export function LogoutButton() {
  const router = useRouter();

  const handleLogout = async () => {
    await logout();
    router.push('/login');
  };

  return (
    <Button onClick={handleLogout}>
      登出
    </Button>
  );
}
