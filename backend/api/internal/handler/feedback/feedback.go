package feedback

import (
	"net/http"
	"strconv"

	"dmh/api/internal/logic/feedback"
	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type CreateFeedbackHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCreateFeedbackHandler(svcCtx *svc.ServiceContext) *CreateFeedbackHandler {
	return &CreateFeedbackHandler{
		svcCtx: svcCtx,
	}
}

func (h *CreateFeedbackHandler) CreateFeedback(w http.ResponseWriter, r *http.Request) {
	var req types.CreateFeedbackReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 从JWT获取用户ID（解析失败时回退为 0）
	userId := getUserIDFromRequest(r)

	l := feedback.NewCreateFeedbackLogic(r.Context(), h.svcCtx)
	resp, err := l.CreateFeedback(&req, userId)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

type ListFeedbackHandler struct {
	svcCtx *svc.ServiceContext
}

func NewListFeedbackHandler(svcCtx *svc.ServiceContext) *ListFeedbackHandler {
	return &ListFeedbackHandler{
		svcCtx: svcCtx,
	}
}

func (h *ListFeedbackHandler) ListFeedback(w http.ResponseWriter, r *http.Request) {
	var req types.ListFeedbackReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 从JWT获取用户ID和角色（解析失败时回退为 0/空字符串）
	userId := getUserIDFromRequest(r)
	userRole := getUserRoleFromRequest(r)

	l := feedback.NewListFeedbackLogic(r.Context(), h.svcCtx)
	resp, err := l.ListFeedback(&req, userId, userRole)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

type GetFeedbackHandler struct {
	svcCtx *svc.ServiceContext
}

func NewGetFeedbackHandler(svcCtx *svc.ServiceContext) *GetFeedbackHandler {
	return &GetFeedbackHandler{
		svcCtx: svcCtx,
	}
}

func (h *GetFeedbackHandler) GetFeedback(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 从JWT获取用户ID和角色
	userId := getUserIDFromRequest(r)
	userRole := getUserRoleFromRequest(r)

	l := feedback.NewGetFeedbackLogic(r.Context(), h.svcCtx)
	resp, err := l.GetFeedback(id, userId, userRole)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

type UpdateFeedbackStatusHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUpdateFeedbackStatusHandler(svcCtx *svc.ServiceContext) *UpdateFeedbackStatusHandler {
	return &UpdateFeedbackStatusHandler{
		svcCtx: svcCtx,
	}
}

func (h *UpdateFeedbackStatusHandler) UpdateFeedbackStatus(w http.ResponseWriter, r *http.Request) {
	var req types.UpdateFeedbackStatusReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 从JWT获取用户ID和角色（角色用于管理员权限判断）
	userId := getUserIDFromRequest(r)
	userRole := getUserRoleFromRequest(r)

	l := feedback.NewUpdateFeedbackStatusLogic(r.Context(), h.svcCtx)
	resp, err := l.UpdateFeedbackStatus(&req, userId, userRole)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

type SubmitSatisfactionSurveyHandler struct {
	svcCtx *svc.ServiceContext
}

func NewSubmitSatisfactionSurveyHandler(svcCtx *svc.ServiceContext) *SubmitSatisfactionSurveyHandler {
	return &SubmitSatisfactionSurveyHandler{
		svcCtx: svcCtx,
	}
}

func (h *SubmitSatisfactionSurveyHandler) SubmitSatisfactionSurvey(w http.ResponseWriter, r *http.Request) {
	var req types.SubmitSatisfactionSurveyReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 从JWT获取用户ID和角色
	userId := getUserIDFromRequest(r)
	userRole := getUserRoleFromRequest(r)

	l := feedback.NewSubmitSatisfactionSurveyLogic(r.Context(), h.svcCtx)
	resp, err := l.SubmitSatisfactionSurvey(&req, userId, userRole)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

type ListFAQHandler struct {
	svcCtx *svc.ServiceContext
}

func NewListFAQHandler(svcCtx *svc.ServiceContext) *ListFAQHandler {
	return &ListFAQHandler{
		svcCtx: svcCtx,
	}
}

func (h *ListFAQHandler) ListFAQ(w http.ResponseWriter, r *http.Request) {
	var req types.ListFAQReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	l := feedback.NewListFAQLogic(r.Context(), h.svcCtx)
	resp, err := l.ListFAQ(&req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

type MarkFAQHelpfulHandler struct {
	svcCtx *svc.ServiceContext
}

func NewMarkFAQHelpfulHandler(svcCtx *svc.ServiceContext) *MarkFAQHelpfulHandler {
	return &MarkFAQHelpfulHandler{
		svcCtx: svcCtx,
	}
}

func (h *MarkFAQHelpfulHandler) MarkFAQHelpful(w http.ResponseWriter, r *http.Request) {
	var req types.MarkFAQHelpfulReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	l := feedback.NewMarkFAQHelpfulLogic(r.Context(), h.svcCtx)
	resp, err := l.MarkFAQHelpful(&req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

type RecordFeatureUsageHandler struct {
	svcCtx *svc.ServiceContext
}

func NewRecordFeatureUsageHandler(svcCtx *svc.ServiceContext) *RecordFeatureUsageHandler {
	return &RecordFeatureUsageHandler{
		svcCtx: svcCtx,
	}
}

func (h *RecordFeatureUsageHandler) RecordFeatureUsage(w http.ResponseWriter, r *http.Request) {
	var req types.RecordFeatureUsageReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 从JWT获取用户ID和角色
	userId := getUserIDFromRequest(r)
	userRole := getUserRoleFromRequest(r)

	l := feedback.NewRecordFeatureUsageLogic(r.Context(), h.svcCtx)
	resp, err := l.RecordFeatureUsage(&req, userId, userRole)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

type GetFeedbackStatisticsHandler struct {
	svcCtx *svc.ServiceContext
}

func NewGetFeedbackStatisticsHandler(svcCtx *svc.ServiceContext) *GetFeedbackStatisticsHandler {
	return &GetFeedbackStatisticsHandler{
		svcCtx: svcCtx,
	}
}

func (h *GetFeedbackStatisticsHandler) GetFeedbackStatistics(w http.ResponseWriter, r *http.Request) {
	var req types.GetFeedbackStatisticsReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 从JWT获取用户角色（仅管理员可访问）
	userRole := getUserRoleFromRequest(r)

	l := feedback.NewGetFeedbackStatisticsLogic(r.Context(), h.svcCtx)
	resp, err := l.GetFeedbackStatistics(&req, userRole)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func getUserIDFromRequest(r *http.Request) int64 {
	userId, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		return 0
	}
	return userId
}

func getUserRoleFromRequest(r *http.Request) string {
	roles, err := middleware.GetUserRolesFromContext(r.Context())
	if err == nil && len(roles) > 0 {
		if hasRole(roles, "platform_admin") {
			return "platform_admin"
		}
		if hasRole(roles, "brand_admin") {
			return "brand_admin"
		}
		return roles[0]
	}

	// 兼容测试中通过 userRole 注入的上下文
	if v := r.Context().Value("userRole"); v != nil {
		if role, ok := v.(string); ok {
			return role
		}
	}

	return ""
}

func hasRole(roles []string, target string) bool {
	for _, role := range roles {
		if role == target {
			return true
		}
	}
	return false
}
