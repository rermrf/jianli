import { useState } from 'react'
import { defaultResume } from '../data/mockResume'
import { loadResumeDraft, resetResumeDraft, saveResumeDraft } from '../lib/storage'
import type { ResumeData } from '../types/resume'

export function useResumeDraft() {
  const [draft, setDraft] = useState<ResumeData>(() => loadResumeDraft())

  function updateDraft(nextDraft: ResumeData) {
    setDraft(nextDraft)
  }

  function saveDraft(nextDraft = draft) {
    saveResumeDraft(nextDraft)
    setDraft(nextDraft)
  }

  function restoreDefaultDraft() {
    resetResumeDraft()
    setDraft(defaultResume)
  }

  return {
    draft,
    setDraft: updateDraft,
    saveDraft,
    restoreDefaultDraft,
  }
}
