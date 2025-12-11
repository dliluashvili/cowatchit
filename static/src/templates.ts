import { escapeHtml, formatTime } from './helpers'
import type { ChatMessage, Participants, Participant } from './types'

export const ParticipantTemplate = (
    userId: string,
    participant: Participant,
    youHost: boolean
) => {
    return `<div id="participant_${userId}" class="flex items-center justify-between p-2 rounded-lg bg-white/5 hover:bg-white/10 transition-colors">
                <span class="text-white text-sm font-medium">${escapeHtml(
                    participant.username
                )}</span>
                <div class="flex items-center gap-2">
                    ${
                        participant.is_host
                            ? `
                        <span class="badge badge-sm bg-yellow-500/20 text-yellow-300">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 mr-0.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <path d="m2 4 3 12h14l3-12-6 7-4-7-4 7-6-7zm3 16h14"></path>
                            </svg>
                        </span>
                    `
                            : youHost
                            ? `
                        <div class="dropdown dropdown-end dropdown-hover">
                            <button class="btn btn-xs btn-ghost text-white/50 hover:text-white hover:bg-white/10" title="Options">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                    <circle cx="12" cy="12" r="1"></circle>
                                    <circle cx="12" cy="5" r="1"></circle>
                                    <circle cx="12" cy="19" r="1"></circle>
                                </svg>
                            </button>
                            <ul class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
                                <li><a onclick="removeParticipant('${userId}')">Remove</a></li>
                                <li><a onclick="muteParticipant('${userId}')">Mute</a></li>
                            </ul>
                        </div>
                    `
                            : ''
                    }
                </div>
            </div>`
}

export const ParticipantsTemplate = (
    participants: Participants,
    youHost: boolean
): string => {
    const participantEntries = Object.entries(participants)

    if (participantEntries.length === 0) {
        return '<div class="text-white/50 text-sm p-4">No participants yet</div>'
    }

    // Sort: host first, then guests
    participantEntries.sort((a, b) => {
        const isHostA = a[1].is_host
        const isHostB = b[1].is_host

        if (isHostA === isHostB) return 0
        return isHostA ? -1 : 1
    })

    const participantHTML = participantEntries
        .map(([userId, participant]) => {
            return ParticipantTemplate(userId, participant, youHost)
        })
        .join('')

    return `<div class="space-y-2">${participantHTML}</div>`
}

export const ChatMessageTemplate = (
    message: ChatMessage,
    currentUserId: string
): string => {
    const isCurrentUser = message.sender_id === currentUserId
    const usernameColor = isCurrentUser ? 'text-purple-300' : 'text-white'

    return `<div class="space-y-1">
        <div class="flex items-center gap-2">
            <div class="flex items-center gap-1">
                <span class="font-medium ${usernameColor}">
                    ${escapeHtml(message.sender_username)}
                </span>
                ${
                    message.is_host
                        ? `<svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-yellow-400 flex-shrink-0" viewBox="0 0 24 24" fill="currentColor">
                        <path d="m2 4 3 12h14l3-12-6 7-4-7-4 7-6-7zm3 16h14"></path>
                    </svg>`
                        : ''
                }
            </div>
            <span class="text-white/50 text-sm">${formatTime(
                message.created_at
            )}</span>
        </div>
        <p class="text-white/80">${escapeHtml(message.content)}</p>
    </div>`
}

export const ChatMessagesTemplate = (
    messages: ChatMessage[],
    currentUserId: string
): string => {
    if (messages.length === 0) {
        return '<div class="text-white/50 text-sm p-4">No messages yet</div>'
    }

    const messagesHTML = messages
        .map((message) => ChatMessageTemplate(message, currentUserId))
        .join('')

    return `<div class="flex flex-col gap-3">${messagesHTML}</div>`
}
