import { Metadata } from "next";
import Image from "next/image";
import Link from "next/link";

import { SignupForm } from "@/app/components/form/signup-form";

export const metadata: Metadata = {
  title: "Sign Up",
  description: "Authentication sign up page",
};

export default function SignupPage() {
  return (
    <>
      <div className="md:hidden relative min-h-screen flex flex-col items-center justify-center">
        {/* Mobile version */}
        <div className="flex gap-2 mb-2">
          <Image
            src="/images/logo-circle-bk.svg"
            height={50}
            width={50}
            alt="Magpie Logo"
            className="block dark:hidden object-contain inset-0" // Mobile Light Image
          />
          <Image
            src="/images/logo-circle-wt.png"
            height={50}
            width={50}
            alt="Magpie Logo"
            className="hidden dark:block object-contain inset-0" // Mobile Dark Image
          />
          <div>
              <div className="text-2xl font-bold">Magpie</div>
              <div className="text-m font-medium">Services at a Glance</div>
            </div>
        </div>
        {/* Center user sign up form */}
        <div className="lg:p-8">
          <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
            <SignupForm />
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

      {/* Desktop and Tablet Version */}
      <div className="container relative hidden md:grid h-screen flex-col items-center justify-center lg:max-w-none lg:grid-cols-2 lg:px-0">

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

        <div className="lg:p-8 flex items-center justify-center">
          <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[400px]">
            <SignupForm />
            <p className="px-8 text-center text-sm text-muted-foreground">
              By creating an account, you agree to our{" "}
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
