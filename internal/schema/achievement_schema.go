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

type GetUserAchievementReq struct {
	Username string `validate:"required,gt=0,lte=100" form:"username"`
	Page     int    `validate:"omitempty,min=1" form:"page"`
	PageSize int    `validate:"omitempty,min=1" form:"page_size"`
	UserID   string `json:"-"`
}

type GetUserAchievementResp struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	BadgeInfo   *GetUserBadgeAwardListResp `json:"badge_info,omitempty"`
	Source      string `json:"source"`
	Reputation  int    `json:"reputation"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
}

type GetUserAchievementSummaryResp struct {
	TotalReputation  int                              `json:"total_reputation"`
	BadgeCount       int                              `json:"badge_count"`
	Badges           []*GetUserBadgeAwardListResp     `json:"badges"`
	RecentActivities []*GetUserAchievementResp        `json:"recent_activities"`
}

type AchievementBadgeEvent struct {
	BadgeID string `json:"badge_id"`
	UserID  string `json:"user_id"`
}

type AchievementReputationEvent struct {
	UserID      string `json:"user_id"`
	Source      string `json:"source"`
	Reputation  int    `json:"reputation"`
	Description string `json:"description"`
}

type CheckConsecutiveLoginReq struct {
	UserID string `json:"-"`
}
