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

package controller

import (
	"github.com/apache/answer/internal/base/handler"
	"github.com/apache/answer/internal/base/middleware"
	"github.com/apache/answer/internal/base/pager"
	"github.com/apache/answer/internal/base/reason"
	"github.com/apache/answer/internal/schema"
	messagecommon "github.com/apache/answer/internal/service/message_common"
	"github.com/gin-gonic/gin"
	"github.com/segmentfault/pacman/errors"
)

// MessageController message controller
type MessageController struct {
	messageCommonService *messagecommon.MessageCommon
}

// NewMessageController new message controller
func NewMessageController(
	messageCommonService *messagecommon.MessageCommon,
) *MessageController {
	return &MessageController{messageCommonService: messageCommonService}
}

// SendMessage send message
// @Summary send message
// @Description send message to another user
// @Tags Message
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.SendMessageReq true "message"
// @Success 200 {object} handler.RespBody{data=schema.GetMessageResp}
// @Router /answer/api/v1/message [post]
func (mc *MessageController) SendMessage(ctx *gin.Context) {
	req := &schema.SendMessageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	if len(req.UserID) == 0 {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	resp, err := mc.messageCommonService.SendMessage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// GetMessageList get message list
// @Summary get message list
// @Description get message list by conversation
// @Tags Message
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "page"
// @Param page_size query int false "page size"
// @Param conversation_id query string false "conversation id"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.GetMessageResp}}
// @Router /answer/api/v1/message/list [get]
func (mc *MessageController) GetMessageList(ctx *gin.Context) {
	req := &schema.GetMessageListReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	if len(req.UserID) == 0 {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	resp, err := mc.messageCommonService.GetMessageList(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if pager.ValPageOutOfRange(resp.Count, req.Page, req.PageSize) {
		handler.HandleResponse(ctx, errors.NotFound(reason.RequestFormatError), nil)
		return
	}
	handler.HandleResponse(ctx, err, resp)
}

// GetConversationList get conversation list
// @Summary get conversation list
// @Description get conversation list for current user
// @Tags Message
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "page"
// @Param page_size query int false "page size"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.GetConversationResp}}
// @Router /answer/api/v1/message/conversation/list [get]
func (mc *MessageController) GetConversationList(ctx *gin.Context) {
	req := &schema.GetConversationListReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	if len(req.UserID) == 0 {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	resp, err := mc.messageCommonService.GetConversationList(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if pager.ValPageOutOfRange(resp.Count, req.Page, req.PageSize) {
		handler.HandleResponse(ctx, errors.NotFound(reason.RequestFormatError), nil)
		return
	}
	handler.HandleResponse(ctx, err, resp)
}

// ReadMessage read message
// @Summary read message
// @Description mark message as read
// @Tags Message
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.ReadMessageReq true "read message"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/message/read [post]
func (mc *MessageController) ReadMessage(ctx *gin.Context) {
	req := &schema.ReadMessageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	if len(req.UserID) == 0 {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	err := mc.messageCommonService.ReadMessage(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// ReadAllMessage read all message
// @Summary read all message
// @Description mark all messages as read
// @Tags Message
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.ReadAllMessageReq true "read all message"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/message/read-all [post]
func (mc *MessageController) ReadAllMessage(ctx *gin.Context) {
	req := &schema.ReadAllMessageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	if len(req.UserID) == 0 {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	err := mc.messageCommonService.ReadAllMessage(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// DeleteMessage delete message
// @Summary delete message
// @Description delete message
// @Tags Message
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.DeleteMessageReq true "delete message"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/message [delete]
func (mc *MessageController) DeleteMessage(ctx *gin.Context) {
	req := &schema.DeleteMessageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	if len(req.UserID) == 0 {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	err := mc.messageCommonService.DeleteMessage(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// BlockUser block user
// @Summary block user
// @Description block user from sending messages
// @Tags Message
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.BlockUserReq true "block user"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/message/block [post]
func (mc *MessageController) BlockUser(ctx *gin.Context) {
	req := &schema.BlockUserReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	if len(req.UserID) == 0 {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	err := mc.messageCommonService.BlockUser(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// UnblockUser unblock user
// @Summary unblock user
// @Description unblock user
// @Tags Message
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.UnblockUserReq true "unblock user"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/message/unblock [post]
func (mc *MessageController) UnblockUser(ctx *gin.Context) {
	req := &schema.UnblockUserReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	if len(req.UserID) == 0 {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	err := mc.messageCommonService.UnblockUser(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// GetBlockedUserList get blocked user list
// @Summary get blocked user list
// @Description get blocked user list
// @Tags Message
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} handler.RespBody{data=[]schema.GetBlockedUserListResp}
// @Router /answer/api/v1/message/blocked [get]
func (mc *MessageController) GetBlockedUserList(ctx *gin.Context) {
	userID := middleware.GetLoginUserIDFromContext(ctx)
	if len(userID) == 0 {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	resp, err := mc.messageCommonService.GetBlockedUserList(ctx, userID)
	handler.HandleResponse(ctx, err, resp)
}

// GetUnreadMessageCount get unread message count
// @Summary get unread message count
// @Description get unread message count
// @Tags Message
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} handler.RespBody{data=schema.GetUnreadMessageCountResp}
// @Router /answer/api/v1/message/unread-count [get]
func (mc *MessageController) GetUnreadMessageCount(ctx *gin.Context) {
	userID := middleware.GetLoginUserIDFromContext(ctx)
	if len(userID) == 0 {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	resp, err := mc.messageCommonService.GetUnreadMessageCount(ctx, userID)
	handler.HandleResponse(ctx, err, resp)
}
