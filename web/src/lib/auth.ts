const AUTH_KEY_STORAGE = 'auth:key'
const REDIRECT_PATH_STORAGE = 'auth:redirect-path'

export function isAuthenticated() {
  return Boolean(sessionStorage.getItem(AUTH_KEY_STORAGE))
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
