import { test, expect } from '@playwright/test';

test.describe('反馈中心 - H5 端测试', () => {

  test.beforeEach(async ({ page }) => {
    await page.goto('/feedback');
  });

  test('场景1: 用户提交反馈', async ({ page }) => {
    // 切换到"提交反馈"标签
    await page.click('text=提交反馈');

    // 填写反馈表单
    await page.selectOption('select[name="category"]', 'other');
    await page.selectOption('select[name="priority"]', 'medium');
    await page.fill('input[name="title"]', '测试反馈标题 - Playwright测试');
    await page.fill('textarea[name="content"]', '这是一个测试反馈内容，用于验证反馈提交功能是否正常工作。');

    // 提交反馈
    await page.click('button:has-text("提交反馈")');

    // 验证成功提示
    await expect(page.locator('text=反馈提交成功')).toBeVisible();

    // 验证切换到"我的反馈"标签
    await expect(page.locator('text=我的反馈')).toHaveClass(/.*active.*/);

    // 验证新提交的反馈出现在列表中
    await expect(page.locator('text=测试反馈标题 - Playwright测试')).toBeVisible();
  });

  test('场景2: 表单验证 - 空表单', async ({ page }) => {
    // 切换到"提交反馈"标签
    await page.click('text=提交反馈');

    // 尝试提交空表单
    await page.click('button:has-text("提交反馈")');

    // 验证必填字段错误提示
    await expect(page.locator('text=标题为必填项')).toBeVisible();
    await expect(page.locator('text=内容为必填项')).toBeVisible();
  });

  test('场景2: 表单验证 - 标题长度', async ({ page }) => {
    // 切换到"提交反馈"标签
    await page.click('text=提交反馈');

    // 填写过短标题
    await page.fill('input[name="title"]', '短');
    await page.fill('textarea[name="content"]', '这是一个有效的反馈内容');
    await page.click('button:has-text("提交反馈")');

    // 验证标题长度验证提示
    await expect(page.locator('text=标题长度应在5-50字符之间')).toBeVisible();

    // 填写过长标题
    await page.fill('input[name="title"]', 'a'.repeat(51));
    await page.click('button:has-text("提交反馈")');

    // 验证标题长度验证提示
    await expect(page.locator('text=标题长度应在5-50字符之间')).toBeVisible();
  });

  test('场景2: 表单验证 - 内容长度', async ({ page }) => {
    // 切换到"提交反馈"标签
    await page.click('text=提交反馈');

    // 填写过短内容
    await page.fill('input[name="title"]', '有效标题');
    await page.fill('textarea[name="content"]', '短');
    await page.click('button:has-text("提交反馈")');

    // 验证内容长度验证提示
    await expect(page.locator('text=内容长度应在10-500字符之间')).toBeVisible();

    // 填写过长内容
    await page.fill('textarea[name="content"]', 'a'.repeat(501));
    await page.click('button:has-text("提交反馈")');

    // 验证内容长度验证提示
    await expect(page.locator('text=内容长度应在10-500字符之间')).toBeVisible();
  });

  test('场景3: 查看我的反馈列表', async ({ page }) => {
    // 切换到"我的反馈"标签
    await page.click('text=我的反馈');

    // 验证反馈列表显示
    await expect(page.locator('.feedback-list')).toBeVisible();

    // 检查列表项包含必要信息
    const feedbackItems = page.locator('.feedback-item');
    const count = await feedbackItems.count();

    if (count > 0) {
      // 检查第一项是否包含标题、类别、状态、提交时间
      const firstItem = feedbackItems.first();
      await expect(firstItem.locator('.feedback-title')).toBeVisible();
      await expect(firstItem.locator('.feedback-category')).toBeVisible();
      await expect(firstItem.locator('.feedback-status')).toBeVisible();
      await expect(firstItem.locator('.feedback-time')).toBeVisible();
    }

    // 测试状态筛选
    await page.selectOption('select[name="statusFilter"]', 'pending');
    await page.waitForTimeout(500);
    // 验证只显示待处理的反馈

    // 测试查看详情
    if (count > 0) {
      await feedbackItems.first().click();
      await expect(page.locator('.feedback-detail')).toBeVisible();
    }
  });

  test('场景4: 查看 FAQ 列表', async ({ page }) => {
    // 切换到"常见问题"标签
    await page.click('text=常见问题');

    // 验证 FAQ 列表显示
    await expect(page.locator('.faq-list')).toBeVisible();

    // 点击 FAQ 问题展开答案
    const faqItems = page.locator('.faq-item');
    const count = await faqItems.count();

    if (count > 0) {
      const firstFaq = faqItems.first();

      // 点击问题展开答案
      await firstFaq.click();
      await expect(firstFaq.locator('.faq-answer')).toBeVisible();

      // 点击"有帮助"按钮
      const helpfulButton = firstFaq.locator('button:has-text("有帮助")');
      const initialCountText = await helpfulButton.textContent();
      const initialCount = parseInt(initialCountText?.match(/\d+/)?.[0] || '0');

      await helpfulButton.click();

      // 验证点赞数增加
      const newCountText = await helpfulButton.textContent();
      const newCount = parseInt(newCountText?.match(/\d+/)?.[0] || '0');
      expect(newCount).toBeGreaterThan(initialCount);
    }
  });

  test('场景5: 从活动列表进入反馈中心', async ({ page }) => {
    // 访问活动列表页面
    await page.goto('/');

    // 点击"帮助反馈"按钮
    await page.click('button:has-text("帮助反馈")');

    // 验证跳转到反馈中心
    await expect(page).toHaveURL(/\/feedback/);

    // 验证反馈中心页面正确加载
    await expect(page.locator('.feedback-center')).toBeVisible();
  });

  test('完整流程: 用户提交反馈并查看状态', async ({ page }) => {
    // 提交反馈
    await page.click('text=提交反馈');
    await page.selectOption('select[name="category"]', 'other');
    await page.selectOption('select[name="priority"]', 'medium');
    const uniqueTitle = `测试反馈 ${Date.now()}`;
    await page.fill('input[name="title"]', uniqueTitle);
    await page.fill('textarea[name="content"]', '这是一个测试反馈内容，用于验证完整流程。');
    await page.click('button:has-text("提交反馈")');

    // 验证提交成功
    await expect(page.locator('text=反馈提交成功')).toBeVisible();

    // 查看我的反馈
    await page.click('text=我的反馈');

    // 验证反馈出现在列表中
    await expect(page.locator(`text=${uniqueTitle}`)).toBeVisible();

    // 验证状态为"待处理"
    const feedbackItem = page.locator(`.feedback-item:has-text("${uniqueTitle}")`);
    await expect(feedbackItem.locator('.feedback-status:has-text("待处理")')).toBeVisible();
  });
});
