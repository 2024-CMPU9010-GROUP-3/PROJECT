import type { Metadata } from "next";
import Providers from "./providers";
import "./globals.css";
import { SessionProvider } from './context/SessionContext';
import { Onborda, OnbordaProvider } from "onborda";
import { steps } from "./components/onboarding/steps";
import { TourCard } from "./components/onboarding/card";


export const metadata: Metadata = {
  title: "Magpie",
  description: "An at a glance view of available public services across Dublin",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="antialiased min-h-screen">
        <SessionProvider>
          <Providers>
            <OnbordaProvider>
              <Onborda
                steps={steps}
                showOnborda={true}
                shadowRgb="0,0,0"
                shadowOpacity="0.8"
                cardComponent={TourCard}
                cardTransition={{ duration: 1, type: "spring" }}
              >
                {children}
              </Onborda>
            </OnbordaProvider>
          </Providers>
        </SessionProvider>
      </body>
    </html>
  );
}
