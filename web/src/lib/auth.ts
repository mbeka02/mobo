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

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let message = "Something went wrong";
    try {
      const body = await response.json();
      message = body.error || body.message || message;
    } catch {
      // ignore parse errors
    }
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
  return handleResponse<User>(res);
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
  return handleResponse<User>(res);
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
