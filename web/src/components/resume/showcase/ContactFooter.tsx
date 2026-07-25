import { Link } from 'react-router-dom'
import type { JobIntention, ResumeProfile } from '../../../types/resume'
import { GlassCard } from './GlassCard'

interface ContactFooterProps {
  profile: ResumeProfile
  jobIntention: JobIntention
  showPdfExport?: boolean
}

export function ContactFooter({
  profile,
  jobIntention,
  showPdfExport = false,
}: ContactFooterProps) {
  const cities = jobIntention.cities.filter(Boolean)

  return (
    <GlassCard className="space-y-4">
      <h2 className="text-lg font-semibold text-[var(--resume-ink-900)]">
        联系与意向
      </h2>

      <div className="grid gap-3 text-sm text-[var(--resume-ink-600)] sm:grid-cols-2">
        {cities.length > 0 ? (
          <p>
            <span className="text-[var(--resume-ink-400)]">意向城市：</span>
            {cities.join(' / ')}
          </p>
        ) : null}
        {jobIntention.availability ? (
          <p>
            <span className="text-[var(--resume-ink-400)]">到岗时间：</span>
            {jobIntention.availability}
          </p>
        ) : null}
        {profile.phone ? (
          <p>
            <span className="text-[var(--resume-ink-400)]">手机：</span>
            <a
              className="outline-none transition hover:text-[var(--resume-brand-600)] focus-visible:ring-2 focus-visible:ring-[var(--resume-brand-500)] focus-visible:ring-offset-2"
              href={`tel:${profile.phone}`}
            >
              {profile.phone}
            </a>
          </p>
        ) : null}
        {profile.email ? (
          <p>
            <span className="text-[var(--resume-ink-400)]">邮箱：</span>
            <a
              className="outline-none transition hover:text-[var(--resume-brand-600)] focus-visible:ring-2 focus-visible:ring-[var(--resume-brand-500)] focus-visible:ring-offset-2"
              href={`mailto:${profile.email}`}
            >
              {profile.email}
            </a>
          </p>
        ) : null}
      </div>

      {showPdfExport ? (
        <div className="md:hidden">
          <Link
            className="resume-btn-primary inline-flex w-full items-center justify-center rounded-xl px-4 py-2.5 text-sm font-medium transition outline-none focus-visible:ring-2 focus-visible:ring-[var(--resume-brand-500)] focus-visible:ring-offset-2"
            to="/print"
          >
            导出 PDF
          </Link>
        </div>
      ) : null}
    </GlassCard>
  )
}
