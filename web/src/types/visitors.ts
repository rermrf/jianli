export type VisitorRange = '7d' | '30d' | 'all'

export interface VisitorStats {
  totalVisits: number
  todayVisits: number
  uniqueVisitors: number
  averageDurationSeconds: number
}

export interface VisitorTrendPoint {
  date: string // YYYY-MM-DD
  count: number
}

export interface VisitorRecord {
  id: number
  ip: string
  country: string
  region: string
  city: string
  isp: string
  device: string
  browser: string
  os: string
  visitTime: string
  duration: number
  pdfExported: boolean
}
