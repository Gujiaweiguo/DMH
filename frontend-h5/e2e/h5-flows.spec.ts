import { test, expect } from '@playwright/test';

test.describe('H5 Campaign Flow', () => {
  test('user can view campaign list', async ({ page }) => {
    await page.goto('http://localhost:3100');
    
    await expect(page.locator('.campaign-list')).toBeVisible();
    await expect(page.locator('.campaign-item').first()).toBeVisible();
  });

  test('user can view campaign detail', async ({ page }) => {
    await page.goto('http://localhost:3100');
    
    await page.locator('.campaign-item').first().click();
    
    await expect(page.locator('.campaign-detail')).toBeVisible();
    await expect(page.locator('text=立即报名')).toBeVisible();
  });

  test('user can register for campaign', async ({ page }) => {
    await page.goto('http://localhost:3100');
    await page.locator('.campaign-item').first().click();
    await page.click('text=立即报名');
    
    const timestamp = Date.now();
    await page.fill('input[name="name"]', 'Test User');
    await page.fill('input[name="phone"]', `138${timestamp.toString().slice(-8)}`);
    
    await page.click('button[type="submit"]');
    
    await expect(page.locator('text=报名成功')).toBeVisible();
  });
});

test.describe('H5 Distributor Flow', () => {
  test('user can access distributor center', async ({ page }) => {
    await page.goto('http://localhost:3100');
    
    await page.click('text=我的');
    await page.click('text=分销中心');
    
    await expect(page.locator('.distributor-center')).toBeVisible();
  });

  test('user can generate poster', async ({ page }) => {
    await page.goto('http://localhost:3100/distributor/center');
    
    await page.click('text=生成海报');
    
    await expect(page.locator('.poster-preview')).toBeVisible();
    await expect(page.locator('text=保存海报')).toBeVisible();
  });

  test('user can view earnings', async ({ page }) => {
    await page.goto('http://localhost:3100/distributor/center');
    
    await expect(page.locator('.earnings-amount')).toBeVisible();
    await expect(page.locator('text=累计收益')).toBeVisible();
  });
});

test.describe('H5 Order Flow', () => {
  test('user can view order list', async ({ page }) => {
    await page.goto('http://localhost:3100');
    
    await page.click('text=我的');
    await page.click('text=我的订单');
    
    await expect(page.locator('.order-list')).toBeVisible();
  });

  test('user can view order detail', async ({ page }) => {
    await page.goto('http://localhost:3100/orders');
    
    if (await page.locator('.order-item').count() > 0) {
      await page.locator('.order-item').first().click();
      await expect(page.locator('.order-detail')).toBeVisible();
    }
  });
});

test.describe('H5 Brand Admin Flow', () => {
  test('brand admin can login', async ({ page }) => {
    await page.goto('http://localhost:3100/brand/login');
    
    await page.fill('input[name="username"]', 'brand_manager');
    await page.fill('input[name="password"]', '123456');
    await page.click('button[type="submit"]');
    
    await expect(page).toHaveURL(/.*brand.*/);
  });

  test('brand admin can view campaigns', async ({ page }) => {
    await page.goto('http://localhost:3100/brand/login');
    await page.fill('input[name="username"]', 'brand_manager');
    await page.fill('input[name="password"]', '123456');
    await page.click('button[type="submit"]');
    
    await page.click('text=活动管理');
    
    await expect(page.locator('.campaign-management')).toBeVisible();
  });
});
