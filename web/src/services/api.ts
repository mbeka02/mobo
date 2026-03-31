const API_BASE = '/api/v1'

export type APIError = {
  status: number
  message: string
  detail?: string
}

/** Type guard for APIError */
export function isAPIError(error: unknown): error is APIError {
  return (
    typeof error === 'object' &&
    error !== null &&
    'status' in error &&
    'message' in error
  )
}

/** Parse the backend error response into an APIError and throw it */
async function handleError(response: Response): Promise<never> {
  const apiError: APIError = {
    status: response.status,
    message: 'An unknown error occurred',
  }

  try {
    const data = await response.json()
    if (data?.message) apiError.message = data.message
    if (data?.detail) apiError.detail = data.detail
  } catch {
    // Response body wasn't valid JSON — keep the fallback message
  }

  throw apiError
}

/** Generic request helper — all requests include credentials for HTTP-only cookies */
async function request<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  })

  if (!response.ok) {
    return handleError(response)
  }

  // Handle 204 No Content (e.g. logout)
  if (response.status === 204) {
    return undefined as T
  }

  return response.json()
}

export const api = {
  get<T>(path: string): Promise<T> {
    return request<T>(path, { method: 'GET' })
  },

  post<T>(path: string, body?: unknown): Promise<T> {
    return request<T>(path, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    })
  },

  patch<T>(path: string, body?: unknown): Promise<T> {
    return request<T>(path, {
      method: 'PATCH',
      body: body ? JSON.stringify(body) : undefined,
    })
  },

  delete<T>(path: string): Promise<T> {
    return request<T>(path, { method: 'DELETE' })
  },
}
