// DOM 요소
const loadingEl = document.getElementById('loading');
const errorEl = document.getElementById('error');
const reportsEl = document.getElementById('reports');

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

        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('리포트를 불러오는데 실패했습니다.');
        }
        const reports = await response.json();
        return reports;
    } catch (error) {
        console.error('Error:', error);
        throw error;
    }
}

// 리포트 카드 생성 함수
function createReportCard(report) {
    const card = document.createElement('div');
    card.className = 'report-card';

    const header = document.createElement('div');
    header.className = 'report-header';

    const title = document.createElement('h2');
    title.className = 'report-title';
    title.textContent = report.name;

    const metadata = document.createElement('div');
    metadata.className = 'report-metadata';

    const date = document.createElement('span');
    date.className = 'report-date';
    date.textContent = report.date;

    const decision = document.createElement('span');
    decision.className = `report-decision ${report.decision.toLowerCase()}`;
    decision.textContent = report.decision;

    const locale = document.createElement('span');
    locale.className = 'report-locale';
    locale.textContent = report.locale;

    metadata.appendChild(date);
    metadata.appendChild(decision);
    metadata.appendChild(locale);

    const content = document.createElement('div');
    content.className = 'report-content';
    content.innerHTML = marked.parse(report.content);

    header.appendChild(title);
    header.appendChild(metadata);
    card.appendChild(header);
    card.appendChild(content);

    return card;
}

// 리포트 목록 렌더링 함수
function renderReports(reports) {
    reportsEl.innerHTML = '';
    reports.forEach(report => {
        const card = createReportCard(report);
        reportsEl.appendChild(card);
    });
}

// 초기 로딩
async function init() {
    try {
        loadingEl.style.display = 'block';
        errorEl.style.display = 'none';
        
        const reports = await fetchReports();
        renderReports(reports);
    } catch (error) {
        errorEl.textContent = error.message;
        errorEl.style.display = 'block';
    } finally {
        loadingEl.style.display = 'none';
    }
}

// 페이지 로드 시 초기화
document.addEventListener('DOMContentLoaded', init); 