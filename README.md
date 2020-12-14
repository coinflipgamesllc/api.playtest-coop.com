# Playtest Co-op API

This API handles the backend behavior for [Playtest Co-op](https://playtest-coop.com). Authentication, games content management, etc.

### Running your own instance

There's not a ton of flexibility right now in terms of external dependencies. PostgreSQL and Mailgun are both required for the moment.

Copy `.env.example` to `.env` and fill in the necessary values. Run via `go run main.go`.
