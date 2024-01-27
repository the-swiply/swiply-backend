package converter

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/chat/internal/domain"
	"github.com/the-swiply/swiply-backend/chat/pkg/api/chat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MessagesToPB(messages []domain.ChatMessage) []*chat.ChatMessage {
	res := make([]*chat.ChatMessage, 0, len(messages))

	for _, msg := range messages {
		res = append(res, &chat.ChatMessage{
			Id:       msg.ID.String(),
			ChatId:   msg.ChatID,
			IdInChat: msg.IDInChat,
			FromId:   msg.From.String(),
			Content:  msg.Content,
			SendTime: timestamppb.New(msg.SendTime),
		})
	}

	return res
}

func ChatsToPB(chats []domain.Chat) []*chat.GenericChat {
	res := make([]*chat.GenericChat, 0, len(chats))

	for _, ch := range chats {
		members := make([]string, 0, len(ch.Members))
		for _, member := range ch.Members {
			members = append(members, member.String())
		}

		res = append(res, &chat.GenericChat{
			Id:      ch.ID,
			Members: members,
		})
	}

	return res
}

func StringsToUUIDs(strings []string) ([]uuid.UUID, error) {
	res := make([]uuid.UUID, 0, len(strings))

	for _, s := range strings {
		su, err := uuid.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("can't parse to uuid string %s: %w", s, err)
		}

		res = append(res, su)
	}

	return res, nil
}
