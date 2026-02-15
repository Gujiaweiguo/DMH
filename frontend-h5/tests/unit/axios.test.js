import { beforeEach, describe, expect, it, vi } from "vitest";

const mockState = vi.hoisted(() => ({
	requestOnFulfilled: null,
	requestOnRejected: null,
	responseOnFulfilled: null,
	responseOnRejected: null,
	toastMock: vi.fn(),
}));

vi.mock("axios", () => ({
	default: {
		create: () => ({
			interceptors: {
				request: {
					use: (onFulfilled, onRejected) => {
						mockState.requestOnFulfilled = onFulfilled;
						mockState.requestOnRejected = onRejected;
					},
				},
				response: {
					use: (onFulfilled, onRejected) => {
						mockState.responseOnFulfilled = onFulfilled;
						mockState.responseOnRejected = onRejected;
					},
				},
			},
		}),
	},
}));

vi.mock("vant", () => ({
	Toast: mockState.toastMock,
}));

import "../../src/utils/axios.js";

const createLocalStorageMock = (initial = {}) => {
	const store = { ...initial };
	return {
		getItem: vi.fn((key) => (key in store ? store[key] : null)),
		setItem: vi.fn((key, value) => {
			store[key] = String(value);
		}),
		removeItem: vi.fn((key) => {
			delete store[key];
		}),
		clear: vi.fn(() => {
			Object.keys(store).forEach((key) => {
				delete store[key];
			});
		}),
	};
};

describe("axios interceptors", () => {
	beforeEach(() => {
		mockState.toastMock.mockReset();
		global.localStorage = createLocalStorageMock();
		global.window = {
			location: {
				pathname: "/brand/dashboard",
				href: "",
			},
		};
	});

	it("injects Bearer token in request headers", () => {
		global.localStorage = createLocalStorageMock({ dmh_token: "token-123" });

		const config = mockState.requestOnFulfilled({ headers: {} });

		expect(config.headers.Authorization).toBe("Bearer token-123");
	});

	it("rejects non-200 business response and shows toast", async () => {
		await expect(mockState.responseOnFulfilled({ data: { code: 500, message: "失败" } })).rejects.toThrow("失败");
		expect(mockState.toastMock).toHaveBeenCalledWith(
			expect.objectContaining({
				message: "失败",
				type: "fail",
			}),
		);
	});

	it("returns res when code is 200", async () => {
		const res = { data: { code: 200, data: "success" } };
		const result = await mockState.responseOnFulfilled(res);
		expect(result).toBe(res.data);
	});

	it("returns res when code is undefined", async () => {
		const res = { data: { data: "success" } };
		const result = await mockState.responseOnFulfilled(res);
		expect(result).toBe(res.data);
	});

	it("shows toast for non-401 errors with response", async () => {
		const err = { response: { status: 500, data: { message: "Server error" } }, message: "Server error" };
		await expect(mockState.responseOnRejected(err)).rejects.toBe(err);
		expect(mockState.toastMock).toHaveBeenCalledWith(
			expect.objectContaining({
				message: "Server error",
				type: "fail",
			}),
		);
	});

	it("shows toast for non-401 errors without message", async () => {
		const err = { response: { status: 500, data: {} }, message: "Server error" };
		await expect(mockState.responseOnRejected(err)).rejects.toBe(err);
		expect(mockState.toastMock).toHaveBeenCalledWith(
			expect.objectContaining({
				message: "请求失败",
				type: "fail",
			}),
		);
	});

	it("redirects to distributor login on 401 for distributor path", async () => {
		global.localStorage = createLocalStorageMock({
			dmh_token: "token-123",
			dmh_user_role: "distributor",
			dmh_user_info: "{}",
		});
		global.window.location.pathname = "/distributor/rewards";

		const err = { response: { status: 401, data: {} }, message: "Unauthorized" };
		await expect(mockState.responseOnRejected(err)).rejects.toBe(err);

		expect(global.localStorage.removeItem).toHaveBeenCalledWith("dmh_token");
		expect(global.localStorage.removeItem).toHaveBeenCalledWith("dmh_user_role");
		expect(global.localStorage.removeItem).toHaveBeenCalledWith("dmh_user_info");
		expect(global.window.location.href).toBe("/distributor/login");
	});

	it("redirects to brand login on 401 for brand user", async () => {
		global.localStorage = createLocalStorageMock({
			dmh_token: "token-123",
			dmh_user_role: "brand_admin",
			dmh_user_info: "{}",
		});
		global.window.location.pathname = "/brand/orders";

		const err = { response: { status: 401, data: {} }, message: "Unauthorized" };
		await expect(mockState.responseOnRejected(err)).rejects.toBe(err);
		expect(global.window.location.href).toBe("/brand/login");
	});

	it("shows timeout and network toasts", async () => {
		await expect(mockState.responseOnRejected({ message: "timeout of 15000ms exceeded" })).rejects.toBeTruthy();
		expect(mockState.toastMock).toHaveBeenCalledWith(
			expect.objectContaining({
				message: "请求超时",
				type: "fail",
			}),
		);

		mockState.toastMock.mockReset();
		await expect(mockState.responseOnRejected({ message: "Network Error" })).rejects.toBeTruthy();
		expect(mockState.toastMock).toHaveBeenCalledWith(
			expect.objectContaining({
				message: "网络连接失败",
				type: "fail",
			}),
		);
	});

	it("passes through request interceptor rejection", async () => {
		const err = new Error("request failed");
		await expect(mockState.requestOnRejected(err)).rejects.toBe(err);
	});
});
