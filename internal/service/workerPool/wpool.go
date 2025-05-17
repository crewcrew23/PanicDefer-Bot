package workerpool

import (
	"log/slog"
	"net/http"
	dbmodel "service-healthz-checker/internal/model/dbModel"
	"service-healthz-checker/internal/service"
	"strings"
	"time"
)

type Job struct {
	Item *dbmodel.Service
}

type Result struct {
	Item *dbmodel.Service
	Err  error
}

func worker(id int, jobs <-chan Job, results chan<- Result, log *slog.Logger) {
	log.Info("Start Worker", slog.Int("ID", id))
	for job := range jobs {
		start := time.Now()
		job.Item.Url = normalizeURL(job.Item.Url)
		res, err := http.Get(job.Item.Url)
		duration := time.Since(start)

		if err != nil {
			log.Warn("Request failed", slog.String("url", job.Item.Url), slog.String("error", err.Error()))
			results <- Result{Item: job.Item, Err: err}
			continue
		}

		res.Body.Close()
		job.Item.LastStatus = res.StatusCode
		job.Item.ResponseTimeMs = int(duration.Microseconds())
		results <- Result{Item: job.Item, Err: nil}
	}
}

func RunPool(service *service.PingService, log *slog.Logger, concurrency int, interval time.Duration) {
	jobs := make(chan Job)
	results := make(chan Result)

	for w := 0; w < concurrency; w++ {
		go worker(w, jobs, results, log)
	}

	for {
		data := service.DataForPing()
		if data == nil {
			time.Sleep(interval)
			continue
		}

		go func() {
			for _, item := range data {
				jobs <- Job{Item: item}
			}
		}()

		var updated []*dbmodel.Service
		for i := 0; i < len(data); i++ {
			results := <-results
			updated = append(updated, results.Item)
		}

		service.UpdateData(updated)
		time.Sleep(interval)
	}
}

func normalizeURL(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}

	if isReachable("https://" + url) {
		return "https://" + url
	}

	return "http://" + url
}

func isReachable(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		return false
	}
	res.Body.Close()

	return res.StatusCode < 500
}
