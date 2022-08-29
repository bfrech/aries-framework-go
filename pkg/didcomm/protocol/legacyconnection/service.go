/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package legacyconnection

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/hyperledger/aries-framework-go/pkg/common/log"
	"github.com/hyperledger/aries-framework-go/pkg/common/model"
	"github.com/hyperledger/aries-framework-go/pkg/crypto"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/common/service"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/dispatcher"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/decorator"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/protocol/mediator"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/transport"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	vdrapi "github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/internal/logutil"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/hyperledger/aries-framework-go/pkg/store/connection"
	didstore "github.com/hyperledger/aries-framework-go/pkg/store/did"
	"github.com/hyperledger/aries-framework-go/pkg/vdr"
	"github.com/hyperledger/aries-framework-go/spi/storage"
)

var logger = log.New("aries-framework/legacyconnection/service")

const (
	// LegacyConnection connection protocol.
	LegacyConnection = "legacyconnection"
	// PIURI is the connection protocol identifier URI.
	PIURI = "https://didcomm.org/connections/1.0"
	// InvitationMsgType defines the legacy-connection invite message type.
	InvitationMsgType = PIURI + "/invitation"
	// RequestMsgType defines the legacy-connection request message type.
	RequestMsgType = PIURI + "/request"
	// ResponseMsgType defines the legacy-connection response message type.
	ResponseMsgType = PIURI + "/response"
	// AckMsgType defines the legacy-connection ack message type.
	AckMsgType             = "https://didcomm.org/notification/1.0/ack"
	routerConnsMetadataKey = "routerConnections"
)

const (
	myNSPrefix = "my"
	// TODO: https://github.com/hyperledger/aries-framework-go/issues/556 It will not be constant, this namespace
	//  will need to be figured with verification key
	theirNSPrefix = "their"
	// InvitationRecipientKey is map key constant.
	InvitationRecipientKey = "invRecipientKey"
)

// message type to store data for eventing. This is retrieved during callback.
type message struct {
	Msg           service.DIDCommMsgMap
	ThreadID      string
	Options       *options
	NextStateName string
	ConnRecord    *connection.Record
	// err is used to determine whether callback was stopped
	// e.g the user received an action event and executes Stop(err) function
	// in that case `err` is equal to `err` which was passing to Stop function
	err error
}

// provider contains dependencies for the Connection protocol and is typically created by using aries.Context().
type provider interface {
	OutboundDispatcher() dispatcher.Outbound
	StorageProvider() storage.Provider
	ProtocolStateStorageProvider() storage.Provider
	DIDConnectionStore() didstore.ConnectionStore
	Crypto() crypto.Crypto
	KMS() kms.KeyManager
	VDRegistry() vdrapi.Registry
	Service(id string) (interface{}, error)
	KeyType() kms.KeyType
	KeyAgreementType() kms.KeyType
	MediaTypeProfiles() []string
}

// stateMachineMsg is an internal struct used to pass data to state machine.
type stateMachineMsg struct {
	service.DIDCommMsg
	connRecord *connection.Record
	options    *options
}

type options struct {
	publicDID         string
	routerConnections []string
	label             string
}

// Service for Connection protocol.
type Service struct {
	service.Action
	service.Message
	ctx                *context
	callbackChannel    chan *message
	connectionRecorder *connection.Recorder
	connectionStore    didstore.ConnectionStore
	initialized        bool
}

type context struct {
	outboundDispatcher dispatcher.Outbound
	crypto             crypto.Crypto
	kms                kms.KeyManager
	connectionRecorder *connection.Recorder
	connectionStore    didstore.ConnectionStore
	vdRegistry         vdrapi.Registry
	routeSvc           mediator.ProtocolService
	doACAPyInterop     bool
	keyType            kms.KeyType
	keyAgreementType   kms.KeyType
	mediaTypeProfiles  []string
}

