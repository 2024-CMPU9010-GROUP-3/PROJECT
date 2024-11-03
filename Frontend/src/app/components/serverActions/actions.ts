'use server';

import { redirect } from 'next/navigation';
import { z } from 'zod';
import { SignupFormSchema, FormState } from '@/lib/interfaces/definitions';
import { deleteSession, createSession } from '@/lib/session';

// define the input schema
const loginSchema = z.object({
  username: z.string().min(1, { message: "Username is required" }),
  password: z.string().min(8, { message: "Password must be at least 8 characters" }), // ensure consistent with backend
});

// public function: handle API request
async function handleApiRequest(url: string, method: string, body: any) {
  try {
    const res = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
      credentials: 'include', // send and receive cookies
    });

    if (!res.ok) {
      const errorText = await res.text();
      throw new Error(errorText);
    }

    return await res.json();
  } catch (error) {
    console.error('API request error:', error);
    throw error;
  }
}

// Login Server Action
export async function login(formData: FormData) {
  const parsedData = loginSchema.safeParse({
    username: formData.get('username'),
    password: formData.get('password'),
  });

  if (!parsedData.success) {
    return { errors: parsedData.error.flatten().fieldErrors };
  }

  const { username, password } = parsedData.data;

  try {
    const data = await handleApiRequest(
      `${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/auth/User/login`,
      'POST',
      { username, password }
    );

    if (data.errors) {
      return { errors: data.errors };
    }

    // create session
    await createSession(data.token);

    // login successful, redirect to home page
    redirect('/');
  } catch (error) {
    return {
      errors: {
        username: ['login error, please try again later'],
      },
    };
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
  try {
    await handleApiRequest(`${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/auth/logout`, 'POST', {});
  } catch (error) {
    console.error('logout error:', error);
  }

  // whether the backend request is successful, clear the local session
  await deleteSession()
}
