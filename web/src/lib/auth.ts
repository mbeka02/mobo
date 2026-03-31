const API_BASE = "/api/v1";

export interface User {
  id: string;
  full_name: string;
  email: string;
  created_at: string;
}

export interface AuthError {
  status: number;
  message: string;
}

function isAuthError(error: unknown): error is AuthError {
  return (
    typeof error === "object" &&
    error !== null &&
    "status" in error &&
    "message" in error
  );
}

/**
 * User-friendly error messages by status code per action context.
 * We never expose raw `detail` from the API to the user.
 */
const friendlyErrors: Record<string, Record<number, string>> = {
  login: {
    400: "Please check your email and password format.",
    401: "Invalid email or password. Please try again.",
    409: "This account uses Google sign-in. Please use the Google button above.",
    500: "We're having trouble signing you in. Please try again later.",
  },
  signup: {
    400: "Please check your details and try again.",
    409: "An account with this email already exists. Try signing in instead.",
    500: "We couldn't create your account right now. Please try again later.",
  },
  default: {
    401: "Your session has expired. Please sign in again.",
    403: "You don't have permission to do that.",
    500: "Something went wrong on our end. Please try again later.",
  },
};

function getFriendlyMessage(status: number, context: string): string {
  const contextErrors = friendlyErrors[context] || friendlyErrors.default;
  return (
    contextErrors[status] ||
    friendlyErrors.default[status] ||
    "Something went wrong. Please try again."
  );
}

async function handleResponse<T>(
  response: Response,
  context: string = "default"
): Promise<T> {
  if (!response.ok) {
    const message = getFriendlyMessage(response.status, context);
    const error: AuthError = { status: response.status, message };
    throw error;
  }
  return response.json();
}

export async function login(
  email: string,
  password: string
): Promise<User> {
  const res = await fetch(`${API_BASE}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify({ email, password }),
  });
  return handleResponse<User>(res, "login");
}

export async function signup(
  full_name: string,
  email: string,
  password: string
): Promise<User> {
  const res = await fetch(`${API_BASE}/auth/signup`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify({ full_name, email, password }),
  });
  return handleResponse<User>(res, "signup");
}

export async function getCurrentUser(): Promise<User> {
  const res = await fetch(`${API_BASE}/me`, {
    credentials: "include",
  });
  return handleResponse<User>(res);
}

export async function logout(): Promise<void> {
  await fetch(`${API_BASE}/auth/logout`, {
    method: "POST",
    credentials: "include",
  });
}

export function redirectToGoogleOAuth(): void {
  window.location.href = `${API_BASE}/auth/google`;
}

export { isAuthError };