// opts are used to provide client properties to Connection service.
type opts interface {
	// PublicDID allows for setting public DID
	PublicDID() string

	// Label allows for setting label
	Label() string

	// RouterConnections allows for setting router connections
	RouterConnections() []string
}

// New return connection service.
func New(prov provider) (*Service, error) {
	svc := Service{}

	err := svc.Initialize(prov)
	if err != nil {
		return nil, err
	}

	return &svc, nil
}

// Initialize initializes the Service. If Initialize succeeds, any further call is a no-op.
func (s *Service) Initialize(p interface{}) error {
	if s.initialized {
		return nil
	}

	prov, ok := p.(provider)
	if !ok {
		return fmt.Errorf("expected provider of type `%T`, got type `%T`", provider(nil), p)
	}

	connRecorder, err := connection.NewRecorder(prov)
	if err != nil {
		return fmt.Errorf("failed to initialize connection recorder: %w", err)
	}

	routeSvcBase, err := prov.Service(mediator.Coordination)
	if err != nil {
		return err
	}

	routeSvc, ok := routeSvcBase.(mediator.ProtocolService)
	if !ok {
		return errors.New("cast service to Route Service failed")
	}

	const callbackChannelSize = 10

	keyType := kms.ED25519Type

	keyAgreementType := kms.X25519ECDHKWType

	mediaTypeProfiles := []string{transport.LegacyDIDCommV1Profile}

	s.ctx = &context{
		outboundDispatcher: prov.OutboundDispatcher(),
		crypto:             prov.Crypto(),
		kms:                prov.KMS(),
		vdRegistry:         prov.VDRegistry(),
		connectionRecorder: connRecorder,
		connectionStore:    prov.DIDConnectionStore(),
		routeSvc:           routeSvc,
		keyType:            keyType,
		keyAgreementType:   keyAgreementType,
		mediaTypeProfiles:  mediaTypeProfiles,
	}

	// TODO channel size - https://github.com/hyperledger/aries-framework-go/issues/246
	s.callbackChannel = make(chan *message, callbackChannelSize)
	s.connectionRecorder = connRecorder
	s.connectionStore = prov.DIDConnectionStore()

	// start the listener
	go s.startInternalListener()

	s.initialized = true

	return nil
}

func retrievingRouterConnections(msg service.DIDCommMsg) []string {
	raw, found := msg.Metadata()[routerConnsMetadataKey]
	if !found {
		return nil
	}

	connections, ok := raw.([]string)
	if !ok {
		return nil
	}

	return connections
}

// HandleInbound handles inbound connection messages.
func (s *Service) HandleInbound(msg service.DIDCommMsg, ctx service.DIDCommContext) (string, error) {
	logger.Debugf("receive inbound message : %s", msg)

	// fetch the thread id
	thID, err := msg.ThreadID()
	if err != nil {
		return "", err
	}

	// valid state transition and get the next state
	next, err := s.nextState(msg.Type(), thID)
	if err != nil {
		return "", fmt.Errorf("handle inbound - next state : %w", err)
	}

	// connection record
	connRecord, err := s.connectionRecord(msg, ctx)
	if err != nil {
		return "", fmt.Errorf("failed to fetch connection record : %w", err)
	}

	logger.Debugf("connection record: %+v", connRecord)

	internalMsg := &message{
		Options:       &options{routerConnections: retrievingRouterConnections(msg)},
		Msg:           msg.Clone(),
		ThreadID:      thID,
		NextStateName: next.Name(),
		ConnRecord:    connRecord,
	}

	go func(msg *message, aEvent chan<- service.DIDCommAction) {
		if err = s.handle(msg, aEvent); err != nil {
			logutil.LogError(logger, LegacyConnection, "processMessage", err.Error(),
				logutil.CreateKeyValueString("msgType", msg.Msg.Type()),
				logutil.CreateKeyValueString("msgID", msg.Msg.ID()),
				logutil.CreateKeyValueString("connectionID", msg.ConnRecord.ConnectionID))
		}

		logutil.LogDebug(logger, LegacyConnection, "processMessage", "success",
			logutil.CreateKeyValueString("msgType", msg.Msg.Type()),
			logutil.CreateKeyValueString("msgID", msg.Msg.ID()),
			logutil.CreateKeyValueString("connectionID", msg.ConnRecord.ConnectionID))
	}(internalMsg, s.ActionEvent())

	logutil.LogDebug(logger, LegacyConnection, "handleInbound", "success",
		logutil.CreateKeyValueString("msgType", msg.Type()),
		logutil.CreateKeyValueString("msgID", msg.ID()),
		logutil.CreateKeyValueString("connectionID", internalMsg.ConnRecord.ConnectionID))

	return connRecord.ConnectionID, nil
}

