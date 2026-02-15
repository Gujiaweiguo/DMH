export const resolveGuardNavigation = (to, token, userRole, userInfoRaw) => {
	if (!to.meta.requiresAuth) {
		return null;
	}

	if (!token) {
		if (to.path.startsWith("/distributor")) {
			return "/distributor/login";
		}
		return "/brand/login";
	}

	if (to.meta.role && userRole !== to.meta.role) {
		if (to.meta.role === "distributor") {
			return "/distributor/login";
		}
		if (to.meta.role === "brand_admin") {
			return "/brand/login";
		}
		return "/";
	}

	if (to.meta.hasBrand) {
		let userInfo = {};
		try {
			userInfo = JSON.parse(userInfoRaw || "{}");
		} catch (_error) {
			userInfo = {};
		}

		const brandIds = userInfo.brandIds || [];
		if (!Array.isArray(brandIds) || brandIds.length === 0) {
			return "/brand/login";
		}
	}

	return null;
};
