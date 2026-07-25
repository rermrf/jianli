import { useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { monogramFromName } from '../../../lib/resumePresentation'
import type { ResumeProfile } from '../../../types/resume'
import { GlassCard } from './GlassCard'

interface ResumeHeroProps {
  profile: ResumeProfile
  desiredCities?: string[]
  tagline: string
}

export function ResumeHero({
  profile,
  desiredCities = [],
  tagline,
}: ResumeHeroProps) {
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

  const chips = [
    profile.age ? `${profile.age}岁` : '',
    profile.gender,
    profile.education,
    profile.experience,
    profile.location,
  ].filter(Boolean)

  const monogram = monogramFromName(profile.name)

  return (
    <GlassCard strong className="space-y-5">
      <div className="flex flex-col items-center gap-5 text-center md:flex-row md:items-start md:text-left">
        <button
          aria-label="头像入口"
          className="shrink-0 rounded-full outline-none transition focus-visible:ring-2 focus-visible:ring-[var(--resume-brand-500)] focus-visible:ring-offset-2"
          onClick={handleAvatarClick}
          type="button"
        >
          {profile.avatarUrl ? (
            <img
              alt={`${profile.name}头像`}
              className="h-24 w-24 rounded-full object-cover"
              src={profile.avatarUrl}
            />
          ) : (
            <div
              aria-hidden="true"
              className="flex h-24 w-24 items-center justify-center rounded-full bg-gradient-to-br from-[var(--resume-brand-500)] to-[var(--resume-brand-600)] text-3xl font-semibold text-white"
            >
              {monogram}
            </div>
          )}
        </button>

        <div className="min-w-0 flex-1 space-y-3">
          <div className="space-y-1">
            <h1 className="text-2xl font-semibold tracking-tight text-[var(--resume-ink-900)] md:text-3xl">
              {profile.name}
            </h1>
            <p className="text-base font-medium text-[var(--resume-brand-600)]">
              {profile.title}
            </p>
            {tagline ? (
              <p className="text-sm text-[var(--resume-ink-600)]">{tagline}</p>
            ) : null}
          </div>

          {chips.length > 0 ? (
            <div className="flex flex-wrap justify-center gap-2 md:justify-start">
              {chips.map((chip) => (
                <span
                  key={chip}
                  className="rounded-full bg-[var(--resume-brand-50)] px-3 py-1 text-xs text-[var(--resume-ink-600)]"
                >
                  {chip}
                </span>
              ))}
            </div>
          ) : null}

          <div className="flex flex-col items-center gap-1 text-sm text-[var(--resume-ink-600)] md:items-start">
            {desiredCities.length > 0 ? (
              <p>{`意向城市：${desiredCities.join(' / ')}`}</p>
            ) : null}
            {profile.phone ? (
              <a
                className="outline-none transition hover:text-[var(--resume-brand-600)] focus-visible:ring-2 focus-visible:ring-[var(--resume-brand-500)] focus-visible:ring-offset-2"
                href={`tel:${profile.phone}`}
              >
                {profile.phone}
              </a>
            ) : null}
            {profile.email ? (
              <a
                className="outline-none transition hover:text-[var(--resume-brand-600)] focus-visible:ring-2 focus-visible:ring-[var(--resume-brand-500)] focus-visible:ring-offset-2"
                href={`mailto:${profile.email}`}
              >
                {profile.email}
              </a>
            ) : null}
          </div>
        </div>
      </div>
    </GlassCard>
  )
}
