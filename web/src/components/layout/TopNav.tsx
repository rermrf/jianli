import { Link, NavLink } from 'react-router-dom'
import { isAuthenticated } from '../../lib/auth'
import { Button } from '../common/Button'

const protectedNavItems = [
  { to: '/edit', label: '编辑' },
  { to: '/visitors', label: '访客' },
]

export function TopNav() {
  const authed = isAuthenticated()

  return (
    <header className="hidden rounded-3xl border border-slate-200 bg-white px-6 py-4 shadow-[var(--shadow-card)] md:flex md:items-center md:justify-between">
      <Link className="flex items-center gap-3 text-sm font-semibold text-slate-900" to="/">
        <span className="flex h-8 w-8 items-center justify-center rounded-full bg-brand-500 text-xs text-white">
          W
        </span>
        温庆京的简历
      </Link>
      <nav className="flex items-center gap-5 text-sm text-slate-500">
        <NavLink
          className={({ isActive }) =>
            isActive ? 'text-brand-600' : 'transition hover:text-slate-900'
          }
          to="/"
        >
          简历
        </NavLink>
        {authed
          ? protectedNavItems.map((item) => (
              <NavLink
                key={item.to}
                className={({ isActive }) =>
                  isActive ? 'text-brand-600' : 'transition hover:text-slate-900'
                }
                to={item.to}
              >
                {item.label}
              </NavLink>
            ))
          : null}
      </nav>
      <Link to="/print">
        <Button className="px-3 py-2 text-xs" variant="primary">
          导出 PDF
        </Button>
      </Link>
    </header>
  )
}
