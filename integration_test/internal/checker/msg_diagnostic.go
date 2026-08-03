package checker

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	sdkUtils "github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/protocol/msg"
)

func diagnoseMessageCount(ctx context.Context, testSDK *sdk.TestSDK, total, correct int) {
	conversations, err := testSDK.SDK.Conversation().GetAllConversationList(ctx)
	if err != nil {
		fmt.Printf(">>> DIAG userID=%s GetAllConversationList err=%v\n", testSDK.UserID, err)
		return
	}

	diagnosticCtx := ccontext.WithOperationID(testSDK.SDK.Context(), sdkUtils.OperationIDGenerator())
	var resp *msg.GetConversationsHasReadAndMaxSeqResp
	for attempt := 0; attempt < 3; attempt++ {
		resp, err = api.GetConversationsHasReadAndMaxSeq.Invoke(
			diagnosticCtx,
			&msg.GetConversationsHasReadAndMaxSeqReq{UserID: testSDK.UserID},
		)
		if err == nil {
			break
		}
		if attempt < 2 {
			time.Sleep(time.Second)
		}
	}
	if err != nil {
		fmt.Printf(">>> DIAG userID=%s GetConversationsHasReadAndMaxSeq err=%v\n", testSDK.UserID, err)
		return
	}

	localByID := make(map[string]int32, len(conversations))
	localMaxSeqByID := make(map[string]int64, len(conversations))
	var localSum int64
	for _, conversation := range conversations {
		localByID[conversation.ConversationID] = conversation.UnreadCount
		localMaxSeqByID[conversation.ConversationID] = conversation.MaxSeq
		localSum += int64(conversation.UnreadCount)
	}

	conversationIDs := make([]string, 0, len(resp.Seqs))
	for conversationID := range resp.Seqs {
		conversationIDs = append(conversationIDs, conversationID)
	}
	sort.Strings(conversationIDs)

	var serverSum int64
	for _, conversationID := range conversationIDs {
		seq := resp.Seqs[conversationID]
		serverUnread := max(seq.MaxSeq-seq.HasReadSeq, 0)
		serverSum += serverUnread
		fmt.Printf(">>> DIAG userID=%s conv=%s local_unread=%d local_max_seq=%d server_unread=%d server_max_seq=%d server_has_read_seq=%d\n",
			testSDK.UserID, conversationID, localByID[conversationID], localMaxSeqByID[conversationID],
			serverUnread, seq.MaxSeq, seq.HasReadSeq)
		delete(localByID, conversationID)
	}

	for conversationID, localUnread := range localByID {
		fmt.Printf(">>> DIAG userID=%s conv=%s local_unread=%d local_max_seq=%d server_state=missing\n",
			testSDK.UserID, conversationID, localUnread, localMaxSeqByID[conversationID])
	}

	classification := classifyUnreadMismatch(localSum, serverSum, int64(correct))
	fmt.Printf(">>> DIAG SUMMARY userID=%s checker_total=%d local_sum=%d server_sum=%d correct=%d classification=%s\n",
		testSDK.UserID, total, localSum, serverSum, correct, classification)
}

func classifyUnreadMismatch(local, server, correct int64) string {
	switch {
	case server == correct && local != server:
		return "local_unread_state"
	case local == server && server != correct:
		return "server_state_or_expectation"
	case local != server:
		return "mixed_local_and_server_state"
	default:
		return "no_mismatch"
	}
}
