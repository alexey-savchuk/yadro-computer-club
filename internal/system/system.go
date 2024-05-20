package system

import (
	"fmt"
	"sort"
	"time"
)

type ComputerClubSystem struct {
	numTables        int
	openTime         time.Time
	closeTime        time.Time
	cost             int
	currentTime      time.Time
	currentEvent     Event
	gotFirstEvent    bool
	knownClients     map[string]struct{}
	waitingClients   map[string]struct{}
	clientQueue      []string
	clientPerTable   map[int]string
	tablePerClient   map[string]int
	startPlayTime    map[int]time.Time
	profitPerTable   map[int]int
	durationPerTable map[int]time.Duration
}

type eventDispatcher func(args ...any) []Event

type StatisticsRecord struct {
	TableNum     int
	Profit       int
	TimeOccupied time.Duration
}

func NewComputerClubSystem(
	numTables int,
	openTime time.Time,
	closeTime time.Time,
	cost int,
) *ComputerClubSystem {
	return &ComputerClubSystem{
		numTables:        numTables,
		cost:             cost,
		openTime:         openTime,
		closeTime:        closeTime,
		gotFirstEvent:    false,
		knownClients:     make(map[string]struct{}),
		waitingClients:   make(map[string]struct{}),
		clientPerTable:   make(map[int]string),
		tablePerClient:   make(map[string]int),
		startPlayTime:    make(map[int]time.Time),
		profitPerTable:   make(map[int]int),
		durationPerTable: make(map[int]time.Duration),
	}
}

func (s *ComputerClubSystem) Process(e Event) ([]Event, error) {
	if !s.gotFirstEvent {
		s.currentTime = e.Time
		s.gotFirstEvent = true
	}
	if e.Time.Before(s.currentTime) {
		return nil, fmt.Errorf(
			"event %q time %s is before last event time %s",
			e, e.Time.Format("15:04"), s.currentTime.Format("15:04"),
		)
	}
	s.currentTime = e.Time
	s.currentEvent = e

	var dispatcher eventDispatcher
	switch e.ID {
	case InputEventEnter:
		dispatcher = s.clientEnter
	case InputEventTakeTable:
		tableNum := e.Body[1].(int)
		if tableNum < 1 || tableNum > s.numTables {
			return nil, fmt.Errorf(
				"event %q table number %d is greater than number of tables %d",
				e, tableNum, s.numTables,
			)
		}
		dispatcher = s.clientTakeTable
		dispatcher = s.clintUnknownMiddleware(dispatcher)
	case InputEventWait:
		dispatcher = s.clientWait
		dispatcher = s.clintUnknownMiddleware(dispatcher)
	case InputEventLeave:
		dispatcher = s.clientLeave
		dispatcher = s.clintUnknownMiddleware(dispatcher)
	default:
		return nil, fmt.Errorf("unprocessable event %q", e)
	}

	dispatcher = s.closeClubMiddleware(dispatcher)
	return dispatcher(e.Body...), nil
}

func (s *ComputerClubSystem) IsClubClose() bool {
	return !(!s.currentTime.Before(s.openTime) && s.currentTime.Before(s.closeTime))
}

func (s *ComputerClubSystem) CloseClub() ([]Event, error) {
	if !s.currentTime.Before(s.closeTime) {
		return nil, fmt.Errorf("club already closed")
	}
	s.currentTime = s.closeTime
	dispatcher := s.closeClubMiddleware(func(args ...any) []Event {
		return []Event{}
	})
	return dispatcher(s.closeTime), nil
}

func (s *ComputerClubSystem) GetStatistics() []StatisticsRecord {
	records := make([]StatisticsRecord, 0, s.numTables)

	for i := 1; i <= s.numTables; i++ {
		tableProfit := 0
		if profit, ok := s.profitPerTable[i]; ok {
			tableProfit = profit
		}

		occupiedDuration := time.Duration(0)
		if duration, ok := s.durationPerTable[i]; ok {
			occupiedDuration = duration
		}

		record := StatisticsRecord{
			TableNum:     i,
			Profit:       tableProfit,
			TimeOccupied: occupiedDuration,
		}
		records = append(records, record)
	}

	return records
}

func (s *ComputerClubSystem) clientEnter(args ...any) []Event {
	name := args[0].(string)

	events := make([]Event, 0)

	if s.isClientInClub(name) {
		events = append(events, Event{
			Time: s.currentTime,
			ID:   OutputEventError,
			Body: []any{"YouShallNotPass"},
		})
		return events
	}
	if s.IsClubClose() {
		events = append(events, Event{
			Time: s.currentTime,
			ID:   OutputEventError,
			Body: []any{"NotOpenYet"},
		})
		return events
	}

	s.knownClients[name] = struct{}{}
	return events
}

