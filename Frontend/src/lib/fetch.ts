import { deleteSession } from './session' // import deleteSession function

export async function authenticatedFetch(url: string, options: RequestInit = {}) {
  const res = await fetch(url, {
    ...options,
    credentials: 'include', // send cookie automatically
    headers: {
      ...options.headers,
      'Content-Type': 'application/json',
    },
  })

  if (res.status === 401) {
    // if not authenticated, delete session and redirect to login page
    await deleteSession()
  }

  return res
}
