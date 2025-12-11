import {
    type CreateRoomBody,
    type HttpGetMeResponse,
    type HttpSuccessResponse,
    type SignInBody,
    type SignUpBody,
} from './types'

export const signIn = async (
    url: string,
    body: SignInBody
): Promise<HttpSuccessResponse> => {
    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
    })

    const data = await response.json()

    return data
}

export const signUp = async (
    url: string,
    body: SignUpBody
): Promise<HttpSuccessResponse> => {
    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
    })

    const data = await response.json()

    return data
}

export const getMe = async (): Promise<HttpGetMeResponse> => {
    const url = '/user/me'

    const response = await fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    })

    const data = await response.json()

    return data
}

export const createRoom = async (
    body: CreateRoomBody
): Promise<HttpSuccessResponse> => {
    const url = '/create-room'

    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
    })

    const data = await response.json()

    return data
}
