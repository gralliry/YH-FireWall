<h1 align="center">YH-FireWall Web Panel</h1>

<p align="center">
  <img alt="Vue" src="https://img.shields.io/badge/Vue-3.5-brightgreen?style=flat-square&logo=vuedotjs">
  <img alt="TypeScript" src="https://img.shields.io/badge/TypeScript-5.9-blue?style=flat-square&logo=typescript">
  <img alt="Vite" src="https://img.shields.io/badge/Vite-7.1-purple?style=flat-square&logo=vite">
  <img alt="License" src="https://img.shields.io/badge/license-MIT-green?style=flat-square">
</p>

<p align="center">
  The web-based management interface for <a href="https://github.com/gralliry/YH-FireWall">YH-FireWall</a>,<br/>
  a lightweight firewall powered by <b>iptables</b> + <b>NFQUEUE</b>.
</p>

---

## Features

- **Rule Management** — CRUD operations on firewall rules with inline editing, column visibility toggles, and server-side persistence
- **Connection Monitor** — view and close active TCP/UDP connections (IPv4/IPv6) with process metadata
- **Interface Overview** — list network interfaces with MAC addresses, MTU, flags, and addresses
- **Config Editor** — live YAML config file editor with line numbers and <kbd>Ctrl+S</kbd> save
- **Dark Mode** — built-in light/dark theme toggle via VueUse
- **Hash Routing** — works without server-side URL rewrite (`/#/rule`, `/#/connection`, etc.)

---

## Tech Stack

| Layer | Choice |
|-------|--------|
| Framework | Vue 3 (Composition API, `<script setup>`) |
| Language | TypeScript |
| Build | Vite 7 |
| UI | Element Plus |
| Router | Vue Router (hash history) |
| HTTP | Axios |
| State | Pinia |
| Editor | CodeMirror 6 (YAML mode) |

---

## Project Structure

```
yfw-front/
├── index.html
├── package.json
├── vite.config.ts
└── src/
    ├── main.ts                  # App entry, global config
    ├── App.vue                  # Root layout (header, nav, router-view)
    ├── style.css                # Global styles
    ├── api/
    │   └── instance.ts          # Axios instance with interceptors
    ├── assets/                  # Static assets
    ├── components/
    │   ├── RuleTable.vue         # Rule CRUD page
    │   ├── ConnectionTable.vue   # Active connections page
    │   ├── InterfaceTable.vue    # Network interfaces page
    │   └── ConfigArea.vue        # YAML config editor page
    └── router/
        └── router.ts            # Route definitions (hash history)
```

---

## Getting Started

### Prerequisites

- Node.js ≥ 18
- A running [YH-FireWall](https://github.com/gralliry/YH-FireWall) instance with the web API enabled

### Install

```bash
npm install
```

### Development

```bash
npm run dev
```

In dev mode, API requests are proxied to `http://172.19.237.114:8080/api` with Basic auth (`admin:admin`).  
Edit `src/main.ts` to point to your own firewall instance.

### Production Build

```bash
npm run build
```

The output lands in `dist/` — serve it behind the firewall's web server or any static file server.  
In production the frontend expects the API at the same origin under `/api`.

---

## Pages

| Path | Page | Description |
|------|------|-------------|
| `/rule` | Rule Table | View, add, edit, delete firewall rules. Toggle column visibility, sort by priority, filter by group/accept/enable. |
| `/connection` | Connection Table | List active connections with PID, executable, user, addresses, and TCP state. Close ESTABLISHED connections. |
| `/interface` | Interface Table | Show network interfaces: index, name, MAC, MTU, flags, and assigned addresses. |
| `/config` | Config Editor | Edit the firewall's `config.yaml` in a textarea with line numbers. Refresh from server, save with <kbd>Ctrl+S</kbd>. |

---

## API Endpoints

The frontend talks to the YH-FireWall web API:

| Method | Path | Purpose |
|--------|------|---------|
| `GET` | `/api/rule` | List all rules |
| `POST` | `/api/rule` | Create a rule |
| `PUT` | `/api/rule/:id` | Update a rule |
| `DELETE` | `/api/rule/:id` | Delete a rule |
| `GET` | `/api/connection` | List active connections |
| `DELETE` | `/api/connection/:id` | Close a connection |
| `GET` | `/api/interface` | List network interfaces |
| `GET` | `/api/config` | Read config file |
| `POST` | `/api/config` | Write config file |

Authentication uses HTTP Basic Auth.

---

## License

[MIT](LICENSE) © Gralliry
