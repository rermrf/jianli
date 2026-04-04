import type {
  VisitorRange,
  VisitorRecord,
  VisitorStats,
  VisitorTrendPoint,
} from '../types/visitors'
import { apiFetch } from './api'
import { getAuthKey } from './auth'

function rangeToDays(range: VisitorRange) {
  switch (range) {
    case '7d':
      return 7
    case '30d':
      return 30
    default:
      return 0
  }
}

function authHeaders() {
  return {
    'X-Auth-Key': getAuthKey() ?? '',
  }
}

export async function fetchVisitorStats(range: VisitorRange) {
  return apiFetch<VisitorStats>(`/api/visitors/stats?days=${rangeToDays(range)}`, {
    headers: authHeaders(),
    method: 'GET',
  })
}

export async function fetchVisitors(range: VisitorRange) {
  return apiFetch<VisitorRecord[]>(
    `/api/visitors?days=${rangeToDays(range)}&page=1&limit=100`,
    {
      headers: authHeaders(),
      method: 'GET',
    },
  )
}

export async function recordVisit() {
  const payload = {
    ip: '',
    city: '未知',
    device: navigator.platform || 'Unknown',
    browser: navigator.userAgent,
    os: navigator.userAgent,
    visitTime: new Date().toISOString(),
    duration: 0,
  }

  const response = await apiFetch<{ id: number }>('/api/visitors', {
    body: JSON.stringify(payload),
    method: 'POST',
  })

  return response.id
}

export function sendVisitDuration(visitorID: number, duration: number) {
  const payload = JSON.stringify({ duration })

  return navigator.sendBeacon(
    `/api/visitors/${visitorID}`,
    new Blob([payload], { type: 'application/json' }),
  )
}

export function buildTrendPoints(records: VisitorRecord[]): VisitorTrendPoint[] {
  const points = [...records]
    .sort((left, right) => left.visitTime.localeCompare(right.visitTime))
    .map((record, index) => ({
      label: record.visitTime,
      value: index + 1,
    }))

  return points.length > 0 ? points : [{ label: '暂无数据', value: 0 }]
}
