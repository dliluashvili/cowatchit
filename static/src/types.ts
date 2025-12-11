import type videojs from 'video.js'

export interface IHttpResponse<T> {
    data: T
    message: string
    status: number
}

export type Gender = 'm' | 'f'

export interface User {
    id: string
    age: number
    username: string
    gender: Gender
}

export interface SuccessResponse {
    success: boolean
}

export type HttpError = Record<string, Array<string>>

export interface SignInBody {
    username: string
    password: string
}

export interface SignUpBody {
    username: string
    email: string
    gender: string
    date_of_birth: string
    password: string
    password_confirmation: string
}

export interface CreateRoomBody {
    title: string
    capacity: number
    description: string
    src: string
    private: boolean
    password?: string
}

export interface HttpSuccessResponse
    extends IHttpResponse<SuccessResponse | HttpError> {}

export interface HttpGetMeResponse extends IHttpResponse<User | HttpError> {}

export function isAuthSuccess(
    response: HttpSuccessResponse
): response is IHttpResponse<SuccessResponse> {
    return (
        [200, 201].includes(response.status) &&
        'success' in response.data &&
        response.data.success === true
    )
}

export function isAuthFailure(
    response: HttpSuccessResponse
): response is IHttpResponse<HttpError> {
    return (
        ![200, 201].includes(response.status) ||
        !('success' in response.data) ||
        response.data.success === false
    )
}

export function isGetMeSuccess(
    response: HttpGetMeResponse
): response is IHttpResponse<User> {
    return response.status === 200
}

export function isGetMeFailure(
    response: HttpGetMeResponse
): response is IHttpResponse<HttpError> {
    return response.status !== 200
}

export function isResponseSuccess(
    response: HttpSuccessResponse
): response is IHttpResponse<SuccessResponse> {
    return (
        response.status === 200 &&
        'success' in response.data &&
        response.data.success === true
    )
}

export function isResponseFailure(
    response: HttpSuccessResponse
): response is IHttpResponse<HttpError> {
    return (
        response.status !== 200 ||
        !('success' in response.data) ||
        response.data.success === false
    )
}

export type Type = 'EVENT' | 'ERROR'

export type State = 'STOP' | 'PAUSED' | 'PLAYING' | 'END'

export type WSEvent =
    | 'HOST_STATE_SEND'
    | 'HOST_STATE_RECEIVED'
    | 'USER_STATE_SEND'
    | 'USER_STATE_RECEIVED'
    | 'CHAT_MESSAGE_SEND'
    | 'CHAT_MESSAGE_RECEIVED'
    | 'ROOM_MESSAGES_REQUEST'
    | 'ROOM_MESSAGES_ANSWER'
    | 'USER_JOIN_REQUEST'
    | 'USER_JOIN_ANSWER'
    | 'USER_JOINT'
    | 'USER_LEFT'
    | 'IDENTIFY'

export interface WSMessage {
    type: Type
    event?: WSEvent | null
    data?: Record<string, any> | null
}

export interface Participant {
    username: string
    is_host: boolean
}

export interface Participants {
    [key: string]: Participant
}

export interface ChatMessage {
    sender_id: string
    sender_username: string
    content: string
    is_host: boolean
    created_at: string
}
