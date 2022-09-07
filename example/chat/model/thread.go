package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/teris-io/shortid"
)

type Thread struct {
	id              string
	uncommitted     []Event
	started, closed bool
	moderator       Member
	members         Members
	messages        int

	settings struct {
		delay time.Duration
	}

	// recalls ...
	recalls struct {
		after    time.Duration
		max, num int
	}

	version int64
}

func NewThread(id string) (*Thread, error) {
	return &Thread{id: id, members: make(Members)}, nil
}

func (t *Thread) ID() string {
	return t.id
}

func (t *Thread) Version() int64 {
	return t.version
}

func (t *Thread) Start(channel string, p Member) error {
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

func (t *Thread) Join(p Member) error {
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

func (t *Thread) Message(p Member, text string) error {
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

func (t *Thread) Leave(p Member) error {
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

func (t *Thread) Kick(moderator, p Member) error {
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

func (t *Thread) Mute(p Member, reason string) error {
	if !t.members.isPresent(p) {
		return nil
	}

	return t.commit(
		ThreadMuted{p, reason},
	)
}

func (t *Thread) Close(moderator Member) error {
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

	t.version++
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

func (t *Thread) String() string {
	if s := len(t.uncommitted); s != 0 {
		return fmt.Sprintf("%s.Thread#%d->%d", t.id, t.version, t.version+int64(len(t.uncommitted)))
	}
	return fmt.Sprintf("%s.Thread#%d", t.id, t.version)
}

func (t *Thread) commit(events ...Event) error {
	t.uncommitted = append(t.uncommitted, events...)
	return nil
}

var (
	ErrNotExists           = errors.New("THREAD: not started or not exists")
	ErrClosed              = errors.New("THREAD: closed")
	ErrModeratorRequired   = errors.New("THREAD: only moderator can do it")
	ErrParticipantNotFound = errors.New("THREAD: participant not found")
	ErrParticipantKicked   = errors.New("THREAD: participant has been kicked")
)

type (
	Event = any

	ThreadStarted struct {
		Moderator Member
		Channel   string
	}
	ThreadJoined struct {
		Participant Member
	}
	ThreadMessage struct {
		ID          string
		Participant Member
		Text        string
	}
	ThreadLeft struct {
		Participant Member
	}
	ThreadKicked struct {
		Participant Member
	}
	ThreadMuted struct {
		Participant Member
		Reason      string
	}
	ThreadClosed struct{}

	recalled struct{}
)
