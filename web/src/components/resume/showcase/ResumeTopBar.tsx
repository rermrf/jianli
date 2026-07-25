import { Link } from 'react-router-dom'
import { monogramFromName } from '../../../lib/resumePresentation'

interface ResumeTopBarProps {
  name?: string
  showPdfExport?: boolean
}

export function ResumeTopBar({
  name = '温庆京',
  showPdfExport = false,
}: ResumeTopBarProps) {
  const monogram = monogramFromName(name)

  return (
    <header className="hidden items-center justify-between md:flex">
      <Link
        className="flex items-center gap-3 rounded-full text-sm font-semibold text-[var(--resume-ink-900)] outline-none transition focus-visible:ring-2 focus-visible:ring-[var(--resume-brand-500)] focus-visible:ring-offset-2"
        to="/"
      >
        <span className="flex h-9 w-9 items-center justify-center rounded-full bg-[var(--resume-brand-500)] text-sm font-semibold text-white">
          {monogram}
        </span>
        温庆京的简历
      </Link>

      {showPdfExport ? (
        <Link
          className="resume-btn-primary inline-flex items-center justify-center rounded-xl px-4 py-2 text-sm font-medium transition outline-none focus-visible:ring-2 focus-visible:ring-[var(--resume-brand-500)] focus-visible:ring-offset-2"
          to="/print"
        >
          导出 PDF
        </Link>
      ) : null}
    </header>
  )
}
