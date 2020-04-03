package esl

// EventType holds type of information of what the type of information exists
type EventType int

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
