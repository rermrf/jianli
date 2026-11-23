import type { VisitorRecord } from '../types/visitors'

export function formatVisitorRegion(record: Pick<VisitorRecord, 'country' | 'region' | 'city'>) {
  const parts = [record.country, record.region, record.city]
    .map((value) => value?.trim())
    .filter((value) => value && value !== '未知' && value !== '0' && value.toLowerCase() !== 'unknown')
  return parts.join(' · ') || '未知'
}

export function formatVisitorDevice(record: Pick<VisitorRecord, 'device' | 'os' | 'osVersion'>) {
  const devices: Record<string, string> = { Desktop: '电脑', Mobile: '手机', Tablet: '平板', Bot: '爬虫' }
  const system = record.os && record.os !== 'Unknown' ? [record.os, record.osVersion].filter(Boolean).join(' ') : ''
  const device = devices[record.device] || (record.device !== 'Unknown' ? record.device : '')
  return [...new Set([system, device].filter(Boolean))].join(' · ') || '未知设备'
}

export function formatVisitorBrowser(record: Pick<VisitorRecord, 'browser' | 'browserVersion'>) {
  if (!record.browser || record.browser === 'Unknown') return '未知浏览器'
  return [record.browser, record.browserVersion].filter(Boolean).join(' ')
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
