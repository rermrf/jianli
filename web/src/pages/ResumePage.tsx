import { ResumeShowcase, ResumeSkeleton } from '../components/resume/showcase'
import { useResumeDraft } from '../hooks/useResumeDraft'
import { useVisitorTracking } from '../hooks/useVisitorTracking'

export function ResumePage() {
  const { draft, loading, siteSettings } = useResumeDraft()
  const showPdfExport = !loading && siteSettings.allowPdfExport
  useVisitorTracking(!loading)

  if (loading) {
    return (
      <main className="resume-public min-h-screen">
        <div className="mx-auto max-w-[1200px] space-y-6 px-4 py-6 md:px-6 md:py-10">
          <ResumeSkeleton />
        </div>
      </main>
    )
  }

  return (
    <main className="resume-public min-h-screen">
      <div className="mx-auto max-w-[1200px] px-4 py-6 md:px-6 md:py-10">
        <ResumeShowcase resume={draft} showPdfExport={showPdfExport} />
      </div>
    </main>
  )
}
