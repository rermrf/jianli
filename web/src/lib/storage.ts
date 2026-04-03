import type { ResumeData } from '../types/resume'
import { apiFetch } from './api'
import { getAuthKey } from './auth'

export async function loadResumeDraft(): Promise<ResumeData> {
  return apiFetch<ResumeData>('/api/resume', {
    method: 'GET',
  })
}

export async function saveResumeDraft(resume: ResumeData): Promise<ResumeData> {
  return apiFetch<ResumeData>('/api/resume', {
    body: JSON.stringify(resume),
    headers: {
      'X-Auth-Key': getAuthKey() ?? '',
    },
    method: 'PUT',
  })
}

export function resetResumeDraft() {
  return undefined
}