func (s *ComputerClubSystem) clientTakeTable(args ...any) []Event {
	name := args[0].(string)
	tableNum := args[1].(int)

	events := make([]Event, 0)

	if _, ok := s.clientPerTable[tableNum]; ok {
		events = append(events, Event{
			Time: s.currentTime,
			ID:   OutputEventError,
			Body: []any{"PlaceIsBusy"},
		})
		return events
	}

	if freeTable, ok := s.tablePerClient[name]; ok {
		s.calcStats(freeTable)
		delete(s.clientPerTable, freeTable)
	}

	s.tablePerClient[name] = tableNum
	s.clientPerTable[tableNum] = name

	s.startPlayTime[tableNum] = s.currentTime

	s.removeFromWaiting(name)

	return []Event{}
}

func (s *ComputerClubSystem) clientWait(args ...any) []Event {
	name := args[0].(string)

	if _, ok := s.waitingClients[name]; !ok {
		s.waitingClients[name] = struct{}{}
		s.clientQueue = append(s.clientQueue, name)
	}

	events := make([]Event, 0)
	if len(s.clientPerTable) < s.numTables {
		events = append(events, Event{
			Time: s.currentTime,
			ID:   OutputEventError,
			Body: []any{
				"ICanWaitNoLonger!",
			},
		})
	}
	if len(s.waitingClients) > s.numTables {
		s.clientCleanup(name)
		events = append(events, Event{
			Time: s.currentTime,
			ID:   OutputEventLeave,
			Body: []any{
				name,
			},
		})
	}

	return events
}

func (s *ComputerClubSystem) clientLeave(args ...any) []Event {
	name := args[0].(string)

	events := make([]Event, 0)

	freeTable, ok := s.tablePerClient[name]
	if ok {
		s.calcStats(freeTable)
	}

	s.clientCleanup(name)

	if !ok {
		return events
	}

	if len(s.waitingClients) == 0 {
		return events
	}

	client := s.clientQueue[0]
	s.clientQueue = s.clientQueue[1:]

	s.removeFromWaiting(client)

	s.clientPerTable[freeTable] = client
	s.tablePerClient[client] = freeTable

	s.startPlayTime[freeTable] = s.currentTime

	events = append(events, Event{
		Time: s.currentTime,
		ID:   OutputEventTakeTable,
		Body: []any{
			client,
			freeTable,
		},
	})

	return events
}

// clintUnknownMiddleware checks if the client is in the club
func (s *ComputerClubSystem) clintUnknownMiddleware(next eventDispatcher) eventDispatcher {
	return func(args ...any) []Event {
		name := args[0].(string)

		events := make([]Event, 0)
		if !s.isClientInClub(name) {
			events = append(events, Event{
				Time: s.currentTime,
				ID:   OutputEventError,
				Body: []any{"ClientUnknown"},
			})
			return events
		}

		return next(args...)
	}
}

// closeClubMiddleware removes all clients from the club
//
// If the incoming event happens after the club is closed, all clients are
// removed from the club before processing the incoming event
func (s *ComputerClubSystem) closeClubMiddleware(next eventDispatcher) eventDispatcher {
	return func(args ...any) []Event {
		events := make([]Event, 0)

		if !s.IsClubClose() {
			events = append(events, s.currentEvent)
			return append(events, next(args...)...)
		}

		eventTime := s.currentTime
		s.currentTime = s.closeTime

		clients := make([]string, 0, len(s.knownClients))
		for name := range s.knownClients {
			clients = append(clients, name)
		}

		sort.Slice(clients, func(i, j int) bool {
			return clients[i] < clients[j]
		})

		for _, name := range clients {
			if freeTable, ok := s.tablePerClient[name]; ok {
				s.calcStats(freeTable)
			}

			events = append(events, Event{
				Time: s.closeTime,
				ID:   OutputEventLeave,
				Body: []any{
					name,
				},
			})
		}
		s.knownClients = make(map[string]struct{})
		s.clientQueue = make([]string, 0)
		s.clientPerTable = make(map[int]string)
		s.tablePerClient = make(map[string]int)
		s.startPlayTime = make(map[int]time.Time)

		s.currentTime = eventTime

		events = append(events, s.currentEvent)
		return append(events, next(args...)...)
	}
}

func (s *ComputerClubSystem) isClientInClub(name string) bool {
	_, ok := s.knownClients[name]
	return ok
}

func (s *ComputerClubSystem) removeFromWaiting(name string) {
	delete(s.waitingClients, name)

	for i, client := range s.clientQueue {
		if client == name {
			s.clientQueue = append(s.clientQueue[:i], s.clientQueue[i+1:]...)
			break
		}
	}
}

func (s *ComputerClubSystem) clientCleanup(name string) {
	delete(s.knownClients, name)
	s.removeFromWaiting(name)

	if table, ok := s.tablePerClient[name]; ok {
		delete(s.tablePerClient, name)
		delete(s.clientPerTable, table)
	}
}

func (s *ComputerClubSystem) calcStats(
	tableNum int,
) {
	startTime := s.startPlayTime[tableNum]
	endTime := s.currentTime

	duration := endTime.Sub(startTime)

	hours := int(duration.Hours())
	if duration.Minutes() > 0 {
		hours++
	}

	s.profitPerTable[tableNum] += hours * s.cost
	s.durationPerTable[tableNum] += duration
}
