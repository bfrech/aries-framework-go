/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package rest

import (
	"github.com/hyperledger/aries-framework-go/cmd/aries-agent-mobile/pkg/wrappers/models"
	"github.com/hyperledger/aries-framework-go/pkg/controller/command/outofband"
)

// OutOfBand contains necessary fields to support its operations.
type OutOfBand struct {
	httpClient httpClient
	endpoints  map[string]*endpoint

	URL   string
	Token string
}

// CreateInvitation creates and saves an out-of-band invitation.
func (oob *OutOfBand) CreateInvitation(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return oob.createRespEnvelope(req, outofband.CreateInvitation)
}

// AcceptInvitation from another agent and return the ID of the new connection records.
func (oob *OutOfBand) AcceptInvitation(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return oob.createRespEnvelope(req, outofband.AcceptInvitation)
}

// Actions returns pending actions that have not yet to be executed or canceled.
func (oob *OutOfBand) Actions(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return oob.createRespEnvelope(req, outofband.Actions)
}

// ActionContinue allows continuing with the protocol after an action event was triggered.
func (oob *OutOfBand) ActionContinue(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return oob.createRespEnvelope(req, outofband.ActionContinue)
}

// ActionStop stops the protocol after an action event was triggered.
func (oob *OutOfBand) ActionStop(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return oob.createRespEnvelope(req, outofband.ActionStop)
}

func (oob *OutOfBand) createRespEnvelope(request *models.RequestEnvelope, endpoint string) *models.ResponseEnvelope {
	return exec(&restOperation{
		url:        oob.URL,
		token:      oob.Token,
		httpClient: oob.httpClient,
		endpoint:   oob.endpoints[endpoint],
		request:    request,
	})
}
