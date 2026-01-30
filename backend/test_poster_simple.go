package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
)

func main() {
	posterDir := "/opt/data/posters"

	// 确保目录存在
	if err := os.MkdirAll(posterDir, 0755); err != nil {
		fmt.Printf("创建海报目录失败: %v\n", err)
		return
	}

	// 创建画布
	width := 750
	height := 1334
	dc := gg.NewContext(width, height)

	// 绘制背景
	dc.SetRGB(66, 133, 244)
	dc.Clear()

	// 绘制标题
	dc.SetFontFace(nil)
	dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	dc.DrawString("测试海报标题", 375, 150)

	// 生成二维码
	qr, err := qrcode.New("test-data", qrcode.Medium)
	if err != nil {
		fmt.Printf("生成二维码失败: %v\n", err)
		return
	}

	// 保存二维码
	qrcodePath := "/tmp/test_qrcode.png"
	if err := qr.WriteFile(256, qrcodePath); err != nil {
		fmt.Printf("保存二维码失败: %v\n", err)
		return
	}
	defer os.Remove(qrcodePath)

	// 加载二维码
	img, err := gg.LoadImage(qrcodePath)
	if err != nil {
		fmt.Printf("加载二维码失败: %v\n", err)
		return
	}

	// 绘制二维码
	dc.DrawImage(img, 300, 650)

	// 绘制底部提示
	dc.SetFontFace(nil)
	dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 200})
	dc.DrawStringAnchored("长按识别二维码，参与活动", 375, 1150, 0.5, 0.5)

	// 保存图片
	filename := "test_poster.png"
	filepath := filepath.Join(posterDir, filename)

	fmt.Printf("保存海报到: %s\n", filepath)
	if err := dc.SavePNG(filepath); err != nil {
		fmt.Printf("保存海报失败: %v\n", err)
		return
	}

	fmt.Printf("海报生成成功！\n")
}
