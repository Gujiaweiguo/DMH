import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    include: ['tests/unit/**/*.{test,spec}.{ts,tsx,js,jsx}'],
    exclude: ['e2e/**', 'node_modules/**', 'dist/**'],
  },
});
