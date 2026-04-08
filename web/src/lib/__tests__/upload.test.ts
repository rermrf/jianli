import { afterEach, describe, expect, it, vi } from 'vitest'
import { createSquareAvatarFile, uploadAvatar } from '../upload'

describe('avatar upload client', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('uploads the avatar file as multipart form-data and returns the url', async () => {
    const fetchSpy = vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ code: 0, data: { url: '/uploads/avatars/avatar-1.png' } }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    const file = new File([new Uint8Array([1, 2, 3])], 'avatar.png', {
      type: 'image/png',
    })

    await expect(uploadAvatar(file, 'resume-key')).resolves.toBe('/uploads/avatars/avatar-1.png')
    expect(fetchSpy).toHaveBeenCalledWith(
      '/api/upload/avatar',
      expect.objectContaining({
        method: 'POST',
        headers: expect.objectContaining({ 'X-Auth-Key': 'resume-key' }),
      }),
    )
  })

  it('creates a cropped jpeg avatar instead of always inflating to png', async () => {
    const drawImage = vi.fn()
    const toBlob = vi.fn((callback: (blob: Blob | null) => void, type?: string, quality?: number) => {
      callback(new Blob([new Uint8Array([1, 2, 3])], { type: type ?? 'image/jpeg' }))
      expect(type).toBe('image/jpeg')
      expect(quality).toBeGreaterThan(0.7)
      expect(quality).toBeLessThanOrEqual(0.92)
    })

    vi.stubGlobal('window', {
      ...window,
      createImageBitmap: vi.fn().mockResolvedValue({ width: 1200, height: 800 }),
    })

    vi.spyOn(document, 'createElement').mockReturnValue({
      width: 0,
      height: 0,
      getContext: vi.fn(() => ({ drawImage })),
      toBlob,
    } as unknown as HTMLCanvasElement)

    const file = new File([new Uint8Array([1, 2, 3])], 'avatar.png', {
      type: 'image/png',
    })

    const result = await createSquareAvatarFile(file)

    expect(result.type).toBe('image/jpeg')
    expect(result.name).toBe('avatar.jpg')
    expect(drawImage).toHaveBeenCalled()
    expect(toBlob).toHaveBeenCalled()
  })
})
