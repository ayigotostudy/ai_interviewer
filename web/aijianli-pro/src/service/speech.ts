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

// ===== 音频工具函数：PCM16 -> WAV =====
function writeString(view: DataView, offset: number, str: string) {
  for (let i = 0; i < str.length; i++) {
    view.setUint8(offset + i, str.charCodeAt(i))
  }
}

function encodeWavFromPcmInt16(
  pcmBuffer: ArrayBuffer,
  wavSampleRate: number,
  wavChannels: number
): ArrayBuffer {
  const bytesPerSample = 2 // 16-bit
  const pcmData = new Int16Array(pcmBuffer)
  const blockAlign = wavChannels * bytesPerSample
  const byteRate = wavSampleRate * blockAlign

  const wavHeaderSize = 44
  const wavBuffer = new ArrayBuffer(wavHeaderSize + pcmData.byteLength)
  const view = new DataView(wavBuffer)

  // "RIFF"
  writeString(view, 0, 'RIFF')
  // file length - 8
  view.setUint32(4, wavHeaderSize + pcmData.byteLength - 8, true)
  // "WAVE"
  writeString(view, 8, 'WAVE')
  // "fmt " chunk
  writeString(view, 12, 'fmt ')
  view.setUint32(16, 16, true) // Subchunk1Size for PCM
  view.setUint16(20, 1, true) // AudioFormat PCM = 1
  view.setUint16(22, wavChannels, true)
  view.setUint32(24, wavSampleRate, true)
  view.setUint32(28, byteRate, true)
  view.setUint16(32, blockAlign, true)
  view.setUint16(34, 16, true) // BitsPerSample
  // "data" chunk
  writeString(view, 36, 'data')
  view.setUint32(40, pcmData.byteLength, true)

  // PCM samples
  const out = new Int16Array(wavBuffer, wavHeaderSize)
  out.set(pcmData)
  return wavBuffer
}

/**
 * 将 PCM16 缓冲区（可超长）裁剪到 limitSeconds，并创建可播放的 WAV Blob URL，立即播放
 */
export function previewPcmAudio(
  audioData: ArrayBuffer,
  sampleRate: number = 16000,
  channels: number = 1,
  limitSeconds: number = 60
): { audio: HTMLAudioElement; url: string; blob: Blob } {
  const bytesPerSample = 2
  const maxSamples = Math.max(1, Math.floor(sampleRate * channels * limitSeconds))
  const maxBytes = maxSamples * bytesPerSample
  const limitedPcmBuffer = audioData.byteLength > maxBytes ? audioData.slice(0, maxBytes) : audioData

  const wavArrayBuffer = encodeWavFromPcmInt16(limitedPcmBuffer, sampleRate, channels)
  const blob = new Blob([wavArrayBuffer], { type: 'audio/wav' })
  const url = URL.createObjectURL(blob)
  const audio = new Audio(url)
  audio.play().catch(() => {})
  return { audio, url, blob }
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
      // 限制音频时长为 60 秒（按 16-bit PCM 计算）
      const bytesPerSample = 2
      const maxSamples = sampleRate * channels * 60
      const maxBytes = maxSamples * bytesPerSample
      const limitedPcmBuffer =
        audioData.byteLength > maxBytes ? audioData.slice(0, maxBytes) : audioData

      // 将 PCM 打包为 WAV
      const wavArrayBuffer = encodeWavFromPcmInt16(limitedPcmBuffer, sampleRate, channels)

      const form = new FormData()
      const wavBlob = new Blob([wavArrayBuffer], { type: 'audio/wav' })
      form.append('audio', wavBlob, 'audio.wav')

      request
        .post('/speech/recognize', form, {
          // 让浏览器自动设置 multipart boundary，不手动设置 Content-Type
          timeout: 60000,
          // @ts-ignore
          maxBodyLength: Infinity,
          // @ts-ignore
          maxContentLength: Infinity,
        })
        .then((response) => {
          console.log('语音识别响应:', response)
          const text: string | undefined = response?.data?.text
          if (text && typeof text === 'string') {
            resolve({ success: true, text })
          } else {
            resolve({ success: false, error: '语音识别失败' })
          }
        })
        .catch((error) => {
          console.error('语音识别请求失败:', error)
          resolve({
            success: false,
            error: error?.response?.data?.error || error?.message || '网络请求失败',
          })
        })
    } catch (error: any) {
      console.error('处理音频数据失败:', error)
      resolve({ success: false, error: '音频数据处理失败' })
    }
  })
}

/**
 * 测试语音识别接口连接
 * @returns 连接状态
 */
export function testSpeechAPI(): Promise<boolean> {
  return new Promise((resolve) => {
    request.get('/speech/health')
      .then(() => resolve(true))
      .catch(() => resolve(false))
  })
}
