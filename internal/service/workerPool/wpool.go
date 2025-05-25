package workerpool

import (
	"fmt"
	"log/slog"
	"net/http"
	dbmodel "service-healthz-checker/internal/model/dbModel"
	"service-healthz-checker/internal/service"
	"service-healthz-checker/internal/service/notification"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Job struct {
	Item *dbmodel.Service
}

type Result struct {
	Item *dbmodel.Service
	Err  error
}

type History struct {
	Item *dbmodel.Service
	Err  error
}

func mainWorker(id int, jobs <-chan Job, results chan<- Result, service *service.PingService, notifier *notification.TGNotifier, log *slog.Logger) {
	log.Info("Start Worker", slog.Int("ID", id))
	for job := range jobs {
		start := time.Now()
		res, err := http.Get(normalizeURL(job.Item.Url))
		duration := time.Since(start)

		if err != nil {
			log.Warn("Request failed", slog.String("url", job.Item.Url), slog.String("error", err.Error()))
			now := time.Now().UTC()
			log.Info("time", slog.String("time", fmt.Sprintf("%v", now.Sub(job.Item.LastErrMsg).Hours())))
			if now.Sub(job.Item.LastErrMsg).Minutes() >= 1 {
				notifier.Send(job.Item.ChatID, job.Item.Url, err.Error())
				job.Item.LastErrMsg = time.Now().UTC()
			}
			results <- Result{Item: job.Item, Err: err}
			continue
		}

		avgResTime, err := service.AvgResTime(job.Item.Id)
		if err == nil {
			abnormalTimeMS := avgResTime * 2.2
			abnormalDuration := time.Duration(abnormalTimeMS * float64(time.Millisecond))
			if duration > abnormalDuration {
				msg := AbnormalTimeMSG(job.Item.Url, avgResTime, duration, abnormalTimeMS)
				notifier.Send(job.Item.ChatID, job.Item.Url, msg)
			}
		}

		res.Body.Close()
		job.Item.LastStatus = res.StatusCode
		job.Item.ResponseTimeMs = int(duration.Milliseconds())
		job.Item.UpdatedAt = time.Now().UTC()
		results <- Result{Item: job.Item, Err: nil}
	}
}

func historyWorker(id int, results <-chan History, service *service.PingService, log *slog.Logger) {
	log.Info("START HISTORY", slog.Int("ID", id))
	var historyBucket []*dbmodel.Service
	maxSize := 100

	timeout := 5 * time.Second
	timer := time.NewTimer(timeout)

	for {
		select {
		case result, ok := <-results:
			if !ok {
				if len(historyBucket) > 0 {
					service.SaveHistory(historyBucket)
				}
				return
			}

			if result.Err != nil {
				log.Error("historyWorker: failed to process result",
					slog.String("error", result.Err.Error()),
					slog.Int("workerID", id))
				continue
			}

			historyBucket = append(historyBucket, result.Item)

			if len(historyBucket) >= maxSize {
				service.SaveHistory(historyBucket)
				historyBucket = nil
				timer.Reset(timeout)
			}

		case <-timer.C:
			if len(historyBucket) > 0 {
				service.SaveHistory(historyBucket)
				historyBucket = nil
			}

			timer.Reset(timeout)
		}
	}
}

func RunMainPool(service *service.PingService, notifier *notification.TGNotifier, log *slog.Logger, mainConcurrency int, historyConcurrency int, interval time.Duration) {
	jobs := make(chan Job, 1000)
	results := make(chan Result, 1000)
	history := make(chan History, 1000)

	for w := 0; w < mainConcurrency; w++ {
		go mainWorker(w, jobs, results, service, notifier, log)
	}

	for w := 0; w < historyConcurrency; w++ {
		go historyWorker(w, history, service, log)
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
			history <- History{Item: results.Item}
			updated = append(updated, results.Item)
		}

		service.UpdateData(updated)
		time.Sleep(interval)
	}
}

func DeleteOldWrites(db *sqlx.DB, log *slog.Logger) {
	log.Info("START DB CLEANER")
	for {
		timeForDelete := time.Now().UTC().Add(-1 * time.Hour)

		tx, err := db.Beginx()
		if err != nil {
			log.Error("Failed to begin transaction", "error", err)
			return
		}
		defer tx.Rollback()
		res, err := tx.Exec(`DELETE FROM history WHERE created_at < $1`, timeForDelete)
		if err != nil {
			log.Info("DB ERR OF CLEAN")
			log.Info("DB ERR", slog.String("err", err.Error()))
			return
		}

		if err := tx.Commit(); err != nil {
			log.Error("Commit failed", "error", err)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			return
		}

		log.Info("DB WAS BE CLEANED")
		time.Sleep(time.Minute)
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

func AbnormalTimeMSG(url string, avgResTime float64, duration time.Duration, abnormalTimeMS float64) string {
	msg := fmt.Sprintf(`
⚠️ *Обнаружено аномальное время ответа!*

Сервис: %s
Среднее время ответа: %.2f мс
Текущее время ответа: %v
Превышение порога: %.2f мс (в %.1f раз больше среднего)

Рекомендуется проверить доступность сервиса.
`,
		url,
		avgResTime,
		duration.Round(time.Millisecond),
		abnormalTimeMS-avgResTime,
		abnormalTimeMS/avgResTime)
	return msg
}
