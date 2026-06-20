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
	AchievementTypeBadge = "badge"
	AchievementTypeReputation = "reputation"
)

const (
	AchievementSourceRegister      = "register"
	AchievementSourceFirstQuestion = "first_question"
	AchievementSourceAnswerAccepted = "answer_accepted"
	AchievementSourceConsecutiveLogin = "consecutive_login"
	AchievementSourceFirstAnswer   = "first_answer"
	AchievementSourceQuestionUpvote = "question_upvote"
	AchievementSourceAnswerUpvote  = "answer_upvote"
)

const MaxUserBadges = 50

type UserAchievement struct {
	ID          string    `xorm:"not null pk autoincr BIGINT(20) id"`
	CreatedAt   time.Time `xorm:"created TIMESTAMP created_at"`
	UpdatedAt   time.Time `xorm:"updated TIMESTAMP updated_at"`
	UserID      string    `xorm:"not null default 0 BIGINT(20) INDEX user_id"`
	AchievementType string `xorm:"not null default '' VARCHAR(32) achievement_type"`
	AchievementID string  `xorm:"not null default 0 BIGINT(20) achievement_id"`
	Source      string    `xorm:"not null default '' VARCHAR(64) source"`
	Reputation  int       `xorm:"not null default 0 INT(11) reputation"`
	Description string    `xorm:"not null default '' VARCHAR(255) description"`
}

func (UserAchievement) TableName() string {
	return "user_achievement"
}

type UserReputationSummary struct {
	UserID         string `xorm:"user_id"`
	TotalReputation int   `xorm:"total_reputation"`
}

func (UserReputationSummary) TableName() string {
	return "user_achievement"
}
