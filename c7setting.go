package c7api

import (
	"encoding/json"
	"fmt"
)

type WinerySettingsResponse struct {
	Settings []WinerySettings `json:"settings"`
}

type WinerySettings struct {
	Id                    string   `json:"id"`
	CompanyName           string   `json:"companyName"`
	BusinessNumber        *string  `json:"businessNumber"`
	Logo                  string   `json:"logo"`
	PrimaryEmail          string   `json:"primaryEmail"`
	PrimaryPhone          string   `json:"primaryPhone"`
	ContinueShopping      string   `json:"continueShopping"`
	ClubList              string   `json:"clubList"`
	ReservationList       string   `json:"reservationList"`
	RouterType            string   `json:"routerType"`
	Url                   string   `json:"url"`
	TimeZone              string   `json:"timeZone"`
	TimeFormat            string   `json:"timeFormat"`
	Currency              string   `json:"currency"`
	DefaultWeightUnit     string   `json:"defaultWeightUnit"`
	DefaultCountryCode    string   `json:"defaultCountryCode"`
	OperatingCountryCodes []string `json:"operatingCountryCodes"`
	OperatingStateCodes   []string `json:"operatingStateCodes"`
	UsesOpusOneFedEx      bool     `json:"usesOpusOneFedEx"`
	EmailLogo             string   `json:"emailLogo"`
	MinimumAge            int      `json:"minimumAge"`
	SetupGuideStatus      string   `json:"setupGuideStatus"`
}

func GetWineryInfoSettings(tenant string, auth string, rl genericRateLimiter) (*WinerySettings, error) {

	// Create url and request
	reqUrl := Endpoints.Setting
	settingsResp, err := RequestWithRetryAndRead("GET", reqUrl, nil, nil, tenant, auth, 2, rl)
	if err != nil {
		return nil, fmt.Errorf("while getting settings: %w", err)
	}

	// unmarshall the settings
	settingsPayload := WinerySettingsResponse{}
	err = json.Unmarshal(*settingsResp, &settingsPayload)
	if err != nil {
		return nil, fmt.Errorf("while unmarshalling settings payload: %w", err)
	}

	if len(settingsPayload.Settings) != 1 {
		return nil, fmt.Errorf("settings payload unexpected length: %d", len(settingsPayload.Settings))
	}

	return &settingsPayload.Settings[0], nil

}
