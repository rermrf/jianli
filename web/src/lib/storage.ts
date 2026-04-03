import { defaultResume } from '../data/mockResume'
import type { ResumeData } from '../types/resume'

const RESUME_DRAFT_KEY = 'resume:draft'

export function loadResumeDraft(): ResumeData {
  const storedDraft = localStorage.getItem(RESUME_DRAFT_KEY)

  if (!storedDraft) {
    return defaultResume
  }

  try {
    return JSON.parse(storedDraft) as ResumeData
  } catch {
    return defaultResume
  }
}

export function saveResumeDraft(resume: ResumeData) {
  localStorage.setItem(RESUME_DRAFT_KEY, JSON.stringify(resume))
}

export function resetResumeDraft() {
  localStorage.removeItem(RESUME_DRAFT_KEY)
}
