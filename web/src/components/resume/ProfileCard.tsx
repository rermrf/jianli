import type { ResumeProfile } from '../../types/resume'
import { SectionCard } from '../common/SectionCard'

interface ProfileCardProps {
  profile: ResumeProfile
}

const baseFacts = (profile: ResumeProfile) => [
  `${profile.age}岁`,
  profile.gender,
  profile.education,
  profile.experience,
]

export function ProfileCard({ profile }: ProfileCardProps) {
  return (
    <SectionCard className="space-y-5">
      <div className="flex flex-col items-center text-center">
        <div className="flex h-20 w-20 items-center justify-center rounded-full bg-slate-100 text-3xl text-slate-400">
          👤
        </div>
        <h1 className="mt-4 text-2xl font-semibold text-slate-900">{profile.name}</h1>
        <p className="mt-2 text-sm text-slate-500">{profile.title}</p>
      </div>
      <div className="flex flex-wrap justify-center gap-2 text-xs text-slate-400">
        {baseFacts(profile).map((item) => (
          <span key={item} className="rounded-full bg-slate-100 px-3 py-1">
            {item}
          </span>
        ))}
      </div>
      <div className="space-y-2 text-sm text-slate-500">
        <p>{profile.location}</p>
        <p>{profile.phone}</p>
        <p>{profile.email}</p>
      </div>
    </SectionCard>
  )
}
