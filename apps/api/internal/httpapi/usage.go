package httpapi

import (
	"net/http"
	"time"

	"github.com/marciojunior/waytutor/apps/api/internal/quota"
)

const adminUsageDays = 14

type adminTotalsBody struct {
	Chars     int64 `json:"chars"`
	Requests  int64 `json:"requests"`
	Summaries int64 `json:"summaries"`
	Accounts  int64 `json:"accounts,omitempty"`
}

type adminLimitsBody struct {
	GlobalChars      int64 `json:"globalChars"`
	AccountChars     int64 `json:"accountChars"`
	MonthlySummaries int64 `json:"monthlySummaries"`
}

type adminCountsBody struct {
	Active  int64 `json:"active"`
	Pending int64 `json:"pending"`
	Paused  int64 `json:"paused"`
}

type adminDayBody struct {
	Day      string `json:"day"`
	Chars    int64  `json:"chars"`
	Requests int64  `json:"requests"`
}

type adminTopBody struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Chars    int64  `json:"chars"`
	Requests int64  `json:"requests"`
}

type adminUsageBody struct {
	Day      string          `json:"day"`
	Period   string          `json:"period"`
	Today    adminTotalsBody `json:"today"`
	Month    adminTotalsBody `json:"month"`
	Limits   adminLimitsBody `json:"limits"`
	Accounts adminCountsBody `json:"accounts"`
	Daily    []adminDayBody  `json:"daily"`
	Top      []adminTopBody  `json:"top"`
}

func (h *api) usage(w http.ResponseWriter, r *http.Request) {
	usage, err := h.store.AdminUsage(r.Context(), quota.Period(time.Now()), adminUsageDays)
	if err != nil {
		respondInternal(w, r, err)
		return
	}

	body := adminUsageBody{
		Day:    usage.Day,
		Period: usage.Period,
		Today: adminTotalsBody{
			Chars:     usage.Today.Chars,
			Requests:  usage.Today.Requests,
			Summaries: usage.Today.Summaries,
			Accounts:  usage.TodayAccounts,
		},
		Month: adminTotalsBody{
			Chars:     usage.Month.Chars,
			Requests:  usage.Month.Requests,
			Summaries: usage.Month.Summaries,
		},
		Limits: adminLimitsBody{
			GlobalChars:      h.limits.GlobalMonthlyChars,
			AccountChars:     h.limits.MonthlyCharsPerAccount,
			MonthlySummaries: h.limits.MonthlySummaries,
		},
		Accounts: adminCountsBody{
			Active:  usage.Active,
			Pending: usage.Pending,
			Paused:  usage.Paused,
		},
		Daily: make([]adminDayBody, 0, len(usage.Daily)),
		Top:   make([]adminTopBody, 0, len(usage.Top)),
	}
	for _, day := range usage.Daily {
		body.Daily = append(body.Daily, adminDayBody{Day: day.Day, Chars: day.Chars, Requests: day.Requests})
	}
	for _, account := range usage.Top {
		body.Top = append(body.Top, adminTopBody{
			Email:    account.Email,
			Name:     account.Name,
			Chars:    account.Chars,
			Requests: account.Requests,
		})
	}

	writeJSON(w, http.StatusOK, body)
}
