package seed

import "encoding/json"

func DefaultResume() json.RawMessage {
	return json.RawMessage(`{
		"profile": {
			"avatarUrl": "",
			"name": "",
			"title": "",
			"age": 0,
			"gender": "",
			"education": "",
			"experience": "",
			"location": "",
			"hometown": "",
			"phone": "",
			"email": ""
		},
		"skills": [],
		"jobIntention": {
			"position": "",
			"cities": [],
			"availability": ""
		},
		"education": [],
		"workExperience": [],
		"projects": [],
		"awards": []
	}`)
}
