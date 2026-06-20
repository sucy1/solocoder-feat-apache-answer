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

package entity

import "time"

const (
	MessageStatusUnread = 1
	MessageStatusRead   = 2
	MessageStatusDeleted = 10
)

type Message struct {
	ID              string    `xorm:"not null pk autoincr BIGINT(20) id"`
	CreatedAt       time.Time `xorm:"created TIMESTAMP created_at"`
	UpdatedAt       time.Time `xorm:"updated TIMESTAMP updated_at"`
	FromUserID      string    `xorm:"not null default 0 BIGINT(20) from_user_id"`
	ToUserID        string    `xorm:"not null default 0 BIGINT(20) INDEX to_user_id"`
	Title           string    `xorm:"not null default '' VARCHAR(200) title"`
	Content         string    `xorm:"not null MEDIUMTEXT content"`
	Status          int       `xorm:"not null default 1 INT(11) status"`
	IsFromSystem    bool      `xorm:"not null default false BOOL is_from_system"`
	ConversationID  string    `xorm:"not null default '' VARCHAR(64) INDEX conversation_id"`
}

func (Message) TableName() string {
	return "message"
}

type MessageBlock struct {
	ID          string    `xorm:"not null pk autoincr BIGINT(20) id"`
	CreatedAt   time.Time `xorm:"created TIMESTAMP created_at"`
	UpdatedAt   time.Time `xorm:"updated TIMESTAMP updated_at"`
	UserID      string    `xorm:"not null default 0 BIGINT(20) user_id"`
	BlockedUserID string  `xorm:"not null default 0 BIGINT(20) blocked_user_id"`
}

func (MessageBlock) TableName() string {
	return "message_block"
}
