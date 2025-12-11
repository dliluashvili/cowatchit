export function isEmail(email: string) {
    if (!minLength(email, 3)) return false

    if (!maxLength(email, 254)) return false

    return !!email.match(
        /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
    )
}

export function isGender(gender: string) {
    return ['f', 'm'].includes(gender)
}

export function minLength(value: string, min: number) {
    return value.length >= min
}

export function maxLength(value: string, max: number) {
    return value.length <= max
}

export function wrongDob(value: string) {
    if (!value.length) {
        return 'date of birth is required'
    }

    const date = new Date(value)
    const currentYear = new Date().getFullYear()
    const minYear = currentYear - 100
    const maxYear = 2007

    if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) {
        return 'invalid date format (use YYYY-MM-DD)'
    }

    if (isNaN(date.getTime())) {
        return 'invalid date'
    }

    const year = date.getFullYear()

    if (year < minYear || year > maxYear) {
        return `year must be between ${minYear} and ${maxYear}`
    }

    return false
}

export function wrongUsernameError(username: string) {
    if (!minLength(username, 4) || !maxLength(username, 15)) {
        return 'username must be between 4 to 15 characters'
    }

    if (!/^[a-zA-Z0-9_]+$/.test(username)) {
        return 'username can only contain letters, numbers, and underscores'
    }

    return false
}

export function wrongPasswordError(password: string) {
    if (!minLength(password, 6) || !maxLength(password, 30)) {
        return 'password must be between 6 to 30 characters'
    }

    return false
}

export function validateVideoTitle(title: string) {
    const regex = /^[a-zA-Z0-9\-:!?'&.,()]+$/

    if (!title || !title.trim()) {
        return 'Movie title is required'
    }

    if (!minLength(title, 1) || !maxLength(title, 200)) {
        return 'Movie title must be between 1 to 200 characters'
    }

    if (!regex.test(title)) {
        return 'Movie title contains invalid characters'
    }

    return false
}
