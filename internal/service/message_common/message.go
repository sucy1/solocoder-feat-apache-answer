/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package messagecommon

import (
	"context"

	"github.com/apache/answer/internal/base/pager"
	"github.com/apache/answer/internal/base/reason"
	"github.com/apache/answer/internal/entity"
	"github.com/apache/answer/internal/repo/message"
	"github.com/apache/answer/internal/schema"
	usercommon "github.com/apache/answer/internal/service/user_common"
	"github.com/segmentfault/pacman/errors"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxContentLen   = 65536
	MaxTitleLen     = 200
	OldMessageDays  = 90
)

// MessageCommon message service
type MessageCommon struct {
	messageRepo      message.MessageRepo
	messageBlockRepo message.MessageBlockRepo
	userCommon       *usercommon.UserCommon
}

func NewMessageCommon(
	messageRepo message.MessageRepo,
	messageBlockRepo message.MessageBlockRepo,
	userCommon *usercommon.UserCommon,
) *MessageCommon {
	return &MessageCommon{
		messageRepo:      messageRepo,
		messageBlockRepo: messageBlockRepo,
		userCommon:       userCommon,
	}
}

func (ms *MessageCommon) SendMessage(ctx context.Context, req *schema.SendMessageReq) (*schema.GetMessageResp, error) {
	if req.UserID == req.ToUserID {
		return nil, errors.BadRequest(reason.RequestFormatError)
	}

	blocked, err := ms.messageBlockRepo.IsBlocked(ctx, req.ToUserID, req.UserID)
	if err != nil {
		return nil, err
	}
	if blocked {
		return nil, errors.BadRequest(reason.UserAccessDenied)
	}

	_, exist, err := ms.userCommon.GetUserBasicInfoByID(ctx, req.ToUserID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.BadRequest(reason.UserNotFound)
	}

	conversationID := message.GetConversationID(req.UserID, req.ToUserID)

	msg := &entity.Message{
		FromUserID:     req.UserID,
		ToUserID:       req.ToUserID,
		Title:          req.Title,
		Content:        req.Content,
		Status:         entity.MessageStatusUnread,
		IsFromSystem:   false,
		ConversationID: conversationID,
	}

	err = ms.messageRepo.AddMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	return ms.getMessageResp(ctx, msg), nil
}

