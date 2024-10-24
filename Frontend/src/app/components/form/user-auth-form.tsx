"use client";

import * as React from "react";
import { useFormState, useFormStatus } from 'react-dom';
import { cn } from "@/lib/utils";
import { Icons } from "@/components/ui/icons";
import { Button } from "@/components/ui/registry/button";
import { Input } from "@/components/ui/registry/input";
import { Label } from "@/components/ui/registry/label";
import { signup } from "@/app/actions";

interface UserAuthFormProps extends React.HTMLAttributes<HTMLDivElement> {}

function SubmitButton() {
  const { pending } = useFormStatus();
  
  return (
    <Button disabled={pending}>
      {pending && (
        <Icons.spinner className="mr-2 h-4 w-4 animate-spin" />
      )}
      Register
    </Button>
  );
}

export function UserAuthForm({ className, ...props }: UserAuthFormProps) {
  const [state, formAction] = useFormState(signup, undefined);
  const [showPassword, setShowPassword] = React.useState<boolean>(false);

  return (
    <div className={cn("grid gap-6", className)} {...props}>
      <form action={formAction}>
        <div className="grid gap-2">
          <div className="flex gap-2">
            <div className="flex-1">
              <Label className="sr-only" htmlFor="first-name">
                名字
              </Label>
              <Input
                id="first-name"
                name="firstName"
                placeholder=" "
                type="text"
                autoCapitalize="none"
                autoCorrect="off"
              />
              {state?.errors && 'firstName' in state.errors && (
                <p className="text-sm text-red-500">{state?.errors?.firstName?.[0]}</p>
              )}
            </div>
            <div className="flex-1">
              <Label className="sr-only" htmlFor="last-name">
                姓氏
              </Label>
              <Input
                id="last-name"
                name="lastName"
                placeholder=" "
                type="text"
                autoCapitalize="none"
                autoCorrect="off"
              />
              {state?.errors && 'lastName' in state.errors && (
                <p className="text-sm text-red-500">{state?.errors?.lastName?.[0]}</p>
              )}
            </div>
          </div>

          <div className="grid gap-1">
            <Label className="sr-only" htmlFor="email">
              Email
            </Label>
            <Input
              id="email"
              name="email"
              placeholder="name@example.com"
              type="email"
              autoCapitalize="none"
              autoComplete="email"
              autoCorrect="off"
            />
            {state?.errors?.email && (
              <p className="text-sm text-red-500">{state.errors.email[0]}</p>
            )}
          </div>

          <div className="grid gap-1">
            <Label className="sr-only" htmlFor="password">
              密码
            </Label>
            <div className="relative">
              <Input
                id="password"
                name="password"
                placeholder="Enter password"
                type={showPassword ? "text" : "password"}
                autoCapitalize="none"
                autoComplete="new-password"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-2 top-1/2 transform -translate-y-1/2"
              >
                {showPassword ? "Hide" : "Show"}
              </button>
            </div>
            {state?.errors && 'password' in state.errors && (
              <div className="text-sm text-red-500">
                <p>Password must:</p>
                <ul className="list-disc pl-4">
                  {state?.errors?.password?.map((error) => (
                    <li key={error}>{error}</li>
                  ))}
                </ul>
              </div>
            )}
          </div>

          {state?.message && (
            <p className="text-sm text-red-500">{state.message}</p>
          )}

          <SubmitButton />
        </div>
      </form>
    </div>
  );
}
