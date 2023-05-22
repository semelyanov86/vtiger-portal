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

type operation struct {
	wg      sync.WaitGroup
	mutex   sync.Mutex
	stats   *domain.Statistics
	limitCh chan struct{}
}

func (s StatisticsService) GetStatistics(ctx context.Context, userModel domain.User) (domain.Statistics, error) {
	var totalErr, openErr, ipError, wrError, closedError error

	statOperation := &operation{
		wg:      sync.WaitGroup{},
		mutex:   sync.Mutex{},
		stats:   &domain.Statistics{},
		limitCh: make(chan struct{}, 1),
	}

	statOperation.limitCh <- struct{}{}

	err := GetFromCache[*domain.Statistics]("stat-"+userModel.Crmid, statOperation.stats, s.cache)
	if err == nil {
		return *statOperation.stats, nil
	}

	if errors.Is(cache.ErrItemNotFound, err) {
		statOperation.wg.Add(5)
		// Total Tickets
		go func() {
			totalErr = s.calcTotalTickets(ctx, userModel, statOperation)
		}()

		// Open tickets
		go func() {
			openErr = s.calcOpenTickets(ctx, userModel, statOperation)
		}()

		// In progress tickets
		go func() {
			ipError = s.calcInProgressTickets(ctx, userModel, statOperation)
		}()

		// Wait For Response tickets
		go func() {
			wrError = s.calcWaitForResponseTickets(ctx, userModel, statOperation)
		}()

		// Closed tickets
		go func() {
			closedError = s.calcClosedTickets(ctx, userModel, statOperation)
		}()

		statOperation.wg.Wait()

		if totalErr != nil {
			return *statOperation.stats, fmt.Errorf("error calculating total tickets: %v", totalErr)
		}
		if openErr != nil {
			return *statOperation.stats, fmt.Errorf("error calculating open tickets: %v", openErr)
		}
		if ipError != nil {
			return *statOperation.stats, fmt.Errorf("error calculating In Progress tickets: %v", ipError)
		}
		if wrError != nil {
			return *statOperation.stats, fmt.Errorf("error calculating In Progress tickets: %v", wrError)
		}
		if closedError != nil {
			return *statOperation.stats, fmt.Errorf("error calculating In Progress tickets: %v", closedError)
		}
		err = StoreInCache[*domain.Statistics]("stat-"+userModel.Crmid, statOperation.stats, CacheStatisticsTtl, s.cache)
		if err != nil {
			return *statOperation.stats, err
		}

		return *statOperation.stats, nil
	} else {
		return *statOperation.stats, e.Wrap("can not convert caches data to stat", err)
	}
}

func (s StatisticsService) calcTotalTickets(ctx context.Context, userModel domain.User, op *operation) error {
	defer op.wg.Done()
	total, err := s.repository.CalcTicketTotal(ctx, userModel)
	<-op.limitCh
	if err != nil {
		return err
	}

	op.mutex.Lock()
	op.stats.Tickets.Total = total
	op.mutex.Unlock()
	return nil
}

func (s StatisticsService) calcOpenTickets(ctx context.Context, userModel domain.User, op *operation) error {
	defer op.wg.Done()
	var openHours float64
	var openDays float64
	var openTotal int
	op.limitCh <- struct{}{}
	tickets, err := s.repository.TicketOpenStat(ctx, userModel)
	<-op.limitCh
	if err != nil {
		return err
	}

	for _, ticket := range tickets {
		openHours += ticket.Hours
		openDays += ticket.Days
		openTotal++
	}
	op.mutex.Lock()
	op.stats.Tickets.Open = openTotal
	op.stats.Tickets.OpenDays = openDays
	op.stats.Tickets.OpenHours = openHours
	op.mutex.Unlock()
	return nil
}

func (s StatisticsService) calcInProgressTickets(ctx context.Context, userModel domain.User, op *operation) error {
	defer op.wg.Done()
	var ipHours float64
	var ipDays float64
	var ipTotal int
	op.limitCh <- struct{}{}
	tickets, err := s.repository.TicketInProgressStat(ctx, userModel)
	<-op.limitCh
	if err != nil {
		return err
	}

	for _, ticket := range tickets {
		ipHours += ticket.Hours
		ipDays += ticket.Days
		ipTotal++
	}
	op.mutex.Lock()
	op.stats.Tickets.InProgress = ipTotal
	op.stats.Tickets.InProgressDays = ipDays
	op.stats.Tickets.InProgressHours = ipHours
	op.mutex.Unlock()
	return nil
}

func (s StatisticsService) calcWaitForResponseTickets(ctx context.Context, userModel domain.User, op *operation) error {
	defer op.wg.Done()
	var wrHours float64
	var wrDays float64
	var wrTotal int
	op.limitCh <- struct{}{}
	tickets, err := s.repository.TicketWaitForResponseStat(ctx, userModel)
	<-op.limitCh
	if err != nil {
		return err
	}

	for _, ticket := range tickets {
		wrHours += ticket.Hours
		wrDays += ticket.Days
		wrTotal++
	}
	op.mutex.Lock()
	op.stats.Tickets.WaitForResponse = wrTotal
	op.stats.Tickets.WaitForResponseDays = wrDays
	op.stats.Tickets.WaitForResponseHours = wrHours
	op.mutex.Unlock()
	return nil
}

func (s StatisticsService) calcClosedTickets(ctx context.Context, userModel domain.User, op *operation) error {
	defer op.wg.Done()
	var wrHours float64
	var wrDays float64
	var wrTotal int
	op.limitCh <- struct{}{}
	tickets, err := s.repository.TicketClosedStat(ctx, userModel)
	<-op.limitCh
	if err != nil {
		return err
	}

	for _, ticket := range tickets {
		wrHours += ticket.Hours
		wrDays += ticket.Days
		wrTotal++
	}
	op.mutex.Lock()
	op.stats.Tickets.Closed = wrTotal
	op.stats.Tickets.ClosedDays = wrDays
	op.stats.Tickets.ClosedHours = wrHours
	op.mutex.Unlock()
	return nil
}
