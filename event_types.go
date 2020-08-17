package esl

// EventType holds type of information of what the type of information exists
type EventType int

// EventContentType is the content type for events
type EventContentType string

// Types of events, known commands have their names
//
const (
	ETEvent            EventType = iota // General identified event
	ETSocketData                        // Unknown header that is not identified as event
	ETAuth                              // Authentication request
	ETAPI                               // API based event
	ETCommandReplay                     // Replay for command request
	ETAPIResponse                       // Response for mod_command API request
	ETLogData                           // Log information returned from a request
	ETDisconnectNotice                  // Event when a connection is disconnected
)

// The event types
const (
	ECTAuthRequest          EventContentType = "auth/request"
	ECTCommandReply         EventContentType = "command/reply"
	ECTAPIResponse          EventContentType = "api/response"
	ECTDisconnectNotice     EventContentType = "text/disconnect-notice"
	ECTRudeRejection        EventContentType = "text/rude-rejection"
	ECTEventPlain           EventContentType = "text/event-plain"
	ECTEventJSON            EventContentType = "text/event-json"
	ECTEventXML             EventContentType = "text/event-xml"
	ECTTextPlain            EventContentType = "text/plain"
	ECTSimpleMessageSummary EventContentType = "application/simple-message-summary"
	ECLogData               EventContentType = "log/data"
)
