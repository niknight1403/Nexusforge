# NexusForge

Go HTTP backend with an Android HTTPS WebView client, a mobile chat UI and free-tier LLM provider fallback.

## Backend (Go 1.23+)

```bash
export GROQ_API_KEY="..."          # optional
export OPENROUTER_API_KEY="..."    # optional
export GEMINI_API_KEY="..."        # optional
export NEXUSFORGE_ACCESS_TOKEN="choose-a-long-random-secret"
go run ./cmd/nexusforge
```

Open http://127.0.0.1:8080 locally. Only providers with configured API keys are used, in order: Groq → OpenRouter → Gemini. Model names can be overridden with `GROQ_MODEL`, `OPENROUTER_MODEL`, `GEMINI_MODEL`. Provider free tiers have quotas, rate limits, and model availability changes; **unlimited free inference is not guaranteed**. The router falls back after failed calls and never prints API keys.

Endpoints: `GET /healthz` (public), `GET /api/providers` and `POST /api/chat` (token-protected when set). A browser stores its access token in sessionStorage only. For public deployment, set `NEXUSFORGE_ADDR=0.0.0.0:8080`, set a strong access token, and put the backend behind HTTPS using a reverse proxy. Do not expose the service publicly without authentication and TLS. No cloud host is automatically provisioned.

## Android

The Android client requires the HTTPS URL of the deployed backend. The APK does not contain the Go backend or any provider credentials. Launch the **Android debug APK** workflow on GitHub Actions, or build with Java 17, Android SDK 35, and Gradle 8.9 using `cd android-app && gradle assembleDebug`. Download the debug artifact after a **successful** run. Debug builds are for testing, not signed production release distribution.

## Verification

Run `go test -race ./...`, `go vet ./...` and `go build ./cmd/nexusforge`. CI workflows: **Go backend checks** and **Android debug APK**. This is a usable single-user starter, not a production multi-agent orchestrator.
