/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package command

import (
	"encoding/json"

	"github.com/hyperledger/aries-framework-go/cmd/aries-agent-mobile/pkg/wrappers/models"
	"github.com/hyperledger/aries-framework-go/pkg/controller/command"
	"github.com/hyperledger/aries-framework-go/pkg/controller/command/mediator"
)

// Mediator contains necessary fields to support its operations.
type Mediator struct {
	handlers map[string]command.Exec
}

// Register registers the agent with the router.
func (m *Mediator) Register(request []byte) *models.ResponseEnvelope {
	args := mediator.RegisterRoute{}

	if err := json.Unmarshal(request, &args); err != nil {
		return &models.ResponseEnvelope{Error: &models.CommandError{Message: err.Error()}}
	}

	response, cmdErr := exec(m.handlers[mediator.RegisterCommandMethod], args)
	if cmdErr != nil {
		return &models.ResponseEnvelope{Error: cmdErr}
	}

	return &models.ResponseEnvelope{Payload: response}
}

// Unregister unregisters the agent with the router.
func (m *Mediator) Unregister(request []byte) *models.ResponseEnvelope {
	response, cmdErr := exec(m.handlers[mediator.UnregisterCommandMethod], request)
	if cmdErr != nil {
		return &models.ResponseEnvelope{Error: cmdErr}
	}

	return &models.ResponseEnvelope{Payload: response}
}

// Connections returns router`s connections.
func (m *Mediator) Connections(request []byte) *models.ResponseEnvelope {
	response, cmdErr := exec(m.handlers[mediator.GetConnectionsCommandMethod], request)
	if cmdErr != nil {
		return &models.ResponseEnvelope{Error: cmdErr}
	}

	return &models.ResponseEnvelope{Payload: response}
}

// Reconnect ends noop message to given mediator connection to re-establish network connection.
func (m *Mediator) Reconnect(request []byte) *models.ResponseEnvelope {
	args := mediator.RegisterRoute{}

	if err := json.Unmarshal(request, &args); err != nil {
		return &models.ResponseEnvelope{Error: &models.CommandError{Message: err.Error()}}
	}

	response, cmdErr := exec(m.handlers[mediator.ReconnectCommandMethod], args)
	if cmdErr != nil {
		return &models.ResponseEnvelope{Error: cmdErr}
	}

	return &models.ResponseEnvelope{Payload: response}
}

// ReconnectAll sends noop message to all mediator connections to re-establish a network connections.
func (m *Mediator) ReconnectAll(request []byte) *models.ResponseEnvelope {
	response, cmdErr := exec(m.handlers[mediator.ReconnectAllCommandMethod], request)
	if cmdErr != nil {
		return &models.ResponseEnvelope{Error: cmdErr}
	}

	return &models.ResponseEnvelope{Payload: response}
}

// Status returns details about pending messages for given connection.
func (m *Mediator) Status(request []byte) *models.ResponseEnvelope {
	args := mediator.StatusRequest{}

	if err := json.Unmarshal(request, &args); err != nil {
		return &models.ResponseEnvelope{Error: &models.CommandError{Message: err.Error()}}
	}

	response, cmdErr := exec(m.handlers[mediator.StatusCommandMethod], args)
	if cmdErr != nil {
		return &models.ResponseEnvelope{Error: cmdErr}
	}

	return &models.ResponseEnvelope{Payload: response}
}

// BatchPickup dispatches pending messages for given connection.
func (m *Mediator) BatchPickup(request []byte) *models.ResponseEnvelope {
	args := mediator.BatchPickupRequest{}

	if err := json.Unmarshal(request, &args); err != nil {
		return &models.ResponseEnvelope{Error: &models.CommandError{Message: err.Error()}}
	}

	response, cmdErr := exec(m.handlers[mediator.BatchPickupCommandMethod], args)
	if cmdErr != nil {
		return &models.ResponseEnvelope{Error: cmdErr}
	}

	return &models.ResponseEnvelope{Payload: response}
}

// RegisterKey registers a new key with the router.
func (m *Mediator) RegisterKey(request []byte) *models.ResponseEnvelope {
	args := mediator.RegisterKey{}

	if err := json.Unmarshal(request, &args); err != nil {
		return &models.ResponseEnvelope{Error: &models.CommandError{Message: err.Error()}}
	}

	response, cmdErr := exec(m.handlers[mediator.RegisterKeyCommandMethod], args)
	if cmdErr != nil {
		return &models.ResponseEnvelope{Error: cmdErr}
	}

	return &models.ResponseEnvelope{Payload: response}
}

// UnregisterKey removes a key from the router.
func (m *Mediator) UnregisterKey(request []byte) *models.ResponseEnvelope {
	args := mediator.RegisterKey{}

	if err := json.Unmarshal(request, &args); err != nil {
		return &models.ResponseEnvelope{Error: &models.CommandError{Message: err.Error()}}
	}

	response, cmdErr := exec(m.handlers[mediator.UnregisterKeyCommandMethod], args)
	if cmdErr != nil {
		return &models.ResponseEnvelope{Error: cmdErr}
	}

	return &models.ResponseEnvelope{Payload: response}
}
