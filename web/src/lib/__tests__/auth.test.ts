import { afterEach, describe, expect, it } from 'vitest'
import {
  consumeRedirectPath,
  isAuthenticated,
  loginWithKey,
  logout,
  setRedirectPath,
} from '../auth'

describe('auth helpers', () => {
  afterEach(() => {
    sessionStorage.clear()
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
})
