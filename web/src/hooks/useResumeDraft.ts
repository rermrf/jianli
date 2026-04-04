import { useEffect, useState } from 'react'
import { defaultResume } from '../data/mockResume'
import { resetResumeDraft, loadResumeDraft, saveResumeDraft } from '../lib/storage'
import type { ResumeData } from '../types/resume'

export function useResumeDraft() {
  const [draft, setDraft] = useState<ResumeData>(defaultResume)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let active = true

    async function fetchDraft() {
      try {
        const nextDraft = await loadResumeDraft()
        if (active) {
          setDraft(nextDraft)
          setError(null)
        }
      } catch (loadError) {
        if (active) {
          setDraft(defaultResume)
          setError(loadError instanceof Error ? loadError.message : '加载失败')
        }
      } finally {
        if (active) {
          setLoading(false)
        }
      }
    }

    void fetchDraft()

    return () => {
      active = false
    }
  }, [])

  function updateDraft(nextDraft: ResumeData) {
    setDraft(nextDraft)
  }

  async function saveDraft(nextDraft = draft) {
    const savedDraft = await saveResumeDraft(nextDraft)
    setDraft(savedDraft)
    setError(null)
    return savedDraft
  }

  function restoreDefaultDraft() {
    resetResumeDraft()
    setDraft(defaultResume)
  }

  return {
    draft,
    error,
    loading,
    setDraft: updateDraft,
    saveDraft,
    restoreDefaultDraft,
  }
}
