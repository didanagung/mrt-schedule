package station

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/didanagung/mrt-schedule/common/client"
	"github.com/didanagung/mrt-schedule/config"
)

type Service interface {
	GetAllStation() (response []StationResponse, err error)
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
