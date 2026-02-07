import axios from "axios";
import { Toast } from "vant";

// 创建axios实例
const service = axios.create({
	baseURL: import.meta.env.VITE_API_BASE_URL || "/api/v1",
	timeout: 15000,
});

// 请求拦截器
service.interceptors.request.use(
	(config) => {
		// 从localStorage获取token
		const token = localStorage.getItem("dmh_token");
		if (token) {
			config.headers.Authorization = `Bearer ${token}`;
		}
		return config;
	},
	(error) => {
		console.error("Request error:", error);
		return Promise.reject(error);
	},
);

// 响应拦截器
service.interceptors.response.use(
	(response) => {
		const res = response.data;

		// 如果返回的状态码不是200，说明接口有问题
		if (res.code !== undefined && res.code !== 200) {
			Toast({
				message: res.message || "请求失败",
				type: "fail",
				duration: 2000,
			});
			return Promise.reject(new Error(res.message || "Error"));
		} else {
			return res;
		}
	},
	(error) => {
		console.error("Response error:", error);

		if (error.response) {
			const { status, data } = error.response;

			// 401 未授权，跳转登录
			if (status === 401) {
				const currentRole = localStorage.getItem("dmh_user_role");
				const isDistributorPath =
					window.location.pathname.startsWith("/distributor");
				localStorage.removeItem("dmh_token");
				localStorage.removeItem("dmh_user_role");
				localStorage.removeItem("dmh_user_info");
				if (currentRole === "distributor" || isDistributorPath) {
					window.location.href = "/distributor/login";
				} else {
					window.location.href = "/brand/login";
				}
				return Promise.reject(error);
			}

			// 其他错误
			const message = data?.message || "请求失败";
			Toast({
				message: message,
				type: "fail",
				duration: 2000,
			});
		} else if (error.message.includes("timeout")) {
			Toast({
				message: "请求超时",
				type: "fail",
				duration: 2000,
			});
		} else if (error.message.includes("Network")) {
			Toast({
				message: "网络连接失败",
				type: "fail",
				duration: 2000,
			});
		}

		return Promise.reject(error);
	},
);

export default service;
