package app

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	chatv1 "github.com/zkxjzmswkwl/go-grpc-testing/gen/go/proto/chat/v1"
	chatv1connect "github.com/zkxjzmswkwl/go-grpc-testing/gen/go/proto/chat/v1/chatv1connect"
	"github.com/zkxjzmswkwl/go-grpc-testing/internal/store/memory"
)

// (csmith): This is pretty silly and not meant to be used in prod.

type ChatService struct {
	chatv1connect.UnimplementedChatServiceHandler
	store *memory.Store
}

func NewChatService(store *memory.Store) *ChatService {
	return &ChatService{store: store}
}

func (s *ChatService) CreateUser(
	ctx context.Context,
	req *connect.Request[chatv1.CreateUserRequest],
) (*connect.Response[chatv1.CreateUserResponse], error) {

	if req.Msg.GetUsername() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("username required"))
	}

	user := &chatv1.User{
		Id:        memory.NewID(),
		Username:  req.Msg.GetUsername(),
		CreatedAt: memory.Now(),
	}
	s.store.Mu.Lock()
	s.store.Users[user.Id] = user
	s.store.Mu.Unlock()

	return connect.NewResponse(&chatv1.CreateUserResponse{
		User: user,
	}), nil
}

func (s *ChatService) CreateServer(
	ctx context.Context,
	req *connect.Request[chatv1.CreateServerRequest],
) (*connect.Response[chatv1.CreateServerResponse], error) {

	if req.Msg.GetName() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("server name required"))
	}

	server := &chatv1.Server{
		Id:        memory.NewID(),
		Name:      req.Msg.GetName(),
		CreatedAt: memory.Now(),
	}

	s.store.Mu.Lock()
	s.store.Servers[server.Id] = server
	s.store.Mu.Unlock()

	return connect.NewResponse(&chatv1.CreateServerResponse{
		Server: server,
	}), nil
}

func (s *ChatService) CreateChannel(
	ctx context.Context,
	req *connect.Request[chatv1.CreateChannelRequest],
) (*connect.Response[chatv1.CreateChannelResponse], error) {

	if req.Msg.GetName() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("channel name required"))
	}

	s.store.Mu.RLock()
	_, exists := s.store.Servers[req.Msg.GetServerId()]
	s.store.Mu.RUnlock()

	if !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("server not found"))
	}

	channel := &chatv1.Channel{
		Id:        memory.NewID(),
		ServerId:  req.Msg.GetServerId(),
		Name:      req.Msg.GetName(),
		CreatedAt: memory.Now(),
	}

	s.store.Mu.Lock()
	s.store.Channels[channel.Id] = channel
	s.store.Mu.Unlock()

	return connect.NewResponse(&chatv1.CreateChannelResponse{
		Channel: channel,
	}), nil
}

func (s *ChatService) SendChannelMessage(
	ctx context.Context,
	req *connect.Request[chatv1.SendChannelMessageRequest],
) (*connect.Response[chatv1.SendChannelMessageResponse], error) {

	if req.Msg.GetContent() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("content required"))
	}

	message := &chatv1.Message{
		Id:           memory.NewID(),
		AuthorUserId: req.Msg.GetAuthorUserId(),
		CreatedAt:    memory.Now(),
		Content:      req.Msg.GetContent(),
		Destination: &chatv1.Message_Channel{
			Channel: &chatv1.ChannelDestination{
				ServerId:  req.Msg.GetServerId(),
				ChannelId: req.Msg.GetChannelId(),
			},
		},
	}

	s.store.Mu.Lock()
	s.store.Messages[req.Msg.GetChannelId()] = append(s.store.Messages[req.Msg.GetChannelId()], message)
	s.store.Mu.Unlock()

	return connect.NewResponse(&chatv1.SendChannelMessageResponse{
		Message: message,
	}), nil
}
