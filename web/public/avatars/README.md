# Avatars folder

Drop any avatar images you want to provide as preinstalled avatars into this folder.

How to use
- Files placed here will be served statically at runtime from `/avatars/<filename>` (for example: `/avatars/avatar1.png`).
- Prefer square images (e.g. 200x200 or 400x400) and common formats (png, jpg, webp).
- Keep filenames URL-safe (no spaces). Example: `avatar_user1.png`.

Notes
- If you build with Vite and serve the app in production, the `public` folder content is usually copied to the root of the built site — `/avatars/...` will work the same way.
- If you prefer bundling images via imports, move them to `src/assets/` and import directly in code.
