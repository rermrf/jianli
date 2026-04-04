import { formatDuration, formatVisitTime, maskIp } from '../../lib/format'
import type { VisitorRecord } from '../../types/visitors'
import { SectionCard } from '../common/SectionCard'

interface VisitorTableProps {
  records: VisitorRecord[]
}

export function VisitorTable({ records }: VisitorTableProps) {
  return (
    <SectionCard className="hidden overflow-hidden md:block">
      <h2 className="mb-4 text-lg font-semibold text-slate-900">访客记录</h2>
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-slate-200 text-left text-sm">
          <thead className="text-slate-400">
            <tr>
              <th className="pb-3 font-medium">IP</th>
              <th className="pb-3 font-medium">地区</th>
              <th className="pb-3 font-medium">设备 / 浏览器</th>
              <th className="pb-3 font-medium">访问时间</th>
              <th className="pb-3 font-medium">停留时长</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 text-slate-600">
            {records.map((record) => (
              <tr key={record.id}>
                <td className="py-3">{maskIp(record.ip)}</td>
                <td className="py-3">{record.city}</td>
                <td className="py-3">
                  {record.device} · {record.browser}
                </td>
                <td className="py-3">{formatVisitTime(record.visitTime)}</td>
                <td className="py-3">{formatDuration(record.duration)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </SectionCard>
  )
}
