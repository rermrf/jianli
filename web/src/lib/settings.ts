import type { SiteSettings } from '../types/siteSettings'
import { apiFetch } from './api'
import { getAuthKey } from './auth'

export async function updateSiteSettings(settings: SiteSettings): Promise<SiteSettings> {
  return apiFetch<SiteSettings>('/api/settings', {
    body: JSON.stringify(settings),
    headers: {
      'X-Auth-Key': getAuthKey() ?? '',
    },
    method: 'PUT',
  })
}
