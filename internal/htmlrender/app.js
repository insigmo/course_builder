// ══ Course App ═══════════════════════════════════════════════════
const b64Data = "__BASE64_DATA__";
const courseData = JSON.parse(decodeURIComponent(escape(window.atob(b64Data))));
const TOTAL_STEPS = __TOTAL_STEPS__;
const COURSE_TITLE = "__COURSE_TITLE_JS__";

let currentLessonIdx = 0, currentStepIdx = 0;
const COURSE_KEY = "course_" + COURSE_TITLE.replace(/\s+/g,"_").replace(/[^\w\u0400-\u04FF]/g,"").slice(0,40);
let completedSteps = new Set(readJSON(COURSE_KEY + "_progress", []));
let savedQuizzes = readJSON(COURSE_KEY + "_quizzes", {});
let openLessons = new Set(courseData.length ? ["0"] : []);

// ── Utils ─────────────────────────────────────────────────────────
function readJSON(k, fb) { try { const r=localStorage.getItem(k); return r?JSON.parse(r):fb; } catch(_){ return fb; } }
function writeJSON(k, v) { try { localStorage.setItem(k, JSON.stringify(v)); } catch(_){} }
function escapeHtml(v) { return String(v).replace(/[&<>"']/g, c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c])); }
function truncateTitle(t) { return t.length>60?t.slice(0,57)+'...':t; }
function stepId(l,s) { return `${l}-${s}`; }
function quizKey(l,s,q) { return `${l}_${s}_${q}`; }
function getNodeByPath(nodeId) {
  const parts=String(nodeId).split('.').map(Number);
  let node=courseData[parts[0]];
  for(let i=1;i<parts.length;i++) node=node&&node.children[parts[i]];
  return node||null;
}
function isValidStep(nodeId,sIdx) { const n=getNodeByPath(nodeId); return !!(n&&n.steps&&n.steps[sIdx]); }

// ── Progress ──────────────────────────────────────────────────────
function markStepCompleted(l,s) {
  const id=stepId(l,s);
  if(!completedSteps.has(id)){completedSteps.add(id);writeJSON(COURSE_KEY+"_progress",[...completedSteps]);}
}
function updateProgressUI() {
  const pct = TOTAL_STEPS ? (completedSteps.size/TOTAL_STEPS*100) : 0;
  const el = document.getElementById('progressText');
  const bar = document.getElementById('progressFill');
  if(el) el.textContent = `Пройдено ${completedSteps.size}/${TOTAL_STEPS}`;
  if(bar) bar.style.width = pct+'%';
}

// ── Sidebar ───────────────────────────────────────────────────────
function renderLessonNode(node, nodeId, depth) {
  const open = openLessons.has(nodeId);
  const ml = depth * 12;
  let stepsHtml = '';
  if (open && node.steps && node.steps.length) {
    stepsHtml = `<div class="lesson-steps">${node.steps.map((step,sIdx)=>{
      const isActive = String(nodeId)===String(currentLessonIdx) && sIdx===currentStepIdx;
      const isDone = completedSteps.has(stepId(nodeId,sIdx));
      return `<button class="step-item ${isActive?'active':''} ${isDone?'completed':''}"
        onclick="loadStepByNodeId('${nodeId}',${sIdx})"
        title="${escapeHtml(step.title)}">
        <span class="step-dot"></span>${escapeHtml(truncateTitle(step.title))}
      </button>`;
    }).join('')}</div>`;
  }
  let childrenHtml = '';
  if (open && node.children && node.children.length) {
    childrenHtml = node.children.map((c,ci)=>renderLessonNode(c,nodeId+'.'+ci,depth+1)).join('');
  }
  return `<section class="lesson" style="margin-left:${ml}px">
    <button class="lesson-header" onclick="toggleLessonNode('${nodeId}')">
      <svg class="lesson-arrow ${open?'open':''}" viewBox="0 0 24 24" fill="currentColor">
        <path d="M8 5l8 7-8 7V5z"/>
      </svg>
      <span class="lesson-name">${escapeHtml(node.lesson)}</span>
    </button>
    ${stepsHtml}${childrenHtml}
  </section>`;
}

function renderSidebar() {
  const sidebar = document.getElementById('sidebar');
  const body = sidebar?.querySelector('.sidebar-body');
  const scrollTop = body ? body.scrollTop : 0;

  const storedTheme = getStoredTheme();
  const lessonsHtml = courseData.map((l,i)=>renderLessonNode(l,String(i),0)).join('');

  sidebar.innerHTML = `
    <div class="sidebar-header">
      <h1 class="course-title">${escapeHtml(COURSE_TITLE)}</h1>
      <div class="progress-row">
        <span class="progress-text" id="progressText">Пройдено 0/${TOTAL_STEPS}</span>
        <div class="progress-bar-wrap"><div class="progress-bar-fill" id="progressFill"></div></div>
      </div>
    </div>
    <div class="sidebar-body">${lessonsHtml}</div>
    <div class="sidebar-footer">
      <button class="sidebar-settings-btn" id="settingsOpenBtn">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
        </svg>
        Настройки
      </button>
    </div>`;

  const newBody = sidebar.querySelector('.sidebar-body');
  if (newBody) newBody.scrollTop = scrollTop;

  updateProgressUI();
  initSettings();
}

function toggleLessonNode(nodeId) {
  openLessons.has(nodeId) ? openLessons.delete(nodeId) : openLessons.add(nodeId);
  renderSidebar();
}

// ── Flat list for prev/next ───────────────────────────────────────
let _flatList = [];
function buildFlatList() {
  _flatList = [];
  function walk(node, nodeId) {
    (node.steps||[]).forEach((_,sIdx)=>_flatList.push({nodeId,sIdx}));
    (node.children||[]).forEach((c,ci)=>walk(c,nodeId+'.'+ci));
  }
  courseData.forEach((l,i)=>walk(l,String(i)));
}
function currentFlatIdx() { return _flatList.findIndex(x=>String(x.nodeId)===String(currentLessonIdx)&&x.sIdx===currentStepIdx); }
function updateNavButtons() {
  const f = currentFlatIdx();
  const prev = document.getElementById('prevBtn');
  const next = document.getElementById('nextBtn');
  if(prev) prev.disabled = f<=0;
  if(next) next.disabled = f>=_flatList.length-1;
}
function navigateFlat(dir) {
  const f=currentFlatIdx(), nx=_flatList[f+dir];
  if(nx) loadStepByNodeId(nx.nodeId,nx.sIdx);
}
function scrollSidebarToActive() {
  const a=document.querySelector('.step-item.active');
  if(a) a.scrollIntoView({behavior:'smooth',block:'nearest'});
}

// ── Load Step ─────────────────────────────────────────────────────
function loadStepByNodeId(nodeId, sIdx) {
  const node = getNodeByPath(nodeId);
  if(!node||!node.steps||!node.steps[sIdx]) return;
  currentLessonIdx=nodeId; currentStepIdx=sIdx;
  openLessons.add(nodeId);
  writeJSON(COURSE_KEY+'_last_viewed',{l:nodeId,s:sIdx});
  markStepCompleted(nodeId,sIdx);
  renderSidebar();

  const step = node.steps[sIdx];
  const quizzes = Array.isArray(step.quizzes)?step.quizzes:[];

  document.getElementById('contentWrap').innerHTML = `
    <div class="content-card">
      <h1 class="step-title-main">${escapeHtml(step.title)}</h1>
      <div class="step-content">${step.html||''}</div>
      ${renderQuizzes(quizzes,nodeId,sIdx)}
      <div class="nav-row">
        <button id="prevBtn" class="nav-btn secondary" onclick="navigateFlat(-1)">← Назад</button>
        <button id="nextBtn" class="nav-btn" onclick="navigateFlat(1)">Вперёд →</button>
      </div>
    </div>`;

  quizzes.forEach((q,qi)=>{
    const k=quizKey(nodeId,sIdx,qi);
    if(Object.prototype.hasOwnProperty.call(savedQuizzes,k)) applyAnswerUI(qi,savedQuizzes[k],q.answer);
  });

  updateNavButtons();
  document.getElementById('main')?.scrollTo({top:0,behavior:'smooth'});
  setTimeout(scrollSidebarToActive,200);
  setTimeout(initAllPlayers,60);
}

// ── Quizzes ───────────────────────────────────────────────────────
function renderQuizzes(quizzes,l,s) {
  if(!quizzes||!quizzes.length) return '';
  return `<section class="quiz-section">
    <div class="quiz-header">
      <h2 class="quiz-title">Проверь себя</h2>
      <button class="quiz-reset" onclick="resetQuizzes('${l}',${s})">Сбросить</button>
    </div>
    <div class="quiz-list">${quizzes.map((q,qi)=>`
      <div class="quiz-card" data-quiz-index="${qi}">
        <p class="quiz-question">${escapeHtml(q.question)}</p>
        <div class="quiz-options">${(q.options||[]).map((o,oi)=>`
          <label class="quiz-option" data-option-index="${oi}">
            <input type="radio" name="q_${l}_${s}_${qi}" onclick="checkAnswer('${l}',${s},${qi},${oi},${q.answer})">
            ${escapeHtml(o)}
          </label>`).join('')}
        </div>
        <p class="quiz-message"></p>
      </div>`).join('')}
    </div>
  </section>`;
}
function applyAnswerUI(qi,si,ci) {
  const card=document.querySelector(`.quiz-card[data-quiz-index="${qi}"]`);
  if(!card) return;
  card.querySelectorAll('.quiz-option').forEach(o=>o.classList.remove('correct','incorrect'));
  card.querySelectorAll('input[type="radio"]').forEach((r,i)=>{r.disabled=true;r.checked=i===si;});
  card.classList.remove('correct','incorrect');
  const msg=card.querySelector('.quiz-message');
  if(si===ci){
    card.classList.add('correct');
    card.querySelector(`.quiz-option[data-option-index="${si}"]`)?.classList.add('correct');
    if(msg)msg.textContent='✓ Правильно!';
  } else {
    card.classList.add('incorrect');
    card.querySelector(`.quiz-option[data-option-index="${si}"]`)?.classList.add('incorrect');
    card.querySelector(`.quiz-option[data-option-index="${ci}"]`)?.classList.add('correct');
    if(msg)msg.textContent='✗ Неверно. Правильный ответ выделен.';
  }
}
function checkAnswer(l,s,q,si,ci){savedQuizzes[quizKey(l,s,q)]=si;writeJSON(COURSE_KEY+'_quizzes',savedQuizzes);applyAnswerUI(q,si,ci);}
function resetQuizzes(l,s){const p=`${l}_${s}_`;Object.keys(savedQuizzes).forEach(k=>{if(k.startsWith(p))delete savedQuizzes[k];});writeJSON(COURSE_KEY+'_quizzes',savedQuizzes);loadStepByNodeId(l,parseInt(s));}

// ── Init ──────────────────────────────────────────────────────────
function init() {
  // Apply theme
  const stored = getStoredTheme();
  applyThemeSetting(stored);

  buildFlatList();
  renderSidebar();

  if(!courseData.length||!TOTAL_STEPS){
    document.getElementById('contentWrap').innerHTML='<div class="empty-state">В курсе нет шагов для отображения.</div>';
    return;
  }

  const saved = readJSON(COURSE_KEY+'_last_viewed', null);
  const nodeId = saved&&saved.l!=null ? String(saved.l) : '0';
  const ss = saved&&Number.isInteger(saved.s) ? saved.s : 0;

  if(isValidStep(nodeId,ss)) loadStepByNodeId(nodeId,ss);
  else if(_flatList.length) loadStepByNodeId(_flatList[0].nodeId,_flatList[0].sIdx);
}

init();
