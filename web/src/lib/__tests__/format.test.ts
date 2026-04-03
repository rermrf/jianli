import { describe, expect, it } from 'vitest'
import { visitorTrendByRange } from '../../data/mockVisitors'
import {
  formatDuration,
  getVisitorRangeData,
  maskIp,
} from '../format'

describe('visitor formatting helpers', () => {
  it('masks visitor IP addresses for list display', () => {
    expect(maskIp('112.17.45.201')).toBe('112.17.***.***')
  })

  it('formats duration values for display', () => {
    expect(formatDuration(135)).toBe('2m 15s')
  })

  it('filters visitor trend data by selected range', () => {
    expect(getVisitorRangeData(visitorTrendByRange, '7d')).toEqual(
      visitorTrendByRange['7d'],
    )
    expect(getVisitorRangeData(visitorTrendByRange, '30d')).toEqual(
      visitorTrendByRange['30d'],
    )
  })
})
