import {deleteSessionFromCookies} from '@/app/components/serverActions/actions';
import {getCookiesAccepted} from '@/lib/cookies';
import { cookies } from 'next/headers';
import {NextResponse} from 'next/server';


const tokenCookieName = "magpie_token";
const uuidCookieName = "magpie_uuid";

export async function GET() {
  const cookieStore = cookies();

  if (!(await getCookiesAccepted())) {
    await deleteSessionFromCookies();
    return NextResponse.json({ token: "", uuid: "" });
  }

  const tokenCookie = cookieStore.get(tokenCookieName);
  const uuidCookie = cookieStore.get(uuidCookieName);

  if (tokenCookie && uuidCookie) {
    console.log("RElOADED SESSION FROM COOKIES >>> ", tokenCookie, uuidCookie)
    return NextResponse.json({
      token: tokenCookie.value,
      uuid: uuidCookie.value
    });
  }

  return NextResponse.json({ token: "", uuid: "" });
}