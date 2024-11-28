"use client";

import { useSession } from "@/app/context/SessionContext";
import { Menu } from "lucide-react";
import Image from "next/image";
import Link from "next/link";

export const Navbar = () => {
  const { sessionToken } = useSession();

  return (
    <nav className="fixed top-0 left-0 right-0 bg-white shadow-sm z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16 items-center">
          <div className="flex items-center">
            <Link href={sessionToken ? "/home" : "/"}>
              <Image
                src="/images/BKlogo.svg"
                alt="Magpie Logo"
                className="h-16 w-16"
                width={500}
                height={500}
              />
            </Link>
          </div>
          <div className="hidden md:flex items-center space-x-8">
            <Link
              href="/landing#about"
              className="text-gray-600 hover:text-gray-900"
            >
              About
            </Link>
            <Link
              href="/landing#features"
              className="text-gray-600 hover:text-gray-900"
            >
              Features
            </Link>
            <Link
              href="/landing#use-cases"
              className="text-gray-600 hover:text-gray-900"
            >
              Use Cases
            </Link>
            <Link
              href="/landing#get-started"
              className="text-gray-600 hover:text-gray-900"
            >
              Get Started
            </Link>
            <Link
              className="bg-black text-white px-4 py-2 rounded-md hover:bg-gray-800"
              href={"/signup"}
            >
              Sign Up
            </Link>
          </div>
          <div className="md:hidden">
            <Menu className="h-6 w-6" />
          </div>
        </div>
      </div>
    </nav>
  );
};
