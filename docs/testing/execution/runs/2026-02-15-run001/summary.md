# Run Summary - 2026-02-15-run001

## Scope

- Environment: Docker Compose (`dmh-api`, `dmh-nginx`, `mysql8`, `redis-dmh`)
- Test groups:
  - Backend full integration: `go test ./test/integration/... -v -count=1`
  - Admin E2E: `npm run test:e2e -- --reporter=line`
  - H5 E2E: `npm run test:e2e -- --reporter=line`

## Result

- Backend integration: PASS (`ok dmh/test/integration 3.998s`)
- Admin E2E: PASS (`21 passed`)
- H5 E2E: PASS (`7 passed`)
- Release gate impact: PASS

## Evidence Files

- `backend-integration.log`
- `admin-e2e.log`
- `h5-e2e.log`
- `docker-ps.log`
- `dmh-api.log`
- `dmh-nginx.log`
