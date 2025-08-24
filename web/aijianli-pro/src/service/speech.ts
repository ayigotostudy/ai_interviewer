import request from '@/plugins/axiosInstance'

// 语音识别接口参数
export interface SpeechRecognitionParams {
  audio_data: ArrayBuffer
  audio_format: string
  sample_rate?: number
  channels?: number
}

// 语音识别响应
export interface SpeechRecognitionResponse {
  success: boolean
  text?: string
  error?: string
}

/**
 * 语音识别接口
 * @param audioData PCM格式的音频数据
 * @param sampleRate 采样率，默认16000
 * @param channels 声道数，默认1
 * @returns 识别结果
 */
export function recognizeSpeech(
  audioData: ArrayBuffer, 
  sampleRate: number = 16000, 
  channels: number = 1
): Promise<SpeechRecognitionResponse> {
  return new Promise((resolve) => {
    try {
      // 将ArrayBuffer转换为Base64字符串
      const uint8Array = new Uint8Array(audioData)
      let binary = ''
      for (let i = 0; i < uint8Array.length; i++) {
        binary += String.fromCharCode(uint8Array[i])
      }
      const base64Audio = btoa(binary)
      
      // 构建请求参数
      const params: SpeechRecognitionParams = {
        audio_data: audioData,
        audio_format: 'pcm',
        sample_rate: sampleRate,
        channels: channels
      }
      
      // 调用后端语音识别接口
      request.post('/api/v1/speech/recognize', {
        audio_data: base64Audio,
        audio_format: 'pcm',
        sample_rate: sampleRate,
        channels: channels
      })
      .then(response => {
        console.log('语音识别响应:', response)
        
        if (response.data.code === 1000) {
          resolve({
            success: true,
            text: response.data.data?.text || '识别结果为空'
          })
        } else {
          resolve({
            success: false,
            error: response.data.msg || '语音识别失败'
          })
        }
      })
      .catch(error => {
        console.error('语音识别请求失败:', error)
        resolve({
          success: false,
          error: error.response?.data?.msg || '网络请求失败'
        })
      })
      
    } catch (error) {
      console.error('处理音频数据失败:', error)
      resolve({
        success: false,
        error: '音频数据处理失败'
      })
    }
  })
}

/**
 * 测试语音识别接口连接
 * @returns 连接状态
 */
export function testSpeechAPI(): Promise<boolean> {
  return new Promise((resolve) => {
    request.get('/api/v1/speech/health')
      .then(() => resolve(true))
      .catch(() => resolve(false))
  })
}
