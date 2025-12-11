// calendar.ts

interface CalendarState {
    currentMonth: number
    currentYear: number
    selectedDate: Date | null
    showYearPicker: boolean
    currentDecadeStart: number
}

const months = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December',
]

export function renderCalendar(selector: string, inputSelector?: string) {
    const container = document.querySelector(selector)
    if (!container) {
        console.error(`Calendar container "${selector}" not found`)
        return
    }

    // Get input field if provided
    const inputField = inputSelector
        ? (document.querySelector(inputSelector) as HTMLInputElement)
        : null

    // Calculate year ranges
    const currentYear = new Date().getFullYear()
    const minYear = currentYear - 100
    const maxYear = 2007

    const state: CalendarState = {
        currentMonth: new Date().getMonth(),
        currentYear: 1990,
        selectedDate: null,
        showYearPicker: false,
        currentDecadeStart: 1990,
    }

    // Get DOM elements
    const monthSelect = container.querySelector(
        '#monthSelect'
    ) as HTMLSelectElement
    const yearBtn = container.querySelector('#yearBtn') as HTMLButtonElement
    const daysGrid = container.querySelector('#daysGrid') as HTMLElement
    const prevBtn = container.querySelector('#prevBtn') as HTMLButtonElement
    const nextBtn = container.querySelector('#nextBtn') as HTMLButtonElement
    const yearPicker = container.querySelector('#yearPicker') as HTMLElement
    const calendarContent = container.querySelector(
        '#calendarContent'
    ) as HTMLElement
    const prevDecade = container.querySelector(
        '#prevDecade'
    ) as HTMLButtonElement
    const nextDecade = container.querySelector(
        '#nextDecade'
    ) as HTMLButtonElement
    const decadeBtn = container.querySelector('#decadeBtn') as HTMLButtonElement
    const yearsGrid = container.querySelector('#yearsGrid') as HTMLElement

    if (
        !monthSelect ||
        !yearBtn ||
        !daysGrid ||
        !prevBtn ||
        !nextBtn ||
        !yearPicker ||
        !calendarContent ||
        !prevDecade ||
        !nextDecade ||
        !decadeBtn ||
        !yearsGrid
    ) {
        console.error('Required calendar elements not found')
        return
    }

    // Populate month dropdown
    months.forEach((month, idx) => {
        const option = document.createElement('option')
        option.value = idx.toString()
        option.textContent = month
        monthSelect.appendChild(option)
    })

    // Generate years for decade (Ant Design style: 3x4 grid)
    function generateYearsForDecade() {
        yearsGrid.innerHTML = ''

        const decadeStart = state.currentDecadeStart
        const decadeEnd = Math.min(decadeStart + 9, maxYear)

        // Update decade button text
        decadeBtn.textContent = `${decadeStart}-${decadeEnd}`

        // Generate only years in the decade (no prev/next)
        for (let year = decadeStart; year <= decadeEnd; year++) {
            const btn = document.createElement('button')
            btn.type = 'button'
            btn.textContent = year.toString()

            // Check if selected
            const isSelected = year === state.currentYear

            if (isSelected) {
                btn.className =
                    'btn btn-sm bg-white/20 text-white border-white/30 h-12'
            } else {
                btn.className =
                    'btn btn-sm btn-ghost text-white hover:bg-white/10 h-12'
            }

            btn.addEventListener('click', () => {
                state.currentYear = year
                state.showYearPicker = false
                calendarContent.classList.remove('hidden')
                yearPicker.classList.add('hidden')
                updateCalendar()
            })

            yearsGrid.appendChild(btn)
        }
    }

    function updateCalendar() {
        daysGrid.innerHTML = ''

        monthSelect.value = state.currentMonth.toString()
        yearBtn.textContent = state.currentYear.toString()

        const firstDay = new Date(
            state.currentYear,
            state.currentMonth,
            1
        ).getDay()
        const daysInMonth = new Date(
            state.currentYear,
            state.currentMonth + 1,
            0
        ).getDate()
        const prevMonthDays = new Date(
            state.currentYear,
            state.currentMonth,
            0
        ).getDate()

        const startDay = firstDay === 0 ? 6 : firstDay - 1

        let dayCount = 1
        let nextMonthDay = 1

        const totalCells = Math.ceil((startDay + daysInMonth) / 7) * 7

        for (let i = 0; i < totalCells; i++) {
            const btn = document.createElement('button')
            btn.type = 'button'

            if (i < startDay) {
                btn.textContent = (prevMonthDays - startDay + i + 1).toString()
                btn.className =
                    'btn btn-sm btn-ghost w-10 h-10 p-0 min-h-0 text-white/60 hover:bg-white/5'
                btn.disabled = true
            } else if (dayCount > daysInMonth) {
                btn.textContent = (nextMonthDay++).toString()
                btn.className =
                    'btn btn-sm btn-ghost w-10 h-10 p-0 min-h-0 text-white/60 hover:bg-white/5'
                btn.disabled = true
            } else {
                const day = dayCount++
                btn.textContent = day.toString()

                const today = new Date()
                const isToday =
                    day === today.getDate() &&
                    state.currentMonth === today.getMonth() &&
                    state.currentYear === today.getFullYear()

                const isSelected =
                    state.selectedDate &&
                    day === state.selectedDate.getDate() &&
                    state.currentMonth === state.selectedDate.getMonth() &&
                    state.currentYear === state.selectedDate.getFullYear()

                if (isSelected) {
                    btn.className =
                        'btn btn-sm w-10 h-10 p-0 min-h-0 bg-white/20 text-white border-white/30 hover:bg-white/30'
                } else if (isToday) {
                    btn.className =
                        'btn btn-sm w-10 h-10 p-0 min-h-0 text-white border border-white/50 hover:bg-white/10'
                } else {
                    btn.className =
                        'btn btn-sm btn-ghost w-10 h-10 p-0 min-h-0 text-white hover:bg-white/10 hover:border-white/30'
                }

                btn.addEventListener('click', () => {
                    state.selectedDate = new Date(
                        state.currentYear,
                        state.currentMonth,
                        day
                    )
                    updateCalendar()

                    const event = new CustomEvent('dateSelected', {
                        detail: { date: state.selectedDate },
                    })
                    container.dispatchEvent(event)
                })
            }

            daysGrid.appendChild(btn)
        }
    }

    // Event listeners
    prevBtn.addEventListener('click', () => {
        if (state.currentMonth === 0) {
            state.currentMonth = 11
            state.currentYear--
        } else {
            state.currentMonth--
        }
        updateCalendar()
    })

    nextBtn.addEventListener('click', () => {
        if (state.currentMonth === 11) {
            state.currentMonth = 0
            state.currentYear++
        } else {
            state.currentMonth++
        }
        updateCalendar()
    })

    monthSelect.addEventListener('change', (e) => {
        state.currentMonth = parseInt((e.target as HTMLSelectElement).value)
        updateCalendar()
    })

    yearBtn.addEventListener('click', () => {
        state.showYearPicker = true
        state.currentDecadeStart = Math.floor(state.currentYear / 10) * 10
        calendarContent.classList.add('hidden')
        yearPicker.classList.remove('hidden')
        generateYearsForDecade()
    })

    prevDecade.addEventListener('click', () => {
        state.currentDecadeStart -= 10
        if (state.currentDecadeStart < Math.floor(minYear / 10) * 10) {
            state.currentDecadeStart = Math.floor(minYear / 10) * 10
        }
        generateYearsForDecade()
    })

    nextDecade.addEventListener('click', () => {
        state.currentDecadeStart += 10
        if (state.currentDecadeStart > Math.floor(maxYear / 10) * 10) {
            state.currentDecadeStart = Math.floor(maxYear / 10) * 10
        }
        generateYearsForDecade()
    })

    // Initial render
    calendarContent.classList.remove('hidden')
    yearPicker.classList.add('hidden')

    // Hide calendar initially if input field is provided
    if (inputField) {
        container.classList.add('hidden')

        // Show calendar on input click
        inputField.addEventListener('click', () => {
            container.classList.remove('hidden')
        })

        // Hide calendar when clicking outside
        document.addEventListener('click', (e) => {
            if (
                !container.contains(e.target as Node) &&
                e.target !== inputField
            ) {
                container.classList.add('hidden')
            }
        })

        // Update input field when date is selected
        container.addEventListener('dateSelected', ((e: CustomEvent) => {
            if (inputField && e.detail.date) {
                const date = e.detail.date as Date

                // ✅ Format as YYYY-MM-DD
                const year = date.getFullYear()
                const month = String(date.getMonth() + 1).padStart(2, '0')
                const day = String(date.getDate()).padStart(2, '0')

                inputField.value = `${year}-${month}-${day}`
                container.classList.add('hidden')
            }
        }) as EventListener)
    }

    updateCalendar()
}
