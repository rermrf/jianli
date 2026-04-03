import { PrintResume } from '../components/resume/PrintResume'
import { useResumeDraft } from '../hooks/useResumeDraft'

export function PrintPage() {
  const { draft, loading } = useResumeDraft()

  if (loading) {
    return <main className="min-h-screen bg-white px-4 py-6 text-slate-500">加载中...</main>
  }

  return (
    <main className="min-h-screen bg-white px-4 py-6 text-slate-900 md:px-6 md:py-10 print:px-0 print:py-0">
      <PrintResume resume={draft} />
    </main>
  )
}
