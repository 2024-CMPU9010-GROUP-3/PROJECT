import { NextResponse } from 'next/server'
import { getUser } from '@/lib/dal'

export async function GET() {
  const user = await getUser()
  
  if (!user) {
    return NextResponse.json({ error: 'Unauthorized' }, { status: 401 })
  }

  return NextResponse.json({ user })
}
