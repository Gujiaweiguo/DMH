package feedback

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFeedbackLogic {
	return &CreateFeedbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateFeedbackLogic) CreateFeedback(req *types.CreateFeedbackReq, userId int64) (*types.FeedbackResp, error) {
	// 参数验证
	if req.Title == "" || req.Content == "" {
		return nil, fmt.Errorf("标题和内容不能为空")
	}

	if req.Category == "" {
		return nil, fmt.Errorf("反馈类别不能为空")
	}

	// 验证评分范围
	if req.Rating != nil && (*req.Rating < 1 || *req.Rating > 5) {
		return nil, fmt.Errorf("评分必须在1-5之间")
	}

	// 设置默认优先级
	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	feedback := &model.UserFeedback{
		UserID:         userId,
		Category:       req.Category,
		Subcategory:    req.Subcategory,
		Rating:         req.Rating,
		Title:          req.Title,
		Content:        req.Content,
		FeatureUseCase: req.FeatureUseCase,
		DeviceInfo:     req.DeviceInfo,
		BrowserInfo:    req.BrowserInfo,
		Priority:       priority,
		Status:         "pending",
	}

	err := l.svcCtx.DB.Create(feedback).Error
	if err != nil {
		return nil, err
	}

	return &types.FeedbackResp{
		Id:             feedback.ID,
		UserId:         feedback.UserID,
		UserName:       "",
		UserRole:       "",
		Category:       feedback.Category,
		Subcategory:    feedback.Subcategory,
		Rating:         feedback.Rating,
		Title:          feedback.Title,
		Content:        feedback.Content,
		FeatureUseCase: feedback.FeatureUseCase,
		DeviceInfo:     feedback.DeviceInfo,
		BrowserInfo:    feedback.BrowserInfo,
		Priority:       feedback.Priority,
		Status:         feedback.Status,
		AssigneeId:     feedback.AssigneeID,
		Response:       feedback.Response,
		ResolvedAt:     feedback.ResolvedAt,
		CreatedAt:      feedback.CreatedAt,
		Tags:           []types.TagResp{},
	}, nil
}

type ListFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFeedbackLogic {
	return &ListFeedbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListFeedbackLogic) ListFeedback(req *types.ListFeedbackReq, userId int64, userRole string) (*types.FeedbackListResp, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	var feedbacks []model.UserFeedback
	var total int64

	query := l.svcCtx.DB.Model(&model.UserFeedback{})

	// 权限控制：普通用户只能查看自己的反馈，管理员可以查看所有
	if userRole != "platform_admin" && userRole != "brand_admin" {
		query = query.Where("user_id = ?", userId)
	}

	// 筛选条件
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Priority != "" {
		query = query.Where("priority = ?", req.Priority)
	}

	// 查询总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err = query.Preload("User").Preload("Assignee").Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&feedbacks).Error
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	feedbackResps := make([]types.FeedbackResp, len(feedbacks))
	for i, f := range feedbacks {
		userName := ""
		userRoleStr := ""
		if f.User != nil {
			userName = f.User.Username
			userRoleStr = getRoleName(f.User.Id)
		}

		feedbackResps[i] = types.FeedbackResp{
			Id:             f.ID,
			UserId:         f.UserID,
			UserName:       userName,
			UserRole:       userRoleStr,
			Category:       f.Category,
			Subcategory:    f.Subcategory,
			Rating:         f.Rating,
			Title:          f.Title,
			Content:        f.Content,
			FeatureUseCase: f.FeatureUseCase,
			DeviceInfo:     f.DeviceInfo,
			BrowserInfo:    f.BrowserInfo,
			Priority:       f.Priority,
			Status:         f.Status,
			AssigneeId:     f.AssigneeID,
			Response:       f.Response,
			ResolvedAt:     f.ResolvedAt,
			CreatedAt:      f.CreatedAt,
			Tags:           []types.TagResp{},
		}
	}

	return &types.FeedbackListResp{
		Total:     total,
		Feedbacks: feedbackResps,
	}, nil
}

// 辅助函数：根据用户ID获取角色名称
func getRoleName(userId int64) string {
	return "" // 暂时返回空，后续可以查询数据库获取角色
}

type GetFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedbackLogic {
	return &GetFeedbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFeedbackLogic) GetFeedback(id int64, userId int64, userRole string) (*types.FeedbackResp, error) {
	var feedback model.UserFeedback
	err := l.svcCtx.DB.Preload("User").Preload("Assignee").First(&feedback, id).Error
	if err != nil {
		return nil, err
	}

	// 权限控制
	if userRole != "platform_admin" && userRole != "brand_admin" && feedback.UserID != userId {
		return nil, fmt.Errorf("无权访问此反馈")
	}

	// 加载标签
	var tags []model.FeedbackTag
	err = l.svcCtx.DB.Raw(`
		SELECT ft.* FROM feedback_tags ft
		INNER JOIN feedback_tag_relations ftr ON ft.id = ftr.tag_id
		WHERE ftr.feedback_id = ?
	`, id).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	tagResps := make([]types.TagResp, len(tags))
	for i, t := range tags {
		tagResps[i] = types.TagResp{
			Id:    t.ID,
			Name:  t.Name,
			Color: t.Color,
		}
	}

	userName := ""
	userRoleStr := ""
	if feedback.User != nil {
		userName = feedback.User.Username
		userRoleStr = getRoleName(feedback.User.Id)
	}

	return &types.FeedbackResp{
		Id:             feedback.ID,
		UserId:         feedback.UserID,
		UserName:       userName,
		UserRole:       userRoleStr,
		Category:       feedback.Category,
		Subcategory:    feedback.Subcategory,
		Rating:         feedback.Rating,
		Title:          feedback.Title,
		Content:        feedback.Content,
		FeatureUseCase: feedback.FeatureUseCase,
		DeviceInfo:     feedback.DeviceInfo,
		BrowserInfo:    feedback.BrowserInfo,
		Priority:       feedback.Priority,
		Status:         feedback.Status,
		AssigneeId:     feedback.AssigneeID,
		Response:       feedback.Response,
		ResolvedAt:     feedback.ResolvedAt,
		CreatedAt:      feedback.CreatedAt,
		Tags:           tagResps,
	}, nil
}

type UpdateFeedbackStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateFeedbackStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFeedbackStatusLogic {
	return &UpdateFeedbackStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFeedbackStatusLogic) UpdateFeedbackStatus(req *types.UpdateFeedbackStatusReq, userId int64, userRole string) (*types.FeedbackResp, error) {
	// 权限验证：只有管理员可以更新反馈状态
	if userRole != "platform_admin" && userRole != "brand_admin" {
		return nil, fmt.Errorf("无权更新反馈状态")
	}

	var feedback model.UserFeedback
	err := l.svcCtx.DB.First(&feedback, req.Id).Error
	if err != nil {
		return nil, err
	}

	// 更新字段
	updates := map[string]interface{}{
		"status": req.Status,
	}

	if req.AssigneeId != nil {
		updates["assignee_id"] = *req.AssigneeId
	}

	if req.Response != "" {
		updates["response"] = req.Response
	}

	// 如果状态为已解决，设置解决时间
	if req.Status == "resolved" {
		now := time.Now()
		updates["resolved_at"] = &now
	}

	err = l.svcCtx.DB.Model(&feedback).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	// 刷新数据
	err = l.svcCtx.DB.Preload("User").Preload("Assignee").First(&feedback, req.Id).Error
	if err != nil {
		return nil, err
	}

	return &types.FeedbackResp{
		Id:             feedback.ID,
		UserId:         feedback.UserID,
		UserName:       "",
		UserRole:       "",
		Category:       feedback.Category,
		Subcategory:    feedback.Subcategory,
		Rating:         feedback.Rating,
		Title:          feedback.Title,
		Content:        feedback.Content,
		FeatureUseCase: feedback.FeatureUseCase,
		DeviceInfo:     feedback.DeviceInfo,
		BrowserInfo:    feedback.BrowserInfo,
		Priority:       feedback.Priority,
		Status:         feedback.Status,
		AssigneeId:     feedback.AssigneeID,
		Response:       feedback.Response,
		ResolvedAt:     feedback.ResolvedAt,
		CreatedAt:      feedback.CreatedAt,
		Tags:           []types.TagResp{},
	}, nil
}

type SubmitSatisfactionSurveyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitSatisfactionSurveyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitSatisfactionSurveyLogic {
	return &SubmitSatisfactionSurveyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitSatisfactionSurveyLogic) SubmitSatisfactionSurvey(req *types.SubmitSatisfactionSurveyReq, userId int64, userRole string) (*types.SatisfactionSurveyResp, error) {
	// 参数验证
	if req.Feature == "" {
		return nil, fmt.Errorf("功能名称不能为空")
	}

	// 验证评分范围
	if req.EaseOfUse != nil && (*req.EaseOfUse < 1 || *req.EaseOfUse > 5) {
		return nil, fmt.Errorf("易用性评分必须在1-5之间")
	}
	if req.Performance != nil && (*req.Performance < 1 || *req.Performance > 5) {
		return nil, fmt.Errorf("性能评分必须在1-5之间")
	}
	if req.Reliability != nil && (*req.Reliability < 1 || *req.Reliability > 5) {
		return nil, fmt.Errorf("稳定性评分必须在1-5之间")
	}
	if req.OverallSatisfaction != nil && (*req.OverallSatisfaction < 1 || *req.OverallSatisfaction > 5) {
		return nil, fmt.Errorf("整体满意度必须在1-5之间")
	}
	if req.WouldRecommend != nil && (*req.WouldRecommend < 1 || *req.WouldRecommend > 5) {
		return nil, fmt.Errorf("推荐意愿必须在1-5之间")
	}

	survey := &model.FeatureSatisfactionSurvey{
		UserID:                 userId,
		UserRole:               userRole,
		Feature:                req.Feature,
		EaseOfUse:              req.EaseOfUse,
		Performance:            req.Performance,
		Reliability:            req.Reliability,
		OverallSatisfaction:    req.OverallSatisfaction,
		WouldRecommend:         req.WouldRecommend,
		MostLiked:              req.MostLiked,
		LeastLiked:             req.LeastLiked,
		ImprovementSuggestions: req.ImprovementSuggestions,
		WouldLikeMoreFeatures:  req.WouldLikeMoreFeatures,
	}

	err := l.svcCtx.DB.Create(survey).Error
	if err != nil {
		return nil, err
	}

	return &types.SatisfactionSurveyResp{
		Id:                     survey.ID,
		UserId:                 survey.UserID,
		UserRole:               survey.UserRole,
		Feature:                survey.Feature,
		EaseOfUse:              survey.EaseOfUse,
		Performance:            survey.Performance,
		Reliability:            survey.Reliability,
		OverallSatisfaction:    survey.OverallSatisfaction,
		WouldRecommend:         survey.WouldRecommend,
		MostLiked:              survey.MostLiked,
		LeastLiked:             survey.LeastLiked,
		ImprovementSuggestions: survey.ImprovementSuggestions,
		WouldLikeMoreFeatures:  survey.WouldLikeMoreFeatures,
		CreatedAt:              survey.CreatedAt,
	}, nil
}

type ListFAQLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListFAQLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFAQLogic {
	return &ListFAQLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListFAQLogic) ListFAQ(req *types.ListFAQReq) (*types.FAQListResp, error) {
	var faqs []model.FAQItem
	var total int64

	query := l.svcCtx.DB.Model(&model.FAQItem{}).Where("is_published = ?", true)

	// 筛选条件
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Keyword != "" {
		query = query.Where("question LIKE ? OR answer LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 查询总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 查询并按排序字段排序
	err = query.Order("sort_order ASC, created_at DESC").Find(&faqs).Error
	if err != nil {
		return nil, err
	}

	faqResps := make([]types.FAQResp, len(faqs))
	for i, f := range faqs {
		faqResps[i] = types.FAQResp{
			Id:              f.ID,
			Category:        f.Category,
			Question:        f.Question,
			Answer:          f.Answer,
			SortOrder:       f.SortOrder,
			ViewCount:       f.ViewCount,
			HelpfulCount:    f.HelpfulCount,
			NotHelpfulCount: f.NotHelpfulCount,
		}
	}

	return &types.FAQListResp{
		Total: total,
		FAQs:  faqResps,
	}, nil
}

type MarkFAQHelpfulLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarkFAQHelpfulLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkFAQHelpfulLogic {
	return &MarkFAQHelpfulLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkFAQHelpfulLogic) MarkFAQHelpful(req *types.MarkFAQHelpfulReq) (*types.FAQResp, error) {
	helpfulColumn, notHelpfulColumn, err := resolveFAQCounterColumns(l.svcCtx.DB)
	if err != nil {
		return nil, err
	}

	counterColumn := ""
	switch req.Type {
	case "helpful":
		if helpfulColumn == "" {
			return nil, fmt.Errorf("faq_items 缺少 helpful 计数字段")
		}
		counterColumn = helpfulColumn
	case "not_helpful":
		if notHelpfulColumn == "" {
			return nil, fmt.Errorf("faq_items 缺少 not_helpful 计数字段")
		}
		counterColumn = notHelpfulColumn
	default:
		return nil, fmt.Errorf("无效的类型: %s", req.Type)
	}

	if err := incrementFAQCounter(l.svcCtx.DB, req.Id, counterColumn); err != nil {
		return nil, err
	}

	faq, err := queryFAQWithCounterAliases(l.svcCtx.DB, req.Id, helpfulColumn, notHelpfulColumn)
	if err != nil {
		return nil, err
	}

	return &types.FAQResp{
		Id:              faq.ID,
		Category:        faq.Category,
		Question:        faq.Question,
		Answer:          faq.Answer,
		SortOrder:       faq.SortOrder,
		ViewCount:       faq.ViewCount,
		HelpfulCount:    faq.HelpfulCount,
		NotHelpfulCount: faq.NotHelpfulCount,
	}, nil
}

func resolveFAQCounterColumns(db *gorm.DB) (helpfulColumn string, notHelpfulColumn string, err error) {
	columns, err := listFAQTableColumns(db)
	if err != nil {
		return "", "", err
	}

	helpfulColumn = pickFAQCounterColumn(columns, []string{"helpful_count", "helpfulCount", "helpful", "helpul_count", "helpulCount"})
	notHelpfulColumn = pickFAQCounterColumn(columns, []string{"not_helpful_count", "notHelpfulCount", "not_helpful", "notHelpful"})

	if helpfulColumn == "" && notHelpfulColumn == "" {
		return "", "", fmt.Errorf("faq_items 计数字段缺失: columns=%v", columns)
	}

	return helpfulColumn, notHelpfulColumn, nil
}

func listFAQTableColumns(db *gorm.DB) ([]string, error) {
	type columnRow struct {
		Name string `gorm:"column:name"`
	}

	var rows []columnRow
	var err error

	if db.Dialector.Name() == "sqlite" {
		err = db.Raw("PRAGMA table_info(faq_items)").Scan(&rows).Error
	} else {
		err = db.Raw("SELECT COLUMN_NAME AS name FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'faq_items'").Scan(&rows).Error
	}
	if err != nil {
		return nil, err
	}

	columns := make([]string, 0, len(rows))
	for _, row := range rows {
		if row.Name != "" {
			columns = append(columns, row.Name)
		}
	}

	return columns, nil
}

func pickFAQCounterColumn(columns []string, candidates []string) string {
	for _, candidate := range candidates {
		for _, column := range columns {
			if strings.EqualFold(column, candidate) {
				return column
			}
		}
	}

	return ""
}

func incrementFAQCounter(db *gorm.DB, faqID int64, column string) error {
	if column != "helpful_count" && column != "helpfulCount" && column != "helpful" && column != "helpul_count" && column != "helpulCount" &&
		column != "not_helpful_count" && column != "notHelpfulCount" && column != "not_helpful" && column != "notHelpful" {
		return fmt.Errorf("不支持的计数字段: %s", column)
	}

	updateSQL := fmt.Sprintf("UPDATE faq_items SET `%s` = COALESCE(`%s`, 0) + 1 WHERE id = ?", column, column)
	result := db.Exec(updateSQL, faqID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func queryFAQWithCounterAliases(db *gorm.DB, faqID int64, helpfulColumn string, notHelpfulColumn string) (*model.FAQItem, error) {
	helpfulExpr := "0 AS helpful_count"
	if helpfulColumn != "" {
		helpfulExpr = fmt.Sprintf("`%s` AS helpful_count", helpfulColumn)
	}

	notHelpfulExpr := "0 AS not_helpful_count"
	if notHelpfulColumn != "" {
		notHelpfulExpr = fmt.Sprintf("`%s` AS not_helpful_count", notHelpfulColumn)
	}

	selectSQL := fmt.Sprintf("id, category, question, answer, sort_order, view_count, %s, %s", helpfulExpr, notHelpfulExpr)

	var faq model.FAQItem
	if err := db.Table("faq_items").Select(selectSQL).Where("id = ?", faqID).Take(&faq).Error; err != nil {
		return nil, err
	}

	return &faq, nil
}

type RecordFeatureUsageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecordFeatureUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordFeatureUsageLogic {
	return &RecordFeatureUsageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecordFeatureUsageLogic) RecordFeatureUsage(req *types.RecordFeatureUsageReq, userId int64, userRole string) (*types.BaseResp, error) {
	usage := &model.FeatureUsageStat{
		UserID:       userId,
		UserRole:     userRole,
		Feature:      req.Feature,
		Action:       req.Action,
		CampaignID:   req.CampaignId,
		Success:      req.Success,
		DurationMs:   req.DurationMs,
		ErrorMessage: req.ErrorMessage,
	}

	err := l.svcCtx.DB.Create(usage).Error
	if err != nil {
		return nil, err
	}

	return &types.BaseResp{
		Code:    200,
		Message: "使用记录已保存",
	}, nil
}

type GetFeedbackStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFeedbackStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedbackStatisticsLogic {
	return &GetFeedbackStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFeedbackStatisticsLogic) GetFeedbackStatistics(req *types.GetFeedbackStatisticsReq, userRole string) (*types.FeedbackStatisticsResp, error) {
	// 权限验证：只有管理员可以查看统计数据
	if userRole != "platform_admin" && userRole != "brand_admin" {
		return nil, fmt.Errorf("无权查看统计数据")
	}

	query := l.svcCtx.DB.Model(&model.UserFeedback{})

	// 时间范围筛选
	if req.StartDate != "" {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate)
	}
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	var totalFeedbacks int64
	err := query.Count(&totalFeedbacks).Error
	if err != nil {
		return nil, err
	}

	// 按类别统计
	type CategoryStat struct {
		Category string
		Count    int64
	}
	var categoryStats []CategoryStat
	err = query.Select("category, COUNT(*) as count").Group("category").Scan(&categoryStats).Error
	if err != nil {
		return nil, err
	}
	byCategory := make(map[string]int64)
	for _, stat := range categoryStats {
		byCategory[stat.Category] = stat.Count
	}

	// 按状态统计
	type StatusStat struct {
		Status string
		Count  int64
	}
	var statusStats []StatusStat
	err = query.Select("status, COUNT(*) as count").Group("status").Scan(&statusStats).Error
	if err != nil {
		return nil, err
	}
	byStatus := make(map[string]int64)
	for _, stat := range statusStats {
		byStatus[stat.Status] = stat.Count
	}

	// 按优先级统计
	type PriorityStat struct {
		Priority string
		Count    int64
	}
	var priorityStats []PriorityStat
	err = query.Select("priority, COUNT(*) as count").Group("priority").Scan(&priorityStats).Error
	if err != nil {
		return nil, err
	}
	byPriority := make(map[string]int64)
	for _, stat := range priorityStats {
		byPriority[stat.Priority] = stat.Count
	}

	// 计算平均评分
	type RatingStat struct {
		AvgRating float64
	}
	var ratingStat RatingStat
	err = query.Select("AVG(rating) as avg_rating").Where("rating IS NOT NULL").Scan(&ratingStat).Error
	if err != nil {
		return nil, err
	}
	averageRating := ratingStat.AvgRating

	// 计算解决率
	var resolvedCount int64
	baseQuery := l.svcCtx.DB.Model(&model.UserFeedback{})
	if req.StartDate != "" {
		baseQuery = baseQuery.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		baseQuery = baseQuery.Where("created_at <= ?", req.EndDate)
	}
	if req.Category != "" {
		baseQuery = baseQuery.Where("category = ?", req.Category)
	}
	baseQuery.Where("status = ?", "resolved").Count(&resolvedCount)
	resolutionRate := 0.0
	if totalFeedbacks > 0 {
		resolutionRate = float64(resolvedCount) / float64(totalFeedbacks)
	}

	// 计算平均解决时间（小时）
	type ResolutionTime struct {
		ResolutionHours float64
	}
	var resolutionTime ResolutionTime
	resolutionHoursExpr := "AVG(ROUND((julianday(resolved_at) - julianday(created_at)) * 24, 2)) as resolution_hours"
	if l.svcCtx.DB.Dialector.Name() != "sqlite" {
		resolutionHoursExpr = "AVG(ROUND(TIMESTAMPDIFF(SECOND, created_at, resolved_at) / 3600, 2)) as resolution_hours"
	}
	err = baseQuery.Model(&model.UserFeedback{}).
		Select(resolutionHoursExpr).
		Where("status = ? AND resolved_at IS NOT NULL", "resolved").
		Scan(&resolutionTime).Error
	avgResolutionTime := 0.0
	if err == nil {
		avgResolutionTime = resolutionTime.ResolutionHours
	}

	// 按评分分布统计
	type RatingDist struct {
		Rating int
		Count  int64
	}
	var ratingDists []RatingDist
	err = query.Select("rating, COUNT(*) as count").Where("rating IS NOT NULL").Group("rating").Scan(&ratingDists).Error
	if err != nil {
		return nil, err
	}
	byRating := make(map[int]int64)
	for _, dist := range ratingDists {
		byRating[dist.Rating] = dist.Count
	}

	return &types.FeedbackStatisticsResp{
		TotalFeedbacks:    totalFeedbacks,
		ByCategory:        byCategory,
		ByStatus:          byStatus,
		ByPriority:        byPriority,
		AverageRating:     averageRating,
		ResolutionRate:    resolutionRate,
		AvgResolutionTime: avgResolutionTime,
		ByRating:          byRating,
	}, nil
}
