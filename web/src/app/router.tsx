import { createBrowserRouter, RouterProvider } from 'react-router-dom'

const router = createBrowserRouter([
  {
    path: '/',
    element: (
      <main className="min-h-screen bg-slate-50 px-4 py-10 text-slate-900">
        <section className="mx-auto flex max-w-5xl flex-col gap-6 rounded-[28px] border border-slate-200 bg-white p-6 shadow-sm md:p-10">
          <div className="flex flex-col gap-2">
            <p className="text-sm font-medium uppercase tracking-[0.24em] text-sky-600">
              Resume Prototype
            </p>
            <h1 className="text-3xl font-semibold tracking-tight md:text-5xl">
              温庆京
            </h1>
            <p className="text-base text-slate-600 md:text-lg">
              Golang 后端工程师
            </p>
          </div>
          <p className="max-w-2xl text-sm leading-7 text-slate-500 md:text-base">
            阶段一正在搭建 React 路由与响应式页面骨架，这里会被完整的简历展示页替换。
          </p>
        </section>
      </main>
    ),
  },
  {
    path: '/edit',
    element: <div>编辑简历</div>,
  },
  {
    path: '/visitors',
    element: <div>访客统计</div>,
  },
  {
    path: '/login',
    element: <div>管理后台</div>,
  },
  {
    path: '/print',
    element: <div>打印预览</div>,
  },
])

export function AppRouter() {
  return <RouterProvider router={router} />
}
