package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"sync"
)

const CacheStatisticsTtl = 500

type StatisticsService struct {
	repository repository.StatisticsCrm
	cache      cache.Cache
}

func NewStatisticsService(repository repository.StatisticsCrm, cache cache.Cache) StatisticsService {
	return StatisticsService{
		repository: repository,
		cache:      cache,
	}
}

func (s StatisticsService) GetStatistics(ctx context.Context, userModel domain.User) (domain.Statistics, error) {
	var stats domain.Statistics
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var totalErr, openErr, ipError, wrError, closedError error
	limitCh := make(chan struct{}, 1)
	limitCh <- struct{}{}

	err := GetFromCache[*domain.Statistics]("stat-"+userModel.Crmid, &stats, s.cache)
	if err == nil {
		return stats, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		wg.Add(5)
		// Total Tickets
		go func() {
			defer wg.Done()
			total, err := s.repository.CalcTicketTotal(ctx, userModel)
			<-limitCh
			if err != nil {
				totalErr = err
				return
			}

			mutex.Lock()
			stats.Tickets.Total = total
			mutex.Unlock()
		}()

		// Open tickets
		go func() {
			defer wg.Done()
			var openHours float64
			var openDays float64
			var openTotal int
			limitCh <- struct{}{}
			tickets, err := s.repository.TicketOpenStat(ctx, userModel)
			<-limitCh
			if err != nil {
				openErr = err
				return
			}

			for _, ticket := range tickets {
				openHours += ticket.Hours
				openDays += ticket.Days
				openTotal++
			}
			mutex.Lock()
			stats.Tickets.Open = openTotal
			stats.Tickets.OpenDays = openDays
			stats.Tickets.OpenHours = openHours
			mutex.Unlock()
		}()

		// In progress tickets
		go func() {
			defer wg.Done()
			var ipHours float64
			var ipDays float64
			var ipTotal int
			limitCh <- struct{}{}
			tickets, err := s.repository.TicketInProgressStat(ctx, userModel)
			<-limitCh
			if err != nil {
				ipError = err
				return
			}

			for _, ticket := range tickets {
				ipHours += ticket.Hours
				ipDays += ticket.Days
				ipTotal++
			}
			mutex.Lock()
			stats.Tickets.InProgress = ipTotal
			stats.Tickets.InProgressDays = ipDays
			stats.Tickets.InProgressHours = ipHours
			mutex.Unlock()
		}()

		// Wait For Response tickets
		go func() {
			defer wg.Done()
			var wrHours float64
			var wrDays float64
			var wrTotal int
			limitCh <- struct{}{}
			tickets, err := s.repository.TicketWaitForResponseStat(ctx, userModel)
			<-limitCh
			if err != nil {
				ipError = err
				return
			}

			for _, ticket := range tickets {
				wrHours += ticket.Hours
				wrDays += ticket.Days
				wrTotal++
			}
			mutex.Lock()
			stats.Tickets.WaitForResponse = wrTotal
			stats.Tickets.WaitForResponseDays = wrDays
			stats.Tickets.WaitForResponseHours = wrHours
			mutex.Unlock()
		}()

		// Closed tickets
		go func() {
			defer wg.Done()
			var wrHours float64
			var wrDays float64
			var wrTotal int
			limitCh <- struct{}{}
			tickets, err := s.repository.TicketClosedStat(ctx, userModel)
			<-limitCh
			if err != nil {
				closedError = err
				return
			}

			for _, ticket := range tickets {
				wrHours += ticket.Hours
				wrDays += ticket.Days
				wrTotal++
			}
			mutex.Lock()
			stats.Tickets.Closed = wrTotal
			stats.Tickets.ClosedDays = wrDays
			stats.Tickets.ClosedHours = wrHours
			mutex.Unlock()
		}()

		wg.Wait()

		if totalErr != nil {
			return stats, fmt.Errorf("error calculating total tickets: %v", totalErr)
		}
		if openErr != nil {
			return stats, fmt.Errorf("error calculating open tickets: %v", openErr)
		}
		if ipError != nil {
			return stats, fmt.Errorf("error calculating In Progress tickets: %v", ipError)
		}
		if wrError != nil {
			return stats, fmt.Errorf("error calculating In Progress tickets: %v", wrError)
		}
		if closedError != nil {
			return stats, fmt.Errorf("error calculating In Progress tickets: %v", closedError)
		}
		err = StoreInCache[*domain.Statistics]("stat-"+userModel.Crmid, &stats, CacheStatisticsTtl, s.cache)
		if err != nil {
			return stats, err
		}

		return stats, nil
	} else {
		return stats, e.Wrap("can not convert caches data to stat", err)
	}
}
