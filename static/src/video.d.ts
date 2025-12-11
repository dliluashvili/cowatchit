import type Video from 'video.js'

declare global {
    interface Window {
        videojs?: typeof Video
    }
}
