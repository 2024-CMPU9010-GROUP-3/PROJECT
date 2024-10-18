import {
    LifeBuoy,
    LogOut,
    Settings,
    User,
    UserRoundX,
} from "lucide-react";

import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuShortcut,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import { useRouter } from "next/navigation"; // 导入 useRouter
import { useAuth } from "@/app/components/AuthContext"; // 导入 AuthContext

type DropdownMenuDemoProps = {
  avatar: React.ReactNode; // 添加 avatar 属性的类型定义
};


export function DropdownMenuDemo() {
    const { isLoggedIn, setIsLoggedIn } = useAuth(); // 使用 AuthContext
    const router = useRouter(); // 使用 useRouter 进行页面跳转
    console.log("isLoggedIn", isLoggedIn);

    const handleLogout = () => {
        localStorage.removeItem("token"); // 移除 token
        setIsLoggedIn(false); // 更新登录状态
        console.log("Logged out!");
    };

    // Function to handle login click
    const handleLoginClick = () => {
        router.push("/login"); // 跳转到登录页面
    };

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                {isLoggedIn ? (
                    <Avatar>
                        <AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
                        <AvatarFallback>CN</AvatarFallback>
                    </Avatar>
                ) : (
                    <UserRoundX className="w-8 h-8 cursor-pointer" />
                )}
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-56">
                {isLoggedIn ? (
                    <>
                        <DropdownMenuLabel>My Account</DropdownMenuLabel>
                        <DropdownMenuSeparator />
                        <DropdownMenuGroup>
                            <DropdownMenuItem>
                                <User className="mr-2 h-4 w-4" />
                                <span>Profile</span>
                                <DropdownMenuShortcut>⇧⌘P</DropdownMenuShortcut>
                            </DropdownMenuItem>
                            <DropdownMenuItem>
                                <Settings className="mr-2 h-4 w-4" />
                                <span>Settings</span>
                                <DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
                            </DropdownMenuItem>
                        </DropdownMenuGroup>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem>
                            <LifeBuoy className="mr-2 h-4 w-4" />
                            <span>Support</span>
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem onClick={handleLogout}>
                            <LogOut className="mr-2 h-4 w-4" />
                            <span>Log out</span>
                            <DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut>
                        </DropdownMenuItem>
                    </>
                ) : (
                    <>
                        <DropdownMenuLabel>Login Required</DropdownMenuLabel>
                        <DropdownMenuSeparator />
                        <DropdownMenuGroup>
                            <DropdownMenuItem onClick={handleLoginClick}>
                                <User className="mr-2 h-4 w-4" />
                                <span>Log in</span>
                                <DropdownMenuShortcut>⇧⌘L</DropdownMenuShortcut>
                            </DropdownMenuItem>
                        </DropdownMenuGroup>
                    </>
                )}
            </DropdownMenuContent>
        </DropdownMenu>
    );
}
