/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package api

import "github.com/hyperledger/aries-framework-go/cmd/aries-agent-mobile/pkg/wrappers/models"

// MessagingController defines methods for the Messaging controller.
type MessagingController interface {

	// RegisterService registers new message service to message handler registrar.
	RegisterService(request []byte) *models.ResponseEnvelope

	// UnregisterService unregisters given message service handler registrar.
	UnregisterService(request []byte) *models.ResponseEnvelope

	// Services returns list of registered service names.
	Services(request []byte) *models.ResponseEnvelope

	// Send sends new message to destination provided.
	Send(request []byte) *models.ResponseEnvelope

	// Reply sends reply to existing message.
	Reply(request []byte) *models.ResponseEnvelope

	// RegisterHTTPService registers new http over didcomm service to message handler registrar.
	RegisterHTTPService(request []byte) *models.ResponseEnvelope
}
