package mobilemsg

import (
	"fmt"
	"github.com/hyperledger/aries-framework-go/pkg/common/log"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/internal/logutil"
)

const (
	// MobileMessageRequestType is mobile message type.
	MobileMessageRequestType = "https://didcomm.org/mobilemessage/1.0/message"

	// error messages.
	errNameAndHandleMandatory = "service name and mobile message handle is mandatory"
	errFailedToDecodeMsg      = "unable to decode incoming DID comm message: %w"

	mobileMessage = "mobileMessage"
)

var logger = log.New("aries-framework/mobilemsg")

// MessageHandle is handle function for mobile message service which gets called by
// `mobilemsg.MessageService` to handle incoming messages.
//
// Args
//
// message : incoming basic message.
// myDID: receiving agent's DID.
// theirDID: sender agent's DID.
//
// Returns
//
// error : handle can return error back to service to notify message dispatcher about failures.
type MessageHandle func(message Message, ctx service.DIDCommContext) error

// NewMessageService creates mobilemessage service which serves
// incoming mobile  messages
//
//
// Args:
//
// name - is name of this message service (this is mandatory argument).
//
// handle - is handle function to which incoming basic message will be sent(this is mandatory argument).
//
// Returns:
//
// MessageService: basic message service,
//
// error: arg validation errors.
func NewMessageService(name string, handle MessageHandle) (*MessageService, error) {
	if name == "" || handle == nil {
		return nil, fmt.Errorf(errNameAndHandleMandatory)
	}

	return &MessageService{
		name:   name,
		handle: handle,
	}, nil
}

// MessageService is message service which transports incoming basic messages to handlers provided.
type MessageService struct {
	name   string
	handle MessageHandle
}

// Name of basic message service.
func (m *MessageService) Name() string {
	return m.name
}

// Accept is acceptance criteria for this mobile message service.
func (m *MessageService) Accept(msgType string, purpose []string) bool {
	return msgType == MobileMessageRequestType
}

// HandleInbound for mobile message service.
func (m *MessageService) HandleInbound(msg service.DIDCommMsg, ctx service.DIDCommContext) (string, error) {
	mobileMsg := Message{}

	err := msg.Decode(&mobileMsg)
	if err != nil {
		return "", fmt.Errorf(errFailedToDecodeMsg, err)
	}

	logutil.LogDebug(logger, mobileMessage, "handleInbound", "received",
		logutil.CreateKeyValueString("msgType", msg.Type()),
		logutil.CreateKeyValueString("msgID", msg.ID()))

	return "", m.handle(mobileMsg, ctx)
}
