# notime.run

CS2 community server browser with native-app-like performance and features.

## Features

### Backend (Go)
- Fetch and cache CS2 server data from Steam API
- RESTful API endpoints for server data
- Efficient server status monitoring
- Rate limiting and error handling
- Server filtering and sorting capabilities

### Frontend (HTMX + Alpine)
- Fast, responsive UI with instant filtering and sorting
- Real-time server status updates
- Offline support
- PWA installable app
- Smooth transitions and animations
- Steam connect link integration

### Possible Features
- Server favorites/history
- Advanced filtering (by region, map, player count, etc.)
- Server ping information
- Rich server details view
- Native sharing capabilities
- Dark/Light theme support
- Steam server browser protocol integration
- Server status notifications

## Technical Stack
- Backend: Go with Gin framework
- Frontend: HTMX + Alpine
- API: Steam Web API
- Development: Docker + CompileDaemon for hot reloading
