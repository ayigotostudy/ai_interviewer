import request from '@/plugins/axiosInstance'

export interface MeetingListResponse {
  code: number;
  msg: string;
  data: Array<{
    id: number;
    created_at: string;
    updated_at: string;
    deleted_at: string | null;
    user_id: number;
    candidate: string;
    position: string;
    job_description: string;
    time: number;
    status: string;
    remark: string;
    resume: string;
    interview_record: string;
    interview_summary: string;
  }>;
}

export interface CreateMeetingParams {
  candidate: string;
  position: string;
  job_description: string;
  time: number;
  status: string;
}

export interface UpdateMeetingParams {
  id: number;
  user_id: number;
  candidate: string;
  position: string;
  job_description: string;
  time: number;
  status: string;
  remark: string;
  interview_record: string;
  interview_summary: string;
}

export interface UploadResumeParams {
  meeting_id: number;
  resume: string;
}

export interface AiInterviewParams {
  meeting_id: number;
  answer: string;
}

export function getMeetingList() {
  return request.get('/meeting/list')
}

export function getMeetingDetail(id: string) {
  return request.get('/meeting', {
    params: { id }
  })
}   

export function createMeeting(data: CreateMeetingParams) { 
  return request.post('/meeting', data)
}

export function updateMeeting(data: UpdateMeetingParams) {
  return request.put('/meeting', data)
}

export function deleteMeeting(id: number) {
  return request.delete('/meeting', {
    data: { id }
  })
}

export function uploadResume(data: UploadResumeParams) {
  return request.post('/meeting/upload_resume', data)
}

export function aiInterview(data: AiInterviewParams) {
  return request.post('/meeting/ai_interview', data)
}

/**
 * 获取面试评价
 * @param meetingId 面试ID
 * @returns 面试评价数据
 */
export function getMeetingEvaluation(meetingId: string | number) {
  return request.get(`/meeting/remark?id=${meetingId}`)
}