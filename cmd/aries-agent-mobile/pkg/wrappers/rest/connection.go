package rest

import (
	"github.com/hyperledger/aries-framework-go/cmd/aries-agent-mobile/pkg/wrappers/models"
	cmdconnection "github.com/hyperledger/aries-framework-go/pkg/controller/command/connection"
)

type Connection struct {
	httpClient httpClient
	endpoints  map[string]*endpoint

	URL   string
	Token string
}

func (ir *Connection) RotateDID(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return ir.createRespEnvelope(req, cmdconnection.RotateDIDCommandMethod)
}

func (ir *Connection) CreateConnectionV2(request []byte) *models.ResponseEnvelope {
	req := &models.RequestEnvelope{Payload: request}
	return ir.createRespEnvelope(req, cmdconnection.CreateV2CommandMethod)
}

func (ir *Connection) createRespEnvelope(request *models.RequestEnvelope, endpoint string) *models.ResponseEnvelope {
	return exec(&restOperation{
		url:        ir.URL,
		token:      ir.Token,
		httpClient: ir.httpClient,
		endpoint:   ir.endpoints[endpoint],
		request:    request,
	})
}
