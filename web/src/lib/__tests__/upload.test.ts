import { describe, expect, it, vi } from 'vitest'
import { uploadAvatar } from '../upload'

describe('avatar upload client', () => {
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
})
