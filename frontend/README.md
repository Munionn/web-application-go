# Frontend (Next.js)

Next.js 15 app with TypeScript and Tailwind CSS. Talks to the Go API in `../server`.

## Setup

```bash
npm install
```

## Run

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000). The app proxies `/api/*` to the backend at `http://localhost:8080` when using `next dev`.

## Scripts

- `npm run dev` — development server (port 3000)
- `npm run build` — production build
- `npm run start` — run production build
- `npm run lint` — ESLint

## Backend

Start the Go API from the repo root:

```bash
cd server && make run
```

Default API base: `http://localhost:8080`.
