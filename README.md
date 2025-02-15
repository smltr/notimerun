# notime.run

A fast, Steam-like server browser for Counter-Strike 2 community servers.

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
- Single-page application design
- Alpine.js for reactive state management
- Fast client-side sorting and filtering
- Steam protocol integration for joining servers
- Responsive layout matching Steam's design

## Technical Stack
- **Backend**: Go 1.21 with Gin framework
- **Frontend**: Alpine.js for reactivity
- **API**: Steam Web API for server data
- **Development**: Docker + CompileDaemon for hot reloading

## Planned Features
- Server favorites
- Server history
- Spectator support
- LAN server discovery
- Detailed server info panel
- Filtering by map, game type, and region
- Server ping information
- Advanced search options
- Player count graphs
- Friend integration

## Development

### Prerequisites
- Docker and Docker Compose
- A Steam Web API key (set in `.env`)

### Running Locally
```bash
make dev     # Build and run with hot reloading
make build   # Build the Docker image
make up      # Start the container
make down    # Stop and clean up
make clean   # Full cleanup
```
