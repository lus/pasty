package cleanup

import (
	"context"
	"github.com/lus/pasty/internal/pastes"
	"github.com/rs/zerolog/log"
	"time"
)

type Task struct {
	Interval    time.Duration
	MaxPasteAge time.Duration
	Repository  pastes.Repository

	running bool
	stop    chan struct{}
}

func (task *Task) Start() {
	if task.running {
		return
	}
	task.stop = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-time.After(task.Interval):
				n, err := task.Repository.DeleteOlderThan(context.Background(), task.MaxPasteAge)
				if err != nil {
					log.Err(err).Msg("Could not clean up expired pastes.")
					continue
				}
				log.Debug().Int("amount", n).Msg("Cleaned up expired pastes.")
			case <-task.stop:
				task.running = false
				return
			}
		}
	}()
	task.running = true
}

func (task *Task) Stop() {
	if !task.running {
		return
	}
	close(task.stop)
}
