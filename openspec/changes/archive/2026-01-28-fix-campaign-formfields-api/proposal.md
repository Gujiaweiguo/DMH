# Change: Fix Campaign FormFields API Data Structure

## Why
The backend API definition for campaign FormFields is inconsistent with the OpenSpec specification. The API currently defines FormFields as `[]string`, but the spec requires complex objects with type, name, label, required, placeholder, and options fields. This prevents proper dynamic form functionality.

## What Changes
- **BREAKING**: Update backend API definition from `[]string` to `[]FormField` objects
- Update Go types and handlers to support complex FormField structure
- Ensure frontend can send and receive proper FormField objects
- Maintain backward compatibility where possible

## Impact
- Affected specs: campaign-management
- Affected code: 
  - `backend/api/dmh.api` - API definitions
  - `backend/api/internal/types/` - Generated types
  - `backend/api/internal/logic/campaign/` - Campaign logic
  - `frontend-h5/src/services/brandApi.js` - API client
  - `frontend-h5/src/views/brand/CampaignEditorVant.vue` - Form editor