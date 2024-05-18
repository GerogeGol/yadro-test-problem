package scan

import (
	"bufio"
	"fmt"
	"io"

	"github.com/GerogeGol/yadro-test-problem/domain/parse"
	memqueue "github.com/GerogeGol/yadro-test-problem/domain/queue/memory"
	"github.com/GerogeGol/yadro-test-problem/domain/service"
	"github.com/GerogeGol/yadro-test-problem/domain/service/event"
	"github.com/GerogeGol/yadro-test-problem/domain/store"
	memstore "github.com/GerogeGol/yadro-test-problem/domain/store/memory"
)

type FileScanner struct {
	*bufio.Scanner
	lastLine string
	lastTime store.DayTime
}

func (s *FileScanner) Scan() bool {
	scanRes := s.Scanner.Scan()
	if !scanRes && s.Scanner.Err() != nil {
		panic("Error while scanning file")
	}
	s.lastLine = s.Scanner.Text()
	return scanRes
}

func (s *FileScanner) ScanTablesCount() (int, error) {

	return parse.TablesCount(s.lastLine)
}

func (s *FileScanner) ScanClubWorkingTime() (openTime store.DayTime, closeTime store.DayTime, err error) {
	return parse.ClubWorkingTime(s.lastLine)
}

func (s *FileScanner) ScanHourCost() (float64, error) {
	return parse.HourCost(s.lastLine)
}

func (s *FileScanner) ScanInputEvent() (e event.InputEvent, err error) {
	e, err = parse.InputEvent(s.lastLine)
	if err != nil {
		return
	}
	if e.Time().Compare(s.lastTime.Time) == -1 {
		err = fmt.Errorf("inconsisten time in events")
		return
	}
	s.lastTime = e.Time()
	return
}

func ScanInputData(r io.Reader, b io.Writer) (string, error) {
	scanner := &FileScanner{Scanner: bufio.NewScanner(r)}

	scanner.Scan()
	tablesCount, err := scanner.ScanTablesCount()
	if err != nil {
		return scanner.lastLine, err
	}
	b.Write([]byte(fmt.Sprintln(tablesCount)))

	scanner.Scan()
	openTime, closeTime, err := scanner.ScanClubWorkingTime()
	if err != nil {
		return scanner.lastLine, err
	}
	b.Write([]byte(fmt.Sprintln(openTime, closeTime)))

	scanner.Scan()
	hourCost, err := scanner.ScanHourCost()
	if err != nil {
		return scanner.lastLine, err
	}
	b.Write([]byte(fmt.Sprintln(hourCost)))

	cc := service.NewComputerClub(tablesCount, hourCost, openTime, closeTime, memstore.NewStore(), memqueue.NewQueue())
	s := service.NewService(cc)

	b.Write([]byte(fmt.Sprintln(openTime)))
	for scanner.Scan() {
		e, err := scanner.ScanInputEvent()
		if err != nil {
			return scanner.lastLine, err
		}

		outEvent := s.ServeEvent(e)
		b.Write([]byte(fmt.Sprintln(e)))
		if event.IsEmpty(outEvent) {
			continue
		}

		errEvent, ok := outEvent.(*event.ErrorEvent)
		if !ok {
			b.Write([]byte(fmt.Sprintln(outEvent)))
			continue
		}

		switch errEvent.Err() {
		case service.ClientUnknown, service.PlaceIsBusy, service.NotOpenYet, service.ICanWaitNoLonger, service.YouShallNotPass:
			b.Write([]byte(fmt.Sprintln(errEvent)))
			continue
		default:
			return scanner.lastLine, errEvent.Err()
		}
	}
	leaveEvents, err := s.Close()
	if err != nil {
		panic(err)
	}
	for _, e := range leaveEvents {
		b.Write([]byte(fmt.Sprintln(&e)))
	}
	b.Write([]byte(fmt.Sprintln(closeTime)))

	tableInfos, err := s.Profit()
	if err != nil {
		panic(err)
	}
	for _, info := range tableInfos {
		b.Write([]byte(fmt.Sprintln(&info)))
	}

	return "", nil
}