// Name return service name.
func (s *Service) Name() string {
	return LegacyConnection
}

func findNamespace(msgType string) string {
	namespace := theirNSPrefix
	if msgType == InvitationMsgType || msgType == ResponseMsgType {
		namespace = myNSPrefix
	}

	return namespace
}

// Accept msg checks the msg type.
func (s *Service) Accept(msgType string) bool {
	return msgType == InvitationMsgType ||
		msgType == RequestMsgType ||
		msgType == ResponseMsgType ||
		msgType == AckMsgType
}

// HandleOutbound handles outbound connection messages.
func (s *Service) HandleOutbound(_ service.DIDCommMsg, _, _ string) (string, error) {
	return "", errors.New("not implemented")
}

func (s *Service) nextState(msgType, thID string) (state, error) {
	logger.Debugf("msgType=%s thID=%s", msgType, thID)

	nsThID, err := connection.CreateNamespaceKey(findNamespace(msgType), thID)
	if err != nil {
		return nil, err
	}

	current, err := s.currentState(nsThID)
	if err != nil {
		return nil, err
	}

	logger.Debugf("retrieved current state [%s] using nsThID [%s]", current.Name(), nsThID)

	next, err := stateFromMsgType(msgType)
	if err != nil {
		return nil, err
	}

	logger.Debugf("check if current state [%s] can transition to [%s]", current.Name(), next.Name())

	if !current.CanTransitionTo(next) {
		return nil, fmt.Errorf("invalid state transition: %s -> %s", current.Name(), next.Name())
	}

	return next, nil
}

func (s *Service) handle(msg *message, aEvent chan<- service.DIDCommAction) error { //nolint:funlen,gocyclo
	logger.Debugf("handling msg: %+v", msg)

	next, err := stateFromName(msg.NextStateName)
	if err != nil {
		return fmt.Errorf("invalid state name: %w", err)
	}

	for !isNoOp(next) {
		s.sendMsgEvents(&service.StateMsg{
			ProtocolName: LegacyConnection,
			Type:         service.PreState,
			Msg:          msg.Msg.Clone(),
			StateID:      next.Name(),
			Properties:   createEventProperties(msg.ConnRecord.ConnectionID, msg.ConnRecord.InvitationID),
		})
		logger.Debugf("sent pre event for state %s", next.Name())

		var (
			action           stateAction
			followup         state
			connectionRecord *connection.Record
		)

		connectionRecord, followup, action, err = next.ExecuteInbound(
			&stateMachineMsg{
				DIDCommMsg: msg.Msg,
				connRecord: msg.ConnRecord,
				options:    msg.Options,
			},
			msg.ThreadID,
			s.ctx)

		if err != nil {
			return fmt.Errorf("failed to execute state '%s': %w", next.Name(), err)
		}

		connectionRecord.State = next.Name()
		logger.Debugf("finished execute state: %s", next.Name())

		if err = s.update(msg.Msg.Type(), connectionRecord); err != nil {
			return fmt.Errorf("failed to persist state '%s': %w", next.Name(), err)
		}

		if connectionRecord.State == StateIDCompleted {
			err = s.connectionStore.SaveDID(connectionRecord.TheirDID, connectionRecord.RecipientKeys...)
			if err != nil {
				return fmt.Errorf("save theirDID: %w", err)
			}
		}

		if err = action(); err != nil {
			return fmt.Errorf("failed to execute state action '%s': %w", next.Name(), err)
		}

		logger.Debugf("finish execute state action: '%s'", next.Name())

		prev := next
		next = followup
		haltExecution := false

		// trigger action event based on message type for inbound messages
		if canTriggerActionEvents(connectionRecord.State, connectionRecord.Namespace) {
			logger.Debugf("action event triggered for msg type: %s", msg.Msg.Type())

			msg.NextStateName = next.Name()
			if err = s.sendActionEvent(msg, aEvent); err != nil {
				return fmt.Errorf("handle inbound: %w", err)
			}

			haltExecution = true
		}

		s.sendMsgEvents(&service.StateMsg{
			ProtocolName: LegacyConnection,
			Type:         service.PostState,
			Msg:          msg.Msg.Clone(),
			StateID:      prev.Name(),
			Properties:   createEventProperties(connectionRecord.ConnectionID, connectionRecord.InvitationID),
		})
		logger.Debugf("sent post event for state %s", prev.Name())

		if haltExecution {
			logger.Debugf("halted execution before state=%s", msg.NextStateName)

			break
		}
	}

	return nil
}

