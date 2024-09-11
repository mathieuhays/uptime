package uptime

import (
	"context"
	"github.com/google/uuid"
	"github.com/mathieuhays/uptime/internal/healthcheck"
	"github.com/mathieuhays/uptime/internal/website"
	"log"
	"net/http"
	"slices"
	"sync"
	"time"
)

type CrawlerHealthCheckRepo interface {
	Create(healthCheck healthcheck.HealthCheck) (*healthcheck.HealthCheck, error)
}

type Crawler struct {
	interval time.Duration

	healthCheckRepo healthcheck.Repository
	websiteRepo     website.Repository

	concurrencyControl chan struct{}
}

func NewCrawler(healthCheckRepo healthcheck.Repository, websiteRepo website.Repository, interval time.Duration, maxConcurrency int) *Crawler {
	return &Crawler{
		healthCheckRepo:    healthCheckRepo,
		websiteRepo:        websiteRepo,
		interval:           interval,
		concurrencyControl: make(chan struct{}, maxConcurrency),
	}
}

func (c *Crawler) Start(ctx context.Context) {
	for {
		c.crawl()

		select {
		case <-ctx.Done():
			log.Printf("Cancelling crawler")
			return
		default:
			time.Sleep(c.interval)
		}
	}
}

func (c *Crawler) crawl() {
	log.Printf("Start crawl")
	var websites []website.Website

	items, err := c.websiteRepo.GetWebsitesByLastFetched(time.Now().Add(time.Minute*-5), 2)
	if err == nil {
		items = slices.Concat(websites, items)
	}

	log.Printf("Crawling %d URLs", len(items))

	wg := &sync.WaitGroup{}

	for _, item := range items {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.crawlURL(item)
		}()
	}

	wg.Wait()

	log.Printf("Done crawling")
}

func (c *Crawler) crawlURL(w website.Website) {
	defer func() {
		<-c.concurrencyControl
	}()

	log.Printf("Initiating crawl for URL: %s", w.URL)

	c.concurrencyControl <- struct{}{}

	log.Printf("Executing crawl for URL: %s", w.URL)

	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, w.URL, nil)
	if err != nil {
		log.Printf("request creation error for %s: %s", w.URL, err)
		return
	}

	client := http.Client{}
	res, err := client.Do(req)
	end := time.Now()

	if err != nil {
		log.Printf("request error for %s: %s", w.URL, err)
		return
	}
	defer res.Body.Close()

	healthCheck := healthcheck.HealthCheck{
		ID:           uuid.New(),
		WebsiteID:    w.ID,
		StatusCode:   res.StatusCode,
		ResponseTime: end.Sub(start).Abs(),
		CreatedAt:    time.Now(),
	}

	_, err = c.healthCheckRepo.Create(healthCheck)
	if err != nil {
		log.Printf("error creating healthcheck for %s: %s", w.URL, err)
		return
	}

	err = c.websiteRepo.SetAsFetched(w.ID, time.Now())
	if err != nil {
		log.Printf("error updating website fetch time for %s: %s", w.URL, err)
	}

	log.Printf("Done crawling %s", w.URL)
}
