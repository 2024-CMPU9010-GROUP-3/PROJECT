"use client"

import * as React from "react"

import { cn } from "@/lib/utils"
import { Icons } from "@/components/ui/icons"
import { Button } from "@/components/ui/registry/button"
import { Input } from "@/components/ui/registry/input"
import { Label } from "@/components/ui/registry/label"
import { useRouter } from 'next/compat/router';

interface UserAuthFormProps extends React.HTMLAttributes<HTMLDivElement> {
  username?: string;
  password?: string;
  confirmPassword?: string;
  email?: string; 
}

export function UserAuthForm({ className, ...props }: UserAuthFormProps) {
  const [isLoading, setIsLoading] = React.useState<boolean>(false)
  const [email, setEmail] = React.useState<string>("")
  const [password, setPassword] = React.useState<string>("")
  const [confirmPassword, setConfirmPassword] = React.useState<string>("");
  const [firstName, setFirstName] = React.useState<string>("");
  const [lastName, setLastName] = React.useState<string>("");
  const [errorMessage, setErrorMessage] = React.useState<string | null>(null);
  const router = useRouter()

  async function onSubmit(event: React.SyntheticEvent) {
    event.preventDefault()
    setIsLoading(true)
    setErrorMessage(null); // clear error previous message 

    if (!email.trim() || !password.trim() || !confirmPassword.trim() || !firstName.trim() || !lastName.trim()) {
      console.error('Fields cannot be empty');
      setErrorMessage('Fields cannot be empty'); 
      setIsLoading(false);
      return;
    }

    // passwordComplexity check
    const passwordComplexity = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;
    if (!passwordComplexity.test(password)) {
      console.error('Password does not meet complexity requirements');
      setErrorMessage('Password does not meet complexity requirements');
      setIsLoading(false);
      return;
    }

    // Check if password and confirm password match
    if (password !== confirmPassword) {
      console.error('Passwords do not match');
      setErrorMessage('Passwords do not match');
      setIsLoading(false);
      return;
    }

    console.log('Backend URL:', process.env.NEXT_PUBLIC_BACKEND_URL); // debug        

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/auth/User/`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
          username: email.split('@')[0],
          email, 
          password, 
          firstName,
          lastName
        }),
      })

      const responseText = await response.text(); // aquire response text
      console.log('Response:', responseText); // display response

      if (response.ok) {
        // Handle successful registration
        alert('Sign up successful'); //  display success message
        if (router) {
          router.push('/')
        } else {
          console.error('Router not found')
        }
      } else {
        // Handle errors
        const errorData = await response.json();
        console.error('Registration failed:', errorData);
        setErrorMessage('Registration failed: ' + errorData.message); // display error message
        alert('Sign up failed: ' + errorData.message); // alert user that sign up failed
      }
    } catch (error) {
      console.error('An error occurred', error)
      alert('An error occurred, please try again later.'); // alert user that an error occurred
    } finally {
      setIsLoading(false)
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
            <Input
              id="password"
              placeholder="Enter your password"
              type="password"
              autoCapitalize="none"
              autoComplete="current-password"
              disabled={isLoading}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>
          <div className="grid gap-1"> {/* Confirm Password input */}
            <Label className="sr-only" htmlFor="confirm-password">
              Confirm Password
            </Label>
            <Input
              id="confirm-password"
              placeholder="Confirm your password"
              type="password"
              autoCapitalize="none"
              autoComplete="current-password"
              disabled={isLoading}
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
            />
          </div>
          {errorMessage && ( // error message
            <div className="text-red-500">{errorMessage}</div>
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
  )
}