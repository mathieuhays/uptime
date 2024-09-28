package healthcheck

import (
	"github.com/google/uuid"
	"math"
	"time"
)

type HealthCheck struct {
	ID           uuid.UUID
	WebsiteID    uuid.UUID
	StatusCode   int
	ResponseTime time.Duration
	CreatedAt    time.Time
}

type HealthCheckSummaryItem struct {
	Date    time.Time
	Average int
	Min     int
	Max     int
}

type DateRange struct {
	Start time.Time
	End   time.Time
}

type Repository interface {
	Create(healthCheck HealthCheck) (*HealthCheck, error)
	Get(id uuid.UUID) (*HealthCheck, error)
	GetByWebsiteID(websiteID uuid.UUID, dateRange DateRange) ([]HealthCheck, error)
	GetLatest(websiteID uuid.UUID) (*HealthCheck, error)
}

func SummarizeHealthChecks(items []HealthCheck, interval time.Duration) []HealthCheckSummaryItem {
	// create groups of healthchecks based on date
	var groups [][]HealthCheck
	var group []HealthCheck

	for _, item := range items {
		if len(group) == 0 {
			group = append(group, item)
			continue
		}

		ref := group[0].CreatedAt
		if item.CreatedAt.Sub(ref).Abs() > interval {
			groups = append(groups, group)
			group = []HealthCheck{
				item,
			}
			continue
		}

		group = append(group, item)
	}

	if len(group) > 0 {
		groups = append(groups, group)
	}

	// process each group
	summaryItems := make([]HealthCheckSummaryItem, len(groups))

	for idx, g := range groups {
		item := HealthCheckSummaryItem{
			Date: g[0].CreatedAt,
			Min:  math.MaxInt,
			Max:  0,
		}
		total := 0

		for _, i := range g {
			milli := int(i.ResponseTime.Milliseconds())

			if milli > item.Max {
				item.Max = milli
			}

			if milli < item.Min {
				item.Min = milli
			}

			total += milli
		}

		item.Average = total / len(g)
		summaryItems[idx] = item
	}

	return summaryItems
}
