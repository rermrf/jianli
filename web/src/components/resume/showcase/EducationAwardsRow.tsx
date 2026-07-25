import type { Award, EducationExperience } from '../../../types/resume'
import { GlassCard } from './GlassCard'

interface EducationAwardsRowProps {
  education: EducationExperience[]
  awards: Award[]
}

export function EducationAwardsRow({
  education,
  awards,
}: EducationAwardsRowProps) {
  const hasEducation = education.length > 0
  const hasAwards = awards.length > 0

  if (!hasEducation && !hasAwards) return null

  const dual = hasEducation && hasAwards

  return (
    <div className={dual ? 'grid gap-4 md:grid-cols-2' : 'grid gap-4'}>
      {hasEducation ? (
        <GlassCard className="space-y-4">
          <h2 className="text-lg font-semibold text-[var(--resume-ink-900)]">
            教育经历
          </h2>
          <ul className="space-y-4">
            {education.map((item) => (
              <li
                key={`${item.school}-${item.major}-${item.startDate}`}
                className="space-y-1"
              >
                <h3 className="text-sm font-semibold text-[var(--resume-ink-900)]">
                  {item.school}
                </h3>
                <p className="text-sm text-[var(--resume-ink-600)]">
                  {`${item.major} · ${item.degree}`}
                </p>
                <p className="text-xs text-[var(--resume-ink-400)]">
                  {`${item.startDate} - ${item.endDate}`}
                </p>
              </li>
            ))}
          </ul>
        </GlassCard>
      ) : null}

      {hasAwards ? (
        <GlassCard className="space-y-4">
          <h2 className="text-lg font-semibold text-[var(--resume-ink-900)]">
            荣誉奖项
          </h2>
          <ul className="space-y-3">
            {awards.map((item) => (
              <li
                key={`${item.title}-${item.date}`}
                className="flex items-start justify-between gap-3"
              >
                <h3 className="text-sm font-semibold text-[var(--resume-ink-900)]">
                  {item.title}
                </h3>
                <p className="shrink-0 text-xs text-[var(--resume-ink-400)]">
                  {item.date}
                </p>
              </li>
            ))}
          </ul>
        </GlassCard>
      ) : null}
    </div>
  )
}
