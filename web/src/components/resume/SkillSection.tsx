import { SectionCard } from '../common/SectionCard'
import { Tag } from '../common/Tag'

interface SkillSectionProps {
  skills: string[]
  title?: string
}

export function SkillSection({
  skills,
  title = '个人技能',
}: SkillSectionProps) {
  return (
    <SectionCard className="space-y-4">
      <div className="flex items-center gap-2">
        <span className="text-brand-500">✦</span>
        <h2 className="text-lg font-semibold text-slate-900">{title}</h2>
      </div>
      <div className="flex flex-wrap gap-2">
        {skills.map((skill) => (
          <Tag key={skill}>{skill}</Tag>
        ))}
      </div>
    </SectionCard>
  )
}
