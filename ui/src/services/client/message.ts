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

import useSWR from 'swr';
import qs from 'qs';

import request from '@/utils/request';
import type * as Type from '@/common/interface';
import { tryLoggedAndActivated } from '@/utils/guard';

export const sendMessage = (receiver_id: string, title: string, content: string) => {
  return request.instance.post('/answer/api/v1/message', {
    receiver_id,
    title,
    content,
  });
};

export const useQueryConversationList = (page: number, page_size: number) => {
  const apiUrl = `/answer/api/v1/message/conversation/page?${qs.stringify(
    { page, page_size },
    { skipNulls: true },
  )}`;

  const { data, error, mutate } = useSWR<Type.ListResult<Type.Conversation>>(
    tryLoggedAndActivated().ok ? apiUrl : null,
    request.instance.get,
  );

  return {
    data,
    isLoading: !data && !error,
    error,
    mutate,
  };
};

export const useQueryMessageList = (
  conversation_id: string,
  page: number,
  page_size: number,
) => {
  const apiUrl = `/answer/api/v1/message/page?${qs.stringify(
    { conversation_id, page, page_size },
    { skipNulls: true },
  )}`;

  const { data, error, mutate } = useSWR<Type.ListResult<Type.Message>>(
    conversation_id && tryLoggedAndActivated().ok ? apiUrl : null,
    request.instance.get,
  );

  return {
    data,
    isLoading: !data && !error,
    error,
    mutate,
  };
};

export const readMessage = (message_id: string) => {
  return request.instance.put('/answer/api/v1/message/read', {
    id: message_id,
  });
};

export const readAllMessage = (conversation_id: string) => {
  return request.instance.put('/answer/api/v1/message/read/all', {
    conversation_id,
  });
};

export const deleteMessage = (message_id: string) => {
  return request.instance.delete('/answer/api/v1/message', {
    id: message_id,
  });
};

export const blockUser = (blocked_user_id: string) => {
  return request.instance.post('/answer/api/v1/message/block', {
    blocked_user_id,
  });
};

export const unblockUser = (blocked_user_id: string) => {
  return request.instance.delete('/answer/api/v1/message/block', {
    blocked_user_id,
  });
};

export const useQueryBlockedUserList = () => {
  const apiUrl = '/answer/api/v1/message/block/list';

  const { data, error, mutate } = useSWR<Type.MessageBlock[]>(
    tryLoggedAndActivated().ok ? apiUrl : null,
    request.instance.get,
  );

  return {
    data,
    isLoading: !data && !error,
    error,
    mutate,
  };
};

export const useQueryUnreadMessageCount = () => {
  const apiUrl = '/answer/api/v1/message/unread/count';

  return useSWR<Type.UnreadMessageCount>(
    tryLoggedAndActivated().ok ? apiUrl : null,
    (url) => request.get(url, { ignoreError: '50X' }),
    {
      refreshInterval: 30000,
    },
  );
};
