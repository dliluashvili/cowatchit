interface Window {
    htmx?: {
        config: {
            wsReconnectDelay?: (retryCount: number) => number
        }
        ajax: any
    }
}
