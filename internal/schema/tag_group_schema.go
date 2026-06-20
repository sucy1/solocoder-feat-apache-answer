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

type AddTagGroupReq struct {
	Name      string `validate:"required,gt=0,lte=50" json:"name"`
	SlugName  string `validate:"required,gt=0,lte=50" json:"slug_name"`
	SortOrder int    `validate:"omitempty,min=0" json:"sort_order"`
	UserID    string `json:"-"`
}

type UpdateTagGroupReq struct {
	ID        string `validate:"required" json:"id"`
	Name      string `validate:"omitempty,gt=0,lte=50" json:"name"`
	SlugName  string `validate:"omitempty,gt=0,lte=50" json:"slug_name"`
	SortOrder int    `validate:"omitempty,min=0" json:"sort_order"`
	UserID    string `json:"-"`
}

type RemoveTagGroupReq struct {
	ID     string `validate:"required" json:"id"`
	UserID string `json:"-"`
}

type GetTagGroupReq struct {
	ID   string `validate:"omitempty" form:"id"`
	Name string `validate:"omitempty,gt=0,lte=50" form:"name"`
}

type GetTagGroupResp struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SlugName  string `json:"slug_name"`
	SortOrder int    `json:"sort_order"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type GetTagGroupWithPageReq struct {
	Page     int    `validate:"omitempty,min=1" form:"page"`
	PageSize int    `validate:"omitempty,min=1" form:"page_size"`
	UserID   string `json:"-"`
}

type TagGroupWithTagsResp struct {
	Group *GetTagGroupResp     `json:"group"`
	Tags  []*GetTagBasicResp   `json:"tags"`
}

type UpdateTagGroupReqForTag struct {
	TagID   string `validate:"required" json:"tag_id"`
	GroupID int64  `validate:"required,min=0" json:"group_id"`
	UserID  string `json:"-"`
}
