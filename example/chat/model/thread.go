package model

import (
	"errors"
	"time"

	"github.com/teris-io/shortid"
)

type Thread struct {
	id              string
	uncommitted     []Event
	started, closed bool
	moderator       Person
	members         Participants
	messages        int

	settings struct {
		delay time.Duration
	}

	// recalls ...
	recalls struct {
		after    time.Duration
		max, num int
	}
}

func NewThread(id string) *Thread {
	return &Thread{id: id, members: make(Participants)}
}

func (t *Thread) Start(channel string, p Person) error {
	switch {
	case t.started:
		return nil

	case t.closed:
		return ErrClosed
	}

	return t.commit(
		ThreadStarted{p, channel},
		ThreadJoined{p},
	)
}

func (t *Thread) Join(p Person) error {
	switch {
	case !t.started:
		return ErrNotExists

	case t.closed:
		return ErrClosed

	case t.members.isKicked(p):
		return ErrParticipantKicked

	case t.members.isPresent(p):
		return nil

	}

	return t.commit(
		ThreadJoined{p},
	)
}

func (t *Thread) Message(p Person, text string) error {
	switch {
	case !t.started:
		return ErrNotExists

	case t.closed:
		return ErrClosed

	case !t.members.isPresent(p):
		return ErrParticipantNotFound

	}

	return t.commit(
		ThreadMessage{shortid.MustGenerate(), p, text},
	)
}

func (t *Thread) Reply(on, participant, text string) error {
	return nil
}

func (t *Thread) Leave(p Person) error {
	switch {
	case !t.started:
		return ErrNotExists

	case t.closed:
		return ErrClosed

	case !t.members.isPresent(p):
		return nil
	}

	return t.commit(
		ThreadLeft{p},
	)
}

func (t *Thread) Kick(moderator, p Person) error {
	switch {
	case !t.started:
		return ErrNotExists

	case t.closed:
		return ErrClosed

	case !t.members.isPresent(p):
		return ErrParticipantNotFound

	case t.moderator != moderator:
		return ErrModeratorRequired
	}

	return t.commit(
		ThreadKicked{p},
	)
}

func (t *Thread) Mute(p Person, reason string) error {
	if !t.members.isPresent(p) {
		return nil
	}

	return t.commit(
		ThreadMuted{p, reason},
	)
}

func (t *Thread) Close(moderator Person) error {
	switch {
	case !t.started:
		return ErrNotExists

	case t.closed:
		return nil

	case t.moderator != moderator:
		return ErrModeratorRequired
	}

	return t.commit(
		ThreadClosed{},
	)
}

func (t *Thread) Awake() error {
	if t.closed {
		return nil
	}

	if t.messages > 3 {
		return t.commit(
			ThreadClosed{},
		)
	}

	return t.commit(recalled{})
}

func (t *Thread) AwakeAfter() time.Time {
	if t.recalls.num >= t.recalls.max {
		return time.Time{} // turn off recall
	}

	return time.Now().Add(t.recalls.after)
}

func (t *Thread) Commit(e Event, at time.Time) error {
	if t.settings.delay > 0 {
		time.Sleep(t.settings.delay)
	}

	switch e := e.(type) {
	case ThreadStarted:
		t.started, t.moderator = true, e.Moderator

	case ThreadJoined:
		t.members.join(e.Participant)

	case ThreadMessage:
		t.messages++

	case ThreadMuted:

	case ThreadLeft:
		t.members.remove(e.Participant)

	case ThreadKicked:
		t.members.kick(e.Participant)

	case ThreadClosed:
		t.closed = true

	case recalled:
		t.recalls.num++
	}

	return nil
}

func (t *Thread) Author() string {
	return "tomek@zklanu.pl:6361"
}

func (t *Thread) Uncommitted(clear bool) []Event {
	defer func() {
		if clear {
			t.uncommitted = []Event{}
		}
	}()

	return t.uncommitted
}

func (t *Thread) commit(events ...Event) error {
	t.uncommitted = append(t.uncommitted, events...)
	return nil
}

type Person string

type Participants map[Person]bool

func (p Participants) join(n Person) {
	p[n] = true
}

func (p Participants) isPresent(n Person) bool {
	_, found := p[n]
	return found
}

func (p Participants) isKicked(n Person) bool {
	if active, found := p[n]; !active && found {
		return true
	}

	return false
}

func (p Participants) remove(n Person) {
	delete(p, n)
}

func (p Participants) kick(n Person) {
	p[n] = false
}

func (p Participants) count() int {
	return len(p)
}

var (
	ErrNotExists           = errors.New("THREAD: not started or not exists")
	ErrClosed              = errors.New("THREAD: closed")
	ErrModeratorRequired   = errors.New("THREAD: only moderator can do it")
	ErrParticipantNotFound = errors.New("THREAD: participant not found")
	ErrParticipantKicked   = errors.New("THREAD: participant has been kicked")
)

type (
	// events
	Event = any

	ThreadStarted struct {
		Moderator Person
		Channel   string
	}
	ThreadJoined struct {
		Participant Person
	}
	ThreadMessage struct {
		ID          string
		Participant Person
		Text        string
	}
	ThreadLeft struct {
		Participant Person
	}
	ThreadKicked struct {
		Participant Person
	}
	ThreadMuted struct {
		Participant Person
		Reason      string
	}
	ThreadClosed struct{}

	recalled struct{}
)
