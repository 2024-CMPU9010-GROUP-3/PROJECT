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
  const [username, setUsername] = useState(""); // username status
  const [password, setPassword] = useState(""); // password status
  const [errorMessage, setErrorMessage] = useState<string | null>(null); // error messesage
  const router = useRouter(); // router
  const [isLoading, setIsLoading] = useState(false); //  loading state

  const onSubmit = async (event: React.SyntheticEvent) => {
    event.preventDefault(); // avoid defauld form submission
    setIsLoading(true); // set loading
    setErrorMessage(null); //   clear previous error message

    //  validate fields
    if (!username.trim() || !password.trim()) {
      setErrorMessage("Username and password are required");
      setIsLoading(false);
      return;
    }

    const formData = new FormData();
    formData.append('username', username);
    formData.append('password', password);

    try {
      const result = await login(formData); //  login

      if (result.errors) {
        setErrorMessage(result.errors.username?.[0] || result.errors.password?.[0] || null);
      } else {
        router.push("/"); //  redirect to home
      }
    } catch (error) {
      console.error("An error occurred during login:", error);
      setErrorMessage("Login failed, please try again later.");
    } finally {
      setIsLoading(false); //   set loading to false
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
            <Button type="submit" disabled={isLoading}>
              {isLoading ? "Logging in..." : "Login"}
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
