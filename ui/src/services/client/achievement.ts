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

export const useGetUserAchievementSummary = (userId: string) => {
  const apiUrl = userId
    ? `/answer/api/v1/achievement/user/summary?user_id=${userId}`
    : null;
  const { data, error, mutate } = useSWR<Type.UserAchievementSummary, Error>(
    apiUrl,
    request.instance.get,
  );
  return {
    data,
    isLoading: !data && !error,
    error,
    mutate,
  };
};

export const useGetMyAchievementSummary = () => {
  const apiUrl = '/answer/api/v1/achievement/me/summary';
  const { data, error, mutate } = useSWR<Type.UserAchievementSummary, Error>(
    apiUrl,
    request.instance.get,
  );
  return {
    data,
    isLoading: !data && !error,
    error,
    mutate,
  };
};

export const useGetUserAchievementList = (
  userId: string,
  page: number,
  pageSize: number,
) => {
  const params = {
    user_id: userId,
    page,
    page_size: pageSize,
  };
  const path = userId
    ? `/answer/api/v1/achievement/user/page?${qs.stringify(params, {
        skipNulls: true,
      })}`
    : null;
  const { data, error, mutate } = useSWR<
    Type.ListResult<Type.Achievement>,
    Error
  >(path, (url) =>
    request.get(url, {
      allow404: true,
    }),
  );

  return {
    data,
    isLoading: !data && !error,
    error,
    mutate,
  };
};
