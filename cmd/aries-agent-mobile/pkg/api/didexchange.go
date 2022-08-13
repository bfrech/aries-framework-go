/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package api

import "github.com/hyperledger/aries-framework-go/cmd/aries-agent-mobile/pkg/wrappers/models"

// DIDExchangeController  defines methods for the DIDExchange protocol controller.
type DIDExchangeController interface {

	// CreateInvitation creates a new connection invitation.
	CreateInvitation(request []byte) *models.ResponseEnvelope

	// ReceiveInvitation receives a new connection invitation.
	ReceiveInvitation(request []byte) *models.ResponseEnvelope

	// AcceptInvitation accepts a stored connection invitation.
	AcceptInvitation(request []byte) *models.ResponseEnvelope

	// CreateImplicitInvitation creates implicit invitation using inviter DID.
	CreateImplicitInvitation(request []byte) *models.ResponseEnvelope

	// AcceptExchangeRequest accepts a stored connection request.
	AcceptExchangeRequest(request []byte) *models.ResponseEnvelope

	// QueryConnections queries agent to agent connections.
	QueryConnections(request []byte) *models.ResponseEnvelope

	// QueryConnectionByID fetches a single connection record by connection ID.
	QueryConnectionByID(request []byte) *models.ResponseEnvelope

	// CreateConnection creates a new connection record in completed state and returns the generated connectionID.
	CreateConnection(request []byte) *models.ResponseEnvelope

	// RemoveConnection removes given connection record.
	RemoveConnection(request []byte) *models.ResponseEnvelope
}
