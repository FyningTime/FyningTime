# GitHub Copilot Instructions for FyningTime

## Project Overview

FyningTime is a desktop time tracking application built with Go and the Fyne GUI framework. It helps users track their work time with features like adding/editing/deleting time entries, local SQLite storage, and vacation planning.

## Technology Stack

- **Language**: Go 1.22+
- **GUI Framework**: Fyne v2.7.0
- **Database**: SQLite3 with go-sqlite3 driver
- **Logging**: charmbracelet/log
- **Date Picker**: sdassow/fyne-datepicker

## Project Structure

```
FyningTime/
├── main.go              # Application entry point with setup and window management
├── app/
│   ├── model/          # Data models and business logic (TimeEntry, Settings, etc.)
│   ├── repo/           # Database repository layer
│   ├── service/        # Business services (settings, importer)
│   ├── view/           # UI views (main, calendar, settings, toolbar, vacation)
│   └── widget/         # Custom Fyne widgets
├── .github/
│   └── workflows/      # GitHub Actions CI/CD
└── testdata.sql        # Test data for development
```

## Development Environment Setup

### Prerequisites
- Go 1.22 or higher
- X11 development libraries (Linux): `sudo apt-get install xorg-dev`
- SQLite3

### Building the Application
```bash
go build -v ./...
```

### Running Tests
```bash
go test -v ./...
```

### Running the Application
```bash
go run main.go
```

For development mode with debug logging:
```bash
go run main.go -d
```

## Coding Standards

### General Guidelines
1. **Follow Go conventions**: Use `gofmt` for formatting, follow effective Go principles
2. **Package naming**: Use lowercase, single-word package names
3. **Error handling**: Always handle errors explicitly; use the logging framework for error reporting
4. **Logging**: Use the charmbracelet/log package with appropriate levels (Debug, Info, Error, Fatal)

### Code Style
- Use descriptive variable and function names
- Keep functions focused and small
- Add comments for exported types, functions, and complex logic
- Use uppercase for exported fields in structs (e.g., `DATE`, `ENTRIES`)
- Follow Go's receiver naming conventions (short, consistent names)

### Database Operations
- Always enable foreign keys with `PRAGMA foreign_keys = ON`
- Use prepared statements for SQL queries
- Handle database errors gracefully and log them
- Close database connections in defer statements

### UI Development with Fyne
- Separate UI logic (views) from business logic (models/services)
- Use Fyne's data binding when appropriate
- Handle window lifecycle events (OnClosed)
- Set reasonable default window sizes
- Use system tray integration when available (desktop.App)

### File and Settings Management
- Store application data in user's home directory under `.fyningtime/`
- Use JSON for settings files
- Create default settings if they don't exist
- Always close file handles in defer statements

## Testing Guidelines

Currently, the project has no test files. When adding tests:
- Place test files in the same package as the code being tested
- Name test files with `_test.go` suffix
- Use table-driven tests when appropriate
- Mock external dependencies (database, filesystem)
- Test both success and error cases

## Dependencies Management

- Use `go mod` for dependency management
- Run `go mod tidy` after adding/removing dependencies
- Keep dependencies up to date but test thoroughly after updates
- Document any special dependency requirements

## Keyboard Shortcuts

The application supports these shortcuts:
- **Ctrl+N**: Add new time entry
- **Ctrl+E**: Edit selected time entry
- **Ctrl+Delete**: Delete selected time entry
- **Ctrl+U**: Unselect table item

When adding new shortcuts, update this documentation and the `CreateAppShortcuts` function.

## Common Patterns

### Creating New Models
- Define struct with exported fields for JSON marshaling
- Add constructor function (e.g., `New()`, `NewSettings()`)
- Include JSON tags for serialization
- Add methods for common operations

### Adding New Views
- Create view in `app/view/` package
- Return Fyne container/widget from creation function
- Use the repository pattern for data access
- Handle errors with dialogs for user feedback

### Database Schema Changes
- Update schema in repository initialization
- Consider migration path for existing users
- Test with both new and existing databases
- Update `testdata.sql` if needed

## CI/CD

The project uses GitHub Actions for continuous integration:
- Builds are triggered on push to `main` and `feature/*` branches
- Pull requests to `main` also trigger builds
- CI installs required dependencies (xorg-dev)
- CI runs both build and test commands

## Known Limitations

- The project currently has no automated tests
- SQLite database is local-only (by design)
- Some planned features are not yet implemented (see README)

## Contributing Guidelines

When making changes:
1. Create feature branches from `main`
2. Keep commits focused and well-described
3. Ensure code builds successfully
4. Update documentation for user-facing changes
5. Follow the existing code style and patterns
6. Test manually since automated tests are limited

## Security Considerations

- Never commit database files with user data
- Store sensitive information securely
- Validate all user inputs
- Use parameterized queries to prevent SQL injection
- Handle file operations securely with proper error checking
