package events

type EventType string

const (
	tencentCloudIMEventPrefix = "im."

	BeforeRequest   EventType = tencentCloudIMEventPrefix + "before_request"
	AfterRequest              = tencentCloudIMEventPrefix + "after_request"
	RequestDuration           = tencentCloudIMEventPrefix + "request_duration"
)

const (
	tencentCloudIMKeyPrefix = "im."
	ServiceName             = tencentCloudIMKeyPrefix + "service"
	CommandName             = tencentCloudIMKeyPrefix + "command"
)

//
type KVPair interface {
	Key() string
	Value() interface{}
}

type kvPair struct {
	key   string
	value interface{}
}

func (p kvPair) Key() string {
	return p.key
}

func (p kvPair) Value() interface{} {
	return p.value
}

//
// NOTE(QuantumGhost): implementor of EventListener should not modify kvs slice.
type EventListener interface {
	Event(event EventType)
	EventKV(event EventType, kvs []KVPair)
	Timing(event EventType, nanoseconds int64)
	TimingKv(event EventType, nanoseconds int64, kvs []KVPair)
}

type nullEventListener struct{}

func (nullEventListener) Event(event EventType) {}

func (nullEventListener) EventKV(event EventType, kvs []KVPair) {}

func (nullEventListener) Timing(event EventType, nanoseconds int64) {}

func (nullEventListener) TimingKv(event EventType, nanoseconds int64, kvs []KVPair) {}

func KV(key string, value interface{}) KVPair {
	return kvPair{key: key, value: value}
}

func NewPairs(pairs ...KVPair) []KVPair {
	return pairs
}
