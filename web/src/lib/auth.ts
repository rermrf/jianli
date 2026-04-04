import { apiFetch } from './api'

const AUTH_KEY_STORAGE = 'auth:key'
const REDIRECT_PATH_STORAGE = 'auth:redirect-path'

export function isAuthenticated() {
  return Boolean(sessionStorage.getItem(AUTH_KEY_STORAGE))
}

export function getAuthKey() {
  return sessionStorage.getItem(AUTH_KEY_STORAGE)
}

export function loginWithKey(key: string) {
  sessionStorage.setItem(AUTH_KEY_STORAGE, key)
}

export function logout() {
  sessionStorage.removeItem(AUTH_KEY_STORAGE)
}

export function setRedirectPath(path: string) {
  sessionStorage.setItem(REDIRECT_PATH_STORAGE, path)
}

export function consumeRedirectPath() {
  const redirectPath = sessionStorage.getItem(REDIRECT_PATH_STORAGE)

  if (!redirectPath) {
    return null
  }

  sessionStorage.removeItem(REDIRECT_PATH_STORAGE)

  return redirectPath
}

export async function verifyAuthKey(key: string) {
  const response = await apiFetch<{ valid: boolean }>('/api/auth/verify', {
    body: JSON.stringify({ key }),
    method: 'POST',
  })

  return response.valid
}
