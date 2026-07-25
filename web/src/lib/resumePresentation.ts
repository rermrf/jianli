import type { ResumeData } from '../types/resume'

export function buildTagline(resume: ResumeData): string {
  const { title, experience } = resume.profile
  const position = resume.jobIntention.position?.trim()
  const base = [title, experience ? `${experience}经验` : '']
    .filter(Boolean)
    .join(' · ')
  if (!position) return base
  return `${base} · 意向 ${position}`
}

export function monogramFromName(name: string): string {
  const trimmed = name.trim()
  return trimmed ? [...trimmed][0] : '?'
}
