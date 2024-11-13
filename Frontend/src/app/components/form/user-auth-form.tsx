"use client";

import * as React from "react";

import { cn } from "@/lib/utils";
import { Icons } from "@/components/ui/icons";
import { Button } from "@/components/ui/registry/button";
import { Input } from "@/components/ui/registry/input";
import { Label } from "@/components/ui/registry/label";
import { useRouter } from "next/navigation";
import { signup } from "@/app/components/serverActions/actions"; // Import the signup action

interface UserAuthFormProps extends React.HTMLAttributes<HTMLDivElement> {
  username?: string;
  password?: string;
  confirmPassword?: string;
  email?: string;
}

export function UserAuthForm({
  className,
  ...props
}: Readonly<UserAuthFormProps>) {
  const [isLoading, setIsLoading] = React.useState<boolean>(false);
  const [email, setEmail] = React.useState<string>("");
  const [password, setPassword] = React.useState<string>("");
  const [confirmPassword, setConfirmPassword] = React.useState<string>("");
  const [firstName, setFirstName] = React.useState<string>("");
  const [lastName, setLastName] = React.useState<string>("");
  const [errorMessage, setErrorMessage] = React.useState<string | null>(null);
  const [formErrors, setFormErrors] = React.useState<{
    firstName?: string[];
    lastName?: string[];
    email?: string[];
    password?: string[];
  }>({});
  const [showPassword, setShowPassword] = React.useState<boolean>(false);
  const [showConfirmPassword, setShowConfirmPassword] =
    React.useState<boolean>(false);
  const router = useRouter();

  async function onSubmit(event: React.SyntheticEvent) {
    event.preventDefault();
    setIsLoading(true);
    setErrorMessage(null);
    setFormErrors({});

    if (
      !email.trim() ||
      !password.trim() ||
      !confirmPassword.trim() ||
      !firstName.trim() ||
      !lastName.trim()
    ) {
      setErrorMessage("Fields cannot be empty");
      setIsLoading(false);
      return;
    }

    // Check if password and confirm password match
    if (password !== confirmPassword) {
      console.error("Passwords do not match");
      setErrorMessage("Passwords do not match");
      setIsLoading(false);
      return;
    }

    try {
      const formData = new FormData();
      formData.append("email", email);
      formData.append("password", password);
      formData.append("firstName", firstName);
      formData.append("lastName", lastName);

      const result = await signup({}, formData); // Call the signup action

      if (result.errors) {
        console.error("Registration failed:", result.errors);
        setFormErrors(result.errors);
      } else {
        alert("Sign up successful in user form");
        router.push("/login");
      }
    } catch (error) {
      console.error("An error occurred", error);
      alert("An error occurred, please try again later.");
      setErrorMessage("An error occurred, please try again later.");
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className={cn("grid gap-6", className)} {...props}>
      <form onSubmit={onSubmit}>
        <div className="grid gap-2">
          <div className="flex gap-2">
            <div className="flex-1">
              <Label className="sr-only" htmlFor="first-name">
                First Name
              </Label>
              <Input
                id="first-name"
                placeholder="First Name"
                type="text"
                disabled={isLoading}
                value={firstName}
                onChange={(e) => setFirstName(e.target.value)}
              />
            </div>
            <div className="flex-1">
              <Label className="sr-only" htmlFor="last-name">
                Last Name
              </Label>
              <Input
                id="last-name"
                placeholder="Last Name"
                type="text"
                disabled={isLoading}
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
              />
            </div>
          </div>
          <div className="grid gap-1">
            <Label className="sr-only" htmlFor="email">
              Email
            </Label>
            <Input
              id="email"
              placeholder="name@example.com"
              type="email"
              autoCapitalize="none"
              autoComplete="email"
              autoCorrect="off"
              disabled={isLoading}
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>
          <div className="grid gap-1">
            <Label className="sr-only" htmlFor="password">
              Password
            </Label>
            <div className="relative">
              <Input
                id="password"
                placeholder="Enter your password"
                type={showPassword ? "text" : "password"}
                autoCapitalize="none"
                autoComplete="current-password"
                disabled={isLoading}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-2 top-1/2 transform -translate-y-1/2"
              >
                {showPassword ? "Hide" : "Show"}
              </button>
            </div>
          </div>
          <div className="grid gap-1">
            <Label className="sr-only" htmlFor="confirm-password">
              Confirm Password
            </Label>
            <div className="relative">
              <Input
                id="confirm-password"
                placeholder="Confirm your password"
                type={showConfirmPassword ? "text" : "password"}
                autoCapitalize="none"
                autoComplete="current-password"
                disabled={isLoading}
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
              />
              <button
                type="button"
                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                className="absolute right-2 top-1/2 transform -translate-y-1/2"
              >
                {showConfirmPassword ? "Hide" : "Show"}
              </button>
            </div>
          </div>
          {errorMessage && <div className="text-red-500">{errorMessage}</div>}
          {formErrors.firstName && (
            <ul className="text-red-500 text-sm mt-1">
              {formErrors.firstName.map((error, index) => (
                <li key={index}>{error}</li>
              ))}
            </ul>
          )}
          {formErrors.lastName && (
            <ul className="text-red-500 text-sm mt-1">
              {formErrors.lastName.map((error, index) => (
                <li key={index}>{error}</li>
              ))}
            </ul>
          )}
          {formErrors.email && (
            <ul className="text-red-500 text-sm mt-1">
              {formErrors.email.map((error, index) => (
                <li key={index}>{error}</li>
              ))}
            </ul>
          )}
          {formErrors.password && (
            <ul className="text-red-500 text-sm mt-1">
              {formErrors.password.map((error, index) => (
                <li key={index}>{error}</li>
              ))}
            </ul>
          )}

          <Button disabled={isLoading}>
            {isLoading && (
              <Icons.spinner className="mr-2 h-4 w-4 animate-spin" />
            )}
            Sign Up
          </Button>
        </div>
      </form>
    </div>
  );
}
