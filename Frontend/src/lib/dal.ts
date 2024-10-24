'use server'

import { cache } from 'react'
import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'
import * as jwt_decode from 'jwt-decode'

interface JWTPayload {
  sub: string  // user ID
  exp: number
  iat: number
}

export const verifySession = cache(async () => {
  const cookieStore = cookies()
  const authCookie = cookieStore.get('magpie_auth')

  if (!authCookie?.value) {
    return null
  }

  try {
    const decoded = jwt_decode<JWTPayload>(authCookie.value)
    
    // check if the token is expired
    if (decoded.exp * 1000 < Date.now()) {
      return null
    }

    return {
      isAuth: true,
      userId: decoded.sub
    }
  } catch (error) {
    console.error('Token parse error:', error)
    return null
  }
})

export const getUser = cache(async () => {
  const session = await verifySession()
  if (!session) return null

  try {
    const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/auth/User/${session.userId}`, {
      credentials: 'include',
    })

    if (!res.ok) {
      throw new Error('fetch user info failed')
    }

    const data = await res.json()
    return data.response.content
  } catch (error) {
    console.error('fetch user info failed:', error)
    return null
  }
})
