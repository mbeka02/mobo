import { api, type APIError, isAPIError } from '../services/api'

export interface User {
  id: string
  full_name: string
  email: string
  created_at: string
}

/** Re-export for route-level consumption */
export { type APIError, isAPIError }

/**
 * User-friendly error messages by status code per action context.
 * We never expose raw `detail` from the API to the user.
 */
const friendlyErrors: Record<string, Record<number, string>> = {
  login: {
    400: 'Please check your email and password format.',
    401: 'Invalid email or password. Please try again.',
    409: 'This account uses Google sign-in. Please use the Google button above.',
    500: "We're having trouble signing you in. Please try again later.",
  },
  signup: {
    400: 'Please check your details and try again.',
    409: 'An account with this email already exists. Try signing in instead.',
    500: "We couldn't create your account right now. Please try again later.",
  },
  default: {
    401: 'Your session has expired. Please sign in again.',
    403: "You don't have permission to do that.",
    500: 'Something went wrong on our end. Please try again later.',
  },
}

export function getFriendlyMessage(
  status: number,
  context: string = 'default',
): string {
  const contextErrors = friendlyErrors[context] || friendlyErrors.default
  return (
    contextErrors[status] ||
    friendlyErrors.default[status] ||
    'Something went wrong. Please try again.'
  )
}

export async function login(email: string, password: string): Promise<User> {
  return api.post<User>('/auth/login', { email, password })
}

export async function signup(
  full_name: string,
  email: string,
  password: string,
): Promise<User> {
  return api.post<User>('/auth/signup', { full_name, email, password })
}

export async function getCurrentUser(): Promise<User> {
  return api.get<User>('/me')
}

export async function logout(): Promise<void> {
  return api.post('/auth/logout')
}

export function redirectToGoogleOAuth(): void {
  window.location.href = '/api/v1/auth/google'
}
