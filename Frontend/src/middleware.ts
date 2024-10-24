import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

// public routes
const publicPaths = ['/login', '/signup', '/forgot-password']

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl
  const authCookie = request.cookies.get('magpie_auth')

  // if it's a public route and already logged in, redirect to home page
  if (publicPaths.includes(pathname) && authCookie) {
    return NextResponse.redirect(new URL('/', request.url))
  }

  // if it's not a public route and not logged in, redirect to login page
  if (!publicPaths.includes(pathname) && !authCookie) {
    return NextResponse.redirect(new URL('/login', request.url))
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    /*
     * match all paths except:
     * /api route
     * /_next static files
     * /images static resources
     * /favicon.ico, manifest.json, etc.
     */
    '/((?!api|_next/static|_next/image|.*\\.png$).*)',
  ],
}
