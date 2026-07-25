import type { ProjectExperience } from '../../../types/resume'
import { GlassCard } from './GlassCard'

interface ProjectShowcaseProps {
  projects: ProjectExperience[]
}

export function ProjectShowcase({ projects }: ProjectShowcaseProps) {
  if (projects.length === 0) return null

  return (
    <GlassCard strong className="space-y-5">
      <div className="flex items-center gap-3">
        <h2 className="text-lg font-semibold text-[var(--resume-ink-900)]">
          项目经历
        </h2>
        <span className="rounded-full bg-[var(--resume-brand-50)] px-2.5 py-0.5 text-xs font-medium text-[var(--resume-brand-600)]">
          {projects.length}
        </span>
      </div>

      <div className="flex flex-col gap-4">
        {projects.map((project) => (
          <article
            key={`${project.name}-${project.startDate}-${project.endDate}`}
            className="relative overflow-hidden rounded-2xl border border-white/50 bg-white/40 p-4 pl-5"
          >
            <span
              aria-hidden="true"
              className="absolute inset-y-0 left-0 w-[3px] bg-[var(--resume-brand-500)]"
            />

            <div className="flex flex-col gap-1 sm:flex-row sm:items-start sm:justify-between">
              {project.url ? (
                <a
                  className="text-sm font-semibold text-[var(--resume-brand-600)] outline-none transition hover:underline focus-visible:ring-2 focus-visible:ring-[var(--resume-brand-500)] focus-visible:ring-offset-2"
                  href={project.url}
                  rel="noopener noreferrer"
                  target="_blank"
                >
                  {project.name}
                </a>
              ) : (
                <h3 className="text-sm font-semibold text-[var(--resume-ink-900)]">
                  {project.name}
                </h3>
              )}
              <p className="text-xs text-[var(--resume-ink-400)] sm:text-sm">
                {`${project.startDate} - ${project.endDate}`}
              </p>
            </div>

            {project.description.length > 0 ? (
              <ul className="mt-3 space-y-2 text-sm leading-6 text-[var(--resume-ink-600)]">
                {project.description.map((bullet) => (
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
          </article>
        ))}
      </div>
    </GlassCard>
  )
}