func (s *Service) handleWithoutAction(msg *message) error {
	return s.handle(msg, nil)
}

func createEventProperties(connectionID, invitationID string) *connectionEvent {
	return &connectionEvent{
		connectionID: connectionID,
		invitationID: invitationID,
	}
}

// sendActionEvent triggers the action event. This function stores the state of current processing and passes a callback
// function in the event message.
func (s *Service) sendActionEvent(internalMsg *message, aEvent chan<- service.DIDCommAction) error {
	// save data to support AcceptConnectionRequest APIs (when client will not be able to invoke the callback function)
	err := s.storeEventProtocolStateData(internalMsg)
	if err != nil {
		return fmt.Errorf("send action event : %w", err)
	}

	if aEvent != nil {
		// trigger action event
		aEvent <- service.DIDCommAction{
			ProtocolName: LegacyConnection,
			Message:      internalMsg.Msg.Clone(),
			Continue: func(args interface{}) {
				switch v := args.(type) {
				case opts:
					internalMsg.Options = &options{
						publicDID:         v.PublicDID(),
						label:             v.Label(),
						routerConnections: v.RouterConnections(),
					}
				default:
					// nothing to do
				}

				s.processCallback(internalMsg)
			},
			Stop: func(err error) {
				// sets an error to the message
				internalMsg.err = err
				s.processCallback(internalMsg)
			},
			Properties: createEventProperties(internalMsg.ConnRecord.ConnectionID, internalMsg.ConnRecord.InvitationID),
		}

		logger.Debugf("dispatched action for msg: %+v", internalMsg.Msg)
	}

	return nil
}

// sendEvent triggers the message events.
func (s *Service) sendMsgEvents(msg *service.StateMsg) {
	// trigger the message events
	for _, handler := range s.MsgEvents() {
		handler <- *msg

		logger.Debugf("sent msg event to handler: %+v", msg)
	}
}

// startInternalListener listens to messages in gochannel for callback messages from clients.
func (s *Service) startInternalListener() {
	for msg := range s.callbackChannel {
		// TODO https://github.com/hyperledger/aries-framework-go/issues/242 - retry logic
		// if no error - do handle
		if msg.err == nil {
			msg.err = s.handleWithoutAction(msg)
		}

		// no error - continue
		if msg.err == nil {
			continue
		}
	}
}

// AcceptInvitation accepts/approves connection invitation.
func (s *Service) AcceptInvitation(connectionID, publicDID, label string, routerConnections []string) error {
	return s.accept(connectionID, publicDID, label, StateIDInvited,
		"accept connection invitation", routerConnections)
}

