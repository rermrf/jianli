import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  consumeRedirectPath,
  isAuthenticated,
  loginWithKey,
  logout,
  setRedirectPath,
  verifyAuthKey,
} from '../auth'

describe('auth helpers', () => {
  afterEach(() => {
    sessionStorage.clear()
    vi.restoreAllMocks()
  })

  it('marks the session as authenticated after login', () => {
    expect(isAuthenticated()).toBe(false)

    loginWithKey('resume-key')

    expect(isAuthenticated()).toBe(true)
  })

  it('stores and consumes a redirect target once', () => {
    setRedirectPath('/visitors')

    expect(consumeRedirectPath()).toBe('/visitors')
    expect(consumeRedirectPath()).toBeNull()
  })

  it('clears the auth state on logout', () => {
    loginWithKey('resume-key')

    logout()

    expect(isAuthenticated()).toBe(false)
  })

  it('verifies auth key against the backend API', async () => {
    const fetchSpy = vi
      .spyOn(globalThis, 'fetch')
      .mockResolvedValue(
        new Response(JSON.stringify({ code: 0, data: { valid: true } }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        }),
      )

    await expect(verifyAuthKey('resume-key')).resolves.toBe(true)
    expect(fetchSpy).toHaveBeenCalledWith(
      '/api/auth/verify',
      expect.objectContaining({
        body: JSON.stringify({ key: 'resume-key' }),
        method: 'POST',
      }),
    )
  })
})
