import type {
  VisitorRange,
  VisitorRecord,
  VisitorStats,
  VisitorTrendPoint,
} from '../types/visitors'
import { apiFetch } from './api'
import { getAuthKey } from './auth'

const VISITOR_ID_KEY = 'jianli.visitorID'
const VISITOR_STARTED_AT_KEY = 'jianli.visitorStartedAt'

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

/**
 * For the trend endpoint we never send days=0 — backend rejects it. The
 * "all" range is mapped to 30 days so the chart stays readable.
 */
function rangeToTrendDays(range: VisitorRange) {
  const days = rangeToDays(range)
  return days > 0 ? days : 30
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

/**
 * Fetch the daily-bucketed visit trend for the selected range. Backend
 * returns only days that had visits; the frontend pads missing dates with
 * count=0 so the chart always shows N continuous bars.
 */
export async function fetchVisitorTrend(range: VisitorRange): Promise<VisitorTrendPoint[]> {
  const days = rangeToTrendDays(range)
  const points = await apiFetch<VisitorTrendPoint[]>(
    `/api/visitors/trend?days=${days}`,
    {
      headers: authHeaders(),
      method: 'GET',
    },
  )
  return padTrend(points, days)
}

/**
 * Pad missing dates in a trend series with count=0 so the chart renders a
 * continuous bar for every day in the requested window.
 */
export function padTrend(points: VisitorTrendPoint[], days: number): VisitorTrendPoint[] {
  const byDate = new Map(points.map((p) => [p.date, p.count]))
  const out: VisitorTrendPoint[] = []
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  for (let i = days - 1; i >= 0; i -= 1) {
    const d = new Date(today)
    d.setDate(today.getDate() - i)
    const key = formatDateKey(d)
    out.push({ date: key, count: byDate.get(key) ?? 0 })
  }
  return out
}

function formatDateKey(d: Date): string {
  const year = d.getFullYear()
  const month = `${d.getMonth() + 1}`.padStart(2, '0')
  const day = `${d.getDate()}`.padStart(2, '0')
  return `${year}-${month}-${day}`
}

export async function recordVisit() {
  const payload = {
    visitTime: new Date().toISOString(),
    duration: 0,
  }

  const response = await apiFetch<{ id: number }>('/api/visitors', {
    body: JSON.stringify(payload),
    method: 'POST',
  })

  return response.id
}

/**
 * Persist visitor ID + start timestamp to sessionStorage. Both must be
 * stored together so a navigation from `/` to `/print` preserves the
 * cumulative visit duration (the heartbeat then computes elapsed = now -
 * stored start, not now - this-page mount).
 */
export function storeVisitorID(id: number) {
  try {
    sessionStorage.setItem(VISITOR_ID_KEY, String(id))
    sessionStorage.setItem(VISITOR_STARTED_AT_KEY, String(Date.now()))
  } catch {
    // Incognito mode or storage quota — accept the loss; the visit is
    // still recorded server-side, just not linked across pages.
  }
}

export function readVisitorID(): number | null {
  try {
    const raw = sessionStorage.getItem(VISITOR_ID_KEY)
    if (!raw) return null
    const n = Number(raw)
    return Number.isFinite(n) && n > 0 ? n : null
  } catch {
    return null
  }
}

export function readVisitorStartedAt(): number | null {
  try {
    const raw = sessionStorage.getItem(VISITOR_STARTED_AT_KEY)
    if (!raw) return null
    const n = Number(raw)
    return Number.isFinite(n) && n > 0 ? n : null
  } catch {
    return null
  }
}

/**
 * Send a duration heartbeat. Backend overwrites the row's duration column,
 * so the value sent here must be the cumulative seconds since the visit
 * was first recorded (computed by the caller from the stored start time).
 */
export async function sendVisitDuration(visitorID: number, duration: number) {
  try {
    await apiFetch(`/api/visitors/${visitorID}`, {
      body: JSON.stringify({ duration }),
      method: 'POST',
    })
  } catch {
    // Heartbeats are advisory; the next tick will retry.
  }
}

/**
 * Mark the current visitor as having exported a PDF. Called after a
 * successful PDF download from PrintPage. No-op if no visitor ID is in
 * sessionStorage (e.g., user navigated directly to /print).
 */
export function markPDFExported(visitorID: number) {
  try {
    void fetch(`/api/visitors/${visitorID}/pdf-export`, { method: 'POST' })
  } catch {
    // Best-effort; admin will see "—" in the column if it fails.
  }
}
