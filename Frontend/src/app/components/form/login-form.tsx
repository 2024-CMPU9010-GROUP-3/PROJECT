"use client";

import Link from "next/link";
import { useState, useEffect, Suspense } from "react"; // useState and useEffect
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
import { useRouter, useSearchParams } from "next/navigation"; // useRouter
import { useSession } from "@/app/context/SessionContext";

export function LoginForm() {
  const [usernameOrEmail, setUsernameOrEmail] = useState(""); // allow login with username or email
  const [password, setPassword] = useState(""); // password state
  const [errorMessage, setErrorMessage] = useState<string | null>(null); // error message
  const [isLoading, setIsLoading] = useState(false); // loading state
  const router = useRouter(); // router
  const {
    sessionToken,
    setSessionToken,
    sessionUUID,
    setSessionUUID,
    setIsUserLoggedIn,
  } = useSession();

  // check if user is already logged in
  useEffect(() => {
    (async () => {
      if (sessionToken) {
        // if user is already logged in, redirect to home
        setTimeout(() => {
          router.push("/"); // redirect to home
        }, 0);
      }
    })();
  }, [sessionToken, router]);

  const onSubmit = async (event: React.SyntheticEvent) => {
    event.preventDefault(); // prevent default form submission behavior
    setIsLoading(true); // set loading state

    // check if fields are empty
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
        body: JSON.stringify({ usernameOrEmail: usernameOrEmail, password }), // send username/email and password
      });

      const data = await response.json();

      if (response.ok) {
        // login success, handle logic
        if (data.response.content) {
          const content = data.response.content;
          setSessionToken(content.token);
          setSessionUUID(content.userid);

          if (sessionToken && sessionUUID) {
            setIsUserLoggedIn(true);
          }

          setErrorMessage(null); // clear any error message

          setTimeout(() => {
            router.push("/"); // redirect to home
          }, 0);
        } else {
          setErrorMessage("Login failed: No user id received"); // if no user id, display error message
        }
      } else {
        // handle error case
        setErrorMessage(data.error.errorMsg); // display original error message
      }
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (error) {
      setErrorMessage("An unknown error occurred during login");
    } finally {
      setIsLoading(false); // reset loading state
    }
  };

  return (
    <Card>
      <Suspense>
        <CardHeaderWithSuccess />
      </Suspense>
      <CardContent>
        <form onSubmit={onSubmit} className="mt-4">
          {" "}
          {/* form submission handling */}
          <div className="grid gap-4">
            <div className="grid gap-2">
              <Label htmlFor="usernameOrEmail">Username/Email</Label>
              <Input
                id="usernameOrEmail"
                type="text"
                placeholder="Username/Email*"
                required
                value={usernameOrEmail} // bind usernameOrEmail state
                onChange={(e) => setUsernameOrEmail(e.target.value)} // update state
              />
            </div>
            <div className="grid gap-2">
              <div className="flex items-center">
                <Label htmlFor="password">Password</Label>
                {/* <Link
                  href="#"
                  className="ml-auto inline-block text-sm underline"
                >
                  Forgot your password?
                </Link> NOT YET IMPLEMENTED*/}
              </div>
              <Input
                id="password"
                type="password"
                placeholder="Password*"
                required
                value={password} // bind password state
                onChange={(e) => setPassword(e.target.value)} // update password state
              />
            </div>
            <div className="text-red-500 w-full text-center">
              {errorMessage || "\u00A0"}
            </div>
            {/* display error message */}
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Logging in..." : "Login"}{" "}
              {/* display loading state */}
            </Button>
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

function CardHeaderWithSuccess() {
  const searchParams = useSearchParams();
  const isSignupSuccess = searchParams.get("signup") === "success";
  return (
    <CardHeader>
      {!isSignupSuccess && (
        <CardTitle className="text-2xl">Welcome to Magpie</CardTitle>
      )}
      {isSignupSuccess && (
        <CardTitle className="text-2xl">Signup successful</CardTitle>
      )}
      {!isSignupSuccess && (
        <CardDescription>
          Please log in using your username or email
        </CardDescription>
      )}
      {isSignupSuccess && (
        <CardDescription>
          You can now use your username or email to log in!
        </CardDescription>
      )}
    </CardHeader>
  );
}