// AcceptConnectionRequest accepts/approves connection request.
func (s *Service) AcceptConnectionRequest(connectionID, publicDID, label string, routerConnections []string) error {
	return s.accept(connectionID, publicDID, label, StateIDRequested,
		"accept exchange request", routerConnections)
}

func (s *Service) accept(connectionID, publicDID, label, stateID, errMsg string, routerConnections []string) error {
	msg, err := s.getEventProtocolStateData(connectionID)
	if err != nil {
		return fmt.Errorf("failed to accept invitation for connectionID=%s : %s : %w", connectionID, errMsg, err)
	}

	connRecord, err := s.connectionRecorder.GetConnectionRecord(connectionID)
	if err != nil {
		return fmt.Errorf("%s : %w", errMsg, err)
	}

	if connRecord.State != stateID {
		return fmt.Errorf("current state (%s) is different from "+
			"expected state (%s)", connRecord.State, stateID)
	}

	msg.Options = &options{publicDID: publicDID, label: label, routerConnections: routerConnections}

	return s.handleWithoutAction(msg)
}

func (s *Service) storeEventProtocolStateData(msg *message) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("store protocol state data : %w", err)
	}

	return s.connectionRecorder.SaveEvent(msg.ConnRecord.ConnectionID, bytes)
}

func (s *Service) getEventProtocolStateData(connectionID string) (*message, error) {
	val, err := s.connectionRecorder.GetEvent(connectionID)
	if err != nil {
		return nil, fmt.Errorf("get protocol state data : %w", err)
	}

	msg := &message{}

	err = json.Unmarshal(val, msg)
	if err != nil {
		return nil, fmt.Errorf("get protocol state data : %w", err)
	}

	return msg, nil
}

func (s *Service) processCallback(msg *message) {
	// pass the callback data to internal channel. This is created to unblock consumer go routine and wrap the callback
	// channel internally.
	s.callbackChannel <- msg
}

func isNoOp(s state) bool {
	_, ok := s.(*noOp)
	return ok
}

func (s *Service) currentState(nsThID string) (state, error) {
	connRec, err := s.connectionRecorder.GetConnectionRecordByNSThreadID(nsThID)
	if err != nil {
		if errors.Is(err, storage.ErrDataNotFound) {
			return &null{}, nil
		}

		return nil, fmt.Errorf("cannot fetch state from store: thID=%s err=%w", nsThID, err)
	}

	return stateFromName(connRec.State)
}

func (s *Service) update(msgType string, record *connection.Record) error {
	if (msgType == RequestMsgType && record.State == StateIDRequested) ||
		(msgType == InvitationMsgType && record.State == StateIDInvited) {
		return s.connectionRecorder.SaveConnectionRecordWithMappings(record)
	}

	return s.connectionRecorder.SaveConnectionRecord(record)
}

// CreateConnection saves the record to the connection store and maps TheirDID to their recipient keys in
// the did connection store.
func (s *Service) CreateConnection(record *connection.Record, theirDID *did.Doc) error {
	logger.Debugf("creating connection using record [%+v] and theirDID [%+v]", record, theirDID)

	didMethod, err := vdr.GetDidMethod(theirDID.ID)
	if err != nil {
		return err
	}

	_, err = s.ctx.vdRegistry.Create(didMethod, theirDID, vdrapi.WithOption("store", true))
	if err != nil {
		return fmt.Errorf("vdr failed to store theirDID : %w", err)
	}

	err = s.connectionStore.SaveDIDFromDoc(theirDID)
	if err != nil {
		return fmt.Errorf("failed to save theirDID to the did.ConnectionStore: %w", err)
	}

	err = s.connectionStore.SaveDIDByResolving(record.MyDID)
	if err != nil {
		return fmt.Errorf("failed to save myDID to the did.ConnectionStore: %w", err)
	}

	if isDIDCommV2(record.MediaTypeProfiles) {
		record.DIDCommVersion = service.V2
	} else {
		record.DIDCommVersion = service.V1
	}

	return s.connectionRecorder.SaveConnectionRecord(record)
}

