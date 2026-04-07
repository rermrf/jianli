import { useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import type { ResumeProfile } from '../../types/resume'
import { SectionCard } from '../common/SectionCard'

interface ProfileCardProps {
  desiredCities?: string[]
  profile: ResumeProfile
}

const baseFacts = (profile: ResumeProfile) => [
  `${profile.age}岁`,
  profile.gender,
  profile.education,
  profile.experience,
]

export function ProfileCard({ desiredCities = [], profile }: ProfileCardProps) {
  const navigate = useNavigate()
  const clickTimestampsRef = useRef<number[]>([])

  function handleAvatarClick() {
    const now = Date.now()
    const nextClicks = [...clickTimestampsRef.current, now].filter(
      (timestamp) => now - timestamp <= 1500,
    )

    clickTimestampsRef.current = nextClicks

    if (nextClicks.length >= 3) {
      clickTimestampsRef.current = []
      navigate('/login')
    }
  }

  return (
    <SectionCard className="space-y-5">
      <div className="flex flex-col items-center text-center">
        <button
          aria-label="头像入口"
          className="rounded-full outline-none transition focus-visible:ring-2 focus-visible:ring-brand-500"
          onClick={handleAvatarClick}
          type="button"
        >
          {profile.avatarUrl ? (
            <img
              alt={`${profile.name}头像`}
              className="h-20 w-20 rounded-full object-cover"
              src={profile.avatarUrl}
            />
          ) : (
            <div className="flex h-20 w-20 items-center justify-center rounded-full bg-slate-100 text-3xl text-slate-400">
              👤
            </div>
          )}
        </button>
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
        <p>{`所在地：${profile.location}`}</p>
        {desiredCities.length > 0 ? (
          <p>{`意向城市：${desiredCities.join(' / ')}`}</p>
        ) : null}
        <p>{`手机号：${profile.phone}`}</p>
        <p>{`邮箱：${profile.email}`}</p>
      </div>
    </SectionCard>
  )
}
