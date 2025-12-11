export const getValueByInputName = (
    form: HTMLFormElement,
    name: string
): string => {
    const input = form.querySelector(
        `input[name="${name}"]`
    ) as HTMLInputElement

    if (!input) {
        throw new Error(`Input with name "${name}" not found`)
    }

    return input.value
}

export const getValueByCheckboxName = (
    form: HTMLFormElement,
    name: string
): boolean => {
    const input = form.querySelector(
        `input[name="${name}"]`
    ) as HTMLInputElement

    if (!input) {
        throw new Error(`Input with name "${name}" not found`)
    }

    return input.checked
}

export const getValueByTextareaName = (
    form: HTMLFormElement,
    name: string
): string => {
    const textarea = form.querySelector(
        `textarea[name="${name}"]`
    ) as HTMLTextAreaElement

    if (!textarea) {
        throw new Error(`Textarea with name "${name}" not found`)
    }

    return textarea.value
}

export function escapeHtml(text: string): string {
    const map: { [key: string]: string } = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;',
    }
    return text.replace(/[&<>"']/g, (char) => map[char])
}

export function formatTime(timestamp: string): string {
    const date = new Date(timestamp)
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}