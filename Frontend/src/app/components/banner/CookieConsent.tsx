"use client"

import { useEffect, useState } from "react"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import Link from "next/link";
import { useRouter } from "next/navigation"
import { X } from "lucide-react"
import { logout } from '@/app/components/serverActions/actions';

export function CookieConsent() {
  const [show, setShow] = useState(true)
  const router = useRouter()

  useEffect(() => {
    // check if cookie exists
    const hasCookie = document.cookie.includes('magpie_auth')
    if (!hasCookie) {
      setShow(true)
    } else {
      setShow(false)
    }
  }, [])

  const handleAccept = () => {
    setShow(false)
  }

  const handleDeny = async () => {
    try {
      // call logout operation
      await logout();

      // redirect to login page
      router.push('/login')
    } catch (error) {
      console.error("Error during logout:", error)
    }
  }

  if (!show) return null

  return (
    <Alert className="fixed bottom-4 left-4 right-4 max-w-xl mx-auto bg-white dark:bg-gray-800 shadow-lg">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <AlertDescription className="text-sm">
            We use cookies to enhance your experience on our website. By continuing to use our site, you agree to our use of cookies.
            <Link href="/terms&privacy" className="underline mx-1">
              Privacy Policy and Terms of Service
            </Link>
          </AlertDescription>
          <div className="mt-3 flex gap-2">
            <Button size="sm" onClick={handleAccept}>
              Accept
            </Button>
            <Button size="sm" variant="outline" onClick={handleDeny}>
              Deny
            </Button>
          </div>
        </div>
        <Button
          variant="ghost"
          size="icon"
          className="h-6 w-6"
          onClick={() => setShow(false)}
        >
          <X className="h-4 w-4" />
        </Button>
      </div>
    </Alert>
  )
} 