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
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60

  return `${minutes}m ${remainingSeconds}s`
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
