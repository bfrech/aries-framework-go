/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package rest

import (
	"github.com/hyperledger/aries-framework-go/cmd/aries-agent-mobile/pkg/wrappers/models"
	cmddidexch "github.com/hyperledger/aries-framework-go/pkg/controller/command/didexchange"
)

// DIDExchange contains necessary fields to support its operations.
type DIDExchange struct {
	httpClient httpClient
	endpoints  map[string]*endpoint

	URL   string
	Token string
}

// CreateInvitation creates a new connection invitation.
func (de *DIDExchange) CreateInvitation(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return de.createRespEnvelope(req, cmddidexch.CreateInvitationCommandMethod)
}

// ReceiveInvitation receives a new connection invitation.
func (de *DIDExchange) ReceiveInvitation(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return de.createRespEnvelope(req, cmddidexch.ReceiveInvitationCommandMethod)
}

// AcceptInvitation accepts a stored connection invitation.
func (de *DIDExchange) AcceptInvitation(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return de.createRespEnvelope(req, cmddidexch.AcceptInvitationCommandMethod)
}

// CreateImplicitInvitation creates implicit invitation using inviter DID.
func (de *DIDExchange) CreateImplicitInvitation(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return de.createRespEnvelope(req, cmddidexch.CreateImplicitInvitationCommandMethod)
}

// AcceptExchangeRequest accepts a stored connection request.
func (de *DIDExchange) AcceptExchangeRequest(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return de.createRespEnvelope(req, cmddidexch.AcceptExchangeRequestCommandMethod)
}

// QueryConnections queries agent to agent connections.
func (de *DIDExchange) QueryConnections(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return de.createRespEnvelope(req, cmddidexch.QueryConnectionsCommandMethod)
}

// QueryConnectionByID fetches a single connection record by connection ID.
func (de *DIDExchange) QueryConnectionByID(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return de.createRespEnvelope(req, cmddidexch.QueryConnectionByIDCommandMethod)
}

// CreateConnection creates a new connection record in completed state and returns the generated connectionID.
func (de *DIDExchange) CreateConnection(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return de.createRespEnvelope(req, cmddidexch.CreateConnectionCommandMethod)
}

// RemoveConnection removes given connection record.
func (de *DIDExchange) RemoveConnection(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return de.createRespEnvelope(req, cmddidexch.RemoveConnectionCommandMethod)
}

func (de *DIDExchange) createRespEnvelope(request *models.RequestEnvelope, endpoint string) *models.ResponseEnvelope {
	return exec(&restOperation{
		url:        de.URL,
		token:      de.Token,
		httpClient: de.httpClient,
		endpoint:   de.endpoints[endpoint],
		request:    request,
	})
}
