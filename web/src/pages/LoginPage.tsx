import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button } from '../components/common/Button'
import { SectionCard } from '../components/common/SectionCard'
import { AppShell } from '../components/layout/AppShell'
import { consumeRedirectPath, loginWithKey } from '../lib/auth'

export function LoginPage() {
  const navigate = useNavigate()
  const [keyValue, setKeyValue] = useState('')
  const [error, setError] = useState('')

  function handleLogin() {
    const trimmedKey = keyValue.trim()

    if (!trimmedKey) {
      setError('请输入访问密钥')
      return
    }

    loginWithKey(trimmedKey)
    setError('')
    navigate(consumeRedirectPath() ?? '/edit', { replace: true })
  }

  return (
    <AppShell contentClassName="flex min-h-screen items-center justify-center md:min-h-[calc(100vh-80px)]">
      <SectionCard className="w-full max-w-[420px] space-y-5 p-8 text-center">
        <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-brand-50 text-2xl text-brand-600">
          🔒
        </div>
        <div className="space-y-2">
          <h1 className="text-2xl font-semibold text-slate-900">管理后台</h1>
          <p className="text-sm text-slate-500">请输入访问密钥以继续</p>
        </div>
        <div className="space-y-2 text-left">
          <input
            className="w-full rounded-2xl border border-slate-200 px-4 py-3 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-brand-500"
            onChange={(event) => setKeyValue(event.target.value)}
            placeholder="输入访问密钥"
            type="password"
            value={keyValue}
          />
          {error ? <p className="text-sm text-rose-500">{error}</p> : null}
        </div>
        <Button className="w-full" onClick={handleLogin}>
          验证并登录
        </Button>
      </SectionCard>
    </AppShell>
  )
}
