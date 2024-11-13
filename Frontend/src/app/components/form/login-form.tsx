"use client";

import Link from "next/link";
import { useState, useEffect } from "react"; // useState and useEffect
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
import { getToken, setSession } from "@/lib/session";
import { getCookiesAccepted } from "@/lib/cookies";
import { useAuth } from "@/app/context/AuthContext";

interface LoginData {
  token: string;
  userid: string;
}

export function LoginForm() {
  const [usernameOrEmail, setUsernameOrEmail] = useState(""); // allow login with username or email
  const [password, setPassword] = useState(""); // password state
  const [errorMessage, setErrorMessage] = useState<string | null>(null); // error message
  const [isLoading, setIsLoading] = useState(false); // loading state
  const [sessionData, setSessionData] = useState<LoginData | null>(null); // session data
  const router = useRouter(); // router

  const authContext = useAuth();

  async function handleSession(sessionData: { token: string; userid: string }) {
    // Save session data to local storage

    setSessionData(sessionData);
    // Set session in the application (e.g., cookies, context, etc.)
    if (await getCookiesAccepted()) {
      await setSession(sessionData.token, sessionData.userid);
    } else {
      authContext.setSessionData({
        token: sessionData.token,
        userid: sessionData.userid,
      });
    }
  }

  // check if user is already logged in
  useEffect(() => {
    (async () => {
      const token = await getToken(); // check user login status
      if (token) {
        // if user is already logged in, redirect to home
        router.push("/"); // redirect to home
      }
    })();
  }, [router]);

  const onSubmit = async (event: React.SyntheticEvent) => {
    event.preventDefault();
    setIsLoading(true);

    if (!usernameOrEmail.trim() || !password.trim()) {
      setErrorMessage("Fields cannot be empty");
      setIsLoading(false);
      return;
    }

    try {
      const response = await fetch("/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ usernameOrEmail: usernameOrEmail, password }),
      });

      if (response.ok) {
        const data = await response.json();
        console.log("Login successful:", data);

        if (data.response.content.userid) {
          const sessionData = {
            token: data.response.content.token,
            userid: data.response.content.userid,
          };

          await handleSession(sessionData);
          setErrorMessage(null);
          router.push("/");
        } else {
          setErrorMessage("Login failed: No user id received");
        }
      } else {
        const errorData = await response.text();
        setErrorMessage("Login failed: " + errorData);
      }
    } catch (error) {
      console.error("An error occurred", error);
      setErrorMessage("An error occurred during login");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (typeof window !== "undefined") {
      localStorage.setItem(
        "session",
        JSON.stringify({
          token: sessionData?.token,
          userId: sessionData?.userid,
        })
      );
    }
  }, [sessionData]);

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
          {" "}
          {/* form submission handling */}
          <div className="grid gap-4">
            <div className="grid gap-2">
              <Label htmlFor="usernameOrEmail">Username/Email</Label>
              <Input
                id="usernameOrEmail"
                type="text"
                placeholder="username or m@example.com"
                required
                value={usernameOrEmail} // bind usernameOrEmail state
                onChange={(e) => setUsernameOrEmail(e.target.value)} // update state
              />
            </div>
            <div className="grid gap-2">
              <div className="flex items-center">
                <Label htmlFor="password">Password</Label>
                <Link
                  href="#"
                  className="ml-auto inline-block text-sm underline"
                >
                  Forgot your password?
                </Link>
              </div>
              <Input
                id="password"
                type="password"
                required
                value={password} // bind password state
                onChange={(e) => setPassword(e.target.value)} // update password state
              />
            </div>
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Logging in..." : "Login"}{" "}
              {/* display loading state */}
            </Button>
            {errorMessage && <div className="text-red-500">{errorMessage}</div>}{" "}
            {/* display error message */}
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
