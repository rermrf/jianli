import { useState } from 'react'
import {
  visitorRecordsByRange,
  visitorStatsByRange,
  visitorTrendByRange,
} from '../data/mockVisitors'
import {
  getVisitorRangeData,
  getVisitorRecordsByRange,
} from '../lib/format'
import type { VisitorRange } from '../types/visitors'
import { Button } from '../components/common/Button'
import { AppShell } from '../components/layout/AppShell'
import { TopNav } from '../components/layout/TopNav'
import { StatsCards } from '../components/visitors/StatsCards'
import { TrendChart } from '../components/visitors/TrendChart'
import { VisitorList } from '../components/visitors/VisitorList'
import { VisitorTable } from '../components/visitors/VisitorTable'

const rangeOptions: { label: string; value: VisitorRange }[] = [
  { label: '7天', value: '7d' },
  { label: '30天', value: '30d' },
  { label: '全部', value: 'all' },
]

export function VisitorsPage() {
  const [range, setRange] = useState<VisitorRange>('7d')

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

      <StatsCards stats={visitorStatsByRange[range]} />
      <TrendChart points={getVisitorRangeData(visitorTrendByRange, range)} />
      <VisitorList records={getVisitorRecordsByRange(visitorRecordsByRange, range)} />
      <VisitorTable records={getVisitorRecordsByRange(visitorRecordsByRange, range)} />
    </AppShell>
  )
}
