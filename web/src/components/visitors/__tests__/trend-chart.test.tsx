import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'
import { TrendChart } from '../TrendChart'

const points = [
  { date: '2026-11-22', count: 4, uniqueVisitors: 2 },
  { date: '2026-11-23', count: 0, uniqueVisitors: 0 },
  { date: '2026-11-24', count: 3, uniqueVisitors: 1 },
]

describe('daily visitor chart', () => {
  it('shows daily people counts including zero days and switches to visit counts', async () => {
    const user = userEvent.setup()
    render(<TrendChart points={points} />)
    expect(screen.getByRole('button', { name: '访客人数' })).toHaveAttribute('aria-pressed', 'true')
    const firstDay = screen.getByRole('button', { name: '2026-11-22：2 人，4 次访问' })
    expect(within(firstDay).getByText('2')).toBeInTheDocument()
    expect(within(screen.getByRole('button', { name: '2026-11-23：0 人，0 次访问' })).getByText('0')).toBeInTheDocument()
    expect(screen.getByText('11/24')).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: '访问次数' }))
    expect(within(firstDay).getByText('4')).toBeInTheDocument()
    expect(screen.getByRole('group', { name: '每日访问次数柱状图' })).toBeInTheDocument()
    await user.click(firstDay)
    expect(screen.getByText('2026-11-22：2 位访客，4 次访问')).toBeInTheDocument()
  })

  it('provides daily detail rows sorted from newest to oldest', async () => {
    const user = userEvent.setup()
    render(<TrendChart points={points} />)
    await user.click(screen.getByText('查看每日明细 · 3 天'))
    const table = screen.getByRole('table', { name: '每日访客人数和访问次数' })
    const rows = within(table).getAllByRole('row')
    expect(rows).toHaveLength(4)
    expect(within(rows[1]).getByText('2026-11-24')).toBeInTheDocument()
    expect(within(rows[1]).getByText('1 人')).toBeInTheDocument()
    expect(within(rows[1]).getByText('3 次')).toBeInTheDocument()
    expect(within(rows[2]).getByText('0 人')).toBeInTheDocument()
  })

  it('shows an empty state without inventing visits', () => {
    render(<TrendChart points={[]} />)
    expect(screen.getByText('暂无访问记录')).toBeInTheDocument()
  })
})