func (s *Service) connectionRecord(msg service.DIDCommMsg, ctx service.DIDCommContext) (*connection.Record, error) {
	switch msg.Type() {
	case InvitationMsgType:
		return s.invitationMsgRecord(msg)
	case RequestMsgType:
		return s.requestMsgRecord(msg, ctx)
	case ResponseMsgType:
		return s.responseMsgRecord(msg)
	case AckMsgType:
		return s.fetchConnectionRecord(theirNSPrefix, msg)
	}

	return nil, errors.New("invalid message type")
}

func (s *Service) invitationMsgRecord(msg service.DIDCommMsg) (*connection.Record, error) {
	thID, msgErr := msg.ThreadID()
	if msgErr != nil {
		return nil, msgErr
	}

	invitation := &Invitation{}

	err := msg.Decode(invitation)
	if err != nil {
		return nil, err
	}

	recKey, err := s.ctx.getInvitationRecipientKey(invitation)
	if err != nil {
		return nil, err
	}

	connRecord := &connection.Record{
		ConnectionID:    generateRandomID(),
		ThreadID:        thID,
		State:           stateNameNull,
		InvitationID:    invitation.ID,
		InvitationDID:   invitation.DID,
		ServiceEndPoint: model.NewDIDCommV1Endpoint(invitation.ServiceEndpoint),
		RecipientKeys:   []string{recKey},
		TheirLabel:      invitation.Label,
		Namespace:       findNamespace(msg.Type()),
		DIDCommVersion:  service.V1,
	}

	if err := s.connectionRecorder.SaveConnectionRecord(connRecord); err != nil {
		return nil, err
	}

	return connRecord, nil
}

func (s *Service) requestMsgRecord(msg service.DIDCommMsg, ctx service.DIDCommContext) (*connection.Record, error) {
	var (
		request          Request
		invitationRecKey []string
	)

	err := msg.Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling failed: %w", err)
	}

	invitationID := msg.ParentThreadID()
	if invitationID == "" {
		// try to retrieve key which was assigned at HandleInboundEnvelope method of inbound handler
		if verKey, ok := ctx.All()[InvitationRecipientKey].(string); ok && verKey != "" {
			invitationRecKey = append(invitationRecKey, verKey)
		} else {
			return nil, fmt.Errorf("missing parent thread ID and invitation recipient key"+
				" on connection request with @id=%s", request.ID)
		}
	}

	if request.Connection == nil {
		return nil, fmt.Errorf("missing connection field on connection request with @id=%s", request.ID)
	}

	connRecord := &connection.Record{
		TheirLabel:              request.Label,
		ConnectionID:            generateRandomID(),
		ThreadID:                request.ID,
		State:                   stateNameNull,
		TheirDID:                request.Connection.DID,
		InvitationID:            invitationID,
		InvitationRecipientKeys: invitationRecKey,
		Namespace:               theirNSPrefix,
		DIDCommVersion:          service.V1,
	}

	if err := s.connectionRecorder.SaveConnectionRecord(connRecord); err != nil {
		return nil, err
	}

	return connRecord, nil
}

func (s *Service) responseMsgRecord(payload service.DIDCommMsg) (*connection.Record, error) {
	return s.fetchConnectionRecord(myNSPrefix, payload)
}

func (s *Service) fetchConnectionRecord(nsPrefix string, payload service.DIDCommMsg) (*connection.Record, error) {
	msg := &struct {
		Thread decorator.Thread `json:"~thread,omitempty"`
	}{}

	err := payload.Decode(msg)
	if err != nil {
		return nil, err
	}

	key, err := connection.CreateNamespaceKey(nsPrefix, msg.Thread.ID)
	if err != nil {
		return nil, err
	}

	return s.connectionRecorder.GetConnectionRecordByNSThreadID(key)
}

