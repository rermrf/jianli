import type { VisitorRecord } from '../../types/visitors'
import { formatDuration, formatVisitTime } from '../../lib/format'
import { SectionCard } from '../common/SectionCard'

interface VisitorListProps {
  records: VisitorRecord[]
  showPdfColumn: boolean
}

function formatRegion(record: VisitorRecord): string {
  const joined = [record.country, record.region, record.city].filter(Boolean).join(' · ')
  return joined || '未知'
}

export function VisitorList({ records, showPdfColumn }: VisitorListProps) {
  return (
    <div className="space-y-3 md:hidden">
      {records.map((record) => (
        <SectionCard key={record.id} className="flex items-start gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-slate-100 text-slate-400">
            💻
          </div>
          <div className="flex-1 space-y-1">
            <p className="text-sm font-medium text-slate-800">
              {formatRegion(record)} · {record.browser} · {record.device}
            </p>
            <p className="font-mono text-xs text-slate-500">{record.ip}</p>
            <p className="text-xs text-slate-400">
              {formatVisitTime(record.visitTime)} · 停留 {formatDuration(record.duration)}
              {record.isp ? ` · ${record.isp}` : ''}
            </p>
            {showPdfColumn && record.pdfExported ? (
              <span className="inline-block rounded-full bg-emerald-50 px-2 py-0.5 text-xs text-emerald-600">
                已导出 PDF
              </span>
            ) : null}
          </div>
        </SectionCard>
      ))}
    </div>
  )
}
