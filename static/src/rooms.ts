import * as v from 'valibot'
import {
    getValueByCheckboxName,
    getValueByInputName,
    getValueByTextareaName,
} from './helpers'
import { drawFormErrors, generalError, resetErrors } from './errors'
import {
    isResponseFailure,
    isResponseSuccess,
    type CreateRoomBody,
} from './types'
import { createRoom } from './api'

document.addEventListener('htmx:load', async function () {
    const createRoomForm = document.querySelector('#create-room-form')

    if (createRoomForm) {
        const createRoomSchema = v.pipe(
            v.object({
                title: v.pipe(
                    v.string(),
                    v.check((title) => {
                        // Deny multiple consecutive spaces
                        return !/\s{2,}/u.test(title)
                    }, 'Title cannot have multiple consecutive spaces'),
                    v.minLength(1, 'Title must be at least 1 character'),
                    v.maxLength(200, 'Title must not exceed 200 characters'),
                    v.regex(
                        /^[\p{L}\p{N}\p{P}\s]+$/u,
                        'Title contains invalid characters.'
                    )
                ),
                capacity: v.pipe(
                    v.number(),
                    v.integer('Capacity must be an integer'),
                    v.minValue(2, 'Capacity must be at least 2'),
                    v.maxValue(10, 'Capacity must not exceed 10')
                ),
                description: v.pipe(v.string(), v.trim()),
                src: v.pipe(v.string(), v.url('Must be a valid URL')),
                private: v.pipe(v.boolean()),
                password: v.optional(
                    v.pipe(
                        v.string(),
                        v.minLength(
                            3,
                            'Password must be at least 3 characters'
                        ),
                        v.maxLength(
                            20,
                            'Password must not exceed 20 characters'
                        )
                    )
                ),
            }),
            v.forward(
                v.check((data) => {
                    return (
                        !data.private ||
                        (data.password !== undefined &&
                            data.password.length > 0)
                    )
                }, 'Password is required for private rooms'),
                ['password']
            )
        )

        createRoomForm
            .querySelector(`input[name="private"]`)
            .addEventListener('change', function () {
                const isChecked = this.checked

                const passwordSection =
                    createRoomForm.querySelector(`#password-section`)
                passwordSection.classList.toggle('hidden', !isChecked)

                const publicIcon =
                    createRoomForm.querySelector(`#privacy-icon-public`)
                const privateIcon = createRoomForm.querySelector(
                    `#privacy-icon-private`
                )
                publicIcon.classList.toggle('hidden', isChecked)
                privateIcon.classList.toggle('hidden', !isChecked)

                const labelText = createRoomForm.querySelector(`#privacy-label`)
                labelText.textContent = isChecked
                    ? 'Private Room'
                    : 'Public Room'

                const description =
                    createRoomForm.querySelector(`#privacy-description`)
                description.textContent = isChecked
                    ? 'Only users with the password can join'
                    : 'Anyone can join this room'
            })

        createRoomForm.addEventListener('submit', async function (e: Event) {
            e.preventDefault()
            const form = this

            resetErrors(form)

            const body: CreateRoomBody = {
                title: getValueByInputName(form, 'title'),
                capacity: parseInt(getValueByInputName(form, 'capacity')),
                description: getValueByTextareaName(form, 'description'),
                src: getValueByInputName(form, 'src'),
                private: getValueByCheckboxName(form, 'private'),
            }

            if (body.private) {
                body.password = getValueByInputName(form, 'password')
            }

            const result = v.safeParse(createRoomSchema, body)

            if (result.success) {
                try {
                    const response = await createRoom(body)

                    if (isResponseSuccess) {
                    } else if (isResponseFailure) {
                        if (response.status === 422) {
                            const { data: errors } = response
                            drawFormErrors(form, errors)
                        } else {
                            generalError(form, 'Unknown Error')
                        }
                    }
                } catch (exception) {}
            } else {
                drawFormErrors(form, result.issues)
            }
        })
    }
})
