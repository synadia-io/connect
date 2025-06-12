# Contributing to Synadia Connect

Hello and welcome! We're thrilled that you're interested in contributing to Synadia Connect. This document provides guidelines and information to help you contribute effectively.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. We expect all contributors to:

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

1. **Fork the Repository**: Click the "Fork" button on GitHub to create your own copy
2. **Clone Your Fork**: 
   ```bash
   git clone https://github.com/YOUR-USERNAME/connect.git
   cd connect
   ```
3. **Add Upstream Remote**:
   ```bash
   git remote add upstream https://github.com/synadia-io/connect.git
   ```
4. **Create a Branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Prerequisites

- Go 1.22 or later
- Task (taskfile.dev)
- Docker (for runtime testing)
- NATS Server (for local testing)

### Building the Project

```bash
# Install dependencies
task deps

# Build all components
task build

# Run tests
task test

# Install binaries
task install
```

### Running Locally

1. Start a local NATS server:
   ```bash
   nats-server -js
   ```

2. Build and run the CLI:
   ```bash
   cd connect
   task build
   ./target/connect --help
   ```

## Contribution Process

### Before You Start

1. **Check Existing Issues**: Look for existing issues or discussions related to your idea
2. **Open a Discussion**: For significant changes, open a discussion first to get feedback
3. **Claim an Issue**: Comment on an issue to let others know you're working on it

### Making Changes

1. **Write Clean Code**:
   - Follow Go best practices and idioms
   - Use meaningful variable and function names
   - Keep functions small and focused
   - Add comments for complex logic

2. **Follow Code Style**:
   - Run `go fmt` before committing
   - Use `golangci-lint` to catch common issues
   - Follow existing patterns in the codebase

3. **Write Tests**:
   - Add unit tests for new functionality
   - Ensure all tests pass: `task test`
   - Aim to maintain or improve code coverage
   - Use table-driven tests where appropriate

4. **Update Documentation**:
   - Update relevant documentation
   - Add code comments for exported functions
   - Update README if adding new features

### Commit Guidelines

We value clean git history. Please follow these guidelines:

1. **Write Clear Commit Messages**:
   ```
   Short summary (50 chars or less)
   
   More detailed explanation if needed. Wrap at 72 characters.
   Explain what and why, not how.
   
   Fixes #123
   ```

2. **Keep Commits Atomic**:
   - One logical change per commit
   - Commits should be self-contained
   - Each commit should pass tests

3. **Sign Your Commits**:
   ```bash
   git commit -s -m "Your commit message"
   ```

### Pull Request Process

1. **Update Your Branch**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Push to Your Fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create Pull Request**:
   - Use a clear, descriptive title
   - Reference any related issues
   - Describe what changes you made and why
   - Include testing instructions

4. **PR Template**:
   ```markdown
   ## Description
   Brief description of changes
   
   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Breaking change
   - [ ] Documentation update
   
   ## Testing
   - [ ] Unit tests pass
   - [ ] Integration tests pass
   - [ ] Manual testing completed
   
   ## Checklist
   - [ ] Code follows project style
   - [ ] Self-review completed
   - [ ] Documentation updated
   - [ ] Tests added/updated
   ```

## Testing Guidelines

### Unit Tests

```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "test",
            expected: "TEST",
            wantErr:  false,
        },
        // Add more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := MyFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("MyFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result != tt.expected {
                t.Errorf("MyFunction() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Integration Tests

Use Ginkgo for BDD-style tests:

```go
var _ = Describe("MyComponent", func() {
    Context("when initialized", func() {
        It("should connect successfully", func() {
            component := NewMyComponent()
            err := component.Connect()
            Expect(err).ToNot(HaveOccurred())
        })
    })
})
```

## Project Structure

```
connect/
├── cli/              # CLI implementation
├── client/           # Go client SDK
├── model/            # Data models (auto-generated)
├── spec/             # Connector specifications
├── builders/         # Fluent API builders
├── convert/          # Conversion utilities
├── docs/             # Documentation
└── test/             # Test utilities
```

### Key Packages

- **cli**: Command-line interface implementation
- **client**: SDK for interacting with Connect API
- **model**: Auto-generated models from JSON schemas
- **builders**: Fluent API for building configurations

## Common Tasks

### Adding a New CLI Command

1. Create command file in `cli/`:
   ```go
   type myCommand struct {
       opts *Options
       // command specific fields
   }
   
   func ConfigureMyCommand(app commandHost, opts *Options) {
       c := &myCommand{opts: opts}
       cmd := app.Command("mycommand", "Description")
       cmd.Action(c.run)
   }
   ```

2. Register in `main.go`:
   ```go
   cli.ConfigureMyCommand(ncli, opts)
   ```

### Updating Models

1. Edit JSON schema in `model/schemas/`
2. Run model generation:
   ```bash
   task models:generate
   ```
3. Commit both schema and generated code

### Adding Tests

1. Create test file alongside code: `myfile_test.go`
2. For CLI tests, create helper functions to avoid `os.Exit`
3. Use mocks for external dependencies
4. Run tests: `go test ./...`

## Documentation

### Code Documentation

- Document all exported functions and types
- Use complete sentences
- Include examples for complex functions

```go
// ProcessMessage transforms the input message according to the configured rules
// and returns the processed result. If the message cannot be processed, an error
// is returned along with the original message.
//
// Example:
//   result, err := ProcessMessage(msg, rules)
//   if err != nil {
//       log.Printf("Failed to process: %v", err)
//   }
func ProcessMessage(msg Message, rules []Rule) (Message, error) {
    // Implementation
}
```

### User Documentation

Update relevant files in `docs/`:
- API changes → `docs/api/`
- CLI changes → `docs/cli-reference.md`
- New features → `docs/getting-started.md`

## Debugging Tips

### Enable Debug Logging

```go
import "log/slog"

slog.SetLogLoggerLevel(slog.LevelDebug)
slog.Debug("Debug message", "key", "value")
```

### Common Issues

1. **Model Generation Fails**:
   - Ensure JSON schemas are valid
   - Check you have the required tools: `task models:deps`

2. **Tests Timeout**:
   - Increase timeout in test
   - Check for deadlocks or infinite loops

3. **Build Fails**:
   - Run `go mod tidy`
   - Check Go version: `go version`

## Community

### Getting Help

- **Slack**: Join #connectors on [NATS Slack](https://slack.nats.io)
- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and ideas

### Code Review

All submissions require review. We use GitHub pull requests for this purpose. Expect feedback and be prepared to make adjustments.

## License

By contributing, you agree that your contributions will be licensed under the Apache 2.0 License.

## Recognition

Contributors are recognized in several ways:
- Listed in release notes
- Mentioned in commit history
- Added to CONTRIBUTORS file (for significant contributions)

Thank you for contributing to Synadia Connect! Your efforts help make data connectivity better for everyone.