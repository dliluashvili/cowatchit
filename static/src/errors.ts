import type { BaseIssue } from 'valibot'

export function drawFormErrors(
    form: HTMLFormElement,
    errors: Record<string, any>
) {
    const errorMap = isValibotIssues(errors) ? mapValibotErrors(errors) : errors

    Object.keys(errorMap).forEach((key) => {
        let error: string | string[] = errorMap[key]

        if (Array.isArray(error)) {
            error = error[0]
        }

        const el = form.querySelector(`input[name="${key}"]`)

        const formControl = el?.closest('.form-control')
        const errorLabel = formControl?.querySelector('.label:last-child')
        const errorSpan = errorLabel?.querySelector('.label-text-alt')

        if (errorSpan && errorLabel) {
            errorLabel.classList.remove('hidden')
            errorSpan.textContent = error
        }
    })
}

export function generalError(form: HTMLFormElement, text?: string | null) {
    const alertError = form.querySelector('.alert-error')
    const errorText = alertError?.querySelector('.error-text')

    alertError?.classList.remove('hidden')
    errorText.textContent = text ?? 'An error occurred. Please try again'
}

export function resetErrors(form: HTMLFormElement) {
    const errorLabels = form.querySelectorAll(
        '.label .label-text-alt.input-error'
    )

    errorLabels.forEach((errorSpan) => {
        errorSpan.textContent = ''

        const label = errorSpan.closest('.label')
        if (label) {
            label.classList.add('hidden')
        }
    })

    const alertError = form.querySelector('.alert-error')
    const errorText = alertError?.querySelector('.error-text')

    alertError?.classList.add('hidden')
    if (errorText) errorText.textContent = ''
}

export function mapValibotErrors(issues: BaseIssue<unknown>[]) {
    return issues.reduce((acc, issue) => {
        const path = issue.path?.map((p) => p.key).join('.') || 'root'
        acc[path] = issue.message
        return acc
    }, {} as Record<string, string>)
}

export function isValibotIssues(
    errors: unknown
): errors is BaseIssue<unknown>[] {
    return (
        Array.isArray(errors) &&
        errors.length > 0 &&
        'message' in errors[0] &&
        'type' in errors[0]
    )
}
