import type { Metadata } from "next";
import "./globals.css";
import { UserProvider } from './contexts/UserContext'

export const metadata: Metadata = {
  title: "Magpie - Services at a glance",
  description: "An at a glance view of available public services across Dublin",
};

export default function RootLayout({
  children, 
}: {
  children: React.ReactNode 
}) {
  return (
    <html lang="zh">
      <body>
        <UserProvider>
          {children}
        </UserProvider>
      </body>
    </html>
  )
}
