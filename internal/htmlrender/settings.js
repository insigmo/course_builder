// ══ Settings Modal ══════════════════════════════════════════════
function initSettings() {
  const backdrop = document.getElementById('settingsBackdrop');
  if (!backdrop) return;

  // Open
  document.getElementById('settingsOpenBtn')?.addEventListener('click', () => {
    backdrop.classList.add('open');
  });

  // Close
  backdrop.addEventListener('click', e => {
    if (e.target === backdrop) backdrop.classList.remove('open');
  });
  document.getElementById('settingsCloseBtn')?.addEventListener('click', () => {
    backdrop.classList.remove('open');
  });
  document.addEventListener('keydown', e => {
    if (e.key === 'Escape') backdrop.classList.remove('open');
  });

  // Theme options
  backdrop.querySelectorAll('.theme-option').forEach(opt => {
    opt.addEventListener('click', () => {
      const theme = opt.dataset.theme;
      applyThemeSetting(theme);
      backdrop.querySelectorAll('.theme-option').forEach(o => o.classList.remove('selected'));
      opt.classList.add('selected');
      try { localStorage.setItem(COURSE_KEY + '_theme', theme); } catch(_){}
    });
  });
}

function applyThemeSetting(theme) {
  if (theme === 'system') {
    const sys = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    document.documentElement.setAttribute('data-theme', sys);
  } else {
    document.documentElement.setAttribute('data-theme', theme);
  }
}

function getStoredTheme() {
  try { return localStorage.getItem(COURSE_KEY + '_theme') || 'system'; } catch(_){ return 'system'; }
}
