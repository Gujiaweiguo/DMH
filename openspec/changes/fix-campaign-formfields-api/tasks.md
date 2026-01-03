## 1. Backend API Updates
- [x] 1.1 Define FormField struct in dmh.api
- [x] 1.2 Update CreateCampaignReq to use []FormField
- [x] 1.3 Update CampaignResp to use []FormField
- [x] 1.4 Regenerate types with goctl
- [x] 1.5 Update campaign logic to handle FormField objects
- [x] 1.6 Test API endpoints with new structure

## 2. Frontend Integration
- [x] 2.1 Update brandApi.js to send FormField objects
- [x] 2.2 Verify CampaignEditorVant.vue works with new API
- [ ] 2.3 Test campaign creation with complex form fields
- [ ] 2.4 Test campaign editing and form field updates

## 3. Testing and Validation
- [ ] 3.1 Create test campaign with all field types
- [ ] 3.2 Verify form fields are properly saved and retrieved
- [ ] 3.3 Test form preview functionality
- [ ] 3.4 Validate API responses match OpenSpec format