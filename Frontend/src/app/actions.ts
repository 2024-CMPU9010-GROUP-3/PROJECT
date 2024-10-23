'use server';

import { redirect } from 'next/navigation';
import { z } from 'zod';

// define the input schema
const loginSchema = z.object({
  username: z.string().min(1, { message: "Username is required" }),
  password: z.string().min(6, { message: "Password must be at least 6 characters" }),
});

const signupSchema = z.object({
  email: z.string().email({ message: "Invalid email" }),
  password: z.string().min(8, { message: "Password must contain uppercase, lowercase numbers and special characters" }),
  firstName: z.string().min(1, { message: "First name is required" }),
  lastName: z.string().min(1, { message: "Last name is required" }),
  profilePicture: z.string().optional(), // optional
});

// // check password complexity
// const passwordComplexity = (password: string) => {
//   const complexityRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;
//   return complexityRegex.test(password);
// };

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

  // call backend api to fetch
  const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/auth/User/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ username, password }),
  });

  if (!res.ok) {
    throw new Error('Login failed');
  }

  // redirect after login successfully
  redirect('/');
}

// signup Server Action
export async function signup(formData: FormData) {
  const parsedData = signupSchema.safeParse({
    email: formData.get('email'),
    password: formData.get('password'),
    firstName: formData.get('firstName'),
    lastName: formData.get('lastName'),
    profilePicture: formData.get('profilePicture'), 
  });

  if (!parsedData.success) {
    return { errors: parsedData.error.flatten().fieldErrors };
  }

  const { email, password, firstName, lastName, profilePicture } = parsedData.data; 

  // extract username
  const username = email.split("@")[0]; // Use the first half of your email as username

//   // check password complexity
//   if (!passwordComplexity(password)) {
//     return { errors: { password: "Password must contain uppercase, lowercase numbers and special characters" } };
//   }

  // Print request data for debugging
  console.log('Signup data:', { email, password, firstName, lastName, username, profilePicture });

  // call backend api to fetch
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
  });

  if (!res.ok) {
    const errorMessage = await res.text(); // await for the response text
    throw new Error(`Signup failed: ${errorMessage}`); // throw an error with the error message
  }

  // redirect after signup successfully
  redirect('/');
}
