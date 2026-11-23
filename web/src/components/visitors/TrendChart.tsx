import { useState } from 'react'
import type { VisitorTrendPoint } from '../../types/visitors'
import { SectionCard } from '../common/SectionCard'

interface TrendChartProps {
  points: VisitorTrendPoint[]
}

export function TrendChart({ points }: TrendChartProps) {
  const [metric, setMetric] = useState<'people' | 'visits'>('people')
  const [activeDate, setActiveDate] = useState<string | null>(null)
  const label = metric === 'people' ? '访客人数' : '访问次数'
  const unit = metric === 'people' ? '人' : '次'
  const valueOf = (point: VisitorTrendPoint) => metric === 'people' ? point.uniqueVisitors : point.count
  const maxValue = Math.max(0, ...points.map(valueOf))
  const step = Math.max(1, Math.ceil(maxValue / 4))
  const axisMax = step * 4
  const active = points.find((point) => point.date === activeDate) ?? points.at(-1)
  const totalVisits = points.reduce((sum, point) => sum + point.count, 0)

  return (
    <SectionCard className="space-y-5">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">每日访问统计</h2>
          <p className="mt-1 text-xs text-slate-500">按北京时间统计 · 访客人数按 IP 去重</p>
        </div>
        <div className="flex gap-1 rounded-xl bg-slate-100 p-1" aria-label="统计指标">
          {([{ value: 'people', label: '访客人数' }, { value: 'visits', label: '访问次数' }] as const).map((option) => (
            <button
              key={option.value}
              type="button"
              aria-pressed={metric === option.value}
              onClick={() => setMetric(option.value)}
              className={`rounded-lg px-3 py-2 text-sm transition-colors ${metric === option.value ? 'bg-white font-medium text-brand-600 shadow-sm' : 'text-slate-500 hover:text-slate-800'}`}
            >{option.label}</button>
          ))}
        </div>
      </div>

      {points.length === 0 ? (
        <p className="rounded-2xl bg-slate-50 p-6 text-sm text-slate-500">暂无访问记录</p>
      ) : (
        <>
          <div className="rounded-2xl bg-slate-50 p-4">
            <div className="flex flex-wrap items-center justify-between gap-2 text-xs text-slate-500">
              <span>{points[0].date} 至 {points.at(-1)?.date} · 共 {points.length} 天</span>
              <span>本期 {totalVisits} 次访问</span>
            </div>
            <div className="mt-4 flex gap-3" role="group" aria-label={`每日${label}柱状图`}>
              <div className="relative mt-8 h-[200px] w-7 shrink-0 text-right text-[11px] text-slate-400" aria-hidden="true">
                {[4, 3, 2, 1, 0].map((tick, index) => (
                  <span key={tick} className="absolute right-0 -translate-y-1/2" style={{ top: `${index * 25}%` }}>{tick * step}</span>
                ))}
              </div>
              <div className="min-w-0 flex-1 overflow-x-auto pb-2">
                <div style={{ minWidth: Math.max(420, points.length * 52) }}>
                  <div className="relative pt-8">
                    <div className="pointer-events-none absolute inset-x-0 top-8 h-[200px]" aria-hidden="true">
                      {[0, 1, 2, 3, 4].map((tick) => (
                        <div key={tick} className="absolute inset-x-0 border-t border-slate-200/80" style={{ top: `${tick * 25}%` }} />
                      ))}
                    </div>
                    <div className="relative grid h-[200px] gap-2" style={{ gridTemplateColumns: `repeat(${points.length}, minmax(0, 1fr))` }}>
                      {points.map((point) => (
                        <button
                          key={point.date}
                          type="button"
                          aria-label={`${point.date}：${point.uniqueVisitors} 人，${point.count} 次访问`}
                          title={`${point.date}：${point.uniqueVisitors} 位访客，${point.count} 次访问`}
                          onMouseEnter={() => setActiveDate(point.date)}
                          onFocus={() => setActiveDate(point.date)}
                          onClick={() => setActiveDate(point.date)}
                          className="flex h-full min-w-0 flex-col items-center justify-end rounded-t-lg focus-visible:outline-2 focus-visible:outline-brand-500"
                        >
                          <span className={`mb-1 shrink-0 text-xs font-semibold ${point.date === active?.date ? 'text-brand-600' : 'text-slate-600'}`}>{valueOf(point)}</span>
                          <span
                            className={`w-3/5 shrink-0 rounded-t-md transition-colors ${point.date === active?.date ? 'bg-brand-600' : 'bg-brand-500/80 hover:bg-brand-600'}`}
                            style={{ height: `${valueOf(point) / axisMax * 200}px` }}
                            aria-hidden="true"
                          />
                        </button>
                      ))}
                    </div>
                  </div>
                  <div className="mt-3 grid gap-2 text-center text-xs text-slate-500" style={{ gridTemplateColumns: `repeat(${points.length}, minmax(0, 1fr))` }}>
                    {points.map((point, index) => (
                      <span key={point.date} className={index === points.length - 1 ? 'font-medium text-brand-600' : ''}>
                        {point.date.slice(5).replace('-', '/')}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            </div>
            <div className="mt-4 flex flex-wrap items-center justify-between gap-2 border-t border-slate-200 pt-3 text-xs text-slate-500">
              <span>纵轴：{label}（{unit}）· 横轴：日期</span>
              {active && <span aria-live="polite" className="font-medium text-slate-700">{active.date}：{active.uniqueVisitors} 位访客，{active.count} 次访问</span>}
            </div>
            {totalVisits === 0 && <p className="mt-3 text-center text-sm text-slate-400">所选时间内暂无访问</p>}
          </div>

          <details className="rounded-xl border border-slate-200">
            <summary className="cursor-pointer px-4 py-3 text-sm font-medium text-slate-700">查看每日明细 · {points.length} 天</summary>
            <div className="max-h-80 overflow-auto border-t border-slate-100 px-4">
              <table className="w-full text-left text-sm">
                <caption className="sr-only">每日访客人数和访问次数</caption>
                <thead className="sticky top-0 bg-white text-xs text-slate-500">
                  <tr><th className="py-3 font-medium">日期</th><th className="py-3 text-right font-medium">访客人数</th><th className="py-3 text-right font-medium">访问次数</th></tr>
                </thead>
                <tbody className="divide-y divide-slate-100 text-slate-700">
                  {[...points].reverse().map((point) => (
                    <tr key={point.date}>
                      <td className="whitespace-nowrap py-3">{point.date}</td>
                      <td className="py-3 text-right tabular-nums">{point.uniqueVisitors} 人</td>
                      <td className="py-3 text-right tabular-nums">{point.count} 次</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </details>
        </>
      )}
    </SectionCard>
  )
}
