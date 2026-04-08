import { ApiError } from './api'

interface UploadAvatarResponse {
  url: string
}

const MAX_AVATAR_DIMENSION = 1024
const TARGET_AVATAR_SIZE_BYTES = 900 * 1024
const JPEG_QUALITIES = [0.86, 0.72, 0.6]

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
    const sourceSize = Math.min(bitmap.width, bitmap.height)
    const outputSize = Math.min(sourceSize, MAX_AVATAR_DIMENSION)
    const offsetX = (bitmap.width - sourceSize) / 2
    const offsetY = (bitmap.height - sourceSize) / 2

    const canvas = document.createElement('canvas')
    canvas.width = outputSize
    canvas.height = outputSize
    const context = canvas.getContext('2d')
    if (!context) {
      return file
    }

    context.drawImage(
      bitmap,
      offsetX,
      offsetY,
      sourceSize,
      sourceSize,
      0,
      0,
      outputSize,
      outputSize,
    )

    let blob: Blob | null = null
    for (const quality of JPEG_QUALITIES) {
      blob = await new Promise<Blob | null>((resolve) => {
        canvas.toBlob(resolve, 'image/jpeg', quality)
      })

      if (blob && blob.size <= TARGET_AVATAR_SIZE_BYTES) {
        break
      }
    }

    if (!blob) {
      return file
    }

    return new File([blob], file.name.replace(/\.[^.]+$/, '') + '.jpg', {
      type: 'image/jpeg',
    })
  } catch {
    return file
  }
}
