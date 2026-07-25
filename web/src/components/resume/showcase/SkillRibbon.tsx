import { GlassCard } from './GlassCard'

interface SkillRibbonProps {
  skills: string[]
  title?: string
}

export function SkillRibbon({
  skills,
  title = '个人技能',
}: SkillRibbonProps) {
  if (skills.length === 0) return null

  return (
    <GlassCard className="space-y-4">
      <h2 className="text-lg font-semibold text-[var(--resume-ink-900)]">
        {title}
      </h2>
      <div className="flex flex-wrap gap-2">
        {skills.map((skill) => (
          <span
            key={skill}
            className="rounded-full bg-[var(--resume-brand-50)] px-3 py-1 text-sm text-[var(--resume-ink-600)]"
          >
            {skill}
          </span>
        ))}
      </div>
    </GlassCard>
  )
}
