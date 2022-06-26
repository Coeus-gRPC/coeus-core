package report

import (
	"github.com/google/uuid"
	"time"
)

type Report struct {
	TotalCallNum           int       `json:"totalCallNum"`
	SuccessCallCount       int       `json:"successCallCount"`
	ConcurrencyLevel       int       `json:"concurrencyLevel"`
	TotalTimeConsumption   float64   `json:"totalTimeConsumption"`
	AverageTimeConsumption float64   `json:"averageTimeConsumption"`
	FastestTimeConsumption float64   `json:"fastestTimeConsumption"`
	SlowestTimeConsumption float64   `json:"slowestTimeConsumption"`
	TimeConsumptions       []float64 `json:"timeConsumptions"`
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

func GenerateReport(reporters []Reporter, totalTime time.Duration, concurrencyLevel int) Report {
	report := Report{
		TotalCallNum:         len(reporters),
		TotalTimeConsumption: float64(totalTime.Microseconds()) / 1000.,
		ConcurrencyLevel:     concurrencyLevel,
	}

	for i := range reporters {
		r := reporters[i]
		if r.StatusStr == "OK" {
			report.SuccessCallCount += 1
		}
		report.TimeConsumptions = append(report.TimeConsumptions, r.TimeConsumption)
		report.Messages = append(report.Messages, r.ReturnStr)
	}

	report.AverageTimeConsumption = addArray(report.TimeConsumptions) / float64(report.TotalCallNum)
	report.FastestTimeConsumption = smallestElement(report.TimeConsumptions)
	report.SlowestTimeConsumption = largestElement(report.TimeConsumptions)
	report.RequestPerSecond = float64(report.TotalCallNum) / totalTime.Seconds()

	return report
}

func addArray(array []float64) float64 {
	result := 0.0
	for _, v := range array {
		result += v
	}
	return result
}

func smallestElement(array []float64) float64 {
	smallestNumber := array[0]
	for _, element := range array {
		if element < smallestNumber {
			smallestNumber = element
		}
	}

	return smallestNumber
}

func largestElement(array []float64) float64 {
	largestNumber := array[0]
	for _, element := range array {
		if element > largestNumber {
			largestNumber = element
		}
	}

	return largestNumber
}
