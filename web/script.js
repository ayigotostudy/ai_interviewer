// 粒子动画
class Particle {
    constructor(canvas) {
        this.canvas = canvas;
        this.ctx = canvas.getContext('2d');
        this.particles = [];
        this.init();
    }

    init() {
        this.canvas.width = window.innerWidth;
        this.canvas.height = window.innerHeight;
        
        // 创建粒子
        for (let i = 0; i < 100; i++) {
            this.particles.push({
                x: Math.random() * this.canvas.width,
                y: Math.random() * this.canvas.height,
                size: Math.random() * 2 + 1,
                speedX: Math.random() * 2 - 1,
                speedY: Math.random() * 2 - 1
            });
        }

        this.animate();
        window.addEventListener('resize', () => this.handleResize());
    }

    handleResize() {
        this.canvas.width = window.innerWidth;
        this.canvas.height = window.innerHeight;
    }

    animate() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        
        this.particles.forEach(particle => {
            particle.x += particle.speedX;
            particle.y += particle.speedY;

            if (particle.x < 0 || particle.x > this.canvas.width) particle.speedX *= -1;
            if (particle.y < 0 || particle.y > this.canvas.height) particle.speedY *= -1;

            this.ctx.beginPath();
            this.ctx.arc(particle.x, particle.y, particle.size, 0, Math.PI * 2);
            this.ctx.fillStyle = 'rgba(0, 240, 252, 0.5)';
            this.ctx.fill();
        });

        requestAnimationFrame(() => this.animate());
    }
}

// API 基础配置
const API_BASE_URL = 'http://localhost:8080/api/v1';

// 工具函数
const showError = (message) => {
    const errorDiv = document.createElement('div');
    errorDiv.className = 'error-message';
    errorDiv.textContent = message;
    document.body.appendChild(errorDiv);
    setTimeout(() => errorDiv.remove(), 3000);
};

const showSuccess = (message) => {
    const successDiv = document.createElement('div');
    successDiv.className = 'success-message';
    successDiv.textContent = message;
    document.body.appendChild(successDiv);
    setTimeout(() => successDiv.remove(), 3000);
};

const setLoading = (isLoading) => {
    const loadingElement = document.getElementById('loading');
    if (loadingElement) {
        loadingElement.style.display = isLoading ? 'flex' : 'none';
    }
};

// 用户信息管理
const getToken = () => localStorage.getItem('token');
const setToken = (token) => localStorage.setItem('token', token);
const removeToken = () => localStorage.removeItem('token');
const setUserId = (userId) => localStorage.setItem('userId', userId);
const getUserId = () => localStorage.getItem('userId');
const setUserEmail = (email) => localStorage.setItem('userEmail', email);
const getUserEmail = () => localStorage.getItem('userEmail');

// API 请求函数
const apiRequest = async (endpoint, method = 'GET', data = null) => {
    setLoading(true);
    const headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    };

    // 非登录注册接口都需要携带 token
    if (!endpoint.includes('/user/login') && !endpoint.includes('/user/register')) {
        const token = getToken();
        if (!token) {
            showError('请先登录');
            window.location.href = 'index.html';
            return;
        }
        headers['Authorization'] = `Bearer ${token}`;
    }

    try {
        const url = `${API_BASE_URL}${endpoint}`;
        console.log('Sending request to:', url);
        console.log('Request method:', method);
        console.log('Request headers:', headers);
        console.log('Request data:', data);

        const response = await fetch(url, {
            method,
            headers,
            body: method !== 'GET' ? JSON.stringify(data) : null,
            credentials: 'include',
            mode: 'cors'
        });

        console.log('Response status:', response.status);
        const result = await response.json();
        console.log('Response data:', result);
        
        if (!response.ok) {
            throw new Error(result.msg || `HTTP error! status: ${response.status}`);
        }

        if (result.code !== 1000) {
            throw new Error(result.msg || '请求失败');
        }

        return result.data;
    } catch (error) {
        console.error('API request failed:', error);
        showError(error.message || '网络请求失败，请稍后重试');
        throw error;
    } finally {
        setLoading(false);
    }
};

// 用户认证相关
const register = async (email, password) => {
    try {
        console.log('Registering user:', { email });
        const result = await apiRequest('/user/register', 'POST', { 
            email, 
            password 
        });
        showSuccess('注册成功，请登录');
        return result;
    } catch (error) {
        console.error('Register failed:', error);
        throw error;
    }
};

