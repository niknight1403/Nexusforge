# NexusForge

Standalone Go orchestration backend with Android WebView client.

## Android
The Android app asks for the HTTPS address of a deployed NexusForge server. It does not bundle or run the Go backend. The Android debug APK is built by the GitHub Actions workflow `Android debug APK` (Actions → Run workflow) and appears as a workflow artifact **only after a successful run**.

## Status
The Android client is committed. Backend source synchronization is still pending; the repository is not yet a complete deployable release.
