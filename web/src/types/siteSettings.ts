export interface SiteSettings {
  allowPdfExport: boolean
}

export interface PublicResumePayload<TResume> {
  resume: TResume
  siteSettings: SiteSettings
}
