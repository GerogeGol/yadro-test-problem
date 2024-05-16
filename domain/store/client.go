package store

import (
	"math"
	"time"
)

type Client struct {
	Name         string
	Table        int
	PlayingSince DayTime
}

func (c *Client) PlayingTime(t DayTime) time.Duration {
	return t.Sub(c.PlayingSince.Time)
}

func (c *Client) Payment(t DayTime, moneyPerHour float64) float64 {
	playingTime := c.PlayingTime(t)
	playedHours := math.Round(playingTime.Hours())
	return float64(playedHours) * moneyPerHour
}
