export type VisitorRange = '7d' | '30d' | 'all'

export interface VisitorStats {
  totalVisits: number
  todayVisits: number
  uniqueVisitors: number
  averageDurationSeconds: number
}

export interface VisitorTrendPoint {
  label: string
  value: number
}

export interface VisitorRecord {
  id: number
  ip: string
  city: string
  device: string
  browser: string
  os: string
  visitTime: string
  durationSeconds: number
}

export type VisitorTrendByRange = Record<VisitorRange, VisitorTrendPoint[]>
