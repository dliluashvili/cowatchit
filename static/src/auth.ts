import { drawFormErrors, generalError, resetErrors } from './errors'
import { renderCalendar } from './calendar'
import {
    isAuthFailure,
    isAuthSuccess,
    type SignInBody,
    type SignUpBody,
} from './types'
import { signIn, signUp } from './api'
import {
    isEmail,
    isGender,
    wrongDob,
    wrongPasswordError,
    wrongUsernameError,
} from './validation'

const signInForm = document.getElementById('signin-form') as HTMLFormElement
const signUpForm = document.getElementById('signup-form') as HTMLFormElement

function switchTab(tab) {
    const signinPanel = document.getElementById('signin-panel')
    const signupPanel = document.getElementById('signup-panel')
    const tabs = document.querySelectorAll('.auth-tab')

    tabs.forEach((t) => t.classList.remove('tab-active'))

    if (tab === 'signin') {
        resetErrors(signInForm)

        signinPanel.classList.remove('hidden')
        signupPanel.classList.add('hidden')
        tabs[0].classList.add('tab-active')
    } else {
        resetErrors(signUpForm)

        signinPanel.classList.add('hidden')
        signupPanel.classList.remove('hidden')
        tabs[1].classList.add('tab-active')
    }
}

document.addEventListener('htmx:load', async function () {
    renderCalendar('#calendar', '#date_of_birth')

    const tabs = document.querySelectorAll('.auth-tab')
    tabs[0]?.addEventListener('click', () => switchTab('signin'))
    tabs[1]?.addEventListener('click', () => switchTab('signup'))

    let errors = {}

    if (signInForm) {
        const action = signInForm.getAttribute('action')

        signInForm.addEventListener('submit', async function (evt) {
            evt.preventDefault()

            const form = this

            const btn = this.querySelector("button[type='submit']")

            btn.setAttribute('disabled', 'disabled')

            resetErrors(signInForm)

            errors = {}

            const body: SignInBody = {
                username: (
                    signInForm.querySelector(
                        'input[name="username"]'
                    ) as HTMLInputElement
                ).value,

                password: (
                    signInForm.querySelector(
                        'input[name="password"]'
                    ) as HTMLInputElement
                ).value,
            }

            const wrongPasswordErrorMsg = wrongPasswordError(body.password)

            const wrongUsernameErrorMsg = wrongUsernameError(body.username)

            if (wrongUsernameErrorMsg || wrongPasswordErrorMsg) {
                generalError(form, 'Invalid credentials')
                btn.removeAttribute('disabled')
                return false
            }

            try {
                const response = await signIn(action, body)

                if (isAuthSuccess(response)) {
                    window.location.href = '/rooms'
                } else if (isAuthFailure(response)) {
                    if (response.status === 422) {
                        const { data: errors } = response
                        drawFormErrors(form, errors)
                    } else {
                        generalError(form, 'Invalid credentials')
                    }
                }
            } catch (error) {
                generalError(form)
            } finally {
                btn.removeAttribute('disabled')
            }
        })
    }

    if (signUpForm) {
        const action = signUpForm.getAttribute('action')

        signUpForm.addEventListener('submit', async function (evt) {
            evt.preventDefault()

            const form = this

            const btn = this.querySelector("button[type='submit']")

            btn.setAttribute('disabled', 'disabled')

            resetErrors(form)

            errors = {}

            const body: SignUpBody = {
                username: (
                    form.querySelector(
                        'input[name="username"]'
                    ) as HTMLInputElement
                ).value,

                email: (
                    form.querySelector(
                        'input[name="email"]'
                    ) as HTMLInputElement
                ).value,

                date_of_birth: (
                    form.querySelector(
                        'input[name="date_of_birth"]'
                    ) as HTMLInputElement
                ).value,

                gender: (
                    form.querySelector(
                        'input[name="gender"]:checked'
                    ) as HTMLSelectElement
                )?.value,

                password: (
                    form.querySelector(
                        'input[name="password"]'
                    ) as HTMLSelectElement
                ).value,

                password_confirmation: (
                    form.querySelector(
                        'input[name="password_confirmation"]'
                    ) as HTMLInputElement
                ).value,
            }

            if (!isEmail(body.email)) {
                errors['email'] = ['wrong email format']
            }

            if (!isGender(body.gender)) {
                errors['gender'] = ['wrong gender']
            }

            const wrongPasswordErrorMsg = wrongPasswordError(body.password)

            if (wrongPasswordErrorMsg) {
                errors['password'] = [wrongPasswordErrorMsg]
            }

            const wrongPasswordConfErrorMsg = wrongPasswordError(
                body.password_confirmation
            )

            if (wrongPasswordConfErrorMsg) {
                errors['password_confirmation'] = wrongPasswordConfErrorMsg
            }

            const wrongUsernameErrorMsg = wrongUsernameError(body.username)

            if (wrongUsernameErrorMsg) {
                errors['username'] = [wrongUsernameErrorMsg]
            }

            if (body.password !== body.password_confirmation) {
                errors['password_confirmation'] = ['passwords do not match']
            }

            const wrongDobMsg = wrongDob(body.date_of_birth)

            if (wrongDobMsg) {
                errors['date_of_birth'] = [wrongDobMsg]
            }

            if (Object.keys(errors).length) {
                drawFormErrors(form, errors)
                btn.removeAttribute('disabled')
                return false
            }

            try {
                const response = await signUp(action, body)

                if (isAuthSuccess(response)) {
                    window.location.href = '/rooms'
                } else if (isAuthFailure(response)) {
                    if (response.status === 422) {
                        const { data: errors } = response
                        drawFormErrors(form, errors)
                    } else {
                        generalError(form, 'Invalid credentials')
                    }
                }
            } catch (error) {
                generalError(form)
            } finally {
                btn.removeAttribute('disabled')
            }
        })
    }
})
