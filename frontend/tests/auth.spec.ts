import { test, expect } from '@playwright/test';

test.describe('Registration Page Tests', () => {
  
  test('Registration form contains email, password and confirm password fields', async ({ page }) => {
    // Go to the registration page
    await page.goto('http://localhost:3000/register');

    // Check if the login link exists and is visible
    const loginLink = await page.$('a[href="/login"]');
    expect(loginLink).not.toBeNull();
    expect(await loginLink?.isVisible()).toBeTruthy();

    // Check if the form exists
    const form = await page.$('form');
    expect(form).not.toBeNull();

    // Check if the email input field exists
    const emailField = await page.$('input[name="email"]');
    expect(emailField).not.toBeNull();

    // Check if the password input field exists
    const passwordField = await page.$('input[name="password"]');
    expect(passwordField).not.toBeNull();

    // Check if the confirm password input field exists
    const confirmPasswordField = await page.$('input[name="confirmPassword"]');
    expect(confirmPasswordField).not.toBeNull();
  });

  test('Fill in the registration form and submit', async ({ page }) => {
    // Go to the registration page
    await page.goto('http://localhost:3000/register');

    // Fill in the email input field
    await page.fill('input[name="email"]', 'test@example.com');

    // Fill in the password input field
    await page.fill('input[name="password"]', 'Password123');

    // Fill in the confirm password input field
    await page.fill('input[name="confirmPassword"]', 'Password123');

    // Submit the form
    await page.click('button[type="submit"]');

    // Check if the page is redirected to the home page
    await page.waitForURL('http://localhost:3000/');
    expect(page.url()).toBe('http://localhost:3000/');
  });
});

test.describe('Login Page Tests', () => {
  test('Login form contains email, password fields', async ({ page }) => {
    // Go to the login page
    await page.goto('http://localhost:3000/login');

    // Check if the register link exists and is visible
    const registerLink = await page.$('a[href="/register"]');
    expect(registerLink).not.toBeNull();
    expect(await registerLink?.isVisible()).toBeTruthy();

    // Check if the form exists
    const form = await page.$('form');
    expect(form).not.toBeNull();

    // Check if the email input field exists
    const emailField = await page.$('input[name="email"]');
    expect(emailField).not.toBeNull();

    // Check if the password input field exists
    const passwordField = await page.$('input[name="password"]');
    expect(passwordField).not.toBeNull();
  });

  test('Fill in the login form and submit', async ({ page }) => {
    // Go to the login page
    await page.goto('http://localhost:3000/login');

    // Fill in the email input field
    await page.fill('input[name="email"]', 'test@example.com');

    // Fill in the password input field
    await page.fill('input[name="password"]', 'Password123');

    // Submit the form
    await page.click('button[type="submit"]');

    // Check if the page is redirected to the home page
    await page.waitForURL('http://localhost:3000/');
    expect(page.url()).toBe('http://localhost:3000/');
  });
});