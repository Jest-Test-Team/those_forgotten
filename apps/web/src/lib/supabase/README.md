# Supabase Clients

- `client.ts`: browser-side PKCE client for Google OAuth and session persistence.
- `server.ts`: server-safe client for page and route handlers.

These wrappers are intentionally thin so auth wiring can evolve without changing callers.
