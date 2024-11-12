"use client"

import { useEffect, useState } from "react"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import Link from "next/link";
import { useRouter } from "next/navigation"
import { X } from "lucide-react"
import { logout } from '@/app/components/serverActions/actions';
import './animate.css'

export function CookieConsent() {
  const [show, setShow] = useState(false)
  const [isPageLoaded, setIsPageLoaded] = useState(false)
  const router = useRouter()

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
    if (isPageLoaded) {
      // get cookie
      const cookieString = document.cookie;
      const cookies = cookieString.split('; ').reduce((acc, cookie) => {
        const [name, value] = cookie.split('=');
        acc[name] = value;
        return acc;
      }, {} as Record<string, string>);

      // check if cookie exists or expired
      const hasCookie = cookies['magpie_auth'];
      const expiryTime = cookies['magpie_auth_expiry']; 

      if (hasCookie && expiryTime) {
        const expiryDate = new Date(parseInt(expiryTime, 10));
        const now = new Date();
        if (now < expiryDate) {
          setShow(false);
        } else {
          setShow(true);
        }
      } else {
        setShow(true);
      }
    }
  }, [isPageLoaded])

  const handleAccept = () => {
    setShow(false)
  }

  const handleDeny = async () => {
    try {
      // call logout operation
      // await logout();
      console.log('user denied cookie consent')
      // redirect to login page
      // router.push('/login')
      setShow(false)
    } catch (error) {
      console.error("Error during logout:", error)
    }
  }

  if (!show) return null

  return (
    <Alert className="fixed bottom-8 left-10 max-w-lg mx-auto bg-white dark:bg-gray-800 shadow-lg p-3 animate-slide-in">
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

