"use client";

import {Button} from "@/components/ui/button";
import {Tooltip, TooltipContent, TooltipProvider, TooltipTrigger} from "@/components/ui/tooltip";
import {BookOpenText, CircleHelp, LogOut} from "lucide-react";
import {useRouter} from "next/navigation";
import {useSession} from "../context/SessionContext";
import {useOnborda} from "onborda";

const Nav = () => {
  const { setSessionToken, setSessionUUID, sessionUUID } = useSession();
  const { startOnborda } = useOnborda();
  const router = useRouter();
  const handleLogout = async () => {
    setSessionToken("");
    setSessionUUID("");
    router.push("/login");
  };

  return (
      // <ProtectedRoute>
        <div className="flex flex-row gap-2">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button className="w-10 h-10 rounded-full bg-white hover:bg-gray-100 p-0" onClick={handleLogout}>
                  <LogOut color="black" className="w-5 h-5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>Logout</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button className="w-10 h-10 rounded-full bg-white hover:bg-gray-100 p-0" onClick={() => {router.push(`/history?userid=${sessionUUID}`)}}>
                  <BookOpenText color="black" className="w-5 h-5"/>
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>Saved Locations</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button className="w-10 h-10 rounded-full bg-white hover:bg-gray-100 p-0" id="onboarding-step-3" onClick={() => startOnborda("general-onboarding")}>
                  <CircleHelp color="black" className="w-5 h-5"/>
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>How do I use Magpie?</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
  );
};

export default Nav;
