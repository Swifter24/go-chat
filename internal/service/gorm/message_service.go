package gorm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go_chat/internal/dao"
	"go_chat/internal/dto/respond"
	"go_chat/internal/model"
	myredis "go_chat/internal/service/redis"
	"go_chat/pkg/constants"
	"go_chat/pkg/zlog"
)

type messageService struct {
}

var MessageService = new(messageService)

// GetMessageList 获取聊天记录
func (m *messageService) GetMessageList(userOneId, userTwoId string) (string, []respond.GetMessageListRespond, int) {
	rspString, err := myredis.GetKeyNilIsErr("message_list_" + userOneId + "_" + userTwoId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			zlog.Info(err.Error())
			zlog.Info(fmt.Sprintf("%s %s", userTwoId, userTwoId))
			var messageList []model.Message
			if res := dao.GormDB.Where("(send_id = ? AND receive_id = ?) OR (send_id = ? AND receive_id = ?)", userOneId, userTwoId, userTwoId, userOneId).Order("created_at ASC").Find(&messageList); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, nil, -1
			}
			var rspList []respond.GetMessageListRespond
			for _, message := range messageList {
				rspList = append(rspList, respond.GetMessageListRespond{
					SendId:     message.SendId,
					SendName:   message.SendName,
					SendAvatar: message.SendAvatar,
					ReceiveId:  message.ReceiveId,
					Content:    message.Content,
					Url:        message.Url,
					Type:       message.Type,
					FileType:   message.FileType,
					FileName:   message.FileName,
					FileSize:   message.FileSize,
					CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
				})
			}
			//rspString, err := json.Marshal(rspList)
			//if err != nil {
			//	zlog.Error(err.Error())
			//}
			//if err := myredis.SetKeyEx("message_list_"+userOneId+"_"+userTwoId, string(rspString), time.Minute*constants.REDIS_TIMEOUT); err != nil {
			//	zlog.Error(err.Error())
			//}
			return "获取聊天记录成功", rspList, 0
		} else {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
	}
	var rsp []respond.GetMessageListRespond
	if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
		zlog.Error(err.Error())
	}
	return "获取聊天记录成功", rsp, 0
}

// GetGroupMessageList 获取群聊消息记录
func (m *messageService) GetGroupMessageList(groupId string) (string, []respond.GetGroupMessageListRespond, int) {
	rspString, err := myredis.GetKeyNilIsErr("group_messagelist_" + groupId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			var messageList []model.Message
			if res := dao.GormDB.Where("receive_id = ?", groupId).Order("created_at ASC").Find(&messageList); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, nil, -1
			}
			var rspList []respond.GetGroupMessageListRespond
			for _, message := range messageList {
				rsp := respond.GetGroupMessageListRespond{
					SendId:     message.SendId,
					SendName:   message.SendName,
					SendAvatar: message.SendAvatar,
					ReceiveId:  message.ReceiveId,
					Content:    message.Content,
					Url:        message.Url,
					Type:       message.Type,
					FileType:   message.FileType,
					FileName:   message.FileName,
					FileSize:   message.FileSize,
					CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
				}
				rspList = append(rspList, rsp)
			}
			//rspString, err := json.Marshal(rspList)
			//if err != nil {
			//	zlog.Error(err.Error())
			//}
			//if err := myredis.SetKeyEx("group_messagelist_"+groupId, string(rspString), time.Minute*constants.REDIS_TIMEOUT); err != nil {
			//	zlog.Error(err.Error())
			//}
			return "获取聊天记录成功", rspList, 0
		} else {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
	}
	var rsp []respond.GetGroupMessageListRespond
	if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
		zlog.Error(err.Error())
	}
	return "获取聊天记录成功", rsp, 0
}
