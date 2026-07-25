import type { VisitorRange, VisitorStats } from '../types/visitors'

// Initial seed values used as placeholders before the real API data
// arrives. Only the stats are still consumed at runtime; trend and record
// fixtures were removed when the page switched to the live trend endpoint.
export const visitorStatsByRange: Record<VisitorRange, VisitorStats> = {
  '7d': {
    totalVisits: 0,
    todayVisits: 0,
    uniqueVisitors: 0,
    averageDurationSeconds: 0,
  },
  '30d': {
    totalVisits: 0,
    todayVisits: 0,
    uniqueVisitors: 0,
    averageDurationSeconds: 0,
  },
  all: {
    totalVisits: 0,
    todayVisits: 0,
    uniqueVisitors: 0,
    averageDurationSeconds: 0,
  },
}
