import type { ResumeData } from '../../types/resume'

interface PrintResumeProps {
  resume: ResumeData
}

/**
 * PrintResume renders a print-optimized, single-column A4 layout.
 *
 * Section order: Profile → Skills → Work → Projects → Education → Awards.
 * Empty sections (zero-length arrays) are hidden entirely, including their
 * heading. The component carries `print-item` / `print-section-title` class
 * hooks that the embedded `@media print` block uses to prevent items from
 * splitting across pages and to keep section headings with their first item.
 *
 * The `id="print-root"` on the root article is the sentinel that the backend
 * chromedp PDF exporter waits for before checking `window.__printReady`.
 */
export function PrintResume({ resume }: PrintResumeProps) {
  const { profile, skills, workExperience, projects, education, awards, jobIntention } = resume

  const profileFacts = [
    profile.age ? `${profile.age}岁` : '',
    profile.gender,
    profile.education,
    profile.experience,
    profile.hometown ? `籍贯：${profile.hometown}` : '',
  ].filter(Boolean)

  const contactLine = [
    profile.location ? `所在地：${profile.location}` : '',
    profile.phone ? `手机号：${profile.phone}` : '',
    profile.email ? `邮箱：${profile.email}` : '',
  ]
    .filter(Boolean)
    .join(' · ')

  return (
    <>
      <PrintCSS />
      <article
        className="mx-auto max-w-[210mm] bg-white p-10 text-slate-900 print:max-w-none print:p-0"
        id="print-root"
      >
        <PrintHeader
          contactLine={contactLine}
          desiredCities={jobIntention.cities}
          profileFacts={profileFacts}
          profile={profile}
        />

        {skills.length > 0 ? (
          <PrintSection title="个人技能">
            <p className="text-[11px] leading-6 text-slate-700">{skills.join(' · ')}</p>
          </PrintSection>
        ) : null}

        {workExperience.length > 0 ? (
          <PrintSection title="工作经历">
            {workExperience.map((item) => (
              <PrintItem
                key={`${item.company}-${item.role}-${item.startDate}`}
                date={`${item.startDate} - ${item.endDate}`}
                subtitle={item.role}
                title={item.company}
              >
                <BulletList items={item.description} />
              </PrintItem>
            ))}
          </PrintSection>
        ) : null}

        {projects.length > 0 ? (
          <PrintSection title="项目经历">
            {projects.map((item) => (
              <PrintItem
                key={`${item.name}-${item.startDate}`}
                date={`${item.startDate} - ${item.endDate}`}
                title={item.name}
                titleHref={item.url}
              >
                <BulletList items={item.description} />
              </PrintItem>
            ))}
          </PrintSection>
        ) : null}

        {education.length > 0 ? (
          <PrintSection title="教育经历">
            {education.map((item) => (
              <PrintItem
                key={`${item.school}-${item.startDate}`}
                date={`${item.startDate} - ${item.endDate}`}
                subtitle={`${item.major} · ${item.degree}`}
                title={item.school}
              />
            ))}
          </PrintSection>
        ) : null}

        {awards.length > 0 ? (
          <PrintSection title="荣誉奖项">
            {awards.map((item) => (
              <PrintItem
                key={`${item.title}-${item.date}`}
                date={item.date}
                title={item.title}
              />
            ))}
          </PrintSection>
        ) : null}
      </article>
    </>
  )
}

function PrintHeader({
  contactLine,
  desiredCities,
  profileFacts,
  profile,
}: {
  contactLine: string
  desiredCities: string[]
  profileFacts: string[]
  profile: ResumeData['profile']
}) {
  return (
    <header className="flex items-start gap-5 border-b border-slate-200 pb-5">
      {profile.avatarUrl ? (
        <img
          alt={`${profile.name}头像`}
          className="h-16 w-16 flex-shrink-0 rounded-full object-cover"
          src={profile.avatarUrl}
        />
      ) : (
        <div className="h-16 w-16 flex-shrink-0 rounded-full bg-slate-100" />
      )}
      <div className="min-w-0 flex-1">
        <h1 className="text-[24px] font-semibold tracking-[0.02em] leading-tight text-slate-900">
          {profile.name}
        </h1>
        {profile.title ? (
          <p className="mt-1 text-[14px] font-medium text-slate-700">{profile.title}</p>
        ) : null}
        {contactLine ? (
          <p className="mt-2 text-[11px] leading-5 text-slate-500">{contactLine}</p>
        ) : null}
        {profileFacts.length > 0 ? (
          <p className="mt-1 text-[11px] leading-5 text-slate-500">{profileFacts.join(' / ')}</p>
        ) : null}
        {desiredCities.length > 0 ? (
          <p className="mt-1 text-[11px] leading-5 text-slate-500">
            {`意向城市：${desiredCities.join(' / ')}`}
          </p>
        ) : null}
      </div>
    </header>
  )
}

function PrintSection({ children, title }: { children: React.ReactNode; title: string }) {
  return (
    <section className="mt-4">
      <h2
        className="print-section-title mb-2 flex items-center gap-2 text-[14px] font-bold text-slate-900"
      >
        <span aria-hidden className="inline-block h-4 w-1 rounded-sm bg-brand-500" />
        {title}
      </h2>
      <div className="space-y-2">{children}</div>
    </section>
  )
}

function PrintItem({
  children,
  date,
  subtitle,
  title,
  titleHref,
}: {
  children?: React.ReactNode
  date: string
  subtitle?: string
  title: string
  titleHref?: string
}) {
  return (
    <div className="print-item">
      <div className="flex items-baseline justify-between gap-3">
        {titleHref ? (
          <a
            className="text-[14px] font-semibold text-brand-600"
            href={titleHref}
            rel="noreferrer"
            target="_blank"
          >
            {title}
            <span aria-hidden className="ml-0.5 align-super text-[9px]">
              ↗
            </span>
          </a>
        ) : (
          <strong className="text-[14px] font-semibold text-slate-900">{title}</strong>
        )}
        {date ? (
          <span className="flex-shrink-0 text-[11px] text-slate-400">{date}</span>
        ) : null}
      </div>
      {subtitle ? (
        <p className="mt-0.5 text-[11px] text-slate-600">{subtitle}</p>
      ) : null}
      {children ? <div className="mt-1">{children}</div> : null}
    </div>
  )
}

function BulletList({ items }: { items: string[] }) {
  if (items.length === 0) {
    return null
  }
  return (
    <ul className="space-y-1 text-[11px] leading-6 text-slate-700">
      {items.map((bullet) => (
        <li key={bullet} className="flex gap-2">
          <span aria-hidden className="mt-[7px] h-1 w-1 flex-shrink-0 rounded-full bg-brand-500" />
          <span>{bullet}</span>
        </li>
      ))}
    </ul>
  )
}

/**
 * Print-only CSS not expressible with Tailwind utilities alone:
 *  - @page paper size and margins (single source of truth for PDF margins;
 *    the backend's PrintToPDF call passes margin=0 + PreferCSSPageSize=true
 *    so these values win)
 *  - page-break / break-inside hints applied via class hooks
 */
function PrintCSS() {
  return (
    <style>{`
      @media print {
        @page {
          size: A4;
          margin: 12mm 14mm;
        }
        .print-item {
          break-inside: avoid;
          page-break-inside: avoid;
        }
        .print-section-title {
          break-after: avoid;
          page-break-after: avoid;
        }
        .print-section-title + * {
          break-before: avoid;
        }
      }
    `}</style>
  )
}
