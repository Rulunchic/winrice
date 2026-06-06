# <project-name>

Short description (1-2 sentences).

## Stack

- **Svelte 5** (runes syntax: `$state`, `$derived`, `$effect`, `$props`, `{@render}`)
- **Vite 5** (dev server, build)
- **TypeScript 5**

Styling: geek-minimal, opaque panels, single accent, monospace. See `design-ui` skill.

## Quick start

```bash
npm install
npm run dev
# open http://localhost:5173
```

## Scripts

| Command           | What it does                          |
| ----------------- | ------------------------------------- |
| `npm run dev`     | Vite dev server with HMR              |
| `npm run build`   | Type-check + production build to `dist/` |
| `npm run preview` | Serve `dist/` for inspection          |
| `npm run check`   | `svelte-check` without building       |

## Structure

```
src/
├── App.svelte     — root component
├── app.css        — global styles + reset
├── theme.css      — CSS variables (palette)
├── main.ts        — entry point, mount(App)
├── lib/           — shared components, helpers
└── vite-env.d.ts  — Vite types
static/           — static assets, copied to dist as-is
```

## License

<MIT / Proprietary / ...>
