import { test, expect } from '@playwright/test';

test.describe('Visual Regression Tests', () => {
  test('home page @vrt', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveScreenshot('home-page.png', { fullPage: true });
  });

  test('items list page @vrt', async ({ page }) => {
    await page.goto('/items');
    await page.waitForSelector('table');
    await expect(page).toHaveScreenshot('items-list.png', { fullPage: true });
  });

  test('staff list page @vrt', async ({ page }) => {
    await page.goto('/staff');
    await page.waitForSelector('table');
    await expect(page).toHaveScreenshot('staff-list.png', { fullPage: true });
  });

  test('stores list page @vrt', async ({ page }) => {
    await page.goto('/stores');
    await page.waitForSelector('table');
    await expect(page).toHaveScreenshot('stores-list.png', { fullPage: true });
  });

  test('sales list page @vrt', async ({ page }) => {
    await page.goto('/sales');
    await page.waitForSelector('table');
    await expect(page).toHaveScreenshot('sales-list.png', { fullPage: true });
  });

  test('create item modal @vrt', async ({ page }) => {
    await page.goto('/items');
    await page.click('button:has-text("New Item")');
    await page.waitForSelector('.modal');
    await expect(page.locator('.modal')).toHaveScreenshot('create-item-modal.png');
  });

  test('responsive design - mobile @vrt', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/');
    await expect(page).toHaveScreenshot('home-page-mobile.png', { fullPage: true });
  });

  test('responsive design - tablet @vrt', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.goto('/');
    await expect(page).toHaveScreenshot('home-page-tablet.png', { fullPage: true });
  });
});

test.describe('Dark Mode Visual Regression', () => {
  test.beforeEach(async ({ page }) => {
    // Set dark mode preference
    await page.emulateMedia({ colorScheme: 'dark' });
  });

  test('home page dark mode @vrt', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveScreenshot('home-page-dark.png', { fullPage: true });
  });

  test('items list dark mode @vrt', async ({ page }) => {
    await page.goto('/items');
    await page.waitForSelector('table');
    await expect(page).toHaveScreenshot('items-list-dark.png', { fullPage: true });
  });
});