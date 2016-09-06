package audio

import (
	"errors"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/internal/testutils"
	v2 "github.com/CyCoreSystems/ari/v2"

	"golang.org/x/net/context"
)

func TestPlayAsyncWithTimeout(t *testing.T) {

	MaxPlaybackTime = 3 * time.Second

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	pb, err := PlayAsync(bus, player, "audio:hello-world")
	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if pb == nil {
		t.Errorf("Expected playback object to be non-nil")
		return
	}

	if pb.Handle() == nil {
		t.Errorf("Expected playback.Handle to be non-nil")
	}

	select {
	case <-pb.StartCh():
		t.Errorf("Unexpected trigger of Start channel")
	case <-pb.StopCh():
		t.Errorf("Unexpected trigger of Stop channel")
	case <-time.After(1 * time.Second):
	}

	// wait for timeout
	<-time.After(MaxPlaybackTime)

	select {
	case <-pb.StartCh():
	default:
		t.Errorf("Expected trigger of start channel after MaxPlaybackTime")
	}

	select {
	case <-pb.StopCh():
	default:
		t.Errorf("Expected trigger of stop channel after MaxPlaybackTime")
	}

	if !isTimeout(pb.Err()) {
		t.Errorf("Expected timeout error, got: '%v'", pb.Err())
	}

}

func TestPlayAsyncQuit(t *testing.T) {

	MaxPlaybackTime = 3 * time.Second

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	pb, err := PlayAsync(bus, player, "audio:hello-world")
	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if pb == nil {
		t.Errorf("Expected playback object to be non-nil")
		return
	}
	if pb.Handle() == nil {
		t.Errorf("Expected playback.Handle to be non-nil")
	}

	pb.Cancel()

	select {
	case <-pb.StartCh():
	case <-time.After(1 * time.Second):
		t.Errorf("Expected trigger of Start channel")
	}

	select {
	case <-pb.StopCh():
	case <-time.After(1 * time.Second):
		t.Errorf("Expected trigger of Stop channel")
	}

	// wait for timeout
	<-time.After(MaxPlaybackTime)

	if err := pb.Err(); err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}
}

func TestPlayAsyncQuitAfterStart(t *testing.T) {

	MaxPlaybackTime = 3 * time.Second

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	pb, err := PlayAsync(bus, player, "audio:hello-world")
	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if pb == nil {
		t.Errorf("Expected playback object to be non-nil")
		return
	}

	if pb.Handle() == nil {
		t.Errorf("Expected playback.Handle to be non-nil")
	}

	bus.Send(playbackStartedGood("pb1"))
	pb.Cancel()

	select {
	case <-pb.StartCh():
	case <-time.After(1 * time.Second):
		t.Errorf("Expected trigger of Start channel")
	}

	select {
	case <-pb.StopCh():
	case <-time.After(1 * time.Second):
		t.Errorf("Expected trigger of Stop channel")
	}

	// wait for timeout
	<-time.After(MaxPlaybackTime)

	if err := pb.Err(); err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}
}

func TestPlayTimeoutStart(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	err := Play(ctx, bus, player, "audio:hello-world")

	if !isTimeout(err) {
		t.Errorf("Expected timeout error, got: '%v'", err)
	}

	if err != nil && err.Error() != "Timeout waiting for start of playback" {
		t.Errorf("Expected timeout waiting for start of playback error, got: '%v'", err)
	}
}

func TestPlayTimeoutStop(t *testing.T) {
	MaxPlaybackTime = 3 * time.Millisecond

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	pb, err := PlayAsync(bus, player, "audio:hello-world")
	if err != nil {
		t.Errorf("Unexpected error playing audio: '%v'", err)
	}

	bus.Send(playbackStartedGood("pb1"))

	select {
	case <-pb.StopCh():
	case <-time.After(3 * time.Second):
		t.Errorf("Expected Stop channel to trigger")
	}

	if !isTimeout(pb.Err()) {
		t.Errorf("Expected timeout error, got: '%v'", err)
	}

	if pb.Err() != nil && pb.Err().Error() != "Timeout waiting for stop of playback" {
		t.Errorf("Expected timeout waiting for stop of playback error, got: '%v'", pb.Err())
	}
}

func TestPlayTimeoutStop100(t *testing.T) {
	for i := 0; i != 100; i++ {
		TestPlayTimeoutStop(t)
	}
}

func TestPlaySuccess(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	pb, err := PlayAsync(bus, player, "audio:hello-world")
	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}
	defer pb.Cancel()

	bus.Send(playbackStartedGood("pb1"))
	bus.Send(playbackFinishedGood("pb1"))

	select {
	case <-pb.StopCh():
	case <-time.After(3 * time.Second):
		t.Errorf("Expected Stop channel to trigger")
	}

	if pb.Err() != nil {
		t.Errorf("Unexpected error: '%v'", pb.Err())
	}
}

func TestPlayNilEvents(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	pb, err := PlayAsync(bus, player, "audio:hello-world")
	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}
	defer pb.Cancel()

	bus.SendTo("PlaybackStarted", nil)
	bus.Send(playbackStartedGood("pb1"))

	select {
	case <-pb.StartCh():
	case <-time.After(3 * time.Second):
		t.Errorf("Expected Start channel to trigger")
	}

	bus.SendTo("PlaybackStarted", nil)
	bus.SendTo("PlaybackFinished", nil)
	bus.Send(playbackFinishedGood("pb1"))

	select {
	case <-pb.StopCh():
	case <-time.After(3 * time.Second):
		t.Errorf("Expected Stop channel to trigger")
	}

	if pb.Err() != nil {
		t.Errorf("Unexpected error: '%v'", pb.Err())
	}
}