const login = async (email, password) => {
    try {
        console.log('Logging in user:', { email });
        const result = await apiRequest('/user/login', 'POST', { 
            email, 
            password 
        });
        setToken(result.token);
        setUserId(result.user_id);
        setUserEmail(email);
        showSuccess('登录成功');
    } catch (error) {
        console.error('Login failed:', error);
        throw error;
    }
};

// 面试管理相关
const createMeeting = async (meetingData) => {
    try {
        const result = await apiRequest('/meeting', 'POST', {
            candidate: meetingData.candidate,
            position: meetingData.position,
            job_description: meetingData.jobDescription,
            time: Date.now(),
            status: '进行中',
            remark: ''
        });
        showSuccess('面试创建成功');
        await loadMeetings();
        return result;
    } catch (error) {
        console.error('Create meeting failed:', error);
        throw error;
    }
};

const updateMeeting = async (meetingId, meetingData) => {
    try {
        const result = await apiRequest('/meeting', 'PUT', {
            id: meetingId,
            candidate: meetingData.candidate,
            position: meetingData.position,
            job_description: meetingData.jobDescription,
            status: meetingData.status,
            remark: meetingData.remark || ''
        });
        showSuccess('面试更新成功');
        await loadMeetings();
        return result;
    } catch (error) {
        console.error('Update meeting failed:', error);
        throw error;
    }
};

const getMeeting = async (meetingId) => {
    try {
        const result = await apiRequest(`/meeting?id=${meetingId}`, 'GET');
        return result;
    } catch (error) {
        console.error('Get meeting failed:', error);
        throw error;
    }
};

const deleteMeeting = async (meetingId) => {
    try {
        const result = await apiRequest(`/meeting?id=${meetingId}`, 'DELETE');
        showSuccess('面试删除成功');
        await loadMeetings();
        return result;
    } catch (error) {
        console.error('Delete meeting failed:', error);
        throw error;
    }
};

// AI面试相关
let currentMeetingId = null;

const showAIInterviewForm = async (meetingId) => {
    currentMeetingId = meetingId;
    const chatWindow = document.getElementById('interviewChat');
    if (!chatWindow) {
        showError('聊天窗口初始化失败');
        return;
    }

    chatWindow.style.display = 'flex';
    chatWindow.innerHTML = ''; // 清空聊天记录

    // 检查是否已上传简历
    try {
        const meeting = await getMeeting(meetingId);
        if (!meeting.resume) {
            // 显示上传简历提示
            addMessageToChat('ai', '请先上传简历，以便开始面试。');
            // 显示简历上传区域
            const resumeInput = document.createElement('div');
            resumeInput.className = 'resume-upload';
            resumeInput.innerHTML = `
                <textarea id="resumeInput" class="cyber-input" placeholder="请粘贴简历内容..."></textarea>
                <button onclick="uploadResumeAndStart()" class="cyber-button">
                    <span class="glow-text">上传简历</span>
                </button>
            `;
            chatWindow.appendChild(resumeInput);
            return;
        }

        // 如果已有简历，直接开始面试
        addMessageToChat('ai', '你好，我是AI面试官。请开始你的面试。');
    } catch (error) {
        showError('获取面试信息失败：' + error.message);
        hideAIInterviewForm();
    }
};

const uploadResumeAndStart = async () => {
    if (!currentMeetingId) {
        showError('面试已结束');
        return;
    }

    const resumeInput = document.getElementById('resumeInput');
    const resume = resumeInput.value.trim();
    
    if (!resume) {
        showError('请输入简历内容');
        return;
    }

    try {
        await uploadResume(currentMeetingId, resume);
        showSuccess('简历上传成功');
        
        // 移除简历上传区域
        const resumeUpload = document.querySelector('.resume-upload');
        if (resumeUpload) {
            resumeUpload.remove();
        }

        // 开始面试
        addMessageToChat('ai', '你好，我是AI面试官。请开始你的面试。');
    } catch (error) {
        showError('上传简历失败：' + error.message);
    }
};

const hideAIInterviewForm = () => {
    const chatWindow = document.getElementById('interviewChat');
    if (chatWindow) {
        chatWindow.style.display = 'none';
    }
    currentMeetingId = null;
};

