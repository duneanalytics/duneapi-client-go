# End-to-End Tests

This directory contains E2E tests that hit the real Dune API.

## Running the tests

Run all E2E tests:
```bash
go test ./e2e/... -v -timeout 120s
```

Run only upload tests:
```bash
go test ./e2e/... -v -run TestUpload -timeout 60s
```

Run only dataset tests:
```bash
go test ./e2e/... -v -run TestDataset -timeout 60s
```

Skip E2E tests in short mode:
```bash
go test ./e2e/... -short
```

## Environment Variables

**Required environment variables** (tests will fail fast if not set):

- `DUNE_API_KEY` - Dune API key
- `DUNE_API_KEY_OWNER_HANDLE` - Namespace for table operations

Example:
```bash
export DUNE_API_KEY=your_api_key_here
export DUNE_API_KEY_OWNER_HANDLE=your_namespace
go test ./e2e/... -v
```

Or inline:
```bash
DUNE_API_KEY=your_key DUNE_API_KEY_OWNER_HANDLE=your_namespace go test ./e2e/... -v
```

## Notes

- Tests create and cleanup their own test tables
- Table names are timestamped to avoid conflicts: `test_uploads_api_{timestamp}`
- Tests require a Plus subscription for upload operations
- Datasets API tests only require a valid API key
