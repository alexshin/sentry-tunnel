# Overview

Simple Sentry-Proxy allows to hide traffic by using your own host and endpoint
instead of global one. Then no Ad-Blockers can block traffic

## Configuration

All configuration you can do by providing env-variables:

- APP_HOST - host you would bind to. `Default: ""`
- APP_PORT - host you would bind to. `Default: "3333"`
- SENTRY_HOST - target sentry host you want to send traffic to
- SENTRY_SCHEMA - schema. `Default: "https"`
- SENTRY_PROJECT_IDS - comma separated IDs of projects you would