import { useEffect, useState } from 'react'
import { visitorStatsByRange } from '../data/mockVisitors'
import type { VisitorRange, VisitorRecord, VisitorStats } from '../types/visitors'
import { Button } from '../components/common/Button'
import { AppShell } from '../components/layout/AppShell'
import { TopNav } from '../components/layout/TopNav'
import { StatsCards } from '../components/visitors/StatsCards'
import { TrendChart } from '../components/visitors/TrendChart'
import { VisitorList } from '../components/visitors/VisitorList'
import { VisitorTable } from '../components/visitors/VisitorTable'
import { buildTrendPoints, fetchVisitors, fetchVisitorStats } from '../lib/visitors'

const rangeOptions: { label: string; value: VisitorRange }[] = [
  { label: '7天', value: '7d' },
  { label: '30天', value: '30d' },
  { label: '全部', value: 'all' },
]

export function VisitorsPage() {
  const [range, setRange] = useState<VisitorRange>('7d')
  const [stats, setStats] = useState<VisitorStats>(visitorStatsByRange['7d'])
  const [records, setRecords] = useState<VisitorRecord[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let active = true

    async function loadVisitors() {
      setLoading(true)
      const [nextStats, nextRecords] = await Promise.all([
        fetchVisitorStats(range),
        fetchVisitors(range),
      ])

      if (active) {
        setStats(nextStats)
        setRecords(nextRecords)
        setLoading(false)
      }
    }

    void loadVisitors()

    return () => {
      active = false
    }
  }, [range])

  return (
    <AppShell contentClassName="space-y-6">
      <TopNav />
      <div className="flex items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-slate-900">访客统计</h1>
          <p className="mt-2 text-sm text-slate-500">查看最近访问趋势与来源设备。</p>
        </div>
        <div className="flex items-center gap-2">
          {rangeOptions.map((option) => (
            <Button
              key={option.value}
              onClick={() => setRange(option.value)}
              variant={range === option.value ? 'primary' : 'secondary'}
            >
              {option.label}
            </Button>
          ))}
        </div>
      </div>

      {loading ? (
        <div className="rounded-3xl border border-slate-200 bg-white p-6 text-sm text-slate-500 shadow-[var(--shadow-card)]">
          加载中...
        </div>
      ) : (
        <>
          <StatsCards stats={stats} />
          <TrendChart points={buildTrendPoints(records)} />
          <VisitorList records={records} />
          <VisitorTable records={records} />
        </>
      )}
    </AppShell>
  )
}
