package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/knadh/listmonk/internal/auth"
	"github.com/knadh/listmonk/models"
	"github.com/labstack/echo/v5"
)

// dashboardCampaignCount is the number of recent campaigns to show on the dashboard.
const dashboardCampaignCount = 6

// dashboardCounts contains stats to render on the dashboard.
type dashboardCounts struct {
	Subscribers struct {
		Total       int `json:"total"`
		Blocklisted int `json:"blocklisted"`
		Orphans     int `json:"orphans"`
	} `json:"subscribers"`

	Lists struct {
		Total       int `json:"total"`
		Private     int `json:"private"`
		Public      int `json:"public"`
		OptinSingle int `json:"optin_single"`
		OptinDouble int `json:"optin_double"`
	} `json:"lists"`

	Campaigns struct {
		Total    int            `json:"total"`
		ByStatus map[string]int `json:"by_status"`
	} `json:"campaigns"`

	Messages int `json:"messages"`
}

type dashboardView struct {
	adminView
	Counts           dashboardCounts
	Charts           types.JSONText
	Campaigns        models.Campaigns
	CacheSlowQueries bool
}

// ViewDashboard renders the admin dashboard (home) page.
func (a *App) ViewDashboard(c *echo.Context) error {
	counts, err := a.core.GetDashboardCounts()
	if err != nil {
		return err
	}

	charts, err := a.core.GetDashboardCharts()
	if err != nil {
		return err
	}

	camps, err := a.getDashboardCampaigns(c)
	if err != nil {
		return err
	}

	var cnt dashboardCounts
	if len(counts) > 0 {
		if err := json.Unmarshal(counts, &cnt); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	data := dashboardView{
		adminView:        newAdminView(c, a.i18n.T("globals.terms.dashboard"), "", "dashboard"),
		Counts:           cnt,
		Charts:           charts,
		Campaigns:        camps,
		CacheSlowQueries: ko.Bool("app.cache_slow_queries"),
	}

	return c.Render(http.StatusOK, "admin-dashboard", data)
}

// GetDashboardCampaigns returns recent campaigns to show on the dashboard,
// with scheduled campaigns on top.
func (a *App) GetDashboardCampaigns(c *echo.Context) error {
	camps, err := a.getDashboardCampaigns(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, okResp{camps})
}

// getDashboardCampaigns queries recent campaigns to show on the dashboard.
func (a *App) getDashboardCampaigns(c *echo.Context) (models.Campaigns, error) {
	user := auth.GetUser(c)

	// Either the user has campaigns:get_all permission and can view all campaigns,
	// or the campaigns are filtered by the lists the user has get|manage access to.
	hasAllPerm := user.HasPerm(auth.PermCampaignsGetAll)
	var permittedLists []int
	if !hasAllPerm {
		hasAllPerm, permittedLists = user.GetPermittedLists(auth.PermTypeGet | auth.PermTypeManage)
	}

	// Get scheduled campaigns to show on top of the dashboard.
	camps, _, err := a.core.QueryCampaigns("", []string{models.CampaignStatusScheduled}, nil, "", models.FieldCreatedAt, models.OrderDesc, hasAllPerm, permittedLists, 0, dashboardCampaignCount)
	if err != nil {
		return nil, err
	}

	// Remaining is non-scheduled campaigns.
	if rem := dashboardCampaignCount - len(camps); rem > 0 {
		others, _, err := a.core.QueryCampaigns("", []string{
			models.CampaignStatusDraft,
			models.CampaignStatusRunning,
			models.CampaignStatusPaused,
			models.CampaignStatusFinished,
			models.CampaignStatusCancelled,
		}, nil, "", models.FieldCreatedAt, models.OrderDesc, hasAllPerm, permittedLists, 0, rem)
		if err != nil {
			return nil, err
		}
		camps = append(camps, others...)
	}

	// The dashboard card only needs status, name and date, so strip the bodies.
	for i := range camps {
		camps[i].Body = ""
		camps[i].BodySource.Valid = false
	}

	return camps, nil
}

// niceDuration formats duration to human-readable string, eg: 1h, 3s etc.
func niceDuration(d time.Duration) string {
	d = d.Round(time.Second)
	var (
		days  = int(d.Hours()) / 24
		hours = int(d.Hours()) % 24
		mins  = int(d.Minutes()) % 60
		secs  = int(d.Seconds()) % 60
	)

	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	case hours > 0:
		return fmt.Sprintf("%dh %dm", hours, mins)
	case mins > 0:
		return fmt.Sprintf("%dm %ds", mins, secs)
	default:
		return fmt.Sprintf("%ds", secs)
	}
}
