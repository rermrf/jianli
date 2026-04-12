import { SectionCard } from '../common/SectionCard'

interface TimelineItem {
  bullets?: string[]
  meta?: string
  title: string
  titleHref?: string
}

interface TimelineSectionProps {
  items: TimelineItem[]
  title: string
}

export function TimelineSection({ items, title }: TimelineSectionProps) {
  return (
    <SectionCard className="space-y-4">
      <div className="flex items-center gap-2">
        <span className="text-brand-500">◌</span>
        <h2 className="text-lg font-semibold text-slate-900">{title}</h2>
      </div>
      <div className="space-y-4">
        {items.map((item) => (
          <article
            key={`${item.title}-${item.meta ?? 'meta'}`}
            className="rounded-2xl border border-slate-200 bg-slate-50 p-4"
          >
            <div className="flex flex-col gap-1 md:flex-row md:items-start md:justify-between">
              {item.titleHref ? (
                <a
                  className="text-sm font-semibold text-brand-600 hover:text-brand-700 hover:underline"
                  href={item.titleHref}
                  rel="noreferrer"
                  target="_blank"
                >
                  {item.title}
                </a>
              ) : (
                <h3 className="text-sm font-semibold text-slate-900">{item.title}</h3>
              )}
              {item.meta ? (
                <p className="text-xs text-slate-400 md:text-sm">{item.meta}</p>
              ) : null}
            </div>
            {item.bullets?.length ? (
              <ul className="mt-3 space-y-2 text-sm leading-6 text-slate-500">
                {item.bullets.map((bullet) => (
                  <li key={bullet} className="flex gap-2">
                    <span className="mt-2 h-1.5 w-1.5 rounded-full bg-brand-300" />
                    <span>{bullet}</span>
                  </li>
                ))}
              </ul>
            ) : null}
          </article>
        ))}
      </div>
    </SectionCard>
  )
}
