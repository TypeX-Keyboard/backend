package ws

import (
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmlock"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gorilla/websocket"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/util/cache"
	"runtime/debug"
)

const (
	// The user connection timeout
	heartbeatExpirationTime = 6 * 60
)

// User login
type login struct {
	UserId uint64
	Client *Client
}

// GetKey Read client data
func (l *login) GetKey() (key string) {
	key = GetUserKey(l.UserId)
	return
}

// Client connections
type Client struct {
	Addr          string          // Client address
	ID            string          // Unique identifier of the connection
	Socket        *websocket.Conn // User connection
	Send          chan *WResponse // data to be sent
	SendClose     bool            // Whether the sending is turned off
	UserId        uint64          // User ID, which is available after the user logs in
	FirstTime     uint64          // First connection event
	HeartbeatTime uint64          // The time of the user's last heartbeat
	LoginTime     uint64          // LoginTime LoginTime LoginTime uint64 // LoginTime LoginTime Log
	isApp         bool            // Whether it's an app or not
	AuthToken     string
	tags          garray.StrArray // tags
	Tokens        garray.StrArray
}

// NewClient initialize
func NewClient(addr string, socket *websocket.Conn, firstTime uint64) (client *Client) {
	client = &Client{
		Addr:          addr,
		ID:            guid.S(),
		Socket:        socket,
		Send:          make(chan *WResponse, 100),
		SendClose:     false,
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}
	return
}

// Read client data
func (c *Client) read() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", string(debug.Stack()), r)
		}
	}()

	defer func() {
		c.close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			return
		}
		ProcessData(c, message)
	}
}

// Write data to the client
func (c *Client) write() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", string(debug.Stack()), r)
		}
	}()
	defer func() {
		clientManager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 发送数据错误 关闭连接
				return
			}
			_ = c.Socket.WriteJSON(message)
		}
	}
}

// SendMsg
func (c *Client) SendMsg(msg *WResponse) {
	if c == nil || c.SendClose {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendMsg stop:", r, string(debug.Stack()))
		}
	}()
	c.Send <- msg
}

// Heartbeat
func (c *Client) Heartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime
	return
}

// IsHeartbeatTimeout
func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}
	return
}

// Close the client
func (c *Client) close() {
	if c.SendClose {
		return
	}
	c.SendClose = true
	ctx := gctx.New()
	gmlock.Lock(consts.ListenTokenPrices)
	defer gmlock.Unlock(consts.ListenTokenPrices)
	get, err := cache.GetCache().Get(ctx, consts.ListenTokenPrices)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	tokens := make(map[string]int)
	err = gconv.Struct(get, &tokens)
	for _, token := range c.Tokens.Slice() {
		tokens[token] = tokens[token] - 1
		if tokens[token] <= 0 {
			delete(tokens, token)
		}
	}
	err = cache.GetCache().Set(ctx, consts.ListenTokenPrices, tokens, 0)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	close(c.Send)
}
