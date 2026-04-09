import { useState } from 'react'
import { Button } from '../components/common/Button'
import { PrintResume } from '../components/resume/PrintResume'
import { useResumeDraft } from '../hooks/useResumeDraft'

export function PrintPage() {
  const { draft, loading, siteSettings } = useResumeDraft()
  const [downloading, setDownloading] = useState(false)
  const [error, setError] = useState('')

  async function handleDownloadPDF() {
    setDownloading(true)
    setError('')

    try {
      const response = await fetch('/api/resume/pdf', { method: 'GET' })
      if (!response.ok) {
        throw new Error('failed to download pdf')
      }

      const blob = await response.blob()
      const objectUrl = URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = objectUrl
      link.download = 'resume.pdf'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      URL.revokeObjectURL(objectUrl)
    } catch {
      setError('PDF 下载失败，请稍后重试')
    } finally {
      setDownloading(false)
    }
  }

  if (loading) {
    return <main className="min-h-screen bg-white px-4 py-6 text-slate-500">加载中...</main>
  }

  return (
    <main className="min-h-screen bg-white px-4 py-6 text-slate-900 md:px-6 md:py-10 print:px-0 print:py-0">
      <div className="mx-auto mb-6 flex max-w-[794px] items-center justify-between gap-4 print:hidden">
        <div>
          <h1 className="text-xl font-semibold text-slate-900">打印预览</h1>
          <p className="mt-2 text-sm text-slate-500">确认版式无误后可直接下载 PDF。</p>
        </div>
        {siteSettings.allowPdfExport ? (
          <Button onClick={handleDownloadPDF} type="button">
            {downloading ? '下载中...' : '下载 PDF'}
          </Button>
        ) : null}
      </div>
      {error ? (
        <div className="mx-auto mb-4 max-w-[794px] rounded-2xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-600 print:hidden">
          {error}
        </div>
      ) : null}
      <PrintResume resume={draft} />
    </main>
  )
}
