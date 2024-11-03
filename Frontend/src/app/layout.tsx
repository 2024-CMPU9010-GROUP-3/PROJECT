import type { Metadata } from "next";
import Providers from "./providers";
import "./globals.css";
import { AuthProvider } from '@/app/context/AuthContext';


export const metadata: Metadata = {
  title: "Magpie - Services at a glance",
  description: "An at a glance view of available public services across Dublin",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`antialiased min-h-screen`}>
        <AuthProvider>
          <Providers>{children}</Providers>
        </AuthProvider>
      </body>
    </html>
  );
}
