import { Metadata } from "next";
import Image from "next/image";
import Link from "next/link";

import { cn } from "@/lib/utils";
import { buttonVariants } from "@/components/ui/registry/button";
import { UserAuthForm } from "@/app/components/form/user-auth-form";
import MagpieLogo from "../components/logo/magpie";

export const metadata: Metadata = {
  title: "Sign Up",
  description: "Authentication sign up page",
};

export default function SignupPage() {
  return (
    <>
      <div className="md:hidden relative min-h-screen flex flex-col items-center justify-center">
        {/* Mobile version */}
        <Image
          src="/images/auth-wt.png"
          width={1280}
          height={843}
          alt="Authentication"
          className="block dark:hidden object-cover absolute inset-0 w-full h-full" // Mobile Light Image
        />
        <Image
          src="/images/auth-bk.png"
          width={1280}
          height={843}
          alt="Authentication"
          className="hidden dark:block object-cover absolute inset-0 w-full h-full" // Mobile Dark Image
        />
        {/* Center user sign up form */}
        <div className="lg:p-8">
          <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
            <div className="flex flex-col space-y-2 text-center">
              <h1 className="text-2xl font-semibold tracking-tight">
                Create an account
              </h1>
              <p className="text-sm text-muted-foreground">
                Enter your email below to create your account
              </p>
            </div>
            <UserAuthForm />
            <p className="px-8 text-center text-sm text-muted-foreground">
              By clicking continue, you agree to our{" "}
              <Link
                href="/terms"
                className="underline underline-offset-4 hover:text-primary"
              >
                Terms and Privacy Policy
              </Link>
              .
            </p>
          </div>
        </div>
      </div>

      {/* Desktop and Tablet Version */}
      <div className="container relative hidden md:grid h-screen flex-col items-center justify-center lg:max-w-none lg:grid-cols-2 lg:px-0">
        <Link
          href="/login"
          className={cn(
            buttonVariants({ variant: "ghost" }),
            "absolute right-4 top-4 md:right-8 md:top-8"
          )}
        >
          Login
        </Link>

        <div className="relative hidden h-full flex-col p-10 text-white dark:border-r lg:flex">
          {/* Desktop Side Section */}
          <div className="absolute inset-0 bg-gradient-to-r from-zinc-700 to-zinc-950" />
          <div className="absolute top-8 left-6" style={{ zIndex: 99 }}>
                <MagpieLogo />
              </div>
          <div className="relative z-20 mt-auto">
            <blockquote className="space-y-2">
              <p className="text-lg">
                &ldquo;Where data meets design, and every street tells a story—shaping cities through the lens of innovation, one space at a time.&rdquo;
              </p>
              <footer className="text-sm">Magpie</footer>
            </blockquote>
          </div>
        </div>

        <div className="lg:p-8">
          <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
            <div className="flex flex-col space-y-2 text-center">
              <h1 className="text-2xl font-semibold tracking-tight">
                Create an account
              </h1>
              <p className="text-sm text-muted-foreground">
                Enter your email below to create your account
              </p>
            </div>
            <UserAuthForm />
            <p className="px-8 text-center text-sm text-muted-foreground">
              By clicking continue, you agree to our{" "}
              <Link
                href="/terms"
                className="underline underline-offset-4 hover:text-primary"
              >
                Terms and Privacy Policy
              </Link>
              .
            </p>
          </div>
        </div>
      </div>
    </>
  );
}
