/*
 *
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 * /
 *
 */

// Package basicv2 provide basic message protocol features
//
// Any incoming message of type "https://didcomm.org/basicmessage/2.0/message" can be handled
// by registering `basicv2.MessageService`.
package basicv2

import (
	"fmt"

	"github.com/hyperledger/aries-framework-go/pkg/common/log"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/internal/logutil"
)

const (
	// MessageRequestType is basic message DIDComm message type.
	MessageRequestType = "https://didcomm.org/basicmessage/2.0/message"

	// error messages.
	errNameAndHandleMandatory = "service name and basic message handle is mandatory"
	errFailedToDecodeMsg      = "unable to decode incoming DID comm message: %w"

	basicV2Message = "basicV2Message"
)

var logger = log.New("aries-framework/basicv2msg")

// MessageHandle is handle function for basic V2 message service which gets called by
// `basicv2.MessageService` to handle incoming messages.
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

// NewMessageService creates basic V2 message service which serves
// incoming basic V2  messages
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

// Accept is acceptance criteria for this basic message service.
func (m *MessageService) Accept(msgType string, purpose []string) bool {
	return msgType == MessageRequestType
}

// HandleInbound for basic message service.
func (m *MessageService) HandleInbound(msg service.DIDCommMsg, ctx service.DIDCommContext) (string, error) {
	basicMsg := Message{}

	err := msg.Decode(&basicMsg)
	if err != nil {
		return "", fmt.Errorf(errFailedToDecodeMsg, err)
	}

	logutil.LogDebug(logger, basicV2Message, "handleInbound", "received",
		logutil.CreateKeyValueString("msgType", msg.Type()),
		logutil.CreateKeyValueString("msgID", msg.ID()))

	return "", m.handle(basicMsg, ctx)
}
