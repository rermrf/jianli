import type { WorkExperience } from '../../../types/resume'
import { GlassCard } from './GlassCard'

interface WorkTimelineProps {
  items: WorkExperience[]
}

export function WorkTimeline({ items }: WorkTimelineProps) {
  if (items.length === 0) return null

  return (
    <GlassCard className="space-y-5">
      <h2 className="text-lg font-semibold text-[var(--resume-ink-900)]">
        工作经历
      </h2>

      <ol className="relative space-y-6 pl-6">
        <span
          aria-hidden="true"
          className="absolute bottom-1 left-[7px] top-1 w-px bg-[var(--resume-brand-50)]"
        />

        {items.map((item) => (
          <li
            key={`${item.company}-${item.role}-${item.startDate}`}
            className="relative"
          >
            <span
              aria-hidden="true"
              className="absolute -left-6 top-1.5 h-3.5 w-3.5 rounded-full border-2 border-white bg-[var(--resume-brand-500)] shadow-[0_0_0_3px_var(--resume-brand-50)]"
            />

            <div className="space-y-1">
              <div className="flex flex-col gap-1 sm:flex-row sm:items-start sm:justify-between">
                <h3 className="text-sm font-semibold text-[var(--resume-ink-900)]">
                  {`${item.company} · ${item.role}`}
                </h3>
                <p className="text-xs text-[var(--resume-ink-400)] sm:text-sm">
                  {`${item.startDate} - ${item.endDate}`}
                </p>
              </div>

              {item.description.length > 0 ? (
                <ul className="mt-2 space-y-2 text-sm leading-6 text-[var(--resume-ink-600)]">
                  {item.description.map((bullet) => (
                    <li key={bullet} className="flex gap-2">
                      <span
                        aria-hidden="true"
                        className="mt-2 h-1.5 w-1.5 shrink-0 rounded-full bg-[var(--resume-brand-500)]"
                      />
                      <span>{bullet}</span>
                    </li>
                  ))}
                </ul>
              ) : null}
            </div>
          </li>
        ))}
      </ol>
    </GlassCard>
  )
}
