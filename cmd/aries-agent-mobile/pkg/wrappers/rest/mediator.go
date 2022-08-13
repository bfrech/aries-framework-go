/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package rest

import (
	"github.com/hyperledger/aries-framework-go/cmd/aries-agent-mobile/pkg/wrappers/models"
	"github.com/hyperledger/aries-framework-go/pkg/controller/command/mediator"
)

// Mediator contains necessary fields to support its operations.
type Mediator struct {
	httpClient httpClient
	endpoints  map[string]*endpoint

	URL   string
	Token string
}

// Register registers the agent with the router.
func (m *Mediator) Register(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return m.createRespEnvelope(req, mediator.RegisterCommandMethod)
}

// Unregister unregisters the agent with the router.
func (m *Mediator) Unregister(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return m.createRespEnvelope(req, mediator.UnregisterCommandMethod)
}

// Connections returns router`s connections.
func (m *Mediator) Connections(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return m.createRespEnvelope(req, mediator.GetConnectionsCommandMethod)
}

// Reconnect sends noop message to given mediator connection to re-establish network connection.
func (m *Mediator) Reconnect(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return m.createRespEnvelope(req, mediator.ReconnectCommandMethod)
}

// ReconnectAll sends noop message to all mediator connections to re-establish a network connections.
func (m *Mediator) ReconnectAll(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return m.createRespEnvelope(req, mediator.ReconnectAllCommandMethod)
}

// Status returns details about pending messages for given connection.
func (m *Mediator) Status(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return m.createRespEnvelope(req, mediator.StatusCommandMethod)
}

// BatchPickup dispatches pending messages for given connection.
func (m *Mediator) BatchPickup(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return m.createRespEnvelope(req, mediator.BatchPickupCommandMethod)
}

// RegisterKey registers a new key with the router.
func (m *Mediator) RegisterKey(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return m.createRespEnvelope(req, mediator.RegisterKeyCommandMethod)
}

// UnregisterKey removes the key from the router.
func (m *Mediator) UnregisterKey(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return m.createRespEnvelope(req, mediator.UnregisterKeyCommandMethod)
}

func (m *Mediator) createRespEnvelope(request *models.RequestEnvelope, endpoint string) *models.ResponseEnvelope {
	return exec(&restOperation{
		url:        m.URL,
		token:      m.Token,
		httpClient: m.httpClient,
		endpoint:   m.endpoints[endpoint],
		request:    request,
	})
}
