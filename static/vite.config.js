import { defineConfig } from 'vite'
import { resolve } from 'path'

export default defineConfig({
    root: '.',
    build: {
        // minify: 'terser',
        outDir: 'dist',
        emptyOutDir: false,
        // terserOptions: {
        //     compress: {
        //         drop_console: true,
        //         drop_debugger: true,
        //     },
        //     mangle: {
        //         properties: false,
        //         toplevel: true,
        //     },
        //     output: {
        //         comments: false,
        //     },
        // },
        rollupOptions: {
            treeshake: true,
            input: {
                app: resolve(__dirname, 'src/app.css'),
                auth: resolve(__dirname, 'src/auth.ts'),
                rooms: resolve(__dirname, 'src/rooms.ts'),
                socket: resolve(__dirname, 'src/socket.ts'),
            },
            output: {
                entryFileNames: 'js/[name].js',
                assetFileNames: 'css/[name].[ext]',
            },
        },
    },
})