const sendAnswerToAI = async () => {
    if (!currentMeetingId) {
        showError('面试已结束');
        return;
    }

    const answerInput = document.getElementById('answerInput');
    const answer = answerInput.value.trim();
    
    if (!answer) {
        showError('请输入回答');
        return;
    }

    try {
        // 显示用户回答
        addMessageToChat('user', answer);
        answerInput.value = ''; // 清空输入框

        // 发送回答到服务器
        const result = await sendAnswer(currentMeetingId, answer);
        
        // 显示AI回复
        addMessageToChat('ai', result.reply);
    } catch (error) {
        showError('发送回答失败：' + error.message);
    }
};

// 修改面试列表项，添加开始面试按钮
const loadMeetings = async () => {
    try {
        const meetings = await apiRequest('/meeting/list', 'GET');
        console.log('Meetings data:', meetings);
        const meetingList = document.getElementById('meetingList');
        
        if (!Array.isArray(meetings)) {
            console.error('Expected meetings to be an array, got:', meetings);
            return;
        }

        meetingList.innerHTML = meetings.map(meeting => `
            <div class="meeting-item" data-id="${meeting.id}">
                <h3>${meeting.candidate}</h3>
                <p>职位：${meeting.position}</p>
                <p>状态：${meeting.status}</p>
                <p>时间：${new Date(meeting.time).toLocaleString()}</p>
                <p>备注：${meeting.remark || '无'}</p>
                <div class="meeting-actions">
                    <button onclick="editMeeting(${meeting.id})" class="cyber-button small">
                        <span class="glow-text">编辑</span>
                    </button>
                    <button onclick="deleteMeeting(${meeting.id})" class="cyber-button small">
                        <span class="glow-text">删除</span>
                    </button>
                    <button onclick="showChat(${meeting.id})" class="cyber-button small">
                        <span class="glow-text">开始面试</span>
                    </button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Load meetings failed:', error);
    }
};

// 简历管理相关
const uploadResume = async (meetingId, resume) => {
    try {
        const result = await apiRequest('/meeting/upload_resume', 'POST', {
            meeting_id: meetingId,
            resume: resume
        });
        showSuccess('简历上传成功');
        return result;
    } catch (error) {
        console.error('Upload resume failed:', error);
        throw error;
    }
};

// AI面试相关
const sendAnswer = async (meetingId, answer) => {
    try {
        const result = await apiRequest('/meeting/ai_interview', 'POST', {
            meeting_id: meetingId,
            answer: answer
        });
        return result;
    } catch (error) {
        console.error('Send answer failed:', error);
        throw error;
    }
};

// UI交互函数
const showRegisterForm = () => {
    // 移除已存在的注册表单
    const existingForm = document.querySelector('.register-form');
    if (existingForm) {
        existingForm.remove();
    }

    const registerForm = document.createElement('div');
    registerForm.className = 'register-form';
    registerForm.innerHTML = `
        <h2>注册账号</h2>
        <form id="registerForm" method="POST">
            <div class="form-group">
                <input type="email" id="registerEmail" name="email" required class="cyber-input">
                <label for="registerEmail">邮箱</label>
                <div class="input-line"></div>
            </div>
            <div class="form-group">
                <input type="password" id="registerPassword" name="password" required class="cyber-input">
                <label for="registerPassword">密码</label>
                <div class="input-line"></div>
            </div>
            <button type="submit" class="cyber-button">
                <span class="glow-text">注册</span>
                <span class="particle-layer"></span>
            </button>
            <div class="form-footer">
                <a href="javascript:void(0)" onclick="document.querySelector('.register-form').remove()" class="cyber-link">返回登录</a>
            </div>
        </form>
    `;
    document.body.appendChild(registerForm);

    // 添加表单提交事件监听器
    const form = registerForm.querySelector('#registerForm');
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const email = document.getElementById('registerEmail').value;
        const password = document.getElementById('registerPassword').value;
        try {
            await register(email, password);
            registerForm.remove();
        } catch (error) {
            showError('注册失败：' + error.message);
        }
    });
};

const showCreateMeetingForm = () => {
    const form = document.getElementById('createMeetingForm');
    form.style.display = 'flex';
};

const hideCreateMeetingForm = () => {
    const form = document.getElementById('createMeetingForm');
    form.style.display = 'none';
};

const addMessageToChat = (type, content) => {
    const chatWindow = document.getElementById('interviewChat');
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${type}`;
    messageDiv.textContent = content;
    chatWindow.appendChild(messageDiv);
    chatWindow.scrollTop = chatWindow.scrollHeight;
};

// 聊天窗口相关变量
let isDragging = false;
let currentX;
let currentY;
let initialX;
let initialY;
let xOffset = 0;
let yOffset = 0;
let chatWindow = document.getElementById("interviewChat");
let chatHeader = document.getElementById("chatHeader");
let chatContent = document.getElementById("chatContent");

// 拖拽相关函数
function dragStart(e) {
    if (e.target === chatHeader || e.target.parentNode === chatHeader) {
        initialX = e.clientX - xOffset;
        initialY = e.clientY - yOffset;

        if (e.target === chatHeader || e.target.parentNode === chatHeader) {
            isDragging = true;
        }
    }
}

function dragEnd(e) {
    initialX = currentX;
    initialY = currentY;
    isDragging = false;
}

function drag(e) {
    if (isDragging) {
        e.preventDefault();
        currentX = e.clientX - initialX;
        currentY = e.clientY - initialY;

        xOffset = currentX;
        yOffset = currentY;

        setTranslate(currentX, currentY, chatWindow);
    }
}

function setTranslate(xPos, yPos, el) {
    el.style.transform = `translate3d(${xPos}px, ${yPos}px, 0)`;
}

// 聊天窗口控制函数
function showChat(meetingId) {
    currentMeetingId = meetingId;
    chatWindow.style.display = 'flex';
    chatContent.innerHTML = ''; // 清空聊天记录
    document.getElementById('answerInput').value = '';
}

function minimizeChat() {
    chatWindow.style.display = 'none';
}

function closeChat() {
    chatWindow.style.display = 'none';
    chatContent.innerHTML = '';
    document.getElementById('answerInput').value = '';
    currentMeetingId = null;
}

// 事件监听器
document.addEventListener('DOMContentLoaded', () => {
    // 初始化粒子动画
    const canvas = document.getElementById('particleCanvas');
    if (canvas) {
        new Particle(canvas);
    }

    // 检查登录状态
    const token = getToken();
    if (!token && window.location.pathname.includes('dashboard.html')) {
        window.location.href = 'index.html';
        return;
    }

    // 登录表单处理
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            try {
                await login(email, password);
                window.location.href = 'dashboard.html';
            } catch (error) {
                showError('登录失败：' + error.message);
            }
        });
    }

    // 注册表单处理
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('registerEmail').value;
            const password = document.getElementById('registerPassword').value;
            try {
                await register(email, password);
                document.querySelector('.register-form').remove();
            } catch (error) {
                showError('注册失败：' + error.message);
            }
        });
    }

    // 创建面试表单处理
    const meetingForm = document.getElementById('meetingForm');
    if (meetingForm) {
        meetingForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const meetingData = {
                candidate: document.getElementById('candidate').value,
                position: document.getElementById('position').value,
                jobDescription: document.getElementById('jobDescription').value
            };
            try {
                await createMeeting(meetingData);
                hideCreateMeetingForm();
                meetingForm.reset();
            } catch (error) {
                showError('创建面试失败：' + error.message);
            }
        });
    }

    // 加载用户信息
    const userEmail = document.getElementById('userEmail');
    if (userEmail) {
        userEmail.textContent = getUserEmail();
    }

    // 加载面试列表
    if (window.location.pathname.includes('dashboard.html')) {
        loadMeetings();
    }

    // 初始化聊天窗口
    const chatWindow = document.getElementById('interviewChat');
    const chatHeader = document.getElementById('chatHeader');
    const answerInput = document.getElementById('answerInput');

    if (chatWindow && chatHeader) {
        // 初始化拖拽事件监听
        chatHeader.addEventListener('mousedown', dragStart);
        document.addEventListener('mousemove', drag);
        document.addEventListener('mouseup', dragEnd);

        // 添加回车发送功能
        if (answerInput) {
            answerInput.addEventListener('keydown', function(e) {
                if (e.key === 'Enter' && e.ctrlKey) {
                    e.preventDefault();
                    sendAnswerToAI();
                }
            });
        }
    }
});

// 全局函数
window.showRegister = showRegisterForm;
window.createMeeting = () => showCreateMeetingForm();
window.editMeeting = async (meetingId) => {
    try {
        const meeting = await getMeeting(meetingId);
        // TODO: 显示编辑表单
    } catch (error) {
        showError('获取面试信息失败：' + error.message);
    }
};
window.deleteMeeting = async (meetingId) => {
    if (confirm('确定要删除这个面试吗？')) {
        try {
            await deleteMeeting(meetingId);
        } catch (error) {
            showError('删除面试失败：' + error.message);
        }
    }
};
window.logout = () => {
    removeToken();
    window.location.href = 'index.html';
}; 