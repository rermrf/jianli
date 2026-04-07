import type { ResumeData } from '../../types/resume'
import { ProfileCard } from './ProfileCard'
import { SkillSection } from './SkillSection'
import { TimelineSection } from './TimelineSection'

interface ResumeDesktopLayoutProps {
  resume: ResumeData
}

export function ResumeDesktopLayout({ resume }: ResumeDesktopLayoutProps) {
  return (
    <div className="hidden gap-6 md:grid md:grid-cols-[320px_minmax(0,1fr)] xl:grid-cols-[360px_minmax(0,1fr)]">
      <aside className="space-y-6">
        <ProfileCard desiredCities={resume.jobIntention.cities} profile={resume.profile} />
        <SkillSection skills={resume.skills} />
        <TimelineSection
          items={resume.education.map((item) => ({
            title: item.school,
            bullets: [`${item.major} · ${item.degree}`],
          }))}
          title="教育经历"
        />
        <TimelineSection
          items={resume.awards.map((item) => ({
            title: item.title,
            meta: item.date,
          }))}
          title="荣誉奖项"
        />
      </aside>
      <div className="space-y-6">
        <TimelineSection
          items={resume.workExperience.map((item) => ({
            title: `${item.company} · ${item.role}`,
            meta: `${item.startDate} - ${item.endDate}`,
            bullets: item.description,
          }))}
          title="工作经历"
        />
        <TimelineSection
          items={resume.projects.map((item) => ({
            title: item.name,
            meta: `${item.startDate} - ${item.endDate}`,
            bullets: item.description,
          }))}
          title="项目经历"
        />
      </div>
    </div>
  )
}
