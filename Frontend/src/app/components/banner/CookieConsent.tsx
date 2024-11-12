"use client"

import { useEffect, useState } from "react"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import Link from "next/link";
import { getCookiesAccepted, setCookiesAccepted } from '@/lib/cookies';
import { deleteSessionFromCookies } from '@/lib/session';

export function CookieConsent() {
  const [show, setShow] = useState(false)
  const [isPageLoaded, setIsPageLoaded] = useState(false)

  useEffect(() => {
    const handlePageLoad = () => {
      setIsPageLoaded(true)
    }

    window.addEventListener('load', handlePageLoad)

    return () => {
      window.removeEventListener('load', handlePageLoad)
    }
  }, [])

  useEffect(() => {
    const checkCookiesAndSession = async () => {
      const cookiesAccepted = await getCookiesAccepted();
      setShow(!cookiesAccepted);
    };

    if (isPageLoaded) {
      checkCookiesAndSession();
    }
  }, [isPageLoaded])

  const handleAccept = async () => {
    await setCookiesAccepted();
    setShow(false);
  }

  const handleDeny = async () => {
    try {
      // delete session info
      await deleteSessionFromCookies();
      // record user denied cookie consent
      console.log('User denied cookie consent');
      // router.push('/login'); // uncomment to redirect to login page
      // hide cookie banner
      setShow(false);
    } catch (error) {
      console.error("Error during handling deny:", error);
    }
  }

  if (!show) return null

  return (
    <Alert className="fixed bottom-8 left-8 max-w-lg mx-auto bg-white dark:bg-gray-800 shadow-lg p-3 animate-slide-in">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <AlertDescription className="text-sm">
            We use cookies to enhance your experience on our website. By continuing to use our site, you agree to our use of cookies.
          </AlertDescription>
          <div className="mt-3 flex justify-between">
            <div>
              <Link href="/terms" className="flex items-left underline mx-1">
                Terms and Privacy
              </Link>
            </div>
            <div className="flex gap-2">
              <Button size="sm" onClick={handleAccept}>
                Accept
              </Button>
              <Button size="sm" variant="outline" onClick={handleDeny}>
                Deny
              </Button>
            </div>
          </div>
        </div>
        {/* <Button
          variant="ghost"
          size="icon"
          className="h-6 w-6"
          onClick={() => setShow(false)}
        >
          <X className="h-4 w-4" />
        </Button> */}
      </div>
    </Alert>
  )
} 

