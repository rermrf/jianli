import { useId, useMemo, useState } from 'react'
import { ApiError } from '../../lib/api'
import { getAuthKey } from '../../lib/auth'
import { createSquareAvatarFile, uploadAvatar } from '../../lib/upload'
import { Button } from '../common/Button'

interface AvatarUploaderProps {
  currentUrl?: string
  onUploaded: (url: string) => Promise<void> | void
}

export function AvatarUploader({ currentUrl, onUploaded }: AvatarUploaderProps) {
  const inputId = useId()
  const [error, setError] = useState('')
  const [previewUrl, setPreviewUrl] = useState('')
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [uploading, setUploading] = useState(false)

  const displayUrl = useMemo(() => previewUrl || currentUrl || '', [previewUrl, currentUrl])

  async function handleConfirmUpload() {
    if (!selectedFile) {
      return
    }

    const authKey = getAuthKey()
    if (!authKey) {
      setError('未登录，无法上传头像')
      return
    }

    setUploading(true)
    setError('')
    try {
      const croppedFile = await createSquareAvatarFile(selectedFile)
      const avatarUrl = await uploadAvatar(croppedFile, authKey)
      await onUploaded(avatarUrl)
      setPreviewUrl(avatarUrl)
      setSelectedFile(null)
    } catch (error) {
      if (error instanceof ApiError && error.status === 413) {
        setError('头像文件过大，请换一张图片或先压缩后再上传')
      } else {
        setError('头像上传失败，请稍后重试')
      }
    } finally {
      setUploading(false)
    }
  }

  return (
    <div className="space-y-4 rounded-2xl border border-slate-200 bg-slate-50 p-4">
      <div className="flex items-center gap-4">
        {displayUrl ? (
          <img
            alt="头像预览"
            className="h-20 w-20 rounded-full object-cover"
            src={displayUrl}
          />
        ) : (
          <div className="flex h-20 w-20 items-center justify-center rounded-full bg-slate-200 text-3xl text-slate-500">
            👤
          </div>
        )}
        <div className="space-y-3">
          <div>
            <p className="text-sm font-medium text-slate-800">头像</p>
            <p className="text-xs text-slate-500">
              支持图片上传，默认裁剪为 1:1 正方形，并自动压缩后再上传。
            </p>
          </div>
          <div className="flex flex-wrap gap-2">
            <label
              className="inline-flex cursor-pointer items-center justify-center rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 transition hover:border-slate-300 hover:bg-slate-50"
              htmlFor={inputId}
            >
              选择头像
            </label>
            {selectedFile ? (
              <Button onClick={handleConfirmUpload} type="button">
                {uploading ? '上传中...' : '确认上传头像'}
              </Button>
            ) : null}
          </div>
        </div>
      </div>
      <input
        accept="image/*"
        className="sr-only"
        id={inputId}
        onChange={(event) => {
          const file = event.target.files?.[0]
          if (!file) {
            return
          }
          if (!file.type.startsWith('image/')) {
            setError('请选择图片文件')
            return
          }
          if (file.size > 5 * 1024 * 1024) {
            setError('图片不能超过 5MB')
            return
          }
          setError('')
          setSelectedFile(file)
          setPreviewUrl(URL.createObjectURL(file))
        }}
        type="file"
      />
      {error ? <p className="text-sm text-rose-500">{error}</p> : null}
    </div>
  )
}
