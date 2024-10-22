"use client";

import Link from "next/link";
import { useState } from "react"; // useState
import { Button } from "@/components/ui/registry/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/registry/card";
import { Input } from "@/components/ui/registry/input";
import { Label } from "@/components/ui/registry/label";
import { useRouter } from "next/navigation"; // useRouter
import { login } from "@/app/actions"; // 导入 login Server Action

export function LoginForm() {
  const [username, setUsername] = useState(""); // 用户名状态
  const [password, setPassword] = useState(""); // 密码状态
  const [errorMessage, setErrorMessage] = useState<string | null>(null); // 错误信息
  const router = useRouter(); // router

  const onSubmit = async (event: React.SyntheticEvent) => {
    event.preventDefault(); // 防止默认表单提交行为

    const formData = new FormData();
    formData.append('username', username);
    formData.append('password', password);

    const result = await login(formData); // 调用 Server Action

    if (result.errors) {
      setErrorMessage(result.errors.username?.[0] || result.errors.password?.[0] || null);
    } else {
      router.push("/"); // 登录成功后重定向
    }
  };

  return (
    <Card className="mx-auto max-w-sm">
      <CardHeader>
        <CardTitle className="text-2xl">Login</CardTitle>
        <CardDescription>
          Enter your username or email below to login to your account
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit}>
          <div className="grid gap-4">
            <div className="grid gap-2">
              <Label htmlFor="username">Username</Label>
              <Input
                id="username"
                type="text"
                required
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
            </div>
            <div className="grid gap-2">
              <div className="flex items-center">
                <Label htmlFor="password">Password</Label>
                <Link href="#" className="ml-auto inline-block text-sm underline">
                  Forgot your password?
                </Link>
              </div>
              <Input
                id="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            <Button type="submit">
              Login
            </Button>
            {errorMessage && <div className="text-red-500">{errorMessage}</div>}
          </div>
        </form>
        <div className="mt-4 text-center text-sm">
          Don&apos;t have an account?{" "}
          <Link href="./signup" className="underline">
            Sign up
          </Link>
        </div>
      </CardContent>
    </Card>
  );
}
