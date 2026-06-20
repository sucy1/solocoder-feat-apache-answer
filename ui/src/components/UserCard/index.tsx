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

import { memo, FC } from 'react';
import { Link } from 'react-router-dom';

import classnames from 'classnames';

import { Avatar, FormatTime } from '@/components';
import { formatCount } from '@/utils';

interface Props {
  data: any;
  time: number;
  preFix?: string;
  isLogged: boolean;
  timelinePath: string;
  className?: string;
  updateTime?: number;
  updateTimePrefix?: string;
  currentUser?: any;
}

const Index: FC<Props> = ({
  data,
  time,
  preFix,
  isLogged,
  timelinePath,
  className = '',
  updateTime = 0,
  updateTimePrefix = '',
  currentUser,
}) => {
  const isAnonymous = data?.is_anonymous_user;
  const isCurrentUser = currentUser?.username === data?.username;
  const isAdmin = currentUser?.role_id === 2;
  const canSeeRealName = isLogged && (isCurrentUser || isAdmin);
  const displayName = isAnonymous && !canSeeRealName ? '匿名用户' : data?.display_name;
  const showUserLink = data?.status !== 'deleted' && !isAnonymous;

  return (
    <div className={classnames('d-flex', className)}>
      {showUserLink ? (
        <Link to={`/users/${data?.username}`}>
          <Avatar
            avatar={data?.avatar}
            size="40px"
            className="me-2 d-none d-md-block"
            searchStr="s=96"
            alt={displayName}
          />

          <Avatar
            avatar={data?.avatar}
            size="24px"
            className="me-2 d-block d-md-none"
            searchStr="s=48"
            alt={displayName}
          />
        </Link>
      ) : (
        <>
          <Avatar
            avatar={isAnonymous ? '' : data?.avatar}
            size="40px"
            className="me-2 d-none d-md-block"
            searchStr="s=96"
            alt={displayName}
          />

          <Avatar
            avatar={isAnonymous ? '' : data?.avatar}
            size="24px"
            className="me-2 d-block d-md-none"
            searchStr="s=48"
            alt={displayName}
          />
        </>
      )}
      <div className="small text-secondary d-flex flex-column">
        <div className="me-1 me-md-0 d-flex align-items-center">
          {showUserLink ? (
            <Link
              to={`/users/${data?.username}`}
              className="me-1 text-break name-ellipsis"
              style={{ maxWidth: '100px' }}>
              {displayName}
              {isAnonymous && canSeeRealName && (
                <span className="text-muted ms-1">(匿名)</span>
              )}
            </Link>
          ) : (
            <span className="me-1 text-break">
              {displayName}
              {isAnonymous && canSeeRealName && (
                <span className="text-muted ms-1">(匿名)</span>
              )}
            </span>
          )}
          {!isAnonymous && (
            <span className="fw-bold" title="Reputation">
              {formatCount(data?.rank)}
            </span>
          )}
        </div>
        {time &&
          (isLogged ? (
            <Link to={timelinePath}>
              <FormatTime
                time={time}
                preFix={preFix}
                className="link-secondary"
              />
              {updateTime > 0 && (
                <>
                  <span className="mx-1 link-secondary">•</span>
                  <FormatTime
                    time={updateTime}
                    preFix={updateTimePrefix}
                    className="link-secondary"
                  />
                </>
              )}
            </Link>
          ) : (
            <>
              <FormatTime time={time} preFix={preFix} />
              {updateTime > 0 && (
                <>
                  <span className="mx-1 link-secondary">•</span>
                  <FormatTime
                    time={updateTime}
                    preFix={updateTimePrefix}
                    className="link-secondary"
                  />
                </>
              )}
            </>
          ))}
      </div>
    </div>
  );
};

export default memo(Index);
