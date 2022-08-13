package api

import "github.com/hyperledger/aries-framework-go/cmd/aries-agent-mobile/pkg/wrappers/models"

type ConnectionController interface {

	// RotateDID rotates the DID of the given connection to the given new DID, using the signing KID for the key in the old
	// DID doc to sign the DID rotation.
	RotateDID(request []byte) *models.ResponseEnvelope

	// CreateConnectionV2 creates a DIDComm V2 connection with the given DID.
	CreateConnectionV2(request []byte) *models.ResponseEnvelope
}
