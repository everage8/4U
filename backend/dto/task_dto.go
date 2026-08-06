package dto

type CreateTaskRequest struct {
	Subject           string `json:"subject" binding:"required"`
	TaskType          string `json:"taskType" binding:"required"`
	Condition         string `json:"condition" binding:"required"`
	ConditionImageURL string `json:"conditionImageUrl"`
	Answer            string `json:"answer" binding:"required"`
	Solution          string `json:"solution"`
	SolutionImageURL  string `json:"solutionImageUrl"`
}

type UpdateTaskRequest struct {
	Subject           string `json:"subject" binding:"required"`
	TaskType          string `json:"taskType" binding:"required"`
	Condition         string `json:"condition" binding:"required"`
	ConditionImageURL string `json:"conditionImageUrl"`
	Answer            string `json:"answer" binding:"required"`
	Solution          string `json:"solution"`
	SolutionImageURL  string `json:"solutionImageUrl"`
}

type AdminTasksListResponse struct {
	Tasks []interface{} `json:"tasks"`
	Total int           `json:"total"`
}
