import { Menu } from "lucide-react";
import Image from "next/image";

export const Navbar = () => (
  <nav className="fixed top-0 left-0 right-0 bg-white shadow-sm z-50">
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div className="flex justify-between h-16 items-center">
        <div className="flex items-center">
          <Image
            src="/api/placeholder/40/40"
            alt="Magpie Logo"
            className="h-10 w-10"
            width={40}
            height={40}
          />
          <span className="ml-2 text-xl font-bold">Magpie</span>
        </div>
        <div className="hidden md:flex items-center space-x-8">
          <a href="#features" className="text-gray-600 hover:text-gray-900">
            Features
          </a>
          <a href="#use-cases" className="text-gray-600 hover:text-gray-900">
            Use Cases
          </a>
          <a href="#why-magpie" className="text-gray-600 hover:text-gray-900">
            Why Magpie
          </a>
          <a href="#get-started" className="text-gray-600 hover:text-gray-900">
            Get Started
          </a>
          <button className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700">
            Sign Up
          </button>
        </div>
        <div className="md:hidden">
          <Menu className="h-6 w-6" />
        </div>
      </div>
    </div>
  </nav>
);
