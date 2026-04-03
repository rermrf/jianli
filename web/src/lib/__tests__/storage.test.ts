import { afterEach, describe, expect, it } from 'vitest'
import { defaultResume } from '../../data/mockResume'
import { loadResumeDraft, resetResumeDraft, saveResumeDraft } from '../storage'

describe('resume draft storage', () => {
  afterEach(() => {
    resetResumeDraft()
  })

  it('returns the seed resume when localStorage contains invalid JSON', () => {
    localStorage.setItem('resume:draft', '{invalid-json')

    expect(loadResumeDraft()).toEqual(defaultResume)
  })

  it('persists and reloads the saved draft', () => {
    const updatedResume = {
      ...defaultResume,
      profile: {
        ...defaultResume.profile,
        name: '测试姓名',
      },
    }

    saveResumeDraft(updatedResume)

    expect(loadResumeDraft()).toEqual(updatedResume)
  })
})
