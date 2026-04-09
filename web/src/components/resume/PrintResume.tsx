import type { ResumeData } from '../../types/resume'

interface PrintResumeProps {
  resume: ResumeData
}

export function PrintResume({ resume }: PrintResumeProps) {
  const profileFacts = [
    `${resume.profile.age}岁`,
    resume.profile.gender,
    resume.profile.education,
    resume.profile.experience,
    `籍贯：${resume.profile.hometown}`,
  ].filter(Boolean)

  return (
    <article className="mx-auto max-w-[794px] space-y-8 rounded-none bg-white p-8 text-slate-900 print:max-w-none print:p-0">
      <header className="border-b border-slate-200 pb-6">
        <div className="flex items-start gap-6">
          {resume.profile.avatarUrl ? (
            <img
              alt={`${resume.profile.name}头像`}
              className="h-20 w-20 rounded-full object-cover"
              src={resume.profile.avatarUrl}
            />
          ) : (
            <div className="flex h-20 w-20 items-center justify-center rounded-full bg-slate-100 text-3xl text-slate-400">
              👤
            </div>
          )}
          <div className="min-w-0 flex-1">
            <h1 className="text-3xl font-semibold tracking-[0.02em]">{resume.profile.name}</h1>
            <p className="mt-2 text-base font-medium text-slate-700">{resume.profile.title}</p>
            <p className="mt-3 text-sm leading-6 text-slate-500">
              {`所在地：${resume.profile.location}`} · {`手机号：${resume.profile.phone}`} · {`邮箱：${resume.profile.email}`}
            </p>
            <p className="mt-2 text-sm leading-6 text-slate-500">{profileFacts.join(' / ')}</p>
            {resume.jobIntention.cities.length > 0 ? (
              <p className="mt-2 text-sm text-slate-500">
                {`意向城市：${resume.jobIntention.cities.join(' / ')}`}
              </p>
            ) : null}
          </div>
        </div>
      </header>

      <section className="space-y-3">
        <h3 className="text-lg font-semibold">个人技能</h3>
        <p className="text-sm leading-7 text-slate-600">{resume.skills.join(' / ')}</p>
      </section>

      <section className="space-y-4">
        <h3 className="text-lg font-semibold">教育经历</h3>
        {resume.education.map((item) => (
          <div key={`${item.school}-${item.startDate}`} className="space-y-2">
            <div className="flex items-start justify-between gap-3">
              <strong>{item.school}</strong>
              <span className="text-sm text-slate-500">
                {item.startDate} - {item.endDate}
              </span>
            </div>
            <p className="text-sm text-slate-600">
              {item.major} / {item.degree}
            </p>
          </div>
        ))}
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

      <section className="space-y-4">
        <h3 className="text-lg font-semibold">荣誉奖项</h3>
        {resume.awards.map((item) => (
          <div key={`${item.title}-${item.date}`} className="flex items-start justify-between gap-3">
            <strong>{item.title}</strong>
            <span className="text-sm text-slate-500">{item.date}</span>
          </div>
        ))}
      </section>
    </article>
  )
}
