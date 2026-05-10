// ══ Settings Modal ══════════════════════════════════════════════
function initSettings() {
  const backdrop = document.getElementById('settingsBackdrop');
  if (!backdrop) return;

  const openBtn = document.getElementById('settingsOpenBtn');
  const closeBtn = document.getElementById('settingsCloseBtn');
  const select = document.getElementById('themeSelect');

  if (openBtn) {
    openBtn.addEventListener('click', () => {
      backdrop.classList.add('open');
      if (select) select.value = getStoredTheme();
    });
  }

  if (closeBtn) {
    closeBtn.addEventListener('click', () => backdrop.classList.remove('open'));
  }

  backdrop.addEventListener('click', e => {
    if (e.target === backdrop) backdrop.classList.remove('open');
  });

  document.addEventListener('keydown', e => {
    if (e.key === 'Escape') backdrop.classList.remove('open');
  });

  if (select) {
    select.addEventListener('change', () => {
      const theme = select.value;
      applyThemeSetting(theme);
      try { localStorage.setItem(COURSE_KEY + '_theme', theme); } catch(_){}
    });
  }
}

function applyThemeSetting(theme) {
  document.documentElement.setAttribute('data-theme', theme);
}

function getStoredTheme() {
  try { return localStorage.getItem(COURSE_KEY + '_theme') || 'dark'; } catch(_){ return 'dark'; }
}