func generateRandomID() string {
	return uuid.New().String()
}

// canTriggerActionEvents true based on role and state.
// 1. Role is invitee and state is invited.
// 2. Role is inviter and state is requested.
func canTriggerActionEvents(stateID, ns string) bool {
	return (stateID == StateIDInvited && ns == myNSPrefix) || (stateID == StateIDRequested && ns == theirNSPrefix)
}

// CreateImplicitInvitation creates implicit invitation. Inviter DID is required, invitee DID is optional.
// If invitee DID is not provided new peer DID will be created for implicit invitation connection request.
//nolint:funlen
func (s *Service) CreateImplicitInvitation(inviterLabel, inviterDID,
	inviteeLabel, inviteeDID string, routerConnections []string) (string, error) {
	logger.Debugf("implicit invitation requested inviterDID[%s] inviteeDID[%s]", inviterDID, inviteeDID)

	docResolution, err := s.ctx.vdRegistry.Resolve(inviterDID)
	if err != nil {
		return "", fmt.Errorf("resolve public did[%s]: %w", inviterDID, err)
	}

	dest, err := service.CreateDestination(docResolution.DIDDocument)
	if err != nil {
		return "", err
	}

	thID := generateRandomID()

	var connRecord *connection.Record

	if accept, e := dest.ServiceEndpoint.Accept(); e == nil && isDIDCommV2(accept) {
		connRecord = &connection.Record{
			ConnectionID:    generateRandomID(),
			ThreadID:        thID,
			State:           stateNameNull,
			InvitationDID:   inviterDID,
			Implicit:        true,
			ServiceEndPoint: dest.ServiceEndpoint,
			RecipientKeys:   dest.RecipientKeys,
			TheirLabel:      inviterLabel,
			Namespace:       findNamespace(InvitationMsgType),
		}
	} else {
		connRecord = &connection.Record{
			ConnectionID:      generateRandomID(),
			ThreadID:          thID,
			State:             stateNameNull,
			InvitationDID:     inviterDID,
			Implicit:          true,
			ServiceEndPoint:   dest.ServiceEndpoint,
			RecipientKeys:     dest.RecipientKeys,
			RoutingKeys:       dest.RoutingKeys,
			MediaTypeProfiles: dest.MediaTypeProfiles,
			TheirLabel:        inviterLabel,
			Namespace:         findNamespace(InvitationMsgType),
		}
	}

	if e := s.connectionRecorder.SaveConnectionRecordWithMappings(connRecord); e != nil {
		return "", fmt.Errorf("failed to save new connection record for implicit invitation: %w", e)
	}

	invitation := &Invitation{
		ID:    uuid.New().String(),
		Label: inviterLabel,
		DID:   inviterDID,
		Type:  InvitationMsgType,
	}

	msg, err := createDIDCommMsg(invitation)
	if err != nil {
		return "", fmt.Errorf("failed to create DIDCommMsg for implicit invitation: %w", err)
	}

	next := &requested{}
	internalMsg := &message{
		Msg:           msg.Clone(),
		ThreadID:      thID,
		NextStateName: next.Name(),
		ConnRecord:    connRecord,
	}
	internalMsg.Options = &options{publicDID: inviteeDID, label: inviteeLabel, routerConnections: routerConnections}

	go func(msg *message, aEvent chan<- service.DIDCommAction) {
		if err = s.handle(msg, aEvent); err != nil {
			logger.Errorf("error from handle for implicit invitation: %s", err)
		}
	}(internalMsg, s.ActionEvent())

	return connRecord.ConnectionID, nil
}

func createDIDCommMsg(invitation *Invitation) (service.DIDCommMsg, error) {
	payload, err := json.Marshal(invitation)
	if err != nil {
		return nil, fmt.Errorf("marshal invitation: %w", err)
	}

	return service.ParseDIDCommMsgMap(payload)
}
