package service

import (
    "chat-service/internal/model"
    "chat-service/internal/adapters/out"
    "log"
)

// ChatService handles the business logic of sending and receiving messages.
type ChatService struct {
    MessageRepo out.MongoMessageRepository
}

// NewChatService creates a new instance of ChatService.
func NewChatService(messageRepo out.MongoMessageRepository) *ChatService {
    return &ChatService{MessageRepo: messageRepo}
}

// SendMessage saves a new message to the database.
func (s *ChatService) SendMessage(message *model.Message) error {
    err := s.MessageRepo.SaveMessage(message)
    if err != nil {
        log.Println("Error saving message:", err)
        return err
    }
    return nil
}

// GetMessages fetches all messages between two users.
func (s *ChatService) GetMessages(sender, receiver string) ([]model.Message, error) {
    messages, err := s.MessageRepo.GetMessagesByUser(sender, receiver)
    if err != nil {
        log.Println("Error fetching messages:", err)
        return nil, err
    }
    return messages, nil
}

func (s *ChatService) GetAllMessages() ([]model.Message, error) {
    messages, err := s.MessageRepo.GetAllMessages()
    if err != nil {
        log.Println("Error fetching messages:", err)
        return nil, err
    }
    return messages, nil
}

func (s *ChatService) DeleteAllMessages() error {
    err := s.MessageRepo.DeleteAllMessages()
    if err != nil {
        log.Println("Error deleting messages:", err)
        return err
    }
    return nil
}
