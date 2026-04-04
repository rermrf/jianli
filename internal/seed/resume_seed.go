package seed

import "encoding/json"

func DefaultResume() json.RawMessage {
	return json.RawMessage(`{
		"profile": {
			"name": "温庆京",
			"title": "Golang 后端工程师",
			"age": 25,
			"gender": "男",
			"education": "本科",
			"experience": "0.9年",
			"location": "浙江杭州",
			"hometown": "江西赣州",
			"phone": "17620096266",
			"email": "3219431643@qq.com"
		},
		"skills": ["Go", "MySQL", "Redis", "Kafka", "Docker", "Gin", "gRPC", "K8s"],
		"jobIntention": {
			"position": "Golang 后端工程师",
			"cities": ["深圳", "杭州", "厦门"],
			"availability": "一周内到岗"
		},
		"education": [{
			"school": "江西财经大学现代经济管理学院",
			"major": "计算机科学与技术",
			"degree": "本科",
			"startDate": "2023.09",
			"endDate": "2025.07"
		}],
		"workExperience": [{
			"company": "杭州云徕科技有限公司",
			"role": "Golang 后端开发",
			"startDate": "2025.05",
			"endDate": "2026.02",
			"description": ["重构项目后台架构", "优化高频查询接口", "参与网关能力建设"]
		}],
		"projects": [{
			"name": "AI Gateway - LLM 统一网关",
			"startDate": "2025.12",
			"endDate": "2026.01",
			"description": ["统一模型协议", "设计鉴权与路由能力"]
		}],
		"awards": [
			{ "date": "2022.09", "title": "国家励志奖学金" },
			{ "date": "2021.09", "title": "国家励志奖学金" }
		]
	}`)
}
