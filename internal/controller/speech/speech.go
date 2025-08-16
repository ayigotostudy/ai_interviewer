package speech

import (
	"ai_jianli_go/config"
	"ai_jianli_go/pkg/speech"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type SpeechController struct {
	recognizer *speech.Recognizer
}

func NewSpeechController() *SpeechController {
	fmt.Println(config.GetSpeechConfig())
	config := speech.Config{
		APIKey:    config.GetSpeechConfig().APIKey,
		APISecret: config.GetSpeechConfig().APISecret,
		AppID:     config.GetSpeechConfig().AppID,
	}

	return &SpeechController{
		recognizer: speech.NewRecognizer(config),
	}
}

// Recognize 处理语音识别请求
func (c *SpeechController) Recognize(ctx *gin.Context) {
	// 获取上传的音频文件
	file, err := ctx.FormFile("audio")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无法获取音频文件",
		})
		return
	}

	// 创建临时文件
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "audio_"+time.Now().Format("20060102150405")+".wav")

	// 保存上传的文件
	if err := ctx.SaveUploadedFile(file, tempFile); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "保存音频文件失败",
		})
		return
	}
	defer os.Remove(tempFile) // 清理临时文件

	// 执行语音识别
	resultData, err := c.recognizer.RecognizeFile(tempFile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "语音识别失败: " + err.Error(),
		})
		return
	}

	// 返回识别结果
	ctx.JSON(http.StatusOK, gin.H{
		"text": resultData,
	})
}
