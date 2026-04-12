import type { ResumeData } from '../../types/resume'
import { ProfileCard } from './ProfileCard'
import { SkillSection } from './SkillSection'
import { TimelineSection } from './TimelineSection'

interface ResumeMobileLayoutProps {
  resume: ResumeData
}

export function ResumeMobileLayout({ resume }: ResumeMobileLayoutProps) {
  return (
    <div className="space-y-4 md:hidden">
      <ProfileCard desiredCities={resume.jobIntention.cities} profile={resume.profile} />
      <SkillSection skills={resume.skills} />
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
          titleHref: item.url,
          meta: `${item.startDate} - ${item.endDate}`,
          bullets: item.description,
        }))}
        title="项目经历"
      />
      <TimelineSection
        items={resume.education.map((item) => ({
          title: `${item.school} · ${item.major}`,
          bullets: [item.degree],
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
    </div>
  )
}
