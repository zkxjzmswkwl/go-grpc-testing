package memory

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	chatv1 "github.com/zkxjzmswkwl/go-grpc-testing/gen/go/proto/chat/v1"
)

// (csmith): This is pretty silly and not meant to be used in prod.
type Store struct {
	Mu sync.RWMutex

	Users    map[string]*chatv1.User
	Servers  map[string]*chatv1.Server
	Channels map[string]*chatv1.Channel
	Messages map[string][]*chatv1.Message
}

func NewStore() *Store {
	return &Store{
		Users:    make(map[string]*chatv1.User),
		Servers:  make(map[string]*chatv1.Server),
		Channels: make(map[string]*chatv1.Channel),
		Messages: make(map[string][]*chatv1.Message),
	}
}

func NewID() string {
	return uuid.NewString()
}

func Now() *timestamppb.Timestamp {
	return timestamppb.New(time.Now())
}
