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

import { FC, useState, useEffect, useRef } from 'react';
import {
  Row,
  Col,
  Badge,
  Button,
  Form,
  Modal,
  Dropdown,
} from 'react-bootstrap';
import { useTranslation } from 'react-i18next';
import { useParams, useNavigate } from 'react-router-dom';

import classNames from 'classnames';

import { usePageTags } from '@/hooks';
import {
  Avatar,
  FormatTime,
  Empty,
  Pagination,
  Icon,
  Modal as ConfirmModal,
} from '@/components';
import { loggedUserInfoStore } from '@/stores';
import {
  useQueryConversationList,
  useQueryMessageList,
  useQueryBlockedUserList,
  sendMessage,
  readAllMessage,
  blockUser,
  unblockUser,
} from '@/services';
import type * as Type from '@/common/interface';
import { floppyNavigation } from '@/utils';

import './index.scss';

const PAGE_SIZE = 20;
const CONVERSATION_PAGE_SIZE = 30;

const Messages: FC = () => {
  const { conversationId } = useParams();
  const navigate = useNavigate();
  const { t } = useTranslation('translation', { keyPrefix: 'messages' });
  const sessionUser = loggedUserInfoStore((state) => state.user);

  const [conversationPage, setConversationPage] = useState(1);
  const [messagePage, setMessagePage] = useState(1);
  const [selectedConversation, setSelectedConversation] =
    useState<Type.Conversation | null>(null);
  const [messageTitle, setMessageTitle] = useState('');
  const [messageContent, setMessageContent] = useState('');
  const [showBlockedModal, setShowBlockedModal] = useState(false);
  const [messages, setMessages] = useState<Type.Message[]>([]);
  const messageListRef = useRef<HTMLDivElement>(null);

  const { data: conversationData, mutate: mutateConversations } =
    useQueryConversationList(conversationPage, CONVERSATION_PAGE_SIZE);

  const { data: messageData, mutate: mutateMessages, isLoading: messageLoading } =
    useQueryMessageList(
      selectedConversation?.conversation_id || conversationId || '',
      messagePage,
      PAGE_SIZE,
    );

  const { data: blockedUsers, mutate: mutateBlockedUsers } =
    useQueryBlockedUserList();

  usePageTags({
    title: t('page_title'),
  });

  useEffect(() => {
    if (conversationId && conversationData?.list) {
      const conv = conversationData.list.find(
        (c) => c.conversation_id === conversationId,
      );
      if (conv) {
        setSelectedConversation(conv);
      }
    }
  }, [conversationId, conversationData]);

  useEffect(() => {
    if (!messageData) {
      return;
    }
    if (messagePage > 1) {
      setMessages([...(messageData?.list || []), ...messages]);
    } else {
      setMessages(messageData?.list || []);
    }
  }, [messageData, messagePage]);

  useEffect(() => {
    if (messageListRef.current) {
      messageListRef.current.scrollTop = messageListRef.current.scrollHeight;
    }
  }, [messages]);

  const handleConversationClick = async (conv: Type.Conversation) => {
    setSelectedConversation(conv);
    setMessagePage(1);
    setMessages([]);
    navigate(`/messages/${conv.conversation_id}`);
    if (conv.unread_count > 0) {
      await readAllMessage(conv.conversation_id);
      mutateConversations();
    }
  };

  const handleSendMessage = async (e) => {
    e.preventDefault();
    if (!selectedConversation || !messageContent.trim()) {
      return;
    }

    try {
      await sendMessage(
        selectedConversation.other_user.id,
        messageTitle,
        messageContent,
      );
      setMessageTitle('');
      setMessageContent('');
      mutateMessages();
      mutateConversations();
    } catch (error) {
      console.error('Failed to send message:', error);
    }
  };

  const handleLoadMoreMessages = () => {
    setMessagePage(messagePage + 1);
  };

  const handleBlockUser = async () => {
    if (!selectedConversation) return;

    ConfirmModal.confirm({
      content: t('block_confirm', {
        username: selectedConversation.other_user.display_name,
      }),
      onConfirm: async () => {
        await blockUser(selectedConversation.other_user.id);
        mutateBlockedUsers();
        mutateConversations();
      },
    });
  };

  const handleUnblockUser = async (blockedUserId: string) => {
    await unblockUser(blockedUserId);
    mutateBlockedUsers();
    mutateConversations();
  };

  const isBlocked = (userId: string) => {
    return blockedUsers?.some((b) => b.blocked_user_id === userId) || false;
  };

  const isCurrentUserBlocked = selectedConversation
    ? isBlocked(selectedConversation.other_user.id)
    : false;

  return (
    <Row className="pt-4 mb-5">
      <Col className="page-main flex-auto">
        <h3 className="mb-4">{t('title')}</h3>

        <div className="messages-container">
          <div className="conversation-list">
            <div className="conversation-list-header d-flex justify-content-between align-items-center">
              <span>{t('conversations')}</span>
              <Dropdown align="end">
                <Dropdown.Toggle
                  as={Button}
                  variant="link"
                  className="p-0 btn-no-border icon-link">
                  <Icon name="three-dots-vertical" className="fs-5" />
                </Dropdown.Toggle>
                <Dropdown.Menu>
                  <Dropdown.Item onClick={() => setShowBlockedModal(true)}>
                    {t('blocked_users')}
                  </Dropdown.Item>
                </Dropdown.Menu>
              </Dropdown>
            </div>
            <div className="conversation-items">
              {conversationData?.list?.length === 0 && <Empty />}
              {conversationData?.list?.map((conv) => {
                if (isBlocked(conv.other_user.id)) {
                  return null;
                }
                return (
                  <div
                    key={conv.conversation_id}
                    className={classNames('conversation-item', {
                      active:
                        selectedConversation?.conversation_id ===
                        conv.conversation_id,
                    })}
                    onClick={() => handleConversationClick(conv)}>
                    <div className="avatar-wrap">
                      <Avatar
                        size="40px"
                        avatar={conv.other_user.avatar}
                        alt={conv.other_user.display_name}
                        searchStr="s=96"
                      />
                    </div>
                    <div className="conversation-content">
                      <div className="conversation-header">
                        <span className="username">
                          {conv.other_user.display_name}
                        </span>
                        <span className="time">
                          <FormatTime
                            time={conv.last_message.created_at}
                            autoUpdate
                          />
                        </span>
                      </div>
                      <div className="last-message">
                        {conv.last_message.content}
                      </div>
                    </div>
                    {conv.unread_count > 0 && (
                      <Badge bg="danger" className="unread-badge">
                        {conv.unread_count}
                      </Badge>
                    )}
                  </div>
                );
              })}
              {(conversationData?.count || 0) >
                CONVERSATION_PAGE_SIZE * conversationPage && (
                <div className="d-flex justify-content-center py-3">
                  <Button
                    variant="link"
                    className="btn-no-border"
                    onClick={() =>
                      setConversationPage(conversationPage + 1)
                    }>
                    {t('load_more')}
                  </Button>
                </div>
              )}
            </div>
          </div>

          <div className="message-panel">
            {selectedConversation ? (
              <>
                <div className="message-header">
                  <div className="user-info">
                    <div className="avatar-wrap">
                      <Avatar
                        size="40px"
                        avatar={selectedConversation.other_user.avatar}
                        alt={selectedConversation.other_user.display_name}
                        searchStr="s=96"
                      />
                    </div>
                    <div className="user-detail">
                      <div className="username">
                        {selectedConversation.other_user.display_name}
                      </div>
                      <small className="text-secondary">
                        @{selectedConversation.other_user.username}
                      </small>
                    </div>
                  </div>
                  <div className="d-flex gap-2">
                    {isCurrentUserBlocked ? (
                      <Button
                        variant="outline-secondary"
                        size="sm"
                        onClick={() =>
                          handleUnblockUser(selectedConversation.other_user.id)
                        }>
                        {t('unblock')}
                      </Button>
                    ) : (
                      <Button
                        variant="outline-danger"
                        size="sm"
                        onClick={handleBlockUser}>
                        {t('block')}
                      </Button>
                    )}
                  </div>
                </div>

                <div className="message-list" ref={messageListRef}>
                  {messageLoading && messages.length === 0 && (
                    <div className="empty-state">{t('loading')}</div>
                  )}
                  {!messageLoading && messages.length === 0 && (
                    <div className="empty-state">{t('no_messages')}</div>
                  )}
                  {(messageData?.count || 0) > PAGE_SIZE * messagePage && (
                    <div className="d-flex justify-content-center py-2">
                      <Button
                        variant="link"
                        className="btn-no-border"
                        onClick={handleLoadMoreMessages}>
                        {t('load_history')}
                      </Button>
                    </div>
                  )}
                  {messages.map((msg) => {
                    const isOwn = msg.from_user_id === sessionUser?.id;
                    return (
                      <div
                        key={msg.id}
                        className={classNames('message-item', {
                          own: isOwn,
                          other: !isOwn,
                          unread: !isOwn && !msg.is_read,
                        })}>
                        <div className="message-avatar">
                          <Avatar
                            size="32px"
                            avatar={
                              isOwn
                                ? sessionUser?.avatar
                                : selectedConversation.other_user.avatar
                            }
                            alt={
                              isOwn
                                ? sessionUser?.display_name
                                : selectedConversation.other_user.display_name
                            }
                            searchStr="s=64"
                          />
                        </div>
                        <div className="message-content">
                          <div className="message-bubble">
                            {msg.title && (
                              <div className="message-title">{msg.title}</div>
                            )}
                            <div className="message-text">{msg.content}</div>
                          </div>
                          <div
                            className={classNames('message-time', {
                              'text-end': isOwn,
                            })}>
                            {!isOwn && (
                              <span
                                className={classNames('me-2', {
                                  'text-primary': msg.is_read,
                                  'text-secondary': !msg.is_read,
                                })}>
                                {msg.is_read ? (
                                  <Icon name="check2-all" size="14px" />
                                ) : (
                                  <Icon name="check2" size="14px" />
                                )}
                              </span>
                            )}
                            <FormatTime time={msg.created_at} autoUpdate />
                          </div>
                        </div>
                      </div>
                    );
                  })}
                </div>

                {!isCurrentUserBlocked ? (
                  <Form className="message-input" onSubmit={handleSendMessage}>
                    <Form.Control
                      type="text"
                      className="input-title mb-2"
                      placeholder={t('title_placeholder')}
                      value={messageTitle}
                      onChange={(e) => setMessageTitle(e.target.value)}
                    />
                    <Form.Control
                      as="textarea"
                      rows={3}
                      placeholder={t('content_placeholder')}
                      value={messageContent}
                      onChange={(e) => setMessageContent(e.target.value)}
                    />
                    <div className="input-actions">
                      <Button
                        type="submit"
                        variant="primary"
                        disabled={!messageContent.trim()}>
                        {t('send')}
                      </Button>
                    </div>
                  </Form>
                ) : (
                  <div className="message-input text-center text-secondary py-4">
                    {t('blocked_cannot_send')}
                  </div>
                )}
              </>
            ) : (
              <div className="empty-state">{t('select_conversation')}</div>
            )}
          </div>
        </div>
      </Col>
      <Col className="page-right-side" />

      <Modal
        show={showBlockedModal}
        onHide={() => setShowBlockedModal(false)}
        centered>
        <Modal.Header closeButton>
          <Modal.Title>{t('blocked_users')}</Modal.Title>
        </Modal.Header>
        <Modal.Body className="blocked-list-modal">
          {blockedUsers?.length === 0 && (
            <div className="text-center text-secondary py-4">
              {t('no_blocked_users')}
            </div>
          )}
          {blockedUsers?.map((blocked) => (
            <div key={blocked.blocked_user_id} className="blocked-user-item">
              <div className="user-info">
                <div className="avatar-wrap">
                  <Avatar
                    size="36px"
                    avatar={blocked.blocked_user_info.avatar}
                    alt={blocked.blocked_user_info.display_name}
                    searchStr="s=64"
                  />
                </div>
                <div>
                  <div className="fw-medium">
                    {blocked.blocked_user_info.display_name}
                  </div>
                  <small className="text-secondary">
                    @{blocked.blocked_user_info.username}
                  </small>
                </div>
              </div>
              <Button
                variant="outline-secondary"
                size="sm"
                onClick={() => handleUnblockUser(blocked.blocked_user_id)}>
                {t('unblock')}
              </Button>
            </div>
          ))}
        </Modal.Body>
        <Modal.Footer>
          <Button
            variant="secondary"
            onClick={() => setShowBlockedModal(false)}>
            {t('close')}
          </Button>
        </Modal.Footer>
      </Modal>
    </Row>
  );
};

export default Messages;
