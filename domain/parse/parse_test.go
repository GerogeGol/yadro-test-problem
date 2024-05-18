package parse_test

import (
	"fmt"
	"testing"

	"github.com/GerogeGol/yadro-test-problem/domain/parse"
	"github.com/GerogeGol/yadro-test-problem/domain/service/event"
	"github.com/GerogeGol/yadro-test-problem/domain/store"
	"github.com/GerogeGol/yadro-test-problem/domain/test"
)

func TestParseTablesCount(t *testing.T) {
	t.Run("correct tables count", func(t *testing.T) {
		count, err := parse.TablesCount("1")
		test.AssertNoError(t, err)
		test.AssertEqual(t, count, 1)
	})

	t.Run("incorrect tables number", func(t *testing.T) {
		_, err := parse.TablesCount("-1")
		test.AssertError(t, err, parse.LessOrEqualZeroError)

		_, err = parse.TablesCount("0")
		test.AssertError(t, err, parse.LessOrEqualZeroError)
	})
}

func TestParseDayTime(t *testing.T) {
	t.Run("correct daytime in format HH:MM", func(t *testing.T) {
		cases := []struct {
			input string
			want  store.DayTime
		}{
			{"09:00", store.NewDayTime(9, 0)},
			{"00:00", store.NewDayTime(0, 0)},
			{"23:59", store.NewDayTime(23, 59)},
			{"15:00", store.NewDayTime(15, 0)},
		}
		for i, c := range cases {
			t.Run(fmt.Sprintf("Case %d: %q", i, c.input), func(t *testing.T) {
				time, err := parse.DayTime(c.input)
				test.AssertNoError(t, err)
				test.AssertEqual(t, time, c.want)
			})
		}
	})
	t.Run("incorrect daytime", func(t *testing.T) {
		cases := []string{
			"9:00",
			"25:00",
			"09:90",
			"24:00",
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("Case %d: %q", i, c), func(t *testing.T) {
				_, err := parse.DayTime(c)
				test.AssertNotNilError(t, err)
			})
		}
	})
}

func TestParseClubWorkingTime(t *testing.T) {
	t.Run("correct open and close time", func(t *testing.T) {
		cases := []struct {
			input     string
			openTime  store.DayTime
			closeTime store.DayTime
		}{
			{"09:00 19:00", store.NewDayTime(9, 0), store.NewDayTime(19, 0)},
			{"00:00 23:59", store.NewDayTime(0, 0), store.NewDayTime(23, 59)},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("Case %d: %s", i, c.input), func(t *testing.T) {
				openTime, closeTime, err := parse.ClubWorkingTime(c.input)
				test.AssertNoError(t, err)
				test.AssertEqual(t, openTime, c.openTime)
				test.AssertEqual(t, closeTime, c.closeTime)
			})
		}
	})
	t.Run("incorrect club working time", func(t *testing.T) {
		cases := []struct {
			input string
			err   error
		}{
			{"9:00 19:00", parse.IncorrectDayTimeFormat},
			{"23:00 09:59", parse.OpenTimeIsAfterCloseTimeError},
			{"00:00  12:59", parse.IncorrectClubWorkingTimeFormat},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("Case %d: %s", i, c.input), func(t *testing.T) {
				_, _, err := parse.ClubWorkingTime(c.input)
				test.AssertError(t, err, c.err)
			})
		}
	})
}

func TestParseHourCost(t *testing.T) {
	t.Run("correct hour cost", func(t *testing.T) {
		count, err := parse.HourCost("1")
		test.AssertNoError(t, err)
		test.AssertEqual(t, count, 1)
	})

	t.Run("incorrect hour cost", func(t *testing.T) {
		_, err := parse.HourCost("-1")
		test.AssertError(t, err, parse.LessOrEqualZeroError)

		_, err = parse.HourCost("0")
		test.AssertError(t, err, parse.LessOrEqualZeroError)
	})
}

func TestParseArriveEvent(t *testing.T) {
	t.Run("parse correct events", func(t *testing.T) {
		cases := []struct {
			input string
			ev    event.ArriveEvent
		}{
			{"08:48 1 client1", *event.NewArrivalEvent(store.NewDayTime(8, 48), "client1")},
			{"12:48 1 client2", *event.NewArrivalEvent(store.NewDayTime(12, 48), "client2")},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("Case: %d, %q", i, c.input), func(t *testing.T) {
				e, err := parse.ArriveEvent(c.input)
				test.AssertNoError(t, err)
				test.AssertEqual(t, e.Time(), c.ev.Time())
				test.AssertEqual(t, e.Id(), c.ev.Id())
				test.AssertEqual(t, e.Client(), c.ev.Client())
			})
		}
	})

	t.Run("incorrect events format", func(t *testing.T) {
		cases := []struct {
			input string
			err   error
		}{
			{"08:48 112 client1", parse.IncorrectEventFormat},
			{"2:48 1 client2", parse.IncorrectDayTimeFormat},
			{"12:48  1  client2", parse.IncorrectEventFormat},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("Case: %d, %q", i, c.input), func(t *testing.T) {
				_, err := parse.ArriveEvent(c.input)
				test.AssertError(t, err, c.err)
			})
		}
	})
}

