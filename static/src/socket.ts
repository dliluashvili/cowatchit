import type {
    ChatMessage,
    WSEvent,
    Participant,
    Participants,
    WSMessage,
    State,
} from './types'
import {
    ChatMessagesTemplate,
    ChatMessageTemplate,
    ParticipantsTemplate,
    ParticipantTemplate,
} from './templates'
import type Player from 'video.js/dist/types/player'

// Define custom event interfaces for HTMX WebSocket events
interface HTMXWebSocketOpenEvent extends Event {
    detail: {
        socketWrapper: WebSocket
    }
}

interface HTMXWebSocketCloseEvent extends Event {
    detail: {
        socket: WebSocket
        reason: string
        code: number
    }
}

interface HTMXWebSocketMessageEvent extends Event {
    detail: {
        message: string
        socket: WebSocket
    }
}

interface HTMXWebSocketSendEvent extends Event {
    detail: {
        message: string
        socket: WebSocket
    }
}

interface HTMXWebSocketErrorEvent extends Event {
    detail: {
        error: Error
        socket: WebSocket
    }
}

document.addEventListener('htmx:load', async function () {
    let socketId: null | string = null
    let authUserId: null | string = null
    let authUsername: null | string = null
    let isRoomHost: null | boolean = null
    let socketWrapper: null | WebSocket = null
    let player: null | Player = null

    const messagesDiv = document.querySelector('#messages') as HTMLDivElement
    const participantsDiv = document.querySelector(
        '#participants-list'
    ) as HTMLDivElement

    const roomId = window?.location?.pathname?.match(/\/rooms\/([^\/]+)/)?.[1]

    if (window.htmx) {
        window.htmx.config.wsReconnectDelay = function (retryCount: number) {
            const baseDelay = 1000
            const maxDelay = 30000
            return Math.min(baseDelay * retryCount, maxDelay)
        }
    }

    // WebSocket Open Event
    document.body.addEventListener(
        'htmx:wsOpen',
        (event: HTMXWebSocketOpenEvent) => {
            if (roomId) {
                socketWrapper = event.detail.socketWrapper

                setTimeout(function () {
                    const wsMessage: WSMessage = {
                        type: 'EVENT',
                        event: 'USER_JOIN_REQUEST',
                        data: {
                            socket_id: socketId,
                            room_id: roomId,
                        },
                    }

                    socketWrapper.send(JSON.stringify(wsMessage))
                }, 200)

                setTimeout(function () {
                    requestChatMessages()
                }, 300)

                setTimeout(function () {
                    registerChatListeners()
                }, 400)
            }
        }
    )

    // WebSocket Close Event
    document.body.addEventListener(
        'htmx:wsClose',
        (event: HTMXWebSocketCloseEvent) => {
            console.log('WebSocket connection closed', {
                reason: event.detail.reason,
                code: event.detail.code,
            })
        }
    )

    // Before Send Event
    document.body.addEventListener(
        'htmx:wsBeforeSend',
        (event: HTMXWebSocketSendEvent) => {
            console.log('About to send WebSocket message', event.detail.message)
        }
    )

    document.body.addEventListener(
        'htmx:wsAfterMessage',
        async (event: HTMXWebSocketSendEvent) => {
            const { message } = event.detail

            try {
                const msg = JSON.parse(message) as WSMessage

                console.log('msgg', msg)

                if (msg.type === 'EVENT') {
                    switch (msg.event) {
                        case 'IDENTIFY':
                            socketId = msg.data.socket_id
                            authUserId = msg.data.auth_id
                            authUsername = msg.data.auth_username
                            break
                        case 'USER_JOIN_ANSWER':
                            isRoomHost = msg.data.is_host
                            document.querySelector('.room-title').textContent =
                                msg.data.title

                            setTimeout(() => {
                                if (msg.data.is_host) {
                                    document
                                        .querySelector('.room-host')
                                        .classList.remove('hidden')
                                    document
                                        .querySelector('.room-guest')
                                        .classList.add('hidden')
                                } else {
                                    document
                                        .querySelector('.room-guest')
                                        .classList.remove('hidden')
                                    document
                                        .querySelector('.room-host')
                                        .classList.add('hidden')
                                }

                                const participants = msg.data.participants

                                participantsDiv.innerHTML =
                                    ParticipantsTemplate(
                                        participants as Participants,
                                        isRoomHost
                                    )

                                document
                                    .querySelector('#video-el')
                                    .addEventListener(
                                        'contextmenu',
                                        function (e) {
                                            e.preventDefault()
                                        }
                                    )

                                document.querySelector(
                                    '.participants'
                                ).textContent =
                                    Object.keys(participants).length.toString()

                                document.querySelector(
                                    '.room-host-username'
                                ).textContent = msg.data.host

                                document
                                    .querySelector('#room-view')
                                    .classList.remove('hidden')
                                document
                                    .querySelector('#join-view')
                                    .classList.add('hidden')

                                player = window.videojs('video-el')

                                player.controls(true)

                                if (!isRoomHost) {
                                    ;(
                                        player as any
                                    ).controlBar.progressControl.hide()
                                }

                                // Set source FIRST
                                player.src({
                                    src: msg.data.src,
                                    type: 'video/mp4',
                                })

                                player.on('play', function () {
                                    if (!player.seeking()) {
                                        let event: WSEvent = 'USER_STATE_SEND'

                                        if (isRoomHost) {
                                            event = 'HOST_STATE_SEND'
                                        }

                                        const state: State = 'PLAYING'

                                        sendStateEvent(event, state)
                                    }
                                })

                                player.on('pause', function () {
                                    if (!player.seeking()) {
                                        let event: WSEvent = 'USER_STATE_SEND'

                                        if (isRoomHost) {
                                            event = 'HOST_STATE_SEND'
                                        }

                                        const state: State = 'PAUSED'

                                        sendStateEvent(event, state)
                                    }
                                })

                                player.on('seeked', function () {
                                    let event: WSEvent = 'USER_STATE_SEND'

                                    if (isRoomHost) {
                                        event = 'HOST_STATE_SEND'
                                    }

                                    const state: State = player.paused()
                                        ? 'PAUSED'
                                        : 'PLAYING'

                                    sendStateEvent(event, state)
                                })
                            }, 300)

                            break
                        case 'CHAT_MESSAGE_RECEIVED':
                            const message: ChatMessage = {
                                sender_id: msg.data.sender_id,
                                sender_username: msg.data.sender_username,
                                is_host: msg.data.is_host,
                                content: msg.data.content,
                                created_at: msg.data.created_at,
                            }

                            setTimeout(function () {
                                const messageHtml = ChatMessageTemplate(
                                    message,
                                    authUserId
                                )

                                messagesDiv.insertAdjacentHTML(
                                    'afterbegin',
                                    messageHtml
                                )

                                messagesDiv.scrollTop = messagesDiv.scrollHeight
                            }, 100)
                            break
                        case 'ROOM_MESSAGES_ANSWER':
                            messagesDiv.innerHTML = ChatMessagesTemplate(
                                msg.data.messages,
                                authUserId
                            )

                            break
                        case 'USER_JOINT':
                            const participant: Participant = {
                                is_host: msg.data.is_host,
                                username: msg.data.username,
                            }

                            const participantHtml = ParticipantTemplate(
                                msg.data.user_id,
                                participant,
                                isRoomHost
                            )

                            let position: 'afterbegin' | 'beforeend' =
                                'beforeend'

                            if (msg.data.is_host) {
                                position = 'afterbegin'
                            }

                            setTimeout(function () {
                                document.querySelector(
                                    '.participants'
                                ).textContent = msg.data.counted_participants

                                participantsDiv.insertAdjacentHTML(
                                    position,
                                    participantHtml
                                )
                            }, 250)
                            break

                        case 'HOST_STATE_RECEIVED':
                            if (msg.data.state === 'PLAYING') {
                                player.play()
                            }

                            if (msg.data.state === 'PAUSED') {
                                player.pause()
                            }

                            player.currentTime(msg.data.current_time_seconds)

                            break
                        case 'USER_LEFT':
                            setTimeout(() => {
                                document.querySelector(
                                    '.participants'
                                ).textContent = msg.data.counted_participants

                                document
                                    .querySelector(
                                        `#participant_${msg.data.user_id}`
                                    )
                                    ?.remove()
                            }, 250)
                            break
                        default:
                            console.log('No such event exists!')
                            break
                    }
                }
            } catch (error) {
                console.error('Error processing WebSocket message:', error)
            }
        }
    )

    function sendStateEvent(event: WSEvent, state: State) {
        const wsMessage: WSMessage = {
            type: 'EVENT',
            event,
            data: {
                state,
                socket_id: socketId,
                room_id: roomId,
                current_time_seconds: player.currentTime(),
            },
        }

        socketWrapper.send(JSON.stringify(wsMessage))
    }

    function requestChatMessages() {
        const wsMessage: WSMessage = {
            type: 'EVENT',
            event: 'ROOM_MESSAGES_REQUEST',
            data: {
                socket_id: socketId,
                room_id: roomId,
            },
        }

        socketWrapper.send(JSON.stringify(wsMessage))
    }

    function registerChatListeners() {
        const chatForm = document.querySelector('#chat-form') as HTMLFormElement

        const contentInput = document.querySelector(
            '#message-content'
        ) as HTMLInputElement

        chatForm.addEventListener('submit', function (e) {
            e.preventDefault()
            const content = contentInput.value.trim()

            const wsMessage: WSMessage = {
                type: 'EVENT',
                event: 'CHAT_MESSAGE_SEND',
                data: {
                    socket_id: socketId,
                    room_id: roomId,
                    content: content,
                },
            }

            contentInput.value = ''

            setTimeout(function () {
                const message: ChatMessage = {
                    sender_id: authUserId,
                    sender_username: authUsername,
                    is_host: isRoomHost,
                    content,
                    created_at: new Date().toISOString(),
                }

                const messageHtml = ChatMessageTemplate(message, authUserId)

                messagesDiv.insertAdjacentHTML('afterbegin', messageHtml)

                messagesDiv.scrollTop = messagesDiv.scrollHeight

                socketWrapper.send(JSON.stringify(wsMessage))
            }, 100)
        })
    }

    // WebSocket Error Event
    document.body.addEventListener(
        'htmx:wsError',
        (event: HTMXWebSocketErrorEvent) => {
            console.error('WebSocket error', event.detail.error)
        }
    )
})
