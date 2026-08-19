"use client";

import { useRouter } from "next/navigation";

import { Button } from "@mui/material";

import { useUserContext } from "@/components/UserProvider";
import { useLogoutContainer } from "@/containers/LogoutContainer";

export function Header() {
  const router = useRouter();
  const { user } = useUserContext();
  const { onLogout } = useLogoutContainer();

  const handleLogout = async () => {
    const success = await onLogout();
    if (success) router.push("/login");
  };

  return (
    <header className="flex items-center justify-between border-b border-gray-800 px-4 py-3">
      <span className="font-semibold">Trabelle</span>
      {user && (
        <div className="flex items-center gap-3 text-sm">
          <span>{user.name} さん</span>
          <Button size="small" variant="outlined" onClick={handleLogout}>
            ログアウト
          </Button>
        </div>
      )}
    </header>
  );
}
