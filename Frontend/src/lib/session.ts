'use server'

import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'

export async function getSession() {
  try {
    const cookieStore = cookies()
    const authToken = cookieStore.get('magpie_auth')
    return authToken?.value
  } catch (error) {
    console.error('Error fetching session:', error)
    return null
  }
}

export async function createSession(token: string) {
  const cookieStore = cookies()
  // set cookie, use the same configuration as the backend
  cookieStore.set('magpie_auth', token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production', // in production, use secure cookies
    sameSite: 'lax',
    path: '/',
  })
}

export async function deleteSession() {
  const cookieStore = cookies()
  cookieStore.delete('magpie_auth')
  redirect('/login')
}
