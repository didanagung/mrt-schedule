package station

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/didanagung/mrt-schedule/common/client"
	"github.com/didanagung/mrt-schedule/config"
)

type Service interface {
	GetAllStation() (response []StationResponse, err error)
	CheckScheduleByStation(id string) (response []ScheduleResponse, err error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *service) GetAllStation() (response []StationResponse, err error) {
	// layer service
	url := config.GetString("URL_STATIONS")

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var stations []Station

	// ini berfungsi untuk menkonversi json dengan bentuk station
	err = json.Unmarshal(byteResponse, &stations)

	for _, item := range stations {
		response = append(response, StationResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	return
}

func (s *service) CheckScheduleByStation(id string) (response []ScheduleResponse, err error) {
	url := config.GetString("URL_STATIONS")

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var schedule []Schedule

	// ini berfungsi untuk menkonversi json dengan bentuk station
	err = json.Unmarshal(byteResponse, &schedule)
	if err != nil {
		return
	}

	// schedule selected by id station
	var scheduleSelected Schedule
	for _, item := range schedule {
		if item.StationId == id {
			scheduleSelected = item
			break
		}
	}

	if scheduleSelected.StationId == "" {
		err = errors.New("station not found.")
		return
	}

	response, err = ConvertDataToResponses(scheduleSelected)
	if err != nil {
		return
	}
	return
}

func ConvertDataToResponses(schedule Schedule) (response []ScheduleResponse, err error) {
	var (
		LebakBulusTripName = "Stasiun Lebak Bulus BSI"
		BundaranHITripName = "Stasiun Bundaran HI Bank DKI"
	)

	scheduleLebakBulus := schedule.ScheduleLebakBulus
	scheduleBundaranHI := schedule.ScheduleBundaranHI

	scheduleLebakBulusParsed, err := ConvertScheduleToTimeFormat(scheduleLebakBulus)
	if err != nil {
		return
	}

	scheduleBundaranHIParsed, err := ConvertScheduleToTimeFormat(scheduleBundaranHI)
	if err != nil {
		return
	}

	// convert to response
	for _, item := range scheduleLebakBulusParsed {
		if item.Format("15:04:05") > time.Now().Format("15:04:05") { //validasi data muncul lebih dari jam sekarang
			response = append(response, ScheduleResponse{
				StationName: LebakBulusTripName,
				Time:        item.Format("15:04:05"),
			})
		}
	}

	for _, item := range scheduleBundaranHIParsed {
		if item.Format("15:04:05") > time.Now().Format("15:04:05") {
			response = append(response, ScheduleResponse{
				StationName: BundaranHITripName,
				Time:        item.Format("15:04:05"),
			})
		}
	}

	return
}

func ConvertScheduleToTimeFormat(schedule string) (response []time.Time, err error) {
	var (
		parsedTime time.Time
		schedules  = strings.Split(schedule, ",") //split yang tadinya string
	)

	for _, item := range schedules {
		trimmedTime := strings.TrimSpace(item) //menghilangkan spasi
		if trimmedTime == "" {
			continue
		}

		parsedTime, err = time.Parse("15:04:05", trimmedTime) //15:04 can manggih naha euy
		if err != nil {
			err = errors.New("invalid time format " + trimmedTime)
			return
		}
		response = append(response, parsedTime)
	}

	return
}
