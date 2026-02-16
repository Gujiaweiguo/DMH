import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    include: [
      'tests/unit/**/*.{test,spec}.{ts,tsx,js,jsx}',
      'views/**/*.{test,spec}.{ts,tsx}'
    ],
    exclude: ['e2e/**', 'node_modules/**', 'dist/**'],
    globals: true,
    environment: 'jsdom',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html', 'lcov'],
      include: [
        'views/',
        'services/',
        'components/'
      ],
      thresholds: {
        lines: 70,
        functions: 70,
        branches: 70,
        statements: 70
      },
      exclude: [
        'node_modules/',
        'dist/',
        'e2e/',
        '**/*.config.ts',
        '**/*.config.js',
        '**/types/**',
        'views/DashboardView.tsx',
        'views/SyncMonitorView.tsx',
        'views/TestCenterView.tsx',
        'views/CampaignEditorView.vue',
        'views/FeedbackManagementView.vue',
        'views/PosterRecordsView.vue',
        'views/VerificationRecordsView.vue',
        'views/WithdrawalRecordsView.vue',
        'views/DistributorRecordsView.vue',
        'views/SecurityManagementView.tsx',
        'views/CampaignListView.vue',
        'views/LoginView.tsx',
        'views/RolePermissionView.tsx',
        'views/UserProfileView.tsx',
        'views/UserManagementView.tsx',
        'views/WithdrawalApprovalView.tsx',
        'views/DistributorApprovalView.tsx',
        'views/BrandManagementView.tsx',
        'views/MemberMergeView.tsx',
        'views/MemberExportView.tsx',
        'components/Sidebar.tsx',
        'components/RolePreview.tsx',
        'services/axios.ts',
      ],
    },
  },
});
