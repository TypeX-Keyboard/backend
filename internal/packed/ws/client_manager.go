package ws

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/util/cache"
	"sync"
)

// ClientManager
type ClientManager struct {
	Clients         map[*Client]bool      // All connections
	ClientsLock     sync.RWMutex          // Read/write lock
	Users           map[string]*Client    // Logged-in user // uuid
	UserLock        sync.RWMutex          // Read/write lock
	Register        chan *Client          // Connection connection processing
	Login           chan *login           // User login processing
	Unregister      chan *Client          // Disconnect Handler
	Broadcast       chan *WResponse       // Broadcast Send data to all members
	PriceBroadcast  chan *PriceResponse   // Broadcast sends data to subscription price members
	ClientBroadcast chan *ClientWResponse // Broadcast Sends data to a client
	TagBroadcast    chan *TagWResponse    // Broadcast Sends data to a tag member
	UserBroadcast   chan *UserWResponse   // Broadcast sends data to all links of a user
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:         make(map[*Client]bool),
		Users:           make(map[string]*Client),
		Register:        make(chan *Client, 1000),
		Login:           make(chan *login, 1000),
		Unregister:      make(chan *Client, 1000),
		Broadcast:       make(chan *WResponse, 1000),
		PriceBroadcast:  make(chan *PriceResponse, 1000),
		ClientBroadcast: make(chan *ClientWResponse, 1000),
		TagBroadcast:    make(chan *TagWResponse, 1000),
		UserBroadcast:   make(chan *UserWResponse, 1000),
	}
	return
}

func GetUserKey(userId uint64) (key string) {
	key = fmt.Sprintf("%s_%d", "ws", userId)
	return
}

func (manager *ClientManager) InClient(client *Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()
	_, ok = manager.Clients[client]
	return
}

func (manager *ClientManager) GetClients() (clients map[*Client]bool) {
	clients = make(map[*Client]bool)
	manager.ClientsRange(func(client *Client, value bool) (result bool) {
		clients[client] = value
		return true
	})
	return
}

func (manager *ClientManager) ClientsRange(f func(client *Client, value bool) (result bool)) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()
	for key, value := range manager.Clients {
		result := f(key, value)
		if result == false {
			return
		}
	}
	return
}

func (manager *ClientManager) GetClientsLen() (clientsLen int) {
	clientsLen = len(manager.Clients)
	return
}

func (manager *ClientManager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	manager.Clients[client] = true
}

func (manager *ClientManager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
}

func (manager *ClientManager) GetUserClient(userId uint64) (client *Client) {
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	userKey := GetUserKey(userId)
	if value, ok := manager.Users[userKey]; ok {
		client = value
	}
	return
}

func (manager *ClientManager) AddUsers(key string, client *Client) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()
	manager.Users[key] = client
}

func (manager *ClientManager) DelUsers(client *Client) (result bool) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()
	key := GetUserKey(client.UserId)
	if value, ok := manager.Users[key]; ok {
		// Determine whether the user is the same
		if value.Addr != client.Addr {
			return
		}
		delete(manager.Users, key)
		result = true
	}
	return
}

// GetUsersLen
func (manager *ClientManager) GetUsersLen() (userLen int) {
	userLen = len(manager.Users)
	return
}

// EventRegister
func (manager *ClientManager) EventRegister(client *Client) {
	manager.AddClients(client)
	client.SendMsg(&WResponse{Event: "connected", Data: g.Map{
		"ID": client.ID,
	}})
}

// EventLogin
func (manager *ClientManager) EventLogin(login *login) {
	client := login.Client
	if manager.InClient(client) {
		userKey := login.GetKey()
		manager.AddUsers(userKey, login.Client)
	}
}

// EventUnregister
func (manager *ClientManager) EventUnregister(client *Client) {
	manager.DelClients(client)
	deleteResult := manager.DelUsers(client)
	if deleteResult == false {
		// Not the currently connected client
		return
	}
	close(client.Send)
}

func (manager *ClientManager) clearTimeoutConnections() {
	currentTime := uint64(gtime.Now().Unix())
	clients := clientManager.GetClients()
	for client := range clients {
		ctx := gctx.New()
		get, err := cache.GetCache().Get(ctx, consts.RedisClientHMACKey)
		if err != nil {
			g.Log().Error(ctx, err)
			go AuthErr(client)
			return
		}
		if get == nil {
			g.Log().Error(ctx, err)
			go AuthErr(client)
			return
		}
		if client.AuthToken != get.String() {
			go AuthErr(client)
			return
		}
		if client.IsHeartbeatTimeout(currentTime) {
			_ = client.Socket.Close()
		}
	}
}

func (manager *ClientManager) ping(ctx context.Context) {
	//定时任务，发送心跳包
	_, _ = gcron.Add(ctx, "0 */1 * * * *", func(ctx context.Context) {
		res := &WResponse{
			Event: Ping,
			Data:  g.Map{},
		}
		SendToAll(res)
	})
	_, _ = gcron.Add(ctx, "*/30 * * * * *", func(ctx context.Context) {
		manager.clearTimeoutConnections()
	})

}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.Register:
			manager.EventRegister(conn)

		case login := <-manager.Login:
			manager.EventLogin(login)

		case conn := <-manager.Unregister:
			manager.EventUnregister(conn)

		case message := <-manager.Broadcast:
			clients := manager.GetClients()
			for conn := range clients {
				conn.SendMsg(message)
			}
		case message := <-manager.PriceBroadcast:
			clients := manager.GetClients()
			for conn := range clients {
				for _, s := range conn.Tokens.Slice() {
					conn.SendMsg(&WResponse{
						Event: message.Event,
						Data:  message.Data[s],
					})
				}
			}
		case message := <-manager.TagBroadcast:
			clients := manager.GetClients()
			for conn := range clients {
				if conn.tags.Contains(message.Tag) {
					conn.SendMsg(message.WResponse)
				}
			}
		case message := <-manager.UserBroadcast:
			clients := manager.GetClients()
			for conn := range clients {
				if conn.UserId == message.UserID {
					conn.SendMsg(message.WResponse)
				}
			}
		case message := <-manager.ClientBroadcast:
			clients := manager.GetClients()
			for conn := range clients {
				if conn.ID == message.ID {
					conn.SendMsg(message.WResponse)
				}
			}
		}

	}
}

func SendPrice(response *PriceResponse) {
	clientManager.PriceBroadcast <- response
}

func SendToAll(response *WResponse) {
	clientManager.Broadcast <- response
}

func SendToClientID(id string, response *WResponse) {
	clientRes := &ClientWResponse{
		ID:        id,
		WResponse: response,
	}
	clientManager.ClientBroadcast <- clientRes
}

func SendToUser(userID uint64, response *WResponse) {
	userRes := &UserWResponse{
		UserID:    userID,
		WResponse: response,
	}
	clientManager.UserBroadcast <- userRes
}

func SendToTag(tag string, response *WResponse) {
	tagRes := &TagWResponse{
		Tag:       tag,
		WResponse: response,
	}
	clientManager.TagBroadcast <- tagRes
}
