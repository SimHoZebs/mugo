# LazyFood Go Server - Agent Guidelines

## Build & Development Commands
- **Build**: `go build -o server ./server/`
- **Run**: `cd server && ./server` (runs on localhost:8080)
- **Generate Swagger docs**: `make swagger-gen` or `swag init ./server/`
- **Install dependencies**: `go mod tidy`

## Code Style Guidelines
- **Imports**: Group standard library, then third-party, then local imports (server/*)
- **API**: Use Huma v2 with Chi router, follow existing endpoint patterns in routes/
- **Environment**: Load .env with godotenv, never commit secrets

## ADK (Agent Development Kit) Guidelines
- **ADK Documentation**: https://google.github.io/adk-docs/ (primary reference)
- **Go ADK API**: https://pkg.go.dev/google.golang.org/adk
- **Agent Creation**: Use llmagent.New() with gemini models (see agents/weather.go:31)
- **Runner Pattern**: Create agents in agents/, initialize runners in runners/ (see runners/echo.go:10)
- **Session Management**: Use session.InMemoryService() for conversation state
- **Tools**: Leverage geminitool.GoogleSearch{} and other built-in tools
- **Model Config**: Use gemini.NewModel() with GOOGLE_API_KEY from environment

---

## Notion Pages Relevant to Mugo (LazyFood / ai-nutrition-tracker)
Below are Notion pages and databases I found that appear directly related to this project. Each entry is the page title followed by the Notion URL.


- Mugo — https://www.notion.so/Mugo-29cee37bd028800ba00cd39ec50dd022
- LazyFood Tasks (database) — https://www.notion.so/6934f50dcef646e6a374a6edc0b44319
- Build single nutrition agent — https://www.notion.so/Build-single-nutrition-agent-fb282cf142ed48d58358871f8c09b129
- LazyFood: Technical Architecture — https://www.notion.so/LazyFood-Technical-Architecture-ce5f20d767774a04bb89a83c4e029bd2
- LazyFood: Market Analysis — https://www.notion.so/LazyFood-Market-Analysis-7843fc72e72441a1afd1b5696490e5d0
- LazyFood: Business Plan — https://www.notion.so/LazyFood-Business-Plan-2b99b481337f4501b3a7df083ad73574
- Add LazyFood project to website — https://www.notion.so/Add-LazyFood-project-to-website-285ee37bd02880c89984eca5cbbde3fd
- LazyFood: October 27, 2025 → November 2, 2025 (weekly note) — https://www.notion.so/LazyFood-October-27-2025-November-2-2025-6d5a994b19694b9f9387b48cea796fed
- Integrate Gemini API for AI cleanup — https://www.notion.so/Integrate-Gemini-API-for-AI-cleanup-5f4892c3f93042fb8ca938c0959fd226
- Integrate Gemini API for AI summaries — https://www.notion.so/Integrate-Gemini-API-for-AI-summaries-90c3efc56a574de1a169557aa68b8829

If you want more pages added (e.g., broader research links, weekly notes, or task pages), tell me the scope and I can append them.
