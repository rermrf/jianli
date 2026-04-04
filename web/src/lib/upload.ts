import { ApiError } from './api'

interface UploadAvatarResponse {
  url: string
}

export async function uploadAvatar(file: File, authKey: string): Promise<string> {
  const body = new FormData()
  body.append('file', file)

  const response = await fetch('/api/upload/avatar', {
    method: 'POST',
    headers: {
      'X-Auth-Key': authKey,
    },
    body,
  })

  const payload = (await response.json()) as {
    code: number
    data?: UploadAvatarResponse
    message?: string
  }

  if (!response.ok || !payload.data) {
    throw new ApiError(payload.message ?? 'upload failed', payload.code, response.status)
  }

  return payload.data.url
}

export async function createSquareAvatarFile(file: File): Promise<File> {
  if (
    typeof document === 'undefined' ||
    typeof window === 'undefined' ||
    typeof window.createImageBitmap !== 'function'
  ) {
    return file
  }

  try {
    const bitmap = await window.createImageBitmap(file)
    const size = Math.min(bitmap.width, bitmap.height)
    const offsetX = (bitmap.width - size) / 2
    const offsetY = (bitmap.height - size) / 2

    const canvas = document.createElement('canvas')
    canvas.width = size
    canvas.height = size
    const context = canvas.getContext('2d')
    if (!context) {
      return file
    }

    context.drawImage(bitmap, offsetX, offsetY, size, size, 0, 0, size, size)

    const blob = await new Promise<Blob | null>((resolve) => {
      canvas.toBlob(resolve, 'image/png')
    })

    if (!blob) {
      return file
    }

    return new File([blob], file.name.replace(/\.[^.]+$/, '') + '.png', {
      type: 'image/png',
    })
  } catch {
    return file
  }
}
