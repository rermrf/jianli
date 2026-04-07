import { useEffect, useRef } from 'react'
import { Link } from 'react-router-dom'
import { Button } from '../components/common/Button'
import { AppShell } from '../components/layout/AppShell'
import { TopNav } from '../components/layout/TopNav'
import { ResumeDesktopLayout } from '../components/resume/ResumeDesktopLayout'
import { ResumeMobileLayout } from '../components/resume/ResumeMobileLayout'
import { useResumeDraft } from '../hooks/useResumeDraft'
import { recordVisit, sendVisitDuration } from '../lib/visitors'

export function ResumePage() {
  const { draft, loading } = useResumeDraft()
  const visitIDRef = useRef<number | null>(null)
  const visitStartRef = useRef<number | null>(null)

  useEffect(() => {
    if (loading) {
      return
    }

    let active = true

    async function trackVisit() {
      try {
        const visitorID = await recordVisit()
        if (!active) {
          return
        }

        visitIDRef.current = visitorID
        visitStartRef.current = Date.now()
      } catch {
        if (!active) {
          return
        }

        visitIDRef.current = null
        visitStartRef.current = null
      }
    }

    function handlePageHide() {
      if (visitIDRef.current === null || visitStartRef.current === null) {
        return
      }

      const durationSeconds = Math.max(
        0,
        Math.round((Date.now() - visitStartRef.current) / 1000),
      )
      sendVisitDuration(visitIDRef.current, durationSeconds)
    }

    void trackVisit()
    window.addEventListener('pagehide', handlePageHide)

    return () => {
      active = false
      window.removeEventListener('pagehide', handlePageHide)
    }
  }, [loading])

  if (loading) {
    return (
      <AppShell contentClassName="space-y-6">
        <TopNav />
        <div className="rounded-3xl border border-slate-200 bg-white p-6 text-sm text-slate-500 shadow-[var(--shadow-card)]">
          加载中...
        </div>
      </AppShell>
    )
  }

  return (
    <AppShell contentClassName="space-y-6">
      <TopNav />
      <ResumeMobileLayout resume={draft} />
      <ResumeDesktopLayout resume={draft} />
      <div className="md:hidden">
        <Link to="/print">
          <Button className="w-full">导出 PDF</Button>
        </Link>
      </div>
    </AppShell>
  )
}
