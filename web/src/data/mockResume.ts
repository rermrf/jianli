import type { ResumeData } from '../types/resume'

export const defaultResume: ResumeData = {
  profile: {
    name: '温庆京',
    title: 'Golang 后端工程师',
    age: 25,
    gender: '男',
    education: '本科',
    experience: '0.9年',
    location: '浙江杭州',
    hometown: '江西赣州',
    phone: '17620096266',
    email: '3219431643@qq.com',
  },
  skills: ['Go', 'MySQL', 'Redis', 'Kafka', 'Docker', 'Gin', 'gRPC', 'K8s'],
  jobIntention: {
    position: 'Golang 后端工程师',
    cities: ['深圳', '杭州', '厦门'],
    availability: '一周内到岗',
  },
  education: [
    {
      school: '江西财经大学现代经济管理学院',
      major: '计算机科学与技术',
      degree: '本科',
      startDate: '2023.09',
      endDate: '2025.07',
    },
  ],
  workExperience: [
    {
      company: '杭州云徕科技有限公司',
      role: 'Golang 后端开发',
      startDate: '2025.05',
      endDate: '2026.02',
      description: [
        '重构项目后台架构，引入 Wire 依赖注入提升模块边界清晰度。',
        '优化高频查询接口，显著降低接口平均响应时间。',
        '参与网关与中间件能力建设，支撑多业务线复用。',
      ],
    },
  ],
  projects: [
    {
      name: 'AI Gateway - LLM 统一网关',
      startDate: '2025.12',
      endDate: '2026.01',
      description: [
        '封装多模型协议差异，统一 OpenAI 与 Anthropic 调用方式。',
        '设计鉴权、计费与路由层能力，支撑多租户接入。',
      ],
    },
    {
      name: 'Kafka 运维消息服务',
      startDate: '2025.01',
      endDate: '2025.06',
      description: [
        '搭建 Kafka 消息消费与监控后台，支撑业务高峰稳定运行。',
      ],
    },
    {
      name: '银税服务中台',
      startDate: '2024.08',
      endDate: '2024.10',
      description: [
        '维护税务相关中台服务，补齐接口治理和数据校验能力。',
      ],
    },
  ],
  awards: [
    { date: '2022.09', title: '国家励志奖学金' },
    { date: '2021.09', title: '国家励志奖学金' },
  ],
}
