import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import type { VisitorRecord } from '../../../types/visitors'
import { VisitorList } from '../VisitorList'
import { VisitorTable } from '../VisitorTable'

describe('visitor record details', () => {
  it('shows full IP, country/province/city, operator, system and browser in both layouts', () => {
    const record: VisitorRecord = {
      id: 1, ip: '39.182.123.234', country: '中国', region: '浙江省', city: '杭州市', isp: '移动',
      device: 'Mobile', os: 'iOS', osVersion: '18.5', browser: 'Chrome', browserVersion: '146.0.7680.90',
      visitTime: new Date().toISOString(), duration: 120, pdfExported: false,
    }
    render(<><VisitorTable records={[record]} showPdfColumn={false} /><VisitorList records={[record]} showPdfColumn={false} /></>)
    expect(screen.getAllByText('39.182.123.234')).toHaveLength(2)
    expect(screen.queryByText('39.182.***.***')).not.toBeInTheDocument()
    expect(screen.getAllByText('中国 · 浙江省 · 杭州市')).toHaveLength(2)
    expect(screen.getByText('移动')).toBeInTheDocument()
    expect(screen.getByText('iOS 18.5 · 手机')).toBeInTheDocument()
    expect(screen.getByText('Chrome 146.0.7680.90')).toBeInTheDocument()
    expect(screen.getByText('iOS 18.5 · 手机 · Chrome 146.0.7680.90')).toBeInTheDocument()
  })
})
