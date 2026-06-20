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

package schema

type SendMessageReq struct {
	ToUserID string `validate:"required" json:"to_user_id"`
	Title    string `validate:"omitempty,lte=200" json:"title"`
	Content  string `validate:"required,gt=0,lte=65536" json:"content"`
	UserID   string `json:"-"`
}

type GetMessageListReq struct {
	Page           int    `validate:"omitempty,min=1" form:"page"`
	PageSize       int    `validate:"omitempty,min=1" form:"page_size"`
	ConversationID string `validate:"omitempty" form:"conversation_id"`
	UserID         string `json:"-"`
}

type GetMessageResp struct {
	ID             string          `json:"id"`
	FromUserID     string          `json:"from_user_id"`
	FromUserInfo   *UserBasicInfo  `json:"from_user_info"`
	ToUserID       string          `json:"to_user_id"`
	ToUserInfo     *UserBasicInfo  `json:"to_user_info"`
	Title          string          `json:"title"`
	Content        string          `json:"content"`
	Status         int             `json:"status"`
	IsFromSystem   bool            `json:"is_from_system"`
	ConversationID string          `json:"conversation_id"`
	CreatedAt      int64           `json:"created_at"`
}

type GetConversationListReq struct {
	Page     int    `validate:"omitempty,min=1" form:"page"`
	PageSize int    `validate:"omitempty,min=1" form:"page_size"`
	UserID   string `json:"-"`
}

type GetConversationResp struct {
	ConversationID   string         `json:"conversation_id"`
	LastMessage      string         `json:"last_message"`
	LastMessageTime  int64          `json:"last_message_time"`
	UnreadCount      int            `json:"unread_count"`
	TargetUserID     string         `json:"target_user_id"`
	TargetUserInfo   *UserBasicInfo `json:"target_user_info"`
}

type ReadMessageReq struct {
	MessageID      string `validate:"required" json:"message_id"`
	ConversationID string `validate:"omitempty" json:"conversation_id"`
	UserID         string `json:"-"`
}

type ReadAllMessageReq struct {
	ConversationID string `validate:"omitempty" json:"conversation_id"`
	UserID         string `json:"-"`
}

type DeleteMessageReq struct {
	MessageID string `validate:"required" json:"message_id"`
	UserID    string `json:"-"`
}

type BlockUserReq struct {
	BlockedUserID string `validate:"required" json:"blocked_user_id"`
	UserID        string `json:"-"`
}

type UnblockUserReq struct {
	BlockedUserID string `validate:"required" json:"blocked_user_id"`
	UserID        string `json:"-"`
}

type GetBlockedUserListResp struct {
	BlockedUserID   string         `json:"blocked_user_id"`
	BlockedUserInfo *UserBasicInfo `json:"blocked_user_info"`
	CreatedAt       int64          `json:"created_at"`
}

type GetUnreadMessageCountResp struct {
	TotalUnread   int `json:"total_unread"`
	SystemUnread  int `json:"system_unread"`
	PrivateUnread int `json:"private_unread"`
}