func TestPlayUnrelatedEvents(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	pb, err := PlayAsync(bus, player, "audio:hello-world")
	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	bus.SendTo("PlaybackStarted", playbackStartedBadMessageType)
	bus.Send(playbackFinishedDifferentPlaybackID)
	bus.Send(playbackStartedDifferentPlaybackID)
	bus.Send(playbackStartedGood("pb1"))

	select {
	case <-pb.StartCh():
	case <-time.After(1 * time.Millisecond):
		t.Errorf("Expected start channel to trigger")
	}

	bus.SendTo("PlaybackFinished", playbackFinishedBadMessageType)
	bus.Send(playbackFinishedDifferentPlaybackID)

	select {
	case <-pb.StopCh():
		t.Errorf("Unexpected stop channel trigger ")
	default:
	}

	bus.Send(playbackFinishedGood("pb1"))

	select {
	case <-pb.StopCh():
	case <-time.After(1 * time.Millisecond):
		t.Errorf("Expected stop channel to trigger")
	}

	if err = pb.Err(); err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}
}

func TestPlayStopBeforeStart(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	go func() {
		bus.Send(playbackFinishedGood("pb1"))
	}()

	pb, err := PlayAsync(bus, player, "audio:hello-world")
	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	select {
	case <-pb.StopCh():
	case <-time.After(1 * time.Second):
		t.Errorf("Expected trigger of stop channel")
	}

	if pb.Err() != nil {
		t.Errorf("Unexpected error: '%v'", pb.Err())
	}
}

func TestContextCancellation(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	cancel()

	err := Play(ctx, bus, player, "audio:hello-world")

	if err == nil {
		t.Errorf("Expected error, got nil")
	} else if err.Error() != "context canceled" { //TODO: should be an interface to cast to here instead of string comparison
		t.Errorf("Expected context cancellation error, got '%v'", err)
	}
}

func TestContextCancellation100(t *testing.T) {
	for i := 0; i != 100; i++ {
		TestContextCancellation(t)
	}
}

func TestContextCancellationAfterStart(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	go func() {
		bus.Send(playbackStartedGood("pb1"))
		cancel()
	}()

	err := Play(ctx, bus, player, "audio:hello-world")

	if err == nil {
		t.Errorf("Expected error, got nil")
	} else if err.Error() != "context canceled" { //TODO: should be an interface to cast to here instead of string comparison
		t.Errorf("Expected context cancellation error, got '%v'", err)
	}
}

func TestContextCancellationAfterStart100(t *testing.T) {
	for i := 0; i != 100; i++ {
		TestContextCancellationAfterStart(t)
	}
}

func TestErrorInPlayer(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(nil, errors.New("Dummy error playing to dummy player"))

	err := Play(ctx, bus, player, "audio:hello-world")

	if err == nil {
		t.Errorf("Expected error, got nil")
	} else if err.Error() != "Dummy error playing to dummy player" {
		t.Errorf("Expected dummy error, got '%v'", err)
	}
}

func TestErrorInDataFetch(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1", failData: true}), nil)

	err := Play(ctx, bus, player, "audio:hello-world")

	if err == nil {
		t.Errorf("Expected error, got nil")
	} else if err.Error() != "Dummy error getting playback data" {
		t.Errorf("Expected dummy error, got '%v'", err)
	}
}

// messages

var channelDtmf = func(dtmf string) v2.Eventer {
	return &v2.ChannelDtmfReceived{
		Event: v2.Event{
			Message: v2.Message{
				Type: "ChannelDtmfReceived",
			},
		},
		Digit: dtmf,
	}
}

var playbackStartedGood = func(id string) v2.Eventer {
	return &v2.PlaybackStarted{
		Event: v2.Event{
			Message: v2.Message{
				Type: "PlaybackStarted",
			},
		},
		Playback: v2.Playback{
			ID: id,
		},
	}
}

var playbackFinishedGood = func(id string) v2.Eventer {
	return &v2.PlaybackFinished{
		Event: v2.Event{
			Message: v2.Message{
				Type: "PlaybackFinished",
			},
		},
		Playback: v2.Playback{
			ID: id,
		},
	}
}

var playbackStartedBadMessageType = &v2.PlaybackStarted{
	Event: v2.Event{
		Message: v2.Message{
			Type: "PlaybackStarted2",
		},
	},
	Playback: v2.Playback{
		ID: "pb1",
	},
}

var playbackFinishedBadMessageType = &v2.PlaybackFinished{
	Event: v2.Event{
		Message: v2.Message{
			Type: "PlaybackFinished2",
		},
	},
	Playback: v2.Playback{
		ID: "pb1",
	},
}

var playbackStartedDifferentPlaybackID = &v2.PlaybackStarted{
	Event: v2.Event{
		Message: v2.Message{
			Type: "PlaybackStarted",
		},
	},
	Playback: v2.Playback{
		ID: "pb2",
	},
}

var playbackFinishedDifferentPlaybackID = &v2.PlaybackFinished{
	Event: v2.Event{
		Message: v2.Message{
			Type: "PlaybackFinished",
		},
	},
	Playback: v2.Playback{
		ID: "pb2",
	},
}

// test playback ari transport

type testPlayback struct {
	id       string
	failData bool
}

func (p *testPlayback) Get(id string) *ari.PlaybackHandle {
	panic("not implemented")
}

func (p *testPlayback) Data(id string) (pd ari.PlaybackData, err error) {
	if p.failData {
		err = errors.New("Dummy error getting playback data")
	}
	pd.ID = p.id
	return
}

func (p *testPlayback) Control(id string, op string) error {
	panic("not implemented")
}

func (p *testPlayback) Stop(id string) error {
	panic("not implemented")
}

func isTimeout(err error) bool {

	type timeout interface {
		Timeout() bool
	}

	te, ok := err.(timeout)
	return ok && te.Timeout()
}
