package command

import (
	"encoding/json"
	"github.com/hyperledger/aries-framework-go/cmd/aries-agent-mobile/pkg/wrappers/models"
	"github.com/hyperledger/aries-framework-go/pkg/controller/command"
	cmdconnection "github.com/hyperledger/aries-framework-go/pkg/controller/command/connection"
)

type Connection struct {
	handlers map[string]command.Exec
}

// RotateDID rotates the DID of the given connection to the given new DID, using the signing KID for the key in the old
// DID doc to sign the DID rotation.
func (i *Connection) RotateDID(request []byte) *models.ResponseEnvelope {
	args := cmdconnection.RotateDIDRequest{}

	if err := json.Unmarshal(request, &args); err != nil {
		return &models.ResponseEnvelope{Error: &models.CommandError{Message: err.Error()}}
	}

	response, cmdErr := exec(i.handlers[cmdconnection.RotateDIDCommandMethod], args)
	if cmdErr != nil {
		return &models.ResponseEnvelope{Error: cmdErr}
	}

	return &models.ResponseEnvelope{Payload: response}
}

// CreateConnectionV2 creates a DIDComm V2 connection with the given DID.
func (i *Connection) CreateConnectionV2(request []byte) *models.ResponseEnvelope {
	args := cmdconnection.CreateConnectionRequest{}

	if err := json.Unmarshal(request, &args); err != nil {
		return &models.ResponseEnvelope{Error: &models.CommandError{Message: err.Error()}}
	}

	response, cmdErr := exec(i.handlers[cmdconnection.CreateV2CommandMethod], args)
	if cmdErr != nil {
		return &models.ResponseEnvelope{Error: cmdErr}
	}

	return &models.ResponseEnvelope{Payload: response}
}

// UpdateTheirDIDForConnection updates TheirDID for the connection ID
func (i *Connection) UpdateTheirDIDForConnection(request []byte) *models.ResponseEnvelope {
	args := cmdconnection.UpateTheirDIDRequest{}

	if err := json.Unmarshal(request, &args); err != nil {
		return &models.ResponseEnvelope{Error: &models.CommandError{Message: err.Error()}}
	}

	response, cmdErr := exec(i.handlers[cmdconnection.UpdateTheirDIDForConnectionCommandMethod], args)
	if cmdErr != nil {
		return &models.ResponseEnvelope{Error: cmdErr}
	}

	return &models.ResponseEnvelope{Payload: response}
}
