import Image from "next/image";
import Link from "next/link";

import { LoginForm } from "@/app/components/form/login-form";


export default function LoginPage() {


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
        {/* Center LoginForm */}
        <div className="relative z-10 p-4 w-full max-w-sm bg-white">
          <LoginForm />
          <p className="text-center text-sm text-muted-foreground mt-4">
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

      {/* Desktop and Tablet Version */}
      <div className="container relative hidden md:grid h-screen flex-col items-center justify-center lg:max-w-none lg:grid-cols-2 lg:px-0">
        <div className="relative hidden h-full flex-col p-10 text-white dark:border-r lg:flex">
          {/* Desktop Side Section */}
            <Image className="object-cover z-0" src="/images/map_screenshot.png" alt="Map Screenshot" fill></Image>
          <div className="absolute inset-0 bg-gradient-to-r from-zinc-900 to-zinc-900/70 z-10"/>
          <Link href="/" className="relative inline-flex items-center gap-2 z-20">
            <Image src="/images/logo-wt.svg" alt="Logo" width={100} height={100} className="mr-2 h-32 w-32" />
            <div>
              <div className="text-6xl font-bold">Magpie</div>
              <div className="text-2xl font-medium">Services at a Glance</div>
            </div>
          </Link>
        </div>

        <div className="lg:p-8 flex items-center justify-center">
          <div className="w-full max-w-sm flex flex-col justify-center space-y-6 sm:w-[400px]">
            <LoginForm />
            <p className="px-8 text-center text-sm text-muted-foreground">
              By using Magpie, you agree to our{" "}
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

