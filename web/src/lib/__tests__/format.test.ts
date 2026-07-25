import { describe, expect, it } from 'vitest'
import { formatDuration, formatVisitTime } from '../format'

describe('visitor formatting helpers', () => {
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
