import { describe, expect, it } from 'vitest'
import { formatDuration, formatVisitorBrowser, formatVisitorDevice, formatVisitorRegion, formatVisitTime } from '../format'

describe('visitor formatting helpers', () => {
  it('shows the operating system and device category', () => {
    expect(formatVisitorDevice({ os: 'iOS', osVersion: '18.5', device: 'Mobile' })).toBe('iOS 18.5 · 手机')
    expect(formatVisitorDevice({ os: 'Android', osVersion: '14', device: 'Mobile' })).toBe('Android 14 · 手机')
    expect(formatVisitorDevice({ os: 'macOS', osVersion: '14.5', device: 'Desktop' })).toBe('macOS 14.5 · 电脑')
    expect(formatVisitorDevice({ os: 'Linux', device: 'Desktop' })).toBe('Linux · 电脑')
    expect(formatVisitorDevice({ os: 'Windows', osVersion: '10/11', device: 'Desktop' })).toBe('Windows 10/11 · 电脑')
    expect(formatVisitorDevice({ os: 'Windows', device: 'Windows' })).toBe('Windows')
  })

  it('shows available browser versions and supports older records without versions', () => {
    expect(formatVisitorBrowser({ browser: 'Edge', browserVersion: '146.0.3856.84' })).toBe('Edge 146.0.3856.84')
    expect(formatVisitorBrowser({ browser: 'Chrome' })).toBe('Chrome')
    expect(formatVisitorBrowser({ browser: 'Unknown' })).toBe('未知浏览器')
  })

  it('shows known geographic fields without unknown placeholders', () => {
    expect(formatVisitorRegion({ country: '中国', region: '浙江', city: '未知' })).toBe('中国 · 浙江')
    expect(formatVisitorRegion({ country: '中国', region: '0', city: '' })).toBe('中国')
    expect(formatVisitorRegion({ country: '', region: '', city: '未知' })).toBe('未知')
  })

  it('formats duration values for display', () => {
    expect(formatDuration(135)).toBe('2m 15s')
  })

  it('shows a friendlier label while visit duration is not yet available', () => {
    expect(formatDuration(0)).toBe('停留中')
  })

  it('formats ISO visit times for display', () => {
    expect(formatVisitTime('2026-04-04T02:46:53.634Z')).toBe('2026-04-04 10:46')
  })
})
