package checker

import (
	"context"
	"sync"

	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
)

// CheckMessageNum check message num.
func CheckMessageNum(ctx context.Context) error {
	var diagnoseOnce sync.Once

	c := &CounterChecker[*sdk.TestSDK, string]{
		CheckName:      "checkMessageNum",
		CheckerKeyName: "userID",
		GoroutineLimit: config.ErrGroupCommonLimit,
		GetTotalCount: func(ctx context.Context, t *sdk.TestSDK) (int, error) {
			totalNum, err := t.SDK.Conversation().GetTotalUnreadMsgCount(ctx)
			if err != nil {
				return 0, err
			}
			return int(totalNum), nil
		},
		CalCorrectCount: ExpectedMessageCount,
		LoopSlice:       sdk.TestSDKs,
		GetKey: func(t *sdk.TestSDK) string {
			return t.UserID
		},
		OnFail: func(ctx context.Context, t *sdk.TestSDK, total, correct int) {
			diagnoseOnce.Do(func() {
				diagnoseMessageCount(ctx, t, total, correct)
			})
		},
	}

	return c.LoopCheck(ctx)
}

func ExpectedMessageCount(userID string) int {
	userNum := utils.MustGetUserNum(userID)
	createdLargeGroupNum := vars.LargeGroupNum / vars.LoginUserNum
	remainder := vars.LargeGroupNum % vars.LoginUserNum
	groupMsgNum := vars.GroupMessageNum*vars.LoginUserNum*vars.LargeGroupNum + vars.LargeGroupNum
	commonUserMsgNum := min(vars.LoginUserNum, vars.SuperUserNum) * vars.SingleMessageNum

	var result int
	if utils.IsSuperUser(userID) {
		result = groupMsgNum + vars.UserNum - 1 - userNum
		if userNum < vars.LoginUserNum {
			result += vars.SingleMessageNum * (vars.LoginUserNum - 1)
			result -= vars.GroupMessageNum*vars.LargeGroupNum + createdLargeGroupNum
		} else {
			result += vars.SingleMessageNum * vars.LoginUserNum
		}
	} else {
		result = commonUserMsgNum + groupMsgNum
		if userNum < vars.LoginUserNum {
			result -= vars.GroupMessageNum*vars.LargeGroupNum + createdLargeGroupNum
		}
	}

	commonGroupNum := calCommonGroup(userNum) * (vars.GroupMessageNum + 1)
	if utils.IsNumLogin(userNum) {
		commonGroupNum -= vars.CommonGroupNum * (vars.GroupMessageNum + 1)
	}
	result += commonGroupNum
	if userNum < remainder {
		result--
	}
	return result
}
