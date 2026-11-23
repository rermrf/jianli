import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { fetchVisitorTrend, padTrend } from '../visitors'

describe('daily visitor trends', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2026-11-23T16:04:00Z'))
  })
  afterEach(() => {
    vi.useRealTimers()
    vi.restoreAllMocks()
  })

  it('uses Beijing calendar days and fills missing dates with both zero metrics', () => {
    const points = padTrend([
      { date: '2026-11-18', count: 4, uniqueVisitors: 2 },
      { date: '2026-11-24', count: 3, uniqueVisitors: 1 },
    ], 7)
    expect(points).toHaveLength(7)
    expect(points[0]).toEqual({ date: '2026-11-18', count: 4, uniqueVisitors: 2 })
    expect(points[5]).toEqual({ date: '2026-11-23', count: 0, uniqueVisitors: 0 })
    expect(points[6]).toEqual({ date: '2026-11-24', count: 3, uniqueVisitors: 1 })
  })

  it('requests all history instead of silently limiting it to 30 days', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(JSON.stringify({
      code: 0,
      data: [{ date: '2026-10-20', count: 1, uniqueVisitors: 1 }, { date: '2026-11-24', count: 2, uniqueVisitors: 1 }],
    }), { status: 200, headers: { 'Content-Type': 'application/json' } }))
    const points = await fetchVisitorTrend('all')
    expect(fetchMock).toHaveBeenCalledWith('/api/visitors/trend?days=0', expect.anything())
    expect(points).toHaveLength(36)
    expect(points[0].date).toBe('2026-10-20')
    expect(points[35].date).toBe('2026-11-24')
  })

  it('returns an empty all-history series when no visits exist', () => {
    expect(padTrend([], 0)).toEqual([])
  })
})
