# findservers.net

A fast, Steam-like server browser for Counter-Strike 2 community servers. My hope is that a no-nonsense, ad-free tool like this will help the custom server community grow.

## Current Features

### Server List
- Live community server list from Steam API
- Auto-refresh every 30 seconds
- Double-click to connect directly to servers
- Sort by server name, player count, or map
- Steam-like UI with classic green theme

### Backend
- Efficient in-memory server caching
- Filters out Valve official servers
- RESTful API endpoint for server data
- Automatic server list updates
- Debug logging of server data

### Frontend
- Single-page application design, homepage is the app
- Alpine.js for reactive state management
- Fast client-side sorting and filtering
- Steam protocol integration for joining servers
- Classic steam inspired design

## Technical Stack
- **Backend**: Go 1.21 with Gin framework
- **Frontend**: Alpine.js for reactivity
- **API**: Steam Web API for server data
- **Development**: Docker + CompileDaemon for hot reloading

## Planned Features
- Server favorites
- Server history
- Detailed server info panel, possibly including user ping
- More advanced filtering
- Friend/steam integration

## Development

### Lifecycle
The current goal is to develop this app bit by bit as live software without major releases.

### Prerequisites
- Docker and Docker Compose
- A Steam Web API key (set in `.env` as `STEAM_API_KEY`), can be obtained from [Steam Web API](https://steamcommunity.com/dev/apikey)

### Running Locally
Local developement is containerized with Docker. Make sure Docker desktop is installed and use the following commands as listed in the `makefile`

```bash
make dev     # Build and run with hot reloading
make build   # Build the Docker image
make up      # Start the container
make down    # Stop and clean up
make clean   # Full cleanup
```

### Deployment
The app is currently deployed on fly.io. See fly.toml for configuration details.
