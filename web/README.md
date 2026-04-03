# 在线简历前端

阶段一前端原型，基于 React + TypeScript + React Router + Tailwind CSS。

## 当前范围

- 简历展示页 `/`
- 简历编辑页 `/edit`
- 访客统计页 `/visitors`
- Key 登录页 `/login`
- 打印页 `/print`

## 本地命令

```bash
npm install
npm run dev
npm run test
npm run build
npm run lint
```

## 说明

- 当前阶段使用本地 mock 数据
- 登录态保存在 `sessionStorage`
- 简历草稿保存在 `localStorage`
- 还未接入真实后端和 PDF 导出
