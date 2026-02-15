import { describe, expect, it } from "vitest";
import { resolveGuardNavigation } from "../../src/router/guard.js";

describe("router auth guard resolver", () => {
	it("allows route without auth requirement", () => {
		const to = { path: "/", meta: {} };
		expect(resolveGuardNavigation(to, null, null, null)).toBeNull();
	});

	it("redirects distributor path to distributor login when no token", () => {
		const to = { path: "/distributor/rewards", meta: { requiresAuth: true } };
		expect(resolveGuardNavigation(to, null, null, null)).toBe("/distributor/login");
	});

	it("redirects brand path to brand login when no token", () => {
		const to = { path: "/brand/dashboard", meta: { requiresAuth: true } };
		expect(resolveGuardNavigation(to, null, null, null)).toBe("/brand/login");
	});

	it("redirects when required role mismatches", () => {
		const to = { path: "/distributor/rewards", meta: { requiresAuth: true, role: "distributor" } };
		expect(resolveGuardNavigation(to, "token", "brand_admin", null)).toBe("/distributor/login");
	});

	it("redirects brand pages when user has no brand ids", () => {
		const to = { path: "/brand/orders", meta: { requiresAuth: true, hasBrand: true } };
		expect(resolveGuardNavigation(to, "token", "brand_admin", JSON.stringify({ brandIds: [] }))).toBe(
			"/brand/login",
		);
	});

	it("allows brand pages with valid brand ids", () => {
		const to = { path: "/brand/orders", meta: { requiresAuth: true, hasBrand: true } };
		expect(
			resolveGuardNavigation(to, "token", "brand_admin", JSON.stringify({ brandIds: [1, 2] })),
		).toBeNull();
	});

	it("handles malformed user info json and redirects safely", () => {
		const to = { path: "/brand/orders", meta: { requiresAuth: true, hasBrand: true } };
		expect(resolveGuardNavigation(to, "token", "brand_admin", "{bad-json")).toBe("/brand/login");
	});

	it("redirects brand admin role mismatch to brand login", () => {
		const to = { path: "/brand/orders", meta: { requiresAuth: true, role: "brand_admin" } };
		expect(resolveGuardNavigation(to, "token", "participant", null)).toBe("/brand/login");
	});

	it("redirects unknown role mismatch to home", () => {
		const to = { path: "/special/path", meta: { requiresAuth: true, role: "special_role" } };
		expect(resolveGuardNavigation(to, "token", "brand_admin", null)).toBe("/");
	});
});
