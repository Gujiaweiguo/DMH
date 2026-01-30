package poster

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"time"

	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
)

// Service 海报生成服务
type Service struct {
	posterDir string // 海报存储目录
	baseURL   string // 海报访问URL前缀
}

// NewService 创建海报服务
func NewService(posterDir, baseURL string) *Service {
	// 确保目录存在
	if err := os.MkdirAll(posterDir, 0755); err != nil {
		panic(fmt.Sprintf("创建海报目录失败: %v", err))
	}

	return &Service{
		posterDir: posterDir,
		baseURL:   baseURL,
	}
}

// GenerateCampaignPoster 生成活动专属海报
func (s *Service) GenerateCampaignPoster(campaignName, campaignDesc, distributorName, qrcodeData string) (string, error) {
	fmt.Printf("[PosterService] GenerateCampaignPoster called: name=%s, desc=%s, distributor=%s\n", campaignName, campaignDesc, distributorName)

	// 1. 创建画布（海报尺寸：750x1334 px，即微信朋友圈图片尺寸）
	width := 750
	height := 1334
	dc := gg.NewContext(width, height)

	dc.SetRGB(66, 133, 244)
	dc.Clear()

	// 使用纯色背景（渐变需要更复杂的处理）
	dc.SetColor(color.RGBA{R: 66, G: 133, B: 244, A: 255})
	dc.Clear()

	// 3. 绘制标题
	if err := s.drawTitle(dc, campaignName); err != nil {
		return "", err
	}

	// 4. 绘制描述
	if err := s.drawDescription(dc, campaignDesc); err != nil {
		return "", err
	}

	// 5. 绘制分销商信息
	if err := s.drawDistributorInfo(dc, distributorName); err != nil {
		return "", err
	}

	// 6. 生成二维码
	qrcodePath, err := s.generateQRCode(qrcodeData)
	if err != nil {
		return "", err
	}

	// 7. 绘制二维码
	if err := s.drawQRCode(dc, qrcodePath); err != nil {
		return "", err
	}

	// 8. 绘制底部提示文字
	if err := s.drawFooter(dc); err != nil {
		return "", err
	}

	// 9. 保存图片
	filename := fmt.Sprintf("poster_%d.png", generateUniqueID())
	filepath := filepath.Join(s.posterDir, filename)
	fmt.Printf("[PosterService] Saving poster to: %s\n", filepath)
	if err := dc.SavePNG(filepath); err != nil {
		fmt.Printf("[PosterService] Failed to save poster: %v\n", err)
		return "", fmt.Errorf("保存海报失败: %w", err)
	}
	fmt.Printf("[PosterService] Poster saved successfully\n")

	// 10. 清理临时二维码文件
	os.Remove(qrcodePath)

	// 11. 返回URL
	posterURL := fmt.Sprintf("%s/posters/%s", s.baseURL, filename)
	return posterURL, nil
}

// GenerateDistributorPoster 生成通用分销商海报
func (s *Service) GenerateDistributorPoster(distributorName string, campaignCount int) (string, error) {
	// 1. 创建画布
	width := 750
	height := 1334
	dc := gg.NewContext(width, height)

	// 2. 绘制背景
	dc.SetColor(color.RGBA{R: 99, G: 102, B: 241, A: 255})
	dc.Clear()

	// 3. 绘制标题
	if err := s.drawTitle(dc, "推广中心"); err != nil {
		return "", err
	}

	// 4. 绘制分销商信息
	if err := s.drawDistributorInfo(dc, distributorName); err != nil {
		return "", err
	}

	// 5. 绘制活动统计
	infoText := fmt.Sprintf("管理活动：%d 个", campaignCount)
	if err := s.drawCampaignInfo(dc, infoText); err != nil {
		return "", err
	}

	// 6. 绘制底部提示
	if err := s.drawFooter(dc); err != nil {
		return "", err
	}

	// 7. 保存图片
	filename := fmt.Sprintf("distributor_%d.png", generateUniqueID())
	filepath := filepath.Join(s.posterDir, filename)
	if err := dc.SavePNG(filepath); err != nil {
		return "", fmt.Errorf("保存海报失败: %w", err)
	}

	// 8. 返回URL
	posterURL := fmt.Sprintf("%s/posters/%s", s.baseURL, filename)
	return posterURL, nil
}

// createGradient 创建渐变（简化版）
func (s *Service) createGradient(width, height int) []color.Color {
	gradient := make([]color.Color, height)
	for y := 0; y < height; y++ {
		ratio := float64(y) / float64(height)
		r := uint8(66 + (99-66)*ratio)
		g := uint8(133 + (102-133)*ratio)
		b := uint8(244 + (241-244)*ratio)
		gradient[y] = color.RGBA{R: r, G: g, B: b, A: 255}
	}
	return gradient
}

func (s *Service) drawTitle(dc *gg.Context, title string) error {
	// 使用默认字体（系统字体）
	if err := dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 48); err != nil {
		// 如果系统字体不可用，使用最小尺寸
		dc.SetRGB(255, 255, 255)
	}

	w, _ := dc.MeasureString(title)
	x := (750 - w) / 2

	dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	dc.DrawString(title, x, 150)

	return nil
}

func (s *Service) drawDescription(dc *gg.Context, desc string) error {
	if desc == "" {
		return nil
	}

	dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 28)

	dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	dc.DrawStringAnchored(desc, 375, 280, 0.5, 0.5)

	return nil
}

func (s *Service) drawDistributorInfo(dc *gg.Context, distributorName string) error {
	dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	dc.DrawRoundedRectangle(50, 400, 650, 150, 20)
	dc.Fill()

	dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 36)
	dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	dc.DrawStringAnchored("推广员: "+distributorName, 375, 475, 0.5, 0.5)

	return nil
}

// drawQRCode 绘制二维码
func (s *Service) drawQRCode(dc *gg.Context, qrcodePath string) error {
	img, err := gg.LoadImage(qrcodePath)
	if err != nil {
		return fmt.Errorf("加载二维码失败: %w", err)
	}

	// 绘制二维码（居中）
	dc.DrawImage(img, 300, 650)

	return nil
}

func (s *Service) drawCampaignInfo(dc *gg.Context, info string) error {
	dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	dc.DrawRoundedRectangle(50, 450, 650, 100, 20)
	dc.Fill()

	dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 32)
	dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	dc.DrawStringAnchored(info, 375, 500, 0.5, 0.5)

	return nil
}

func (s *Service) drawFooter(dc *gg.Context) error {
	dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 28)
	dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 200})
	dc.DrawStringAnchored("长按识别二维码，参与活动", 375, 1150, 0.5, 0.5)

	return nil
}

// generateQRCode 生成二维码
func (s *Service) generateQRCode(data string) (string, error) {
	// 生成二维码（中等纠错级别）
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("生成二维码失败: %w", err)
	}

	// 保存到临时文件
	tmpPath := fmt.Sprintf("/tmp/qrcode_%d.png", generateUniqueID())
	if err := qr.WriteFile(256, tmpPath); err != nil {
		return "", fmt.Errorf("保存二维码失败: %w", err)
	}

	return tmpPath, nil
}

// GenerateQRCodeAsBase64 生成二维码并返回base64编码
func (s *Service) GenerateQRCodeAsBase64(data string) (string, error) {
	// 简化处理：实际应该使用png.Encode
	// 这里返回一个占位符
	return "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==", nil
}

func generateUniqueID() int64 {
	return int64(time.Now().UnixNano())
}
