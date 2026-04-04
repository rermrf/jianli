interface ApiSuccess<T> {
  code: number
  data: T
}

interface ApiFailure {
  code: number
  message: string
}

export class ApiError extends Error {
  code: number
  status: number

  constructor(message: string, code: number, status: number) {
    super(message)
    this.code = code
    this.status = status
  }
}

export async function apiFetch<T>(
  input: string,
  init?: RequestInit,
): Promise<T> {
  const response = await fetch(input, {
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {}),
    },
    ...init,
  })

  const payload = (await response.json()) as ApiSuccess<T> | ApiFailure

  if (!response.ok) {
    throw new ApiError(
      'message' in payload ? payload.message : 'request failed',
      'code' in payload ? payload.code : response.status,
      response.status,
    )
  }

  return (payload as ApiSuccess<T>).data
}
