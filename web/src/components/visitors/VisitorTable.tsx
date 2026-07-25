import { formatDuration, formatVisitTime } from '../../lib/format'
import type { VisitorRecord } from '../../types/visitors'
import { SectionCard } from '../common/SectionCard'

interface VisitorTableProps {
  records: VisitorRecord[]
  showPdfColumn: boolean
}

function formatRegion(record: VisitorRecord): string {
  const joined = [record.country, record.region, record.city].filter(Boolean).join(' · ')
  return joined || '未知'
}

export function VisitorTable({ records, showPdfColumn }: VisitorTableProps) {
  return (
    <SectionCard className="hidden overflow-hidden md:block">
      <h2 className="mb-4 text-lg font-semibold text-slate-900">访客记录</h2>
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-slate-200 text-left text-sm">
          <thead className="text-slate-400">
            <tr>
              <th className="pb-3 font-medium">IP</th>
              <th className="pb-3 font-medium">地区</th>
              <th className="pb-3 font-medium">ISP</th>
              <th className="pb-3 font-medium">设备 / 浏览器</th>
              <th className="pb-3 font-medium">访问时间</th>
              <th className="pb-3 font-medium">停留时长</th>
              {showPdfColumn ? <th className="pb-3 font-medium">导出 PDF</th> : null}
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 text-slate-600">
            {records.map((record) => (
              <tr key={record.id}>
                <td className="py-3 font-mono text-xs text-slate-700">{record.ip}</td>
                <td className="py-3">{formatRegion(record)}</td>
                <td className="py-3">{record.isp || '—'}</td>
                <td className="py-3">
                  {record.device} · {record.browser}
                </td>
                <td className="py-3">{formatVisitTime(record.visitTime)}</td>
                <td className="py-3">{formatDuration(record.duration)}</td>
                {showPdfColumn ? (
                  <td className="py-3 text-center">
                    {record.pdfExported ? (
                      <span className="text-emerald-600">✓</span>
                    ) : (
                      <span className="text-slate-300">—</span>
                    )}
                  </td>
                ) : null}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </SectionCard>
  )
}
