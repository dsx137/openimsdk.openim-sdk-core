package conversation_msg

import (
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
)

func TestApplyUnreadOnlyChanges_whenConversationExists(t *testing.T) {
	// Given: an existing visible conversation and an unread-only notification.
	local := map[string]*model_struct.LocalConversation{
		"group": {
			ConversationID:    "group",
			UnreadCount:       10,
			LatestMsg:         "latest text message",
			LatestMsgSendTime: 100,
		},
	}
	unreadOnly := map[string]*model_struct.LocalConversation{
		"group": {ConversationID: "group", UnreadCount: 1},
	}
	changed := make(map[string]*model_struct.LocalConversation)

	// When: the unread-only changes are applied.
	hidden := applyUnreadOnlyChanges(local, unreadOnly, changed)

	// Then: unread advances without replacing the latest message projection.
	conversation := changed["group"]
	if conversation == nil {
		t.Fatal("expected the existing conversation to be changed")
	}
	if conversation.UnreadCount != 11 {
		t.Fatalf("expected unread count 11, got %d", conversation.UnreadCount)
	}
	if conversation.LatestMsg != "latest text message" || conversation.LatestMsgSendTime != 100 {
		t.Fatalf("latest message changed: %+v", conversation)
	}
	if len(hidden) != 0 {
		t.Fatalf("expected no hidden conversations, got %d", len(hidden))
	}
}

func TestApplyUnreadOnlyChanges_whenConversationAlreadyChanged(t *testing.T) {
	// Given: a normal conversation update and an unread-only notification in one batch.
	changed := map[string]*model_struct.LocalConversation{
		"group": {ConversationID: "group", UnreadCount: 2, LatestMsg: "new text", LatestMsgSendTime: 200},
	}
	unreadOnly := map[string]*model_struct.LocalConversation{
		"group": {ConversationID: "group", UnreadCount: 1},
	}

	// When: the unread-only changes are applied.
	applyUnreadOnlyChanges(nil, unreadOnly, changed)

	// Then: the delta is applied once to the already changed conversation.
	if changed["group"].UnreadCount != 3 {
		t.Fatalf("expected unread count 3, got %d", changed["group"].UnreadCount)
	}
}

func TestApplyUnreadOnlyChanges_whenConversationDoesNotExist(t *testing.T) {
	// Given: an unread-only notification without a local conversation row.
	unreadOnly := map[string]*model_struct.LocalConversation{
		"group": {ConversationID: "group", GroupID: "group", UnreadCount: 1},
	}

	// When: the unread-only changes are applied.
	hidden := applyUnreadOnlyChanges(nil, unreadOnly, make(map[string]*model_struct.LocalConversation))

	// Then: a hidden conversation row preserves the unread count without a latest message.
	if len(hidden) != 1 {
		t.Fatalf("expected one hidden conversation, got %d", len(hidden))
	}
	if hidden[0].UnreadCount != 1 || hidden[0].LatestMsg != "" || hidden[0].LatestMsgSendTime != 0 {
		t.Fatalf("unexpected hidden conversation: %+v", hidden[0])
	}
}
