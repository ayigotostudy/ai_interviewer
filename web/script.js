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

    const token = getToken();
    if (token) {
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
            body: data ? JSON.stringify(data) : null,
            credentials: 'same-origin',
            mode: 'cors'
        });

        console.log('Response status:', response.status);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result = await response.json();
        console.log('Response data:', result);
        
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
        await apiRequest('/user/register', 'POST', { email, password });
        showSuccess('注册成功，请登录');
    } catch (error) {
        console.error('Register failed:', error);
        throw error;
    }
};

const login = async (email, password) => {
    try {
        console.log('Logging in user:', { email });
        const result = await apiRequest('/user/login', 'POST', { email, password });
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
            user_id: getUserId(),
            ...meetingData
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
            user_id: getUserId(),
            ...meetingData
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
        const result = await apiRequest('/meeting', 'GET', { id: meetingId });
        return result.data;
    } catch (error) {
        console.error('Get meeting failed:', error);
        throw error;
    }
};

const deleteMeeting = async (meetingId) => {
    try {
        const result = await apiRequest('/meeting', 'DELETE', { id: meetingId });
        showSuccess('面试删除成功');
        await loadMeetings();
        return result;
    } catch (error) {
        console.error('Delete meeting failed:', error);
        throw error;
    }
};

const loadMeetings = async () => {
    try {
        const result = await apiRequest('/meeting/list', 'GET');
        const meetingList = document.getElementById('meetingList');
        
        meetingList.innerHTML = result.data.map(meeting => `
            <div class="meeting-item" data-id="${meeting.id}">
                <h3>${meeting.candidate}</h3>
                <p>职位：${meeting.position}</p>
                <p>状态：${meeting.status}</p>
                <div class="meeting-actions">
                    <button onclick="editMeeting(${meeting.id})" class="cyber-button small">
                        <span class="glow-text">编辑</span>
                    </button>
                    <button onclick="deleteMeeting(${meeting.id})" class="cyber-button small">
                        <span class="glow-text">删除</span>
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
            user_id: getUserId(),
            meeting_id: meetingId,
            resume
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
            user_id: getUserId(),
            meeting_id: meetingId,
            answer
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
    const form = document.createElement('div');
    form.className = 'meeting-form';
    form.innerHTML = `
        <h2>创建面试</h2>
        <form id="createMeetingForm">
            <div class="form-group">
                <input type="text" id="candidate" required class="cyber-input">
                <label for="candidate">候选人姓名</label>
            </div>
            <div class="form-group">
                <input type="text" id="position" required class="cyber-input">
                <label for="position">应聘职位</label>
            </div>
            <div class="form-group">
                <textarea id="jobDescription" required class="cyber-input"></textarea>
                <label for="jobDescription">职位描述</label>
            </div>
            <button type="submit" class="cyber-button">
                <span class="glow-text">创建</span>
            </button>
        </form>
    `;
    document.body.appendChild(form);
};

const addMessageToChat = (type, content) => {
    const chatWindow = document.getElementById('interviewChat');
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${type}`;
    messageDiv.textContent = content;
    chatWindow.appendChild(messageDiv);
    chatWindow.scrollTop = chatWindow.scrollHeight;
};

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
    const createMeetingForm = document.getElementById('createMeetingForm');
    if (createMeetingForm) {
        createMeetingForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const meetingData = {
                candidate: document.getElementById('candidate').value,
                position: document.getElementById('position').value,
                job_description: document.getElementById('jobDescription').value,
                time: Date.now(),
                status: '进行中'
            };
            try {
                await createMeeting(meetingData);
                document.querySelector('.meeting-form').remove();
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