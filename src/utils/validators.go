package utils

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ParseMeterIDs(param string) ([]int, error) {
	ids := strings.Split(param, ",")
	var result []int
	for _, id := range ids {
		parsedID, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		result = append(result, parsedID)
	}
	return result, nil
}

func IsValidPeriod(kindPeriod string) bool {
	validPeriods := map[string]bool{"daily": true, "weekly": true, "monthly": true}
	return validPeriods[kindPeriod]
}

func IsValidDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02"
	return time.Parse(layout, dateStr)
}

func ValidateQueryParams(query url.Values) ([]int, time.Time, time.Time, string, int, string) {

    meterIDsParam := query.Get("meters_ids")
    meterIDs, err := ParseMeterIDs(meterIDsParam)
    if err != nil {
        return nil, time.Time{}, time.Time{}, "", http.StatusBadRequest, "❌ Invalid meter_ids format. Expected comma-separated integers."
    }

    startDateStr := query.Get("start_date")
    startDate, err := IsValidDate(startDateStr)
    if err != nil {
        return nil, time.Time{}, time.Time{}, "", http.StatusBadRequest, "❌ Invalid start_date format. Expected YYYY-MM-DD."
    }

    endDateStr := query.Get("end_date")
    endDate, err := IsValidDate(endDateStr)
    if err != nil {
        return nil, time.Time{}, time.Time{}, "", http.StatusBadRequest, "❌ Invalid end_date format. Expected YYYY-MM-DD."
    }

    if startDate.After(endDate) {
        return nil, time.Time{}, time.Time{}, "", http.StatusBadRequest, "❌ start_date cannot be later than end_date."
    }

    kindPeriod := query.Get("kind_period")
    if !IsValidPeriod(kindPeriod) {
        return nil, time.Time{}, time.Time{}, "", http.StatusBadRequest, "❌ Invalid kind_period. Only 'daily', 'weekly', or 'monthly' are allowed."
    }

    return meterIDs, startDate, endDate, kindPeriod, http.StatusOK, ""
}