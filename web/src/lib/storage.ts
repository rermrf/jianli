import type { ResumeData } from '../types/resume'
import type { PublicResumePayload } from '../types/siteSettings'
import { apiFetch } from './api'
import { getAuthKey } from './auth'

function isPublicResumePayload(
  payload: ResumeData | PublicResumePayload<ResumeData>,
): payload is PublicResumePayload<ResumeData> {
  return typeof payload === 'object' && payload !== null && 'resume' in payload && 'siteSettings' in payload
}

export async function loadResumeDraft(): Promise<PublicResumePayload<ResumeData>> {
  const payload = await apiFetch<ResumeData | PublicResumePayload<ResumeData>>('/api/resume', {
    method: 'GET',
  })

  if (isPublicResumePayload(payload)) {
    return payload
  }

  return {
    resume: payload,
    siteSettings: {
      allowPdfExport: true,
    },
  }
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