func (ms *MessageCommon) GetMessageList(ctx context.Context, req *schema.GetMessageListReq) (*pager.PageModel, error) {
	if req.Page <= 0 {
		req.Page = DefaultPage
	}
	if req.PageSize <= 0 {
		req.PageSize = DefaultPageSize
	}

	messages, total, err := ms.messageRepo.GetMessageList(ctx, req.UserID, req.ConversationID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0)
	for _, msg := range messages {
		userIDs = append(userIDs, msg.FromUserID, msg.ToUserID)
	}

	userMap, err := ms.userCommon.BatchUserBasicInfoByID(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	respList := make([]*schema.GetMessageResp, 0, len(messages))
	for _, msg := range messages {
		resp := ms.getMessageResp(ctx, msg)
		resp.FromUserInfo = userMap[msg.FromUserID]
		resp.ToUserInfo = userMap[msg.ToUserID]
		respList = append(respList, resp)
	}

	return &pager.PageModel{
		Count: total,
		List:  respList,
	}, nil
}

func (ms *MessageCommon) GetConversationList(ctx context.Context, req *schema.GetConversationListReq) (*pager.PageModel, error) {
	if req.Page <= 0 {
		req.Page = DefaultPage
	}
	if req.PageSize <= 0 {
		req.PageSize = DefaultPageSize
	}

	conversations, total, err := ms.messageRepo.GetConversationList(ctx, req.UserID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0)
	for _, conv := range conversations {
		targetUserID := ms.getTargetUserID(req.UserID, conv.FromUserID, conv.ToUserID)
		userIDs = append(userIDs, targetUserID)
	}

	userMap, err := ms.userCommon.BatchUserBasicInfoByID(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	unreadCountMap, err := ms.messageRepo.GetUnreadCountByConversation(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	respList := make([]*schema.GetConversationResp, 0, len(conversations))
	for _, conv := range conversations {
		targetUserID := ms.getTargetUserID(req.UserID, conv.FromUserID, conv.ToUserID)

		unreadCount := int(unreadCountMap[conv.ConversationID])

		resp := &schema.GetConversationResp{
			ConversationID:  conv.ConversationID,
			LastMessage:     conv.Content,
			LastMessageTime: conv.CreatedAt.Unix(),
			UnreadCount:     unreadCount,
			TargetUserID:    targetUserID,
			TargetUserInfo:  userMap[targetUserID],
		}
		respList = append(respList, resp)
	}

	return &pager.PageModel{
		Count: total,
		List:  respList,
	}, nil
}

func (ms *MessageCommon) ReadMessage(ctx context.Context, req *schema.ReadMessageReq) error {
	if len(req.ConversationID) > 0 {
		return ms.messageRepo.MarkConversationAsRead(ctx, req.UserID, req.ConversationID)
	}

	msg, exist, err := ms.messageRepo.GetMessageByID(ctx, req.MessageID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.BadRequest(reason.ObjectNotFound)
	}
	if msg.ToUserID != req.UserID {
		return errors.BadRequest(reason.UserAccessDenied)
	}
	if msg.Status != entity.MessageStatusUnread {
		return nil
	}

	return ms.messageRepo.UpdateMessageStatus(ctx, req.MessageID, entity.MessageStatusRead)
}

func (ms *MessageCommon) ReadAllMessage(ctx context.Context, req *schema.ReadAllMessageReq) error {
	if len(req.ConversationID) > 0 {
		return ms.messageRepo.MarkConversationAsRead(ctx, req.UserID, req.ConversationID)
	}
	return ms.messageRepo.MarkConversationAsRead(ctx, req.UserID, "")
}

func (ms *MessageCommon) DeleteMessage(ctx context.Context, req *schema.DeleteMessageReq) error {
	msg, exist, err := ms.messageRepo.GetMessageByID(ctx, req.MessageID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.BadRequest(reason.ObjectNotFound)
	}
	if msg.FromUserID != req.UserID && msg.ToUserID != req.UserID {
		return errors.BadRequest(reason.UserAccessDenied)
	}

	return ms.messageRepo.DeleteMessage(ctx, req.MessageID)
}

func (ms *MessageCommon) BlockUser(ctx context.Context, req *schema.BlockUserReq) error {
	if req.UserID == req.BlockedUserID {
		return errors.BadRequest(reason.RequestFormatError)
	}

	_, exist, err := ms.userCommon.GetUserBasicInfoByID(ctx, req.BlockedUserID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.BadRequest(reason.UserNotFound)
	}

	blocked, err := ms.messageBlockRepo.IsBlocked(ctx, req.UserID, req.BlockedUserID)
	if err != nil {
		return err
	}
	if blocked {
		return nil
	}

	block := &entity.MessageBlock{
		UserID:        req.UserID,
		BlockedUserID: req.BlockedUserID,
	}
	return ms.messageBlockRepo.BlockUser(ctx, block)
}

func (ms *MessageCommon) UnblockUser(ctx context.Context, req *schema.UnblockUserReq) error {
	if req.UserID == req.BlockedUserID {
		return errors.BadRequest(reason.RequestFormatError)
	}

	return ms.messageBlockRepo.UnblockUser(ctx, req.UserID, req.BlockedUserID)
}

func (ms *MessageCommon) GetBlockedUserList(ctx context.Context, userID string) ([]*schema.GetBlockedUserListResp, error) {
	blocks, err := ms.messageBlockRepo.GetBlockedUserList(ctx, userID)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(blocks))
	for _, block := range blocks {
		userIDs = append(userIDs, block.BlockedUserID)
	}

	userMap, err := ms.userCommon.BatchUserBasicInfoByID(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	respList := make([]*schema.GetBlockedUserListResp, 0, len(blocks))
	for _, block := range blocks {
		resp := &schema.GetBlockedUserListResp{
			BlockedUserID:   block.BlockedUserID,
			BlockedUserInfo: userMap[block.BlockedUserID],
			CreatedAt:       block.CreatedAt.Unix(),
		}
		respList = append(respList, resp)
	}

	return respList, nil
}

func (ms *MessageCommon) GetUnreadMessageCount(ctx context.Context, userID string) (*schema.GetUnreadMessageCountResp, error) {
	total, system, private, err := ms.messageRepo.GetUnreadMessageCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &schema.GetUnreadMessageCountResp{
		TotalUnread:   int(total),
		SystemUnread:  int(system),
		PrivateUnread: int(private),
	}, nil
}

func (ms *MessageCommon) CleanOldMessages(ctx context.Context) error {
	return ms.messageRepo.CleanOldMessages(ctx, OldMessageDays)
}

func (ms *MessageCommon) getMessageResp(ctx context.Context, msg *entity.Message) *schema.GetMessageResp {
	return &schema.GetMessageResp{
		ID:             msg.ID,
		FromUserID:     msg.FromUserID,
		ToUserID:       msg.ToUserID,
		Title:          msg.Title,
		Content:        msg.Content,
		Status:         msg.Status,
		IsRead:         msg.Status == entity.MessageStatusRead,
		IsFromSystem:   msg.IsFromSystem,
		ConversationID: msg.ConversationID,
		CreatedAt:      msg.CreatedAt.Unix(),
	}
}

func (ms *MessageCommon) getTargetUserID(currentUserID, fromUserID, toUserID string) string {
	if currentUserID == fromUserID {
		return toUserID
	}
	return fromUserID
}
