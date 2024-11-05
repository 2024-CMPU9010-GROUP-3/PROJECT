import {
    LifeBuoy,
    LogOut,
    Settings,
    User,
    UserRoundX,
  } from "lucide-react";
  
  import { useAuth } from "@/app/context/AuthContext"; // Import useAuth to get auth context
  import { logout } from "@/app/components/serverActions/actions"; // Import logout action
  import { useRouter } from "next/navigation";
  import { useEffect } from "react";
  import { getSession } from "@/lib/session";
  
  import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuShortcut,
    DropdownMenuTrigger,
  } from "@/components/ui/dropdown-menu";
  
  import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
  
  export function DropdownMenuDemo() {
    const { isLoggedIn, setIsLoggedIn } = useAuth(); // Use useAuth to get isLoggedIn state
    const router = useRouter();
  
    useEffect(() => {
      // check session
      const checkSession = async () => {
        const session = await getSession();
        if (session) {
          setIsLoggedIn(true);
        }
      };
      checkSession();
    }, [setIsLoggedIn]);
  
    // Function to handle login
    const handleLogin = () => {
      router.push("/login");
      console.log("Login clicked!");
    };
  
    // Function to handle logout
    const handleLogout = async () => {
      await logout(); // Call logout action
      console.log("Logged out!");
    };
  
    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          {/* Conditionally render Avatar or UserRoundX icon based on isLoggedIn state */}
          {isLoggedIn ? (
            <Avatar>
              <AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
              <AvatarFallback>CN</AvatarFallback>
            </Avatar>
          ) : (
            <UserRoundX className="w-8 h-8 cursor-pointer" />
          )}
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56">
          {/* Show different dropdown items based on login status */}
          {isLoggedIn ? (
            <>
              <DropdownMenuLabel>My Account</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuGroup>
                <DropdownMenuItem>
                  <User className="mr-2 h-4 w-4" />
                  <span>Profile</span>
                  <DropdownMenuShortcut>⇧⌘P</DropdownMenuShortcut>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Settings className="mr-2 h-4 w-4" />
                  <span>Settings</span>
                  <DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
                </DropdownMenuItem>
              </DropdownMenuGroup>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <LifeBuoy className="mr-2 h-4 w-4" />
                <span>Support</span>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={handleLogout}>
                <LogOut className="mr-2 h-4 w-4" />
                <span>Log out</span>
                <DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut>
              </DropdownMenuItem>
            </>
          ) : (
            <>
              <DropdownMenuLabel>Login Required</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuGroup>
                <DropdownMenuItem onClick={handleLogin}>
                  <User className="mr-2 h-4 w-4" />
                  <span>Log in</span>
                  <DropdownMenuShortcut>⇧⌘L</DropdownMenuShortcut>
                </DropdownMenuItem>
              </DropdownMenuGroup>
            </>
          )}
        </DropdownMenuContent>
      </DropdownMenu>
    );
  }
  
  