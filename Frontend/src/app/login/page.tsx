import { Metadata } from "next";
import Image from "next/image";
import Link from "next/link";

import { cn } from "@/lib/utils";
import { buttonVariants } from "@/components/ui/registry/button";
import { LoginForm } from "@/app/components/form/login-form";

export const metadata: Metadata = {
  title: "Login",
  description: "Authentication login page",
};

export default function AuthenticationPage() {
  return (
    <>
      <div className="md:hidden relative">
        {/* Add relative positioning */}
        <Image
          src="/images/auth-wt.png"
          width={1280}
          height={843}
          alt="Authentication"
          className="block dark:hidden object-cover absolute inset-0" // Make image cover container and position absolutely
        />
        <Image
          src="/images/auth-bk.png"
          width={1280}
          height={843}
          alt="Authentication"
          className="hidden dark:block object-cover absolute inset-0" // Same for dark image
        />
      </div>

      <div className="container relative hidden h-[800px] flex-col items-center justify-center md:grid lg:max-w-none lg:grid-cols-2 lg:px-0">
        <Link
          href="/examples/authentication"
          className={cn(
            buttonVariants({ variant: "ghost" }),
            "absolute right-4 top-4 md:right-8 md:top-8"
          )}
        >
          Test
        </Link>

        <div className="relative hidden h-full flex-col p-10 text-white dark:border-r lg:flex">
          {/* Ensure this div is the full height */}
          <div className="absolute inset-0 bg-gradient-to-r from-zinc-900 to-zinc-700" />
          <Link href="/" className="relative inline-flex items-center gap-2">
            <div className="relative z-20 flex items-center text-lg font-medium">
              <Image
                src="/images/BKlogo.svg"
                alt="Logo"
                width={30} // Set the desired width
                height={30} // Set the desired height
                className="mr-2 h-12 w-12"
              />
              Magpie
            </div>
          </Link>
          <div className="relative z-20 mt-auto">
            <blockquote className="space-y-2">
              <p className="text-lg">
                &ldquo;This library has saved me countless hours of work and
                helped me deliver stunning designs to my clients faster than
                ever before.&rdquo;
              </p>
              <footer className="text-sm">Sofia Davis</footer>
            </blockquote>
          </div>
        </div>

        <div className="lg:p-8">
          <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
            <LoginForm />
            <p className="px-8 text-center text-sm text-muted-foreground">
              By clicking continue, you agree to our{" "}
              <Link
                href="/terms"
                className="underline underline-offset-4 hover:text-primary"
              >
                Terms of Service
              </Link>{" "}
              and{" "}
              <Link
                href="/privacy"
                className="underline underline-offset-4 hover:text-primary"
              >
                Privacy Policy
              </Link>
              .
            </p>
          </div>
        </div>
      </div>
    </>
  );
}

