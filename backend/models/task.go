package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Subject           string             `bson:"subject" json:"subject"`
	TaskType          string             `bson:"task_type" json:"taskType"`
	Number            int                `bson:"number" json:"number"`
	Condition         string             `bson:"condition" json:"condition"`
	ConditionImageURL string             `bson:"condition_image_url" json:"conditionImageUrl"`
	Answer            string             `bson:"answer" json:"answer"`
	Solution          string             `bson:"solution" json:"solution"`
	SolutionImageURL  string             `bson:"solution_image_url" json:"solutionImageUrl"`
}

type TaskResponse struct {
	ID                string `json:"id"`
	Subject           string `json:"subject"`
	TaskType          string `json:"taskType"`
	Number            int    `json:"number"`
	Condition         string `json:"condition"`
	ConditionImageURL string `json:"conditionImageUrl"`
	Answer            string `json:"answer"`
	Solution          string `json:"solution"`
	SolutionImageURL  string `json:"solutionImageUrl"`
}

func (t *Task) ToResponse() TaskResponse {
	return TaskResponse{
		ID:                t.ID.Hex(),
		Subject:           t.Subject,
		TaskType:          t.TaskType,
		Number:            t.Number,
		Condition:         t.Condition,
		ConditionImageURL: t.ConditionImageURL,
		Answer:            t.Answer,
		Solution:          t.Solution,
		SolutionImageURL:  t.SolutionImageURL,
	}
}

func ToResponseList(tasks []Task) []TaskResponse {
	out := make([]TaskResponse, 0, len(tasks))
	for i := range tasks {
		if tasks[i].ID.IsZero() {
			continue
		}
		out = append(out, tasks[i].ToResponse())
	}
	return out
}
