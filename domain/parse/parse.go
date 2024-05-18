package parse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/GerogeGol/yadro-test-problem/domain/service/event"
	"github.com/GerogeGol/yadro-test-problem/domain/store"
)

var LessOrEqualZeroError = fmt.Errorf("tables count could not be less or equal 0")
var IncorrectDayTimeFormat = fmt.Errorf("incorrect DayTime format. Should be 'XX:XX'")
var IncorrectClubWorkingTimeFormat = fmt.Errorf("incorrect club working time format. Should be 'XX:XX XX:XX'")
var OpenTimeIsAfterCloseTimeError = fmt.Errorf("open time could not be after close time")
var IncorrectEventFormat = fmt.Errorf("incorrect event format")
var IncorrectClientNameFormat = fmt.Errorf("incorrect client name format. should contain only a..z letters, 0..9 numbers, '_' and '-'")

func TablesCount(s string) (int, error) {
	tablesCount, err := positiveNumber(s)
	if err != nil {
		return 0, fmt.Errorf("parse.TablesCount: %w", err)
	}
	return tablesCount, nil
}

func HourCost(s string) (float64, error) {
	hourCost, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parse.HourCost: %w", err)
	}
	if hourCost <= 0 {
		return 0, fmt.Errorf("parse.HourCost: %w", LessOrEqualZeroError)
	}
	return hourCost, nil
}

func DayTime(s string) (time store.DayTime, err error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 || len(parts[0]) != 2 || len(parts[1]) != 2 {
		err = fmt.Errorf("parse.DayTime: %w", IncorrectDayTimeFormat)
		return
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		err = fmt.Errorf("parse.DayTime: %w", err)
		return
	}

	if hours < 0 || hours >= 24 {
		err = fmt.Errorf("parse.DayTime: %w", IncorrectDayTimeFormat)
		return
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		err = fmt.Errorf("parse.DayTime: %w", err)
		return
	}

	if minutes < 0 || minutes > 60 {
		err = fmt.Errorf("parse.DayTime: %w", IncorrectDayTimeFormat)
		return
	}
	return store.NewDayTime(hours, minutes), nil
}

func ClubWorkingTime(s string) (openTime store.DayTime, closeTime store.DayTime, err error) {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		err = fmt.Errorf("parse.ClubWorkingTime: %w", IncorrectClubWorkingTimeFormat)
		return
	}

	openTime, err = DayTime(parts[0])
	if err != nil {
		err = fmt.Errorf("parse.ClubWorkingTime: %w", err)
		return

	}
	closeTime, err = DayTime(parts[1])
	if err != nil {
		err = fmt.Errorf("parse.ClubWorkingTime: %w", err)
		return
	}

	if openTime.Compare(closeTime.Time) == 1 {
		err = fmt.Errorf("parse.ClubWorkingTime: %w", OpenTimeIsAfterCloseTimeError)
		return
	}

	return
}

func ArriveEvent(s string) (e *event.ArriveEvent, err error) {
	parts := strings.Split(s, " ")
	if len(parts) != 3 {
		err = fmt.Errorf("parse.ArriveEvent: %w", IncorrectEventFormat)
		return
	}

	t, id, client, err := inputEvent(s)
	if err != nil {
		return
	}

	if id != event.ArrivalEventId {
		err = fmt.Errorf("parse.ArriveEvent: cant parse id %w", IncorrectEventFormat)
		return
	}

	return event.NewArrivalEvent(t, client), nil
}

func SitDownEvent(s string) (e *event.SitDownEvent, err error) {
	parts := strings.Split(s, " ")
	if len(parts) != 4 {
		err = fmt.Errorf("parse.SitDownEvent: %w", IncorrectEventFormat)
		return
	}
	t, id, client, err := inputEvent(s)
	if err != nil {
		err = fmt.Errorf("parse.SitDownEvent: %w", err)
		return
	}
	if id != event.SitDownEventId {
		err = fmt.Errorf("parse.SitDownEvent: incorrect id %w", IncorrectEventFormat)
		return
	}

	tableNumber, err := positiveNumber(parts[3])
	if err != nil {
		err = fmt.Errorf("parse.SitDownEvent: cant parse table number %w", err)
		return
	}

	return event.NewSitDownEvent(t, client, tableNumber), nil
}

func WaitEvent(s string) (e *event.WaitEvent, err error) {
	parts := strings.Split(s, " ")
	if len(parts) != 3 {
		err = fmt.Errorf("parse.WaitEvent: %w", IncorrectEventFormat)
		return
	}
	t, id, client, err := inputEvent(s)
	if err != nil {
		err = fmt.Errorf("parse.WaitEvent: %w", err)
		return
	}
	if id != event.WaitEventId {
		err = fmt.Errorf("parse.WaitEvent: incorrect id %w", IncorrectEventFormat)
		return
	}

	return event.NewWaitEvent(t, client), nil
}

func LeaveEvent(s string) (e *event.LeaveEvent, err error) {
	parts := strings.Split(s, " ")
	if len(parts) != 3 {
		err = fmt.Errorf("parse.LeaveEvent: %w", IncorrectEventFormat)
		return
	}
	t, id, client, err := inputEvent(s)
	if err != nil {
		err = fmt.Errorf("parse.LeaveEvent: %w", err)
		return
	}
	if id != event.LeaveEventId {
		err = fmt.Errorf("parse.LeaveEvent: incorrect id %w", IncorrectEventFormat)
		return
	}

	return event.NewLeaveEvent(t, client), nil
}

func InputEvent(s string) (e event.InputEvent, err error) {
	s = strings.Trim(s, " ")
	_, id, _, err := inputEvent(s)
	if err != nil {
		return event.EmptyInputEvent, err
	}

	var errParse error
	switch id {
	case event.ArrivalEventId:
		e, errParse = ArriveEvent(s)
	case event.SitDownEventId:
		e, errParse = SitDownEvent(s)
	case event.WaitEventId:
		e, errParse = WaitEvent(s)
	case event.LeaveEventId:
		e, errParse = LeaveEvent(s)
	default:
		return event.EmptyInputEvent, IncorrectEventFormat
	}

	if errParse != nil {
		return event.EmptyInputEvent, errParse
	}

	return e, err
}

func positiveNumber(s string) (int, error) {
	n, err := strconv.Atoi(s)

	if err != nil {
		return 0, err
	}
	if n <= 0 {
		return 0, LessOrEqualZeroError
	}

	return n, nil
}

func clientName(s string) (string, error) {
	match, _ := regexp.MatchString(`^[a-z | 1-9 |\_|\-]+$`, s)
	if !match {
		return "", IncorrectClientNameFormat
	}
	return s, nil
}

func inputEvent(s string) (t store.DayTime, id int, client string, err error) {
	parts := strings.Split(s, " ")
	if len(parts) < 3 {
		return
	}

	t, err = DayTime(parts[0])
	if err != nil {
		return
	}

	id, err = positiveNumber(parts[1])
	if err != nil {
		return
	}

	client, err = clientName(parts[2])
	if err != nil {
		return
	}
	return t, id, client, nil
}
