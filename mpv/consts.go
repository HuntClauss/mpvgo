package mpv

import "errors"

// Format represents supported formats by mpv API
// used for retrieving and setting options and properties.
type Format int

const (
	FormatNone      Format = 0 // FormatNone is used to represent invalid or in some cases empty values.
	FormatString    Format = 1 // FormatString is used for string values
	FormatOsdString Format = 2 // FormatOsdString is used for OSD string values
	FormatFlag      Format = 3 // FormatFlag is used for boolean values
	FormatBoolean   Format = 3 // FormatBoolean is alias for FormatFlag
	FormatInt64     Format = 4 // FormatInt64 is used for integer values
	FormatDouble    Format = 5 // FormatDouble is used for floating point values
	FormatFloat64   Format = 5 // FormatFloat64 is alias for FormatDouble
	FormatNode      Format = 6 // FormatNode is used for Node values
	FormatNodeArray Format = 7 // FormatNodeArray is used for NodeList values
	FormatNodeMap   Format = 8 // FormatNodeMap is used for NodeMap values
	FormatByteArray Format = 9 // FormatByteArray is used for byte array values
)

type EventID int

const (
	EventNone             EventID = 0
	EventShutdown         EventID = 1
	EventLogMessage       EventID = 2
	EventGetPropertyReply EventID = 3
	EventSetPropertyReply EventID = 4
	EventCommandReply     EventID = 5
	EventStartFile        EventID = 6
	EventEndFile          EventID = 7
	EventFileLoaded       EventID = 8
	EventIdle             EventID = 11 // Deprecated: use ObserveProperty("idle-activity") instead
	EventPause            EventID = 12 // Deprecated: use ObserveProperty("pause") instead
	EventUnpause          EventID = 13 // Deprecated: use ObserveProperty("unpause") instead
	EventTick             EventID = 14 // Deprecated: use for example ObserveProperty("playback-time") instead
	EventClientMessage    EventID = 16
	EventVideoReconfig    EventID = 17
	EventAudioReconfig    EventID = 18
	EventSeek             EventID = 20
	EventPlaybackRestart  EventID = 21
	EventPropertyChange   EventID = 22
	EventQueueOverflow    EventID = 24
	EventHook             EventID = 25
)

type LogLevel int

const (
	LogLevelNone  LogLevel = 0  // "no"    - disable absolutely all messages
	LogLevelFatal LogLevel = 10 // "fatal" - critical/aborting errors
	LogLevelError LogLevel = 20 // "error" - simple errors
	LogLevelWarn  LogLevel = 30 // "warn"  - possible problems
	LogLevelInfo  LogLevel = 40 // "info"  - informational message
	LogLevelV     LogLevel = 50 // "v"     - noisy informational message
	LogLevelDebug LogLevel = 60 // "debug" - very noisy technical information
	LogLevelTrace LogLevel = 70 // "trace" - extremely noisy
)

func (l LogLevel) String() string {
	return []string{"no", "fatal", "error", "warn", "info", "v", "debug", "trace"}[l/10]
}

type EndFileReason int

const (
	EndFileReasonEof      EndFileReason = 0
	EndFileReasonStop     EndFileReason = 2
	EndFileReasonQuit     EndFileReason = 3
	EndFileReasonError    EndFileReason = 4
	EndFileReasonRedirect EndFileReason = 5
)

// Error this is not used anywhere, but maybe someone finds use case for it
//
// Every error code is wrapped with go standard error package
type Error int

const (
	ErrSuccess             = 0
	ErrEventQueueFull      = -1
	ErrNomem               = -2
	ErrUninitialized       = -3
	ErrInvalidParameter    = -4
	ErrOptionNotFound      = -5
	ErrOptionFormat        = -6
	ErrOptionError         = -7
	ErrPropertyNotFound    = -8
	ErrPropertyFormat      = -9
	ErrPropertyUnavailable = -10
	ErrPropertyError       = -11
	ErrCommand             = -12
	ErrLoadingFailed       = -13
	ErrAoInitFailed        = -14
	ErrVoInitFailed        = -15
	ErrNothingToPlay       = -16
	ErrUnknownFormat       = -17
	ErrUnsupported         = -18
	ErrNotImplemented      = -19
	ErrGeneric             = -20
)

func (e Error) Err() error {
	return []error{
		nil,
		errors.New("event ringbuffer is full and can't receive any events"),
		errors.New("memory allocation failed"),
		errors.New("mpv core is not initialized"),
		errors.New("invalid or unsupported parameter value"),
		errors.New("option does not exists"),
		errors.New("unsupported FORMAT of option"),
		errors.New("provided option value could not be parsed"),
		errors.New("accessed property does not exist"),
		errors.New("usage of unsupported FORMAT"),
		errors.New("property exists, but is currently unavailable"),
		errors.New("something went wrong when setting or getting a property"),
		errors.New("something went wrong while running a command"),
		errors.New("something went wrong while loading"),
		errors.New("initialization of audio output failed"),
		errors.New("initialization of video output failed"),
		errors.New("there is no audio or video data to play"),
		errors.New("cannot identify file format"),
		errors.New("system requirements not fulfilled"),
		errors.New("function is not implemented"),
		errors.New("unknown error occurred"),
	}[-e]
}
