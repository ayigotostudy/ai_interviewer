package rag

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
	RedisOptions  *redis.Client
}

func GetDefaultRedisMemory() *redisMemory {
	return NewRedisMemory(RedisMemoryConfig{
		MaxWindowSize: 6,
		RedisOptions: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	})
}

func NewRedisMemory(cfg RedisMemoryConfig) *redisMemory {
	if cfg.RedisOptions == nil {
		cfg.RedisOptions = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
	}

	return &redisMemory{
		client:        cfg.RedisOptions,
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

	LastConversationsKnowledge string `json:"last_conversations_knowledge"`
	RoundCount                 int    `json:"round_count"` // 对话轮数
}

func (c *Conversation) Append(msg ...*schema.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Messages = append(c.Messages, msg...)
	c.RoundCount++ // 增加轮数
	c.save()
}

func (c *Conversation) GetRoundCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.RoundCount
}

func (c *Conversation) SetRoundCount(count int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.RoundCount = count
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

// 修改load方法，使其加载整个Conversation对象，包括RoundCount
func (c *Conversation) load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ctx := context.Background()
	data, err := c.client.Get(ctx, "conversation:"+c.ID).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get conversation: %w", err)
	}

	// 反序列化整个Conversation对象，包括RoundCount
	conversationData := &Conversation{}
	if err := json.Unmarshal([]byte(data), conversationData); err != nil {
		return fmt.Errorf("failed to unmarshal conversation: %w", err)
	}

	// 保留原有指针和配置
	c.Messages = conversationData.Messages
	c.RoundCount = conversationData.RoundCount
	c.LastConversationsKnowledge = conversationData.LastConversationsKnowledge
	
	return nil
}

// 修改save方法，使其保存整个Conversation对象，包括RoundCount
func (c *Conversation) save() {
	c.mu.Lock()
	defer c.mu.Unlock()

	ctx := context.Background()
	// 创建一个副本，不包含client和maxWindowSize等不需要序列化的字段
	serializableData := &Conversation{
		Messages:                  c.Messages,
		RoundCount:                c.RoundCount,
		LastConversationsKnowledge: c.LastConversationsKnowledge,
	}

	data, err := json.Marshal(serializableData)
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
