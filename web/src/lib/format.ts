import type {
  VisitorRange,
  VisitorRecord,
  VisitorTrendByRange,
  VisitorTrendPoint,
} from '../types/visitors'

export function maskIp(ip: string) {
  const segments = ip.split('.')

  if (segments.length !== 4) {
    return ip
  }

  return `${segments[0]}.${segments[1]}.***.***`
}

export function formatDuration(seconds: number) {
  if (seconds <= 0) {
    return '停留中'
  }

  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60

  return `${minutes}m ${remainingSeconds}s`
}

export function formatVisitTime(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  const year = date.getFullYear()
  const month = `${date.getMonth() + 1}`.padStart(2, '0')
  const day = `${date.getDate()}`.padStart(2, '0')
  const hours = `${date.getHours()}`.padStart(2, '0')
  const minutes = `${date.getMinutes()}`.padStart(2, '0')

  return `${year}-${month}-${day} ${hours}:${minutes}`
}

export function getVisitorRangeData(
  trendByRange: VisitorTrendByRange,
  range: VisitorRange,
): VisitorTrendPoint[] {
  return trendByRange[range]
}

export function getVisitorRecordsByRange(
  recordsByRange: Record<VisitorRange, VisitorRecord[]>,
  range: VisitorRange,
) {
  return recordsByRange[range]
}
