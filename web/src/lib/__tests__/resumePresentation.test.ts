import { describe, expect, it } from 'vitest'
import { buildTagline, monogramFromName } from '../resumePresentation'
import { defaultResume } from '../../data/mockResume'

describe('resumePresentation', () => {
  it('builds a stable tagline from profile + intention', () => {
    expect(buildTagline(defaultResume)).toBe(
      'Golang 后端工程师 · 0.9年经验 · 意向 Golang 后端工程师',
    )
  })

  it('uses first character as monogram', () => {
    expect(monogramFromName('温庆京')).toBe('温')
    expect(monogramFromName('')).toBe('?')
  })
})
