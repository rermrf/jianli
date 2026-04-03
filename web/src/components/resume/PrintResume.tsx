import type { ResumeData } from '../../types/resume'

interface PrintResumeProps {
  resume: ResumeData
}

export function PrintResume({ resume }: PrintResumeProps) {
  return (
    <article className="mx-auto max-w-[794px] space-y-8 rounded-none bg-white p-8 text-slate-900 print:max-w-none print:p-0">
      <header className="border-b border-slate-200 pb-6">
        <h1 className="text-3xl font-semibold">打印版简历</h1>
        <h2 className="mt-3 text-2xl font-semibold">{resume.profile.name}</h2>
        <p className="mt-2 text-sm text-slate-500">{resume.profile.title}</p>
        <p className="mt-3 text-sm text-slate-500">
          {resume.profile.location} · {resume.profile.phone} · {resume.profile.email}
        </p>
      </header>

      <section className="space-y-3">
        <h3 className="text-lg font-semibold">个人技能</h3>
        <p className="text-sm leading-7 text-slate-600">{resume.skills.join(' / ')}</p>
      </section>

      <section className="space-y-4">
        <h3 className="text-lg font-semibold">工作经历</h3>
        {resume.workExperience.map((item) => (
          <div key={`${item.company}-${item.role}`} className="space-y-2">
            <div className="flex items-start justify-between gap-3">
              <strong>{item.company}</strong>
              <span className="text-sm text-slate-500">
                {item.startDate} - {item.endDate}
              </span>
            </div>
            <p className="text-sm text-slate-600">{item.role}</p>
            <ul className="list-disc space-y-1 pl-5 text-sm leading-6 text-slate-600">
              {item.description.map((bullet) => (
                <li key={bullet}>{bullet}</li>
              ))}
            </ul>
          </div>
        ))}
      </section>

      <section className="space-y-4">
        <h3 className="text-lg font-semibold">项目经历</h3>
        {resume.projects.map((item) => (
          <div key={item.name} className="space-y-2">
            <div className="flex items-start justify-between gap-3">
              <strong>{item.name}</strong>
              <span className="text-sm text-slate-500">
                {item.startDate} - {item.endDate}
              </span>
            </div>
            <ul className="list-disc space-y-1 pl-5 text-sm leading-6 text-slate-600">
              {item.description.map((bullet) => (
                <li key={bullet}>{bullet}</li>
              ))}
            </ul>
          </div>
        ))}
      </section>
    </article>
  )
}
