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
  password: z.string().min(8, { message: "Password must be at least 8 characters" }),
  firstName: z.string().min(1, { message: "First name is required" }),
  lastName: z.string().min(1, { message: "Last name is required" }),
});

// 密码复杂性检查
const passwordComplexity = (password: string) => {
  const complexityRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;
  return complexityRegex.test(password);
};

// 登录 Server Action
export async function login(formData: FormData) {
  const parsedData = loginSchema.safeParse({
    username: formData.get('username'),
    password: formData.get('password'),
  });

  if (!parsedData.success) {
    return { errors: parsedData.error.flatten().fieldErrors };
  }

  const { username, password } = parsedData.data;

  // 调用后端 API 进行登录
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

  // 登录成功后重定向
  redirect('/');
}

// 注册 Server Action
export async function signup(formData: FormData) {
  const parsedData = signupSchema.safeParse({
    email: formData.get('email') as string, // 强制转换为 string
    password: formData.get('password') as string, // 强制转换为 string
    firstName: formData.get('firstName') as string, // 强制转换为 string
    lastName: formData.get('lastName') as string, // 强制转换为 string
  });

  if (!parsedData.success) {
    return { errors: parsedData.error.flatten().fieldErrors };
  }

  const { email, password, firstName, lastName } = parsedData.data;

  // 检查密码复杂性
  if (!passwordComplexity(password)) {
    return { errors: { password: "Password does not meet complexity requirements" } };
  }

  // 调用后端 API 进行注册
  const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/auth/User/signup`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password, firstName, lastName }),
  });

  if (!res.ok) {
    throw new Error('Signup failed');
  }

  // 注册成功后重定向
  redirect('/');
}
