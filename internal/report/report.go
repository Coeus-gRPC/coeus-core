package report

import "github.com/google/uuid"

type Report struct {
	TotalCallNum           int       `json:"totalCallNum"`
	SuccessCallCount       int       `json:"successCallCount"`
	TotalTimeConsumption   float64   `json:"totalTimeConsumption"`
	AverageTimeConsumption float64   `json:"averageTimeConsumption"`
	FastestTimeConsumption float64   `json:"fastestTimeConsumption"`
	SlowestTimeConsumption float64   `json:"slowestTimeConsumption"`
	Distribution           []float64 `json:"distribution"`
	RequestPerSecond       float64   `json:"requestPerSecond"`
	Messages               []string  `json:"messages"`
}

type Reporter struct {
	ID uuid.UUID
	// StatusStr is represented as string b/c status.Code(err) can be easily converted to string
	StatusStr       string
	TimeConsumption float64
	ReturnStr       string
}

func NewReporter() Reporter {
	return Reporter{
		ID:              uuid.New(),
		StatusStr:       "",
		TimeConsumption: -1.0,
		ReturnStr:       "",
	}
}
