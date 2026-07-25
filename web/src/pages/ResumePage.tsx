import { Link } from 'react-router-dom'
import { Button } from '../components/common/Button'
import { AppShell } from '../components/layout/AppShell'
import { TopNav } from '../components/layout/TopNav'
import { ResumeDesktopLayout } from '../components/resume/ResumeDesktopLayout'
import { ResumeMobileLayout } from '../components/resume/ResumeMobileLayout'
import { useResumeDraft } from '../hooks/useResumeDraft'
import { useVisitorTracking } from '../hooks/useVisitorTracking'

export function ResumePage() {
  const { draft, loading, siteSettings } = useResumeDraft()
  const showPdfExport = !loading && siteSettings.allowPdfExport

  // Track this visit + send 15s heartbeats once the resume has loaded.
  // Hook handles cross-page session continuity via sessionStorage.
  useVisitorTracking(!loading)

  if (loading) {
    return (
      <AppShell contentClassName="space-y-6">
        <TopNav showPdfExport={showPdfExport} />
        <div className="rounded-3xl border border-slate-200 bg-white p-6 text-sm text-slate-500 shadow-[var(--shadow-card)]">
          加载中...
        </div>
      </AppShell>
    )
  }

  return (
    <AppShell contentClassName="space-y-6">
      <TopNav showPdfExport={showPdfExport} />
      <ResumeMobileLayout resume={draft} />
      <ResumeDesktopLayout resume={draft} />
      {showPdfExport ? (
        <div className="md:hidden">
          <Link to="/print">
            <Button className="w-full">导出 PDF</Button>
          </Link>
        </div>
      ) : null}
    </AppShell>
  )
}
