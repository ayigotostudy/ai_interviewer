package component

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
)

type RedisMemoryConfig struct {
	MaxWindowSize int
	RedisOptions  *redis.Options
}

func GetDefaultRedisMemory() *redisMemory {
	return NewRedisMemory(RedisMemoryConfig{
		MaxWindowSize: 6,
		RedisOptions: &redis.Options{
			Addr: "localhost:6379",
		},
	})
}

func NewRedisMemory(cfg RedisMemoryConfig) *redisMemory {
	if cfg.RedisOptions == nil {
		cfg.RedisOptions = &redis.Options{
			Addr: "localhost:6379",
		}
	}

	client := redis.NewClient(cfg.RedisOptions)

	return &redisMemory{
		client:        client,
		maxWindowSize: cfg.MaxWindowSize,
		conversations: make(map[string]*Conversation),
	}
}

// simple memory can store messages of each conversation
type redisMemory struct {
	mu            sync.Mutex
	client        *redis.Client
	maxWindowSize int
	conversations map[string]*Conversation
}

func (m *redisMemory) GetConversation(id string, createIfNotExist bool) *Conversation {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.conversations[id]
	if !ok {
		if createIfNotExist {
			con := &Conversation{
				ID:            id,
				Messages:      make([]*schema.Message, 0),
				client:        m.client,
				maxWindowSize: m.maxWindowSize,
			}
			m.conversations[id] = con
		} else {
			con := &Conversation{
				ID:            id,
				Messages:      make([]*schema.Message, 0),
				client:        m.client,
				maxWindowSize: m.maxWindowSize,
			}
			con.load()
			m.conversations[id] = con
		}
	}

	return m.conversations[id]
}

func (c *Conversation) GetLastConversationsKnowledge() string {
	return c.LastConversationsKnowledge
}

func (c *Conversation) SetLastConversationKnowledge(knowledge string) {
	c.LastConversationsKnowledge = knowledge
}

func (m *redisMemory) ListConversations() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx := context.Background()
	keys, err := m.client.Keys(ctx, "conversation:*").Result()
	if err != nil {
		return nil
	}

	ids := make([]string, 0, len(keys))
	for _, key := range keys {
		ids = append(ids, strings.TrimPrefix(key, "conversation:"))
	}

	return ids
}

func (m *redisMemory) DeleteConversation(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx := context.Background()
	if err := m.client.Del(ctx, "conversation:"+id).Err(); err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	delete(m.conversations, id)
	return nil
}

type Conversation struct {
	mu sync.Mutex

	ID            string            `json:"id"`
	Messages      []*schema.Message `json:"messages"`
	client        *redis.Client
	maxWindowSize int

	LastConversationsKnowledge string
}

func (c *Conversation) Append(msg ...*schema.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Messages = append(c.Messages, msg...)
	c.save()
}

func (c *Conversation) GetFullMessages() []*schema.Message {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Messages
}

// get messages with max window size
func (c *Conversation) GetMessages() []*schema.Message {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.Messages) > c.maxWindowSize {
		return c.Messages[len(c.Messages)-c.maxWindowSize:]
	}

	return c.Messages
}

func (c *Conversation) load() error {
	ctx := context.Background()
	data, err := c.client.Get(ctx, "conversation:"+c.ID).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get conversation: %w", err)
	}

	var messages []*schema.Message
	if err := json.Unmarshal([]byte(data), &messages); err != nil {
		return fmt.Errorf("failed to unmarshal messages: %w", err)
	}
	c.Messages = messages
	return nil
}

func (c *Conversation) save() {
	ctx := context.Background()
	data, err := json.Marshal(c.Messages)
	if err != nil {
		return
	}
	c.client.Set(ctx, "conversation:"+c.ID, data, 0)
}

func (c *Conversation) String() string {
	content := ""
	for _, v := range c.Messages {
		content += v.String() + "\n"
	}
	return content
}
