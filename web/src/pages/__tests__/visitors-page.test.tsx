import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { loginWithKey } from '../../lib/auth'
import { VisitorsPage } from '../VisitorsPage'

function makeRecords(count: number) {
  // Use a real ISO timestamp from "today" so formatVisitTime renders a
  // stable, predictable string regardless of environment locale.
  const now = new Date().toISOString()
  return Array.from({ length: count }).map((_, index) => ({
    id: index + 1,
    ip: `112.17.45.${index + 10}`,
    country: '中国',
    region: '浙江',
    city: '杭州',
    isp: '电信',
    device: 'Windows',
    browser: 'Chrome',
    os: 'Windows',
    visitTime: now,
    duration: 120,
    pdfExported: false,
  }))
}

function makeStats(totalVisits: number) {
  return {
    totalVisits,
    todayVisits: 23,
    uniqueVisitors: 10,
    averageDurationSeconds: 120,
  }
}

describe('visitors page', () => {
  afterEach(() => {
    sessionStorage.clear()
    vi.restoreAllMocks()
  })

  it('switches between 7-day and 30-day visitor data views', async () => {
    loginWithKey('resume-key')
    vi.spyOn(globalThis, 'fetch').mockImplementation(async (input) => {
      const url = String(input)

      if (url.includes('/api/visitors/stats?days=7')) {
        return new Response(
          JSON.stringify({ code: 0, data: makeStats(128) }),
          { status: 200, headers: { 'Content-Type': 'application/json' } },
        )
      }
      if (url.includes('/api/visitors?days=7')) {
        return new Response(
          JSON.stringify({ code: 0, data: makeRecords(1) }),
          { status: 200, headers: { 'Content-Type': 'application/json' } },
        )
      }
      if (url.includes('/api/visitors/stats?days=30')) {
        return new Response(
          JSON.stringify({ code: 0, data: makeStats(462) }),
          { status: 200, headers: { 'Content-Type': 'application/json' } },
        )
      }
      if (url.includes('/api/visitors?days=30')) {
        return new Response(
          JSON.stringify({ code: 0, data: makeRecords(2) }),
          { status: 200, headers: { 'Content-Type': 'application/json' } },
        )
      }
      if (url.includes('/api/visitors/trend')) {
        // Empty trend payload — page pads to N zero-count bars.
        return new Response(
          JSON.stringify({ code: 0, data: [] }),
          { status: 200, headers: { 'Content-Type': 'application/json' } },
        )
      }
      return new Response(JSON.stringify({ code: 0, data: [] }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    })

    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <VisitorsPage />
      </MemoryRouter>,
    )

    expect(await screen.findByText('128')).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: '30天' }))

    await waitFor(() => {
      expect(screen.getByText('462')).toBeInTheDocument()
    })
  })
})
