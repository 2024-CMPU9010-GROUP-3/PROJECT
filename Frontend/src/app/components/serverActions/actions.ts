'use server';

import { redirect } from 'next/navigation';
import { SignupFormSchema, FormState } from '@/lib/interfaces/definitions';
import {cookies} from 'next/headers';
import { jwtDecode, JwtPayload } from "jwt-decode";

const tokenCookieName = "magpie_token";
const uuidCookieName = "magpie_uuid";

// public function: handle API request
async function handleApiRequest(url: string, method: string, body: Record<string, unknown> ) {
  console.log("HANDLE API REQUEST")
  try {

    // console.log("bearer token: ", await getToken())
    const res = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        // 'Authorization': `Bearer ${await getToken()}`
      },
      body: JSON.stringify(body),
      credentials: 'include', // send and receive cookies
    });

    if (!res.ok) {
      const errorText = await res.text();

      // logout if 401 (unauthorized)
      if (res.status === 401){
        console.log("LOGOUT CALLED")
        await logout();
      }
      throw new Error(errorText);
    }

    return await res.json();
  } catch (error) {
    console.error('API request error:', error);
    throw error;
  }
}

export async function signup(prevState: FormState, formData: FormData) {
  // validate form fields
  const validatedFields = SignupFormSchema.safeParse({
    firstName: formData.get('firstName'),
    lastName: formData.get('lastName'),
    email: formData.get('email'),
    password: formData.get('password'),
    profilePicture: formData.get('profilePicture') || '',
  });

  // if validation fails, return error information
  if (!validatedFields.success) {
    return {
      errors: validatedFields.error.flatten().fieldErrors,
      message: 'please check the input information',
      errorDetails: validatedFields.error.flatten().fieldErrors,
    };
  }

  const { email, password, firstName, lastName, profilePicture } = validatedFields.data
  const username = email.split("@")[0] // use email prefix as username

  try {
    const data = await handleApiRequest(
      `${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/auth/User/signup`,
      'POST',
      { username, email, password, firstName, lastName, profilePicture }
    );

    console.log('Signup successful:', data); // use console.log to check the data

    // signup successful, redirect to home page
    redirect('/');
  } catch (error) {
    return {
      message: `error: ${error instanceof Error ? error.message : 'unknown error'}`,
    };
  }
}

export async function logout() {
  redirect('/');
}

export async function saveSessionToCookies(sessionToken: string, sessionUUID: string) {
  const cookieStore = cookies();

  // Set cookies to store sessionToken and sessionUUID
  let expiryDate = new Date(Date.now() + 86400 * 1000); // 1 day expiry default
  const decoded = jwtDecode<JwtPayload>(sessionToken)
  if(decoded.exp){
    expiryDate = new Date(decoded.exp * 1000);
  }
  cookieStore.set(tokenCookieName, sessionToken, { expires: expiryDate });
  cookieStore.set(uuidCookieName, sessionUUID, { expires: expiryDate });
}

export async function deleteSessionFromCookies() {
  const cookieStore = cookies();
  cookieStore.delete(tokenCookieName);
  cookieStore.delete(uuidCookieName);
}
