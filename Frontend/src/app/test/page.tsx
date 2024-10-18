"use client";

import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import { DropdownMenuDemo } from "@/app/components/avatarmenu/avatarmenu";
import { AuthProvider, useAuth } from "@/app/components/AuthContext";

const TestPage = () => {
    return (
        <AuthProvider>
            <Content /> {/* 将内容提取到一个子组件中 */}
        </AuthProvider>
    );
};

const Content = () => {
    const { isLoggedIn, setIsLoggedIn } = useAuth(); // 使用 AuthContext

    const handleAvatarClick = () => {
        setIsLoggedIn(true); // 点击头像时设置为已登录
    };

    const avatar = isLoggedIn ? (
        <Avatar onClick={handleAvatarClick}> {/* 添加点击事件 */}
            <AvatarImage
                src="https://github.com/shadcn.png"
                alt="@shadcn"
            />
            <AvatarFallback>CN</AvatarFallback>
        </Avatar>
    ) : (
        <Avatar onClick={handleAvatarClick}> {/* 添加点击事件 */}
            <AvatarFallback>DF</AvatarFallback>
        </Avatar>
    );

    return (
        <div>All the about us 
            <div>
                <DropdownMenuDemo /> 
            </div>
            content goes in this
        </div>
    );
};

export default TestPage;
