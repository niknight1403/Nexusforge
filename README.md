# NexusForge

Go chat backend, mobile web UI, Android HTTPS WebView client, and selectable free-tier/local LLM routing.

## Run locally

Requires Go 1.23+. Set one or more of `GROQ_API_KEY`, `OPENROUTER_API_KEY`, `GEMINI_API_KEY` or `OLLAMA_BASE_URL` and run:

```sh
export NEXUSFORGE_ACCESS_TOKEN="choose-a-long-unique-random-secret"
go run ./cmd/nexusforge
```

Open http://127.0.0.1:8080 from the same machine. Providers are tried in the following order: Groq, OpenRouter, Gemini, Ollama. Model overrides: `GROQ_MODEL`, `OPENROUTER_MODEL`, `GEMINI_MODEL`, `OLLAMA_MODEL`. For local Ollama set `OLLAMA_BASE_URL=http://127.0.0.1:11434`, start Ollama and pull a model such as `ollama pull llama3.2`. External Ollama endpoints must use HTTPS; a remote HTTP URL is deliberately rejected. Remote models have quotas, rate limits, and changing availability. No unlimited free cloud inference is guaranteed.

## Docker deployment

```sh
cp .env.example .env
# Edit .env: set a strong random NEXUSFORGE_ACCESS_TOKEN and at least one provider
docker compose up --build -d
curl -f http://127.0.0.1:8080/healthz
```

Compose binds only to loopback on the host. To connect a phone remotely, configure a real HTTPS reverse proxy, authentication, firewall, and domain. The Go server rejects non-loopback binds unless an access token is set. Docker does not automatically provide hosting or Ollama access from a host loopback address.

## API and Android

`GET /healthz` is public; `GET /api/providers` and `POST /api/chat` require the `X-NexusForge-Token` header when a token is set. The browser UI keeps the token in sessionStorage. It does not persist chats on the server. Use HTTPS whenever outside the local machine.

The Android app connects to an existing HTTPS server and contains neither the backend nor provider API keys. Its debug APK is built by GitHub Actions workflow **Android debug APK**; retrieve the artifact from a successful run. A production-signed APK and live model test need separate deployment validation.

## Quality checks

```sh
go test -race ./...
go vet ./...
go build ./cmd/nexusforge
docker build -t nexusforge:local .
```

CI runs Go tests, vet, compilation, image build, Compose config validation and a separate Android debug APK build. The software is a single-user chat foundation, not yet an autonomous multi-agent product.
