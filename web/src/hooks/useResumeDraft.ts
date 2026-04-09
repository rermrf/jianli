import { useEffect, useState } from 'react'
import type { Dispatch, SetStateAction } from 'react'
import { defaultResume } from '../data/mockResume'
import { resetResumeDraft, loadResumeDraft, saveResumeDraft } from '../lib/storage'
import type { ResumeData } from '../types/resume'
import type { SiteSettings } from '../types/siteSettings'

const defaultSiteSettings: SiteSettings = {
  allowPdfExport: true,
}

export function useResumeDraft() {
  const [draft, setDraft] = useState<ResumeData>(defaultResume)
  const [siteSettings, setSiteSettings] = useState<SiteSettings>(defaultSiteSettings)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let active = true

    async function fetchDraft() {
      try {
        const nextDraft = await loadResumeDraft()
        if (active) {
          setDraft(nextDraft.resume)
          setSiteSettings(nextDraft.siteSettings)
          setError(null)
        }
      } catch (loadError) {
        if (active) {
          setDraft(defaultResume)
          setSiteSettings(defaultSiteSettings)
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

  async function saveDraft(nextDraft = draft) {
    const savedDraft = await saveResumeDraft(nextDraft)
    setDraft(savedDraft)
    setError(null)
    return savedDraft
  }

  function restoreDefaultDraft() {
    resetResumeDraft()
    setDraft(defaultResume)
    setSiteSettings(defaultSiteSettings)
  }

  return {
    draft,
    error,
    loading,
    restoreDefaultDraft,
    saveDraft,
    setDraft: setDraft as Dispatch<SetStateAction<ResumeData>>,
    setSiteSettings,
    siteSettings,
  }
}
