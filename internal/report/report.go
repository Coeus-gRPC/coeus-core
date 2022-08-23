package report

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/montanaflynn/stats"
	"io/ioutil"
	"strconv"
	"time"
)

type Report struct {
	ReportID               uuid.UUID          `json:"reportID"`
	TotalCallNum           int                `json:"totalCallNum"`
	SuccessCallCount       int                `json:"successCallCount"`
	ConcurrencyLevel       int                `json:"concurrencyLevel"`
	TotalTimeConsumption   float64            `json:"totalTimeConsumption"`
	AverageTimeConsumption float64            `json:"averageTimeConsumption"`
	FastestTimeConsumption float64            `json:"fastestTimeConsumption"`
	SlowestTimeConsumption float64            `json:"slowestTimeConsumption"`
	TimeConsumptions       []float64          `json:"timeConsumptions"`
	Distribution           map[string]float64 `json:"distribution"`
	RequestPerSecond       float64            `json:"requestPerSecond"`
	Messages               []string           `json:"messages"`
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
	report.Distribution = calculatePercentile(report.TimeConsumptions)

	// fmt.Printf("%v", report)

	// A random id to identify the report.
	report.ReportID = uuid.New()

	return report
}

func WriteReportToFile(report Report, filePath string) error {
	if len(filePath) == 0 {
		fmt.Printf("Report Filepath empty. Print the report struct:\n%v", report)
		return nil
	}

	fileData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, fileData, 0644)

	return err
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

func calculatePercentile(distribution []float64) map[string]float64 {
	distributionMap := make(map[string]float64)
	disData := stats.LoadRawData(distribution)
	percentiles := []float64{5.0, 10.0, 25.0, 50.0, 75.0, 90.0, 95.0, 99.0}

	for _, v := range percentiles {
		// fmt.Printf("Calculating %f\n", v)
		percentile, err := stats.Percentile(disData, v)
		if err != nil {
			continue // to skip the NaN value, which are not encodable in a JSON file
			// println(err.Error())
		}
		distributionMap[strconv.Itoa(int(v))] = percentile
	}

	//println("Finish Calculating percentils")
	return distributionMap
}
