"use server";

import { redirect } from "next/navigation";
import { z } from "zod";
import { SignupFormSchema, FormState } from "@/lib/interfaces/definitions";
import { getToken, deleteSessionFromCookies, setSession } from "@/lib/session";

// define the input schema
const loginSchema = z.object({
  usernameOrEmail: z
    .string()
    .min(1, { message: "Username or Email is required" }),
  password: z
    .string()
    .min(8, { message: "Password must be at least 8 characters" }), // ensure consistent with backend
});

// public function: handle API request
async function handleApiRequest(
  url: string,
  method: string,
  body: Record<string, unknown>
) {
  try {
    console.log("bearer token: ", await getToken());
    const res = await fetch(url, {
      method,
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${await getToken()}`,
      },
      body: JSON.stringify(body),
      credentials: "include", // send and receive cookies
    });

    if (!res.ok) {
      const errorText = await res.text();

      // logout if 401 (unauthorized)
      if (res.status === 401) {
        console.log("LOGOUT CALLED");
        await logout();
      }
      throw new Error(errorText);
    }

    return await res.json();
  } catch (error) {
    console.error(`API request error at ${url} with method ${method}:`, error);
    throw error;
  }
}

// Login Server Action
export async function login(formData: FormData) {
  const parsedData = loginSchema.safeParse({
    usernameOrEmail: formData.get("username"),
    password: formData.get("password"),
  });

  if (!parsedData.success) {
    return { errors: parsedData.error.flatten().fieldErrors };
  }

  const { usernameOrEmail, password } = parsedData.data;

  try {
    console.log({ username: usernameOrEmail, password: password });
    const data = await handleApiRequest(
      `${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/auth/User/login`,
      "POST",
      { usernameOrEmail, password }
    );

    if (data.errors) {
      return { errors: data.errors };
    }

    try {
      // Try to set session in cookies first
      const cookieSession = await setSession(data.token, data.userId);

      if (typeof cookieSession === "undefined") {
        // If cookie session fails (cookies not accepted), use localStorage
        if (typeof window !== "undefined") {
          localStorage.setItem(
            "session",
            JSON.stringify({
              token: data.token,
              userId: data.userId,
            })
          );
          return {
            success: true,
            sessionData: {
              token: data.token,
              userId: data.userId,
            },
          };
        }
      }
    } catch (sessionError) {
      console.error("Failed to set session:", sessionError);
      // Return session data to be handled client-side
      return {
        success: true,
        sessionData: {
          token: data.token,
          userId: data.userId,
        },
      };
    }

    // Only redirect if everything was successful
    redirect("/");
  } catch (error) {
    console.error("Login error:", error);
    return {
      errors: {
        username: ["login error, please try again later"],
      },
    };
  }
}

export async function signup(prevState: FormState, formData: FormData) {
  // validate form fields
  const validatedFields = SignupFormSchema.safeParse({
    firstName: formData.get("firstName"),
    lastName: formData.get("lastName"),
    email: formData.get("email"),
    password: formData.get("password"),
    profilePicture: formData.get("profilePicture") || "",
  });

  // if validation fails, return error information
  if (!validatedFields.success) {
    return {
      errors: validatedFields.error.flatten().fieldErrors,
      message: "please check the input information",
      errorDetails: validatedFields.error.flatten().fieldErrors,
    };
  }

  const { email, password, firstName, lastName, profilePicture } =
    validatedFields.data;
  const username = email.split("@")[0]; // use email prefix as username

  try {
    const data = await handleApiRequest(
      `${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/auth/User/signup`,
      "POST",
      { username, email, password, firstName, lastName, profilePicture }
    );

    console.log("Signup successful:", data); // use console.log to check the data

    // signup successful, redirect to home page
    redirect("/login");
  } catch (error) {
    return {
      message: `error: ${
        error instanceof Error ? error.message : "unknown error"
      }`,
    };
  }
}

export async function logout() {
  // backend does not implement logout route yet
  // try {
  //   await handleApiRequest(`${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/auth/logout`, 'POST', {});
  // } catch (error) {
  //   console.error('logout error:', error);
  // }

  try {
    await deleteSessionFromCookies();
  } catch (error) {
    console.error("Error deleting session from cookies:", error);
  }

  try {
    redirect("/login");
  } catch (error) {
    console.error("Error redirecting to login:", error);
  }
  redirect("/login");
}
