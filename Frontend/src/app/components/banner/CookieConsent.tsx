"use client";

import { CookieIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useEffect, useState } from "react";
import { cn } from "@/lib/utils";

import Link from "next/link";
import { setCookiesAccepted } from "@/lib/cookies";
import {
  deleteSessionFromCookies,
  commitSessionToCookies,
} from "@/lib/session";

export default function CookieConsent({
  variant = "default",
  demo = false,
  onAcceptCallback = async () => {
    await setCookiesAccepted();
    await commitSessionToCookies();
  },
  onDeclineCallback = async () => {
    try {
      await deleteSessionFromCookies();
      console.log("User denied cookie consent");
    } catch (error) {
      console.error("Error during handling deny:", error);
    }
  },
}) {
  const [isOpen, setIsOpen] = useState(false);
  const [hide, setHide] = useState(false);

  const accept = async () => {
    setIsOpen(false);
    document.cookie =
      "cookieConsent=true; expires=Fri, 31 Dec 9999 23:59:59 GMT";
    setTimeout(() => {
      setHide(true);
    }, 700);
    await onAcceptCallback();
  };

  const decline = async () => {
    setIsOpen(false);
    setTimeout(() => {
      setHide(true);
    }, 700);
    await onDeclineCallback();
  };

  useEffect(() => {
    try {
      setIsOpen(true);
      if (document.cookie.includes("cookieConsent=true")) {
        if (!demo) {
          setIsOpen(false);
          setTimeout(() => {
            setHide(true);
          }, 700);
        }
      }
    } catch (e) {
      console.log("Error: ", e);
    }
  }, []);

  return variant != "default" ? (
    <div
      className={cn(
        "fixed z-[200] bottom-0 left-0 right-0 sm:left-12 sm:bottom-4 w-full sm:max-w-md duration-700",
        !isOpen
          ? "transition-[opacity,transform] translate-y-8 opacity-0"
          : "transition-[opacity,transform] translate-y-0 opacity-100",
        hide && "hidden"
      )}
    >
      <div className="dark:bg-card bg-background rounded-md m-3 border border-border shadow-lg">
        <div className="grid gap-2">
          <div className="border-b border-border h-14 flex items-center justify-between p-4">
            <h1 className="text-lg font-medium">We use cookies</h1>
            <CookieIcon className="h-[1.2rem] w-[1.2rem]" />
          </div>
          <div className="p-4">
            <p className="text-sm font-normal text-start">
              We use cookies to ensure you get the best experience on our
              website. For more information on how we use cookies, please see
              our cookie policy.
              <br />
              <br />
              <span className="text-xs">
                By clicking "
                <span className="font-medium opacity-80">Accept</span>
                ", you agree to our use of cookies.
              </span>
              <br />
              <Link href="/terms" className="text-xs underline">
                Learn more.
              </Link>
            </p>
          </div>
          <div className="flex gap-2 p-4 py-5 border-t border-border dark:bg-background/20">
            <Button onClick={accept} className="w-full">
              Accept
            </Button>
            <Button onClick={decline} className="w-full" variant="secondary">
              Decline
            </Button>
          </div>
        </div>
      </div>
    </div>
  ) : (
    <div
      className={cn(
        "fixed z-[200] bottom-0 left-0 right-0 sm:left-12 sm:bottom-4 w-full sm:max-w-md duration-700",
        !isOpen
          ? "transition-[opacity,transform] translate-y-8 opacity-0"
          : "transition-[opacity,transform] translate-y-0 opacity-100",
        hide && "hidden"
      )}
    >
      <div className="m-3 dark:bg-card bg-background border border-border rounded-lg">
        <div className="flex items-center justify-between p-3">
          <h1 className="text-lg font-medium">We use cookies</h1>
          <CookieIcon className="h-[1.2rem] w-[1.2rem]" />
        </div>
        <div className="p-3 -mt-2">
          <p className="text-sm text-left text-muted-foreground">
            We use cookies to ensure you get the best experience on our website.
            For more information on how we use cookies, please see our cookie
            policy.{" "}
            <Link href="/terms" className="flex items-left underline mx-1">
              Terms and Privacy
            </Link>
          </p>
        </div>

        <div className="p-3 flex items-center gap-2 mt-2 border-t">
          <Button onClick={accept} className="w-full h-9 rounded-full">
            accept
          </Button>
          <Button
            onClick={decline}
            className="w-full h-9 rounded-full"
            variant="outline"
          >
            decline
          </Button>
        </div>
      </div>
    </div>
  );
}

// export function CookieConsent() {
//   const [show, setShow] = useState(false)
//   const [isPageLoaded, setIsPageLoaded] = useState(false)

//   useEffect(() => {
//     const handlePageLoad = () => {
//       setIsPageLoaded(true)
//     }

//     window.addEventListener('load', handlePageLoad)

//     return () => {
//       window.removeEventListener('load', handlePageLoad)
//     }
//   }, [])

//   useEffect(() => {
//     const checkCookiesAndSession = async () => {
//       const cookiesAccepted = await getCookiesAccepted();
//       setShow(!cookiesAccepted);
//     };

//     if (isPageLoaded) {
//       checkCookiesAndSession();
//     }
//   }, [isPageLoaded])

// const handleAccept = async () => {
//   await setCookiesAccepted();
//   await commitSessionToCookies();
//   setShow(false);
// }

// const handleDeny = async () => {
//   try {
//     // delete session info
//     await deleteSessionFromCookies();
//     // record user denied cookie consent
//     console.log('User denied cookie consent');
//     // router.push('/login'); // uncomment to redirect to login page
//     // hide cookie banner
//     setShow(false);
//   } catch (error) {
//     console.error("Error during handling deny:", error);
//   }
// }

//   if (!show) return null

//   return (
//     <Alert className="fixed bottom-8 left-20 max-w-lg mx-auto bg-white dark:bg-gray-800 shadow-lg p-3 animate-slide-in">
//       <div className="flex items-start justify-between">
//         <div className="flex-1">
//           <AlertDescription className="text-sm">
//             We use cookies to enhance your experience on our website. By continuing to use our site, you agree to our use of cookies.
//           </AlertDescription>
//           <div className="mt-3 flex justify-between">
//             <div>
//               <Link href="/terms" className="flex items-left underline mx-1">
//                 Terms and Privacy
//               </Link>
//             </div>
//             <div className="flex gap-2">
//               <Button size="sm" onClick={handleAccept}>
//                 Accept
//               </Button>
//               <Button size="sm" variant="outline" onClick={handleDeny}>
//                 Deny
//               </Button>
//             </div>
//           </div>
//         </div>
//         {/* <Button
//           variant="ghost"
//           size="icon"
//           className="h-6 w-6"
//           onClick={() => setShow(false)}
//         >
//           <X className="h-4 w-4" />
//         </Button> */}
//       </div>
//     </Alert>
//   )
// }