func TestParseSitDownEvent(t *testing.T) {
	t.Run("parse correct events", func(t *testing.T) {
		cases := []struct {
			input string
			ev    event.SitDownEvent
		}{
			{"08:48 2 client1 1", *event.NewSitDownEvent(store.NewDayTime(8, 48), "client1", 1)},
			{"12:48 2 client2 2", *event.NewSitDownEvent(store.NewDayTime(12, 48), "client2", 2)},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("Case: %d, %q", i, c.input), func(t *testing.T) {
				e, err := parse.SitDownEvent(c.input)
				test.AssertNoError(t, err)
				test.AssertEqual(t, e.Time(), c.ev.Time())
				test.AssertEqual(t, e.Id(), c.ev.Id())
				test.AssertEqual(t, e.Client(), c.ev.Client())
				test.AssertEqual(t, e.Table(), c.ev.Table())
			})
		}
	})
}

func TestParseWaitEvent(t *testing.T) {
	t.Run("parse correct events", func(t *testing.T) {
		cases := []struct {
			input string
			ev    event.WaitEvent
		}{
			{"08:48 3 client1", *event.NewWaitEvent(store.NewDayTime(8, 48), "client1")},
			{"12:48 3 client2", *event.NewWaitEvent(store.NewDayTime(12, 48), "client2")},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("Case: %d, %q", i, c.input), func(t *testing.T) {
				e, err := parse.WaitEvent(c.input)
				test.AssertNoError(t, err)
				test.AssertEqual(t, e.Time(), c.ev.Time())
				test.AssertEqual(t, e.Id(), c.ev.Id())
				test.AssertEqual(t, e.Client(), c.ev.Client())
			})
		}
	})
}

func TestParseLeaveEvent(t *testing.T) {
	t.Run("parse correct events", func(t *testing.T) {
		cases := []struct {
			input string
			ev    event.LeaveEvent
		}{
			{"08:48 4 client1", *event.NewLeaveEvent(store.NewDayTime(8, 48), "client1")},
			{"12:48 4 client2", *event.NewLeaveEvent(store.NewDayTime(12, 48), "client2")},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("Case: %d, %q", i, c.input), func(t *testing.T) {
				e, err := parse.LeaveEvent(c.input)
				test.AssertNoError(t, err)
				test.AssertEqual(t, e.Time(), c.ev.Time())
				test.AssertEqual(t, e.Id(), c.ev.Id())
				test.AssertEqual(t, e.Client(), c.ev.Client())
			})
		}
	})
}

func TestParseEvent(t *testing.T) {
	t.Run("parse correct events", func(t *testing.T) {
		cases := []struct {
			input string
			ev    event.Event
		}{
			{"08:48 1 client1", *event.NewArrivalEvent(store.NewDayTime(8, 48), "client1")},
			{"08:48 2 client1 2", *event.NewSitDownEvent(store.NewDayTime(8, 48), "client1", 1)},
			{"08:48 3 client1", *event.NewWaitEvent(store.NewDayTime(8, 48), "client1")},
			{"12:48 4 client2", *event.NewLeaveEvent(store.NewDayTime(12, 48), "client2")},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("Case: %d, %q", i, c.input), func(t *testing.T) {
				e, err := parse.InputEvent(c.input)
				test.AssertNoError(t, err)
				test.AssertEqual(t, e.Time(), c.ev.Time())
				test.AssertEqual(t, e.Id(), c.ev.Id())
			})
		}
	})

	t.Run(" incorrect events", func(t *testing.T) {
		cases := []struct {
			input string
		}{
			{"08:48 1 client1 2"},
			{"08:48 2 client1!"},
			{"8:48 3 client1"},
			{"12:8 4 client2"},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("Case: %d, %q", i, c.input), func(t *testing.T) {
				e, err := parse.InputEvent(c.input)
				test.AssertNotNilError(t, err)
				test.AssertEmptyEvent(t, e)
			})
		}
	})
}
