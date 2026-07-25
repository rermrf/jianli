import { buildTagline } from '../../../lib/resumePresentation'
import type { ResumeData } from '../../../types/resume'
import { ContactFooter } from './ContactFooter'
import { EducationAwardsRow } from './EducationAwardsRow'
import { ProjectShowcase } from './ProjectShowcase'
import { ResumeHero } from './ResumeHero'
import { ResumeTopBar } from './ResumeTopBar'
import { SkillRibbon } from './SkillRibbon'
import { WorkTimeline } from './WorkTimeline'

interface ResumeShowcaseProps {
  resume: ResumeData
  showPdfExport: boolean
}

export function ResumeShowcase({ resume, showPdfExport }: ResumeShowcaseProps) {
  const tagline = buildTagline(resume)

  return (
    <div className="space-y-6">
      <ResumeTopBar name={resume.profile.name} showPdfExport={showPdfExport} />
      <ResumeHero
        profile={resume.profile}
        desiredCities={resume.jobIntention.cities}
        tagline={tagline}
      />
      <div className="grid gap-6 md:grid-cols-[1.1fr_0.9fr]">
        <ProjectShowcase projects={resume.projects} />
        <WorkTimeline items={resume.workExperience} />
      </div>
      <SkillRibbon skills={resume.skills} />
      <EducationAwardsRow education={resume.education} awards={resume.awards} />
      <ContactFooter
        profile={resume.profile}
        jobIntention={resume.jobIntention}
        showPdfExport={showPdfExport}
      />
    </div>
  )
}
