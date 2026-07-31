# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog,
and this project adheres to Semantic Versioning.

## [Unreleased]

## [0.1.0] - 2026-07-31

### Added
- Initial public release of Subscription Billing Simulator.
- Go REST API for customer creation, subscription lifecycle management, payment failure simulation, retry attempts, and event history.
- PostgreSQL persistence layer for customers, subscriptions, payments, retry attempts, and events.
- Auto-applied SQL migrations for local startup.
- Docker Compose local environment for PostgreSQL-based development.
- Smoke test script for end-to-end verification of the main billing flow.
- Review-ready verification script for local release checks.
- Idempotent payment failure and retry endpoints using idempotency keys.

### Changed
- Improved service testability by refactoring service-layer logic.
- Improved migration loading and local developer workflow.
- Expanded README with architecture, design decisions, local run instructions, and verification steps.
- Updated .gitignore for local review helper artifacts.
