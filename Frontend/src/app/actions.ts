'use server';

import { redirect } from 'next/navigation';
import { z } from 'zod';
import { SignupFormSchema, FormState } from '@/lib/definitions'
import { deleteSession } from '@/lib/session'
import { verifySession } from '@/lib/dal'

// define the input schema
const loginSchema = z.object({
  username: z.string().min(1, { message: "Username is required" }),
  password: z.string().min(6, { message: "Password must be at least 6 characters" }),
});




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
    const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/auth/User/login`, {
      method: 'POST',
      credentials: 'include', // ensure sending and receiving cookies
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    console.log('Response status:', res.status); // add debug information
    console.log('Response headers:', res.headers); // add debug information

    if (!res.ok) {
      return { 
        errors: { 
          username: ['login failed, please check username and password'] 
        } 
      };
    }

    const data = await res.json();
    console.log('Response data:', data); // add debug information

    if (data.errors) {
      return { errors: data.errors };
    }

    // // assume the response from the backend contains the user ID
    // const userId = data.response.content.userid;
    // localStorage.setItem('userId', userId); // store the user ID in localStorage

    // login successful, redirect to the home page
    redirect('/');
    
  } catch (error) {
    console.error('login error:', error); // add debug information
    return { 
      errors: { 
        username: ['login error, please try again later'] 
      } 
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
  })

  // if validation fails, return error information
  if (!validatedFields.success) {
    return {
      errors: validatedFields.error.flatten().fieldErrors,
      message: 'please check the input information'
    }
  }

  const { email, password, firstName, lastName, profilePicture } = validatedFields.data
  const username = email.split("@")[0] // use email prefix as username

  try {
    // call backend API
    const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/auth/User/signup`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username,
        email,
        password,
        firstName,
        lastName,
        profilePicture
      }),
    })

    if (!res.ok) {
      const errorText = await res.text()
      if (errorText.includes('email already exists')) {
        return {
          errors: { email: ['email already exists'] }
        }
      }
      if (errorText.includes('username already exists')) {
        return {
          errors: { email: ['username already exists'] }
        }
      }
      throw new Error(`signup failed: ${errorText}`)
    }

    // redirect after signup
    redirect('/')
    
  } catch (error) {
    return {
      message: `error: ${error instanceof Error ? error.message : 'unknown error'}`
    }
  }
}

export async function logout() {
  try {
    const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/auth/logout`, {
      method: 'POST',
      credentials: 'include', // send cookie
    })

    if (!res.ok) {
      throw new Error('logout failed')
    }
  } catch (error) {
    console.error('logout error:', error)
  }

  // whether the backend request is successful, clear the local session
  await deleteSession()
}

export async function updateUserProfile(formData: FormData) {
  const session = await verifySession()
  if (!session) {
    return {
      error: 'unauthorized operation'
    }
  }

  // continue processing the update user profile logic
  // ...
}
