package service

import (
	"ai-service/internal/adapters/outbound"
	"ai-service/internal/model"
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	systemPrompt = "You are the in-app assistant for a Facebook-like social platform with friends, posts, notifications, chat, and a slot game. Answer briefly and help the user use the product. Do not invent private user data."
	historyLimit = 16
	maxMessage   = 2000
)

type AIService interface {
	GetInbox(ctx context.Context, userID string) (*model.Conversation, []model.Message, error)
	Chat(ctx context.Context, userID, conversationID, content string) (*model.Conversation, *model.Message, *model.Message, error)
	Clear(ctx context.Context, userID, conversationID string) error
}

type aiService struct {
	conversations outbound.ConversationRepository
	messages      outbound.MessageRepository
	llm           outbound.LLMClient
}

func NewAIService(
	conversations outbound.ConversationRepository,
	messages outbound.MessageRepository,
	llm outbound.LLMClient,
) AIService {
	return &aiService{conversations: conversations, messages: messages, llm: llm}
}

func (s *aiService) GetInbox(ctx context.Context, userID string) (*model.Conversation, []model.Message, error) {
	if userID == "" {
		return nil, nil, errors.New("user_id is required")
	}
	conversation, err := s.conversations.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	messages, err := s.messages.ListByConversation(ctx, conversation.ID.Hex(), 50)
	if err != nil {
		return nil, nil, err
	}
	return conversation, messages, nil
}

func (s *aiService) Chat(ctx context.Context, userID, conversationID, content string) (*model.Conversation, *model.Message, *model.Message, error) {
	content = strings.TrimSpace(content)
	if userID == "" {
		return nil, nil, nil, errors.New("user_id is required")
	}
	if content == "" {
		return nil, nil, nil, errors.New("message is required")
	}
	if utf8.RuneCountInString(content) > maxMessage {
		return nil, nil, nil, errors.New("message is too long")
	}

	var (
		conversation *model.Conversation
		err          error
	)
	if conversationID != "" {
		conversation, err = s.conversations.GetByID(ctx, conversationID, userID)
	} else {
		conversation, err = s.conversations.GetOrCreate(ctx, userID)
	}
	if err != nil {
		return nil, nil, nil, err
	}

	userMessage := &model.Message{
		ID:             primitive.NewObjectID(),
		ConversationID: conversation.ID.Hex(),
		Role:           "user",
		Content:        content,
	}
	if err := s.messages.Create(ctx, userMessage); err != nil {
		return nil, nil, nil, err
	}

	history, err := s.messages.ListByConversation(ctx, conversation.ID.Hex(), historyLimit)
	if err != nil {
		return nil, nil, nil, err
	}

	turns := []outbound.ChatTurn{{Role: "system", Content: systemPrompt}}
	for _, item := range history {
		if item.Role == "user" || item.Role == "assistant" {
			turns = append(turns, outbound.ChatTurn{Role: item.Role, Content: item.Content})
		}
	}

	reply, err := s.llm.Complete(ctx, turns)
	if err != nil || strings.TrimSpace(reply) == "" {
		reply = "I could not reach the model just now. Try again in a moment."
	}

	assistantMessage := &model.Message{
		ID:             primitive.NewObjectID(),
		ConversationID: conversation.ID.Hex(),
		Role:           "assistant",
		Content:        reply,
	}
	if err := s.messages.Create(ctx, assistantMessage); err != nil {
		return nil, nil, nil, err
	}
	_ = s.conversations.Touch(ctx, conversation.ID.Hex())

	return conversation, userMessage, assistantMessage, nil
}

func (s *aiService) Clear(ctx context.Context, userID, conversationID string) error {
	if userID == "" || conversationID == "" {
		return errors.New("conversation_id is required")
	}
	conversation, err := s.conversations.GetByID(ctx, conversationID, userID)
	if err != nil {
		return err
	}
	if err := s.messages.DeleteByConversation(ctx, conversation.ID.Hex()); err != nil {
		return err
	}
	return s.conversations.Delete(ctx, conversation.ID.Hex(), userID)
}
