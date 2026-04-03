import { Link } from 'react-router-dom'
import { Button } from '../components/common/Button'
import { AppShell } from '../components/layout/AppShell'
import { TopNav } from '../components/layout/TopNav'
import { ResumeDesktopLayout } from '../components/resume/ResumeDesktopLayout'
import { ResumeMobileLayout } from '../components/resume/ResumeMobileLayout'
import { useResumeDraft } from '../hooks/useResumeDraft'

export function ResumePage() {
  const { draft } = useResumeDraft()

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
