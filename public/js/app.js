// DOM 요소
const loadingEl = document.getElementById('loading');
const errorEl = document.getElementById('error');
const listViewEl = document.getElementById('list-view');
const gridViewEl = document.getElementById('grid-view');
const reportsListEl = document.getElementById('reports-list');
const reportsGridEl = document.getElementById('grid-view');
const modalEl = document.getElementById('report-modal');
const modalTitleEl = document.getElementById('modal-title');
const modalDateEl = document.getElementById('modal-date');
const modalDecisionEl = document.getElementById('modal-decision');
const modalTypeEl = document.getElementById('modal-type');
const modalLocaleEl = document.getElementById('modal-locale');
const modalContentEl = document.getElementById('modal-content');
const closeBtnEl = document.querySelector('.close-btn');
const reportTypeSelect = document.getElementById('report-type');

// 현재 상태
let currentLocale = 'EN';
let currentView = 'list';
let currentType = 'complete';
let reports = [];

// API 호출 함수
async function fetchReports(params = {}) {
    try {
        let url = '/api/reports';
        if (params.date) {
            url = `/api/reports/date?date=${params.date}`;
        } else if (params.decision) {
            url = `/api/reports/decision?decision=${params.decision}`;
        }
        
        if (params.locale) {
            url += `${url.includes('?') ? '&' : '?'}locale=${params.locale}`;
        }
        
        if (params.type) {
            url += `${url.includes('?') ? '&' : '?'}type=${params.type}`;
        }

        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('리포트를 불러오는데 실패했습니다.');
        }
        return await response.json();
    } catch (error) {
        console.error('Error:', error);
        throw error;
    }
}

// 리스트 뷰 렌더링
function renderListView(reports) {
    reportsListEl.innerHTML = '';
    reports.forEach(report => {
        const row = document.createElement('div');
        row.className = 'list-row';
        
        row.innerHTML = `
            <div class="list-cell">${report.date}</div>
            <div class="list-cell">${report.name}</div>
            <div class="list-cell">
                <span class="report-decision ${report.decision.toLowerCase()}">${report.decision}</span>
            </div>
            <div class="list-cell">
                <span class="report-type">${report.type}</span>
            </div>
            <div class="list-cell">
                <button class="view-btn" onclick="showReportDetail('${report.name}', '${report.date}', '${report.decision}', '${report.locale}', '${report.type}', \`${report.content.replace(/`/g, '\\`')}\`)">상세보기</button>
            </div>
        `;
        
        reportsListEl.appendChild(row);
    });
}

// 그리드 뷰 렌더링
function renderGridView(reports) {
    reportsGridEl.innerHTML = '';
    reports.forEach(report => {
        const card = document.createElement('div');
        card.className = 'report-card';
        
        card.innerHTML = `
            <div class="report-header">
                <h2 class="report-title">${report.name}</h2>
                <div class="report-metadata">
                    <span class="report-date">${report.date}</span>
                    <span class="report-decision ${report.decision.toLowerCase()}">${report.decision}</span>
                    <span class="report-type">${report.type}</span>
                    <span class="report-locale">${report.locale}</span>
                </div>
            </div>
            <div class="report-content">
                ${marked.parse(report.content)}
            </div>
        `;
        
        reportsGridEl.appendChild(card);
    });
}

// 리포트 상세 보기
function showReportDetail(name, date, decision, locale, type, content) {
    modalTitleEl.textContent = name;
    modalDateEl.textContent = date;
    modalDecisionEl.textContent = decision;
    modalDecisionEl.className = `report-decision ${decision.toLowerCase()}`;
    modalTypeEl.textContent = type;
    modalLocaleEl.textContent = locale;
    modalContentEl.innerHTML = marked.parse(content);
    
    modalEl.style.display = 'flex';
}

// 모달 닫기
function closeModal() {
    modalEl.style.display = 'none';
}

// 언어 변경
function changeLocale(locale) {
    currentLocale = locale;
    document.querySelectorAll('.lang-btn').forEach(btn => {
        btn.classList.toggle('active', btn.dataset.lang === locale);
    });
    loadReports();
}

// 리포트 타입 변경
function changeReportType(type) {
    currentType = type;
    loadReports();
}

// 뷰 변경
function changeView(view) {
    currentView = view;
    document.querySelectorAll('.view-btn').forEach(btn => {
        btn.classList.toggle('active', btn.dataset.view === view);
    });
    
    listViewEl.style.display = view === 'list' ? 'block' : 'none';
    gridViewEl.style.display = view === 'grid' ? 'grid' : 'none';
    
    if (view === 'list') {
        renderListView(reports);
    } else {
        renderGridView(reports);
    }
}

// 리포트 로드
async function loadReports() {
    try {
        loadingEl.style.display = 'block';
        errorEl.style.display = 'none';
        
        reports = await fetchReports({ 
            locale: currentLocale,
            type: currentType
        });
        
        if (currentView === 'list') {
            renderListView(reports);
        } else {
            renderGridView(reports);
        }
    } catch (error) {
        errorEl.textContent = error.message;
        errorEl.style.display = 'block';
    } finally {
        loadingEl.style.display = 'none';
    }
}

// 이벤트 리스너
document.addEventListener('DOMContentLoaded', () => {
    // 언어 선택 버튼
    document.querySelectorAll('.lang-btn').forEach(btn => {
        btn.addEventListener('click', () => changeLocale(btn.dataset.lang));
    });
    
    // 리포트 타입 선택
    reportTypeSelect.addEventListener('change', (e) => changeReportType(e.target.value));
    
    // 뷰 선택 버튼
    document.querySelectorAll('.view-btn').forEach(btn => {
        if (btn.dataset.view) {
            btn.addEventListener('click', () => changeView(btn.dataset.view));
        }
    });
    
    // 모달 닫기
    closeBtnEl.addEventListener('click', closeModal);
    modalEl.addEventListener('click', (e) => {
        if (e.target === modalEl) {
            closeModal();
        }
    });
    
    // 초기 로드
    loadReports();
}); 