import { afterEach, describe, expect, it, vi } from 'vitest'
import { defaultResume } from '../../data/mockResume'
import { loadResumeDraft, resetResumeDraft, saveResumeDraft } from '../storage'

describe('resume draft storage', () => {
  afterEach(() => {
    resetResumeDraft()
    vi.restoreAllMocks()
  })

  it('loads the resume and site settings from the backend API', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(
        JSON.stringify({
          code: 0,
          data: {
            resume: defaultResume,
            siteSettings: {
              allowPdfExport: false,
            },
          },
        }),
        {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        },
      ),
    )

    await expect(loadResumeDraft()).resolves.toEqual({
      resume: defaultResume,
      siteSettings: {
        allowPdfExport: false,
      },
    })
  })

  it('saves the resume through the backend API', async () => {
    const updatedResume = {
      ...defaultResume,
      profile: {
        ...defaultResume.profile,
        name: '测试姓名',
      },
    }
    const fetchSpy = vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ code: 0, data: updatedResume }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    await expect(saveResumeDraft(updatedResume)).resolves.toEqual(updatedResume)

    expect(fetchSpy).toHaveBeenCalledWith(
      '/api/resume',
      expect.objectContaining({
        body: JSON.stringify(updatedResume),
        method: 'PUT',
      }),
    )
  })
})
