// ══ Custom Video Player ══════════════════════════════════════════
const VP_SPEEDS = [0.75, 1, 1.25, 1.5, 2];

function vpFmt(s) {
  s = Math.floor(s || 0);
  const h = Math.floor(s / 3600), m = Math.floor((s % 3600) / 60), sec = s % 60;
  return h > 0
    ? `${h}:${String(m).padStart(2,'0')}:${String(sec).padStart(2,'0')}`
    : `${String(m).padStart(2,'0')}:${String(sec).padStart(2,'0')}`;
}
function vpSavePos(k, t) { try { localStorage.setItem('vp_'+k, String(t)); } catch(_){} }
function vpLoadPos(k) { try { return parseFloat(localStorage.getItem('vp_'+k))||0; } catch(_){ return 0; } }

function initVideoPlayer(wrapper) {
  const video = wrapper.querySelector('video');
  if (!video) return;

  wrapper.classList.add('vp-wrap');
  const vpKey = (video.querySelector('source')?.src || video.src || '').split('/').pop().replace(/\?.*$/,'');
  const origSources = [...video.querySelectorAll('source')].map(s => ({src: s.src, type: s.type}));
  const origSrc = video.src;

  wrapper.innerHTML = `
    <video tabindex="-1"></video>
    <div class="vp-overlay"><div class="vp-ripple" id="vpRipple"><svg viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg></div></div>
    <div class="vp-end-overlay" id="vpEndOverlay">
      <button class="vp-end-next" id="vpEndNextBtn" title="Следующий урок">
        <svg viewBox="0 0 24 24" fill="currentColor" width="28" height="28"><path d="M6 18l8.5-6L6 6v12zm2.5-6 5.5 3.9V8.1L8.5 12zM16 6h2v12h-2z"/></svg>
        <svg class="vp-end-ring" viewBox="0 0 44 44">
          <circle class="vp-end-ring-bg" cx="22" cy="22" r="19"/>
          <circle class="vp-end-ring-fill" id="vpEndRingFill" cx="22" cy="22" r="19"/>
        </svg>
      </button>
      <div class="vp-end-info">
        <button class="vp-end-label" id="vpEndLabel">Следующий урок</button>
        <span class="vp-end-countdown" id="vpEndCountdown">откроется через 3...</span>
      </div>
      <label class="vp-end-auto">
        <span class="vp-end-toggle-wrap">
          <input type="checkbox" id="vpAutoNext" class="vp-end-toggle-input">
          <span class="vp-end-toggle-track"><span class="vp-end-toggle-thumb"></span></span>
        </span>
        <span class="vp-end-auto-label" id="vpAutoLabel">Автопереход включён</span>
      </label>
    </div>
    <div class="vp-toast" id="vpToast"></div>
    <div class="vp-controls">
      <div class="vp-progress" id="vpProg">
        <div class="vp-progress-buf" id="vpBuf"></div>
        <div class="vp-progress-fill" id="vpFill"></div>
        <div class="vp-progress-thumb" id="vpThumb"></div>
        <div class="vp-time-tip" id="vpTip">0:00</div>
      </div>
      <div class="vp-row">
        <button class="vp-btn vp-step-btn" id="vpPrevStepBtn" title="Предыдущий урок (Shift+←)">
          <svg viewBox="0 0 24 24" fill="currentColor"><path d="M6 6h2v12H6zm3.5 6 8.5 6V6z"/></svg>
        </button>
        <button class="vp-btn" id="vpPlayBtn" title="Пауза/Воспроизведение (Пробел)">
          <svg id="vpPlayIco" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
        </button>
        <button class="vp-btn vp-step-btn" id="vpNextStepBtn" title="Следующий урок (Shift+→)">
          <svg viewBox="0 0 24 24" fill="currentColor"><path d="M6 18l8.5-6L6 6v12zm2.5-6 5.5 3.9V8.1L8.5 12zM16 6h2v12h-2z"/></svg>
        </button>
        <button class="vp-btn" id="vpRewBtn" title="−10 сек (←)">
          <svg viewBox="0 0 24 24"><path d="M12.5 8c-2.65 0-5.05 1-6.9 2.6L2 7v9h9l-3.62-3.62c1.39-1.16 3.16-1.88 5.12-1.88 3.54 0 6.55 2.31 7.6 5.5l2.37-.78C21.08 11.03 17.15 8 12.5 8z"/><text x="12" y="20" font-size="6" text-anchor="middle" fill="currentColor" font-family="sans-serif">10</text></svg>
        </button>
        <button class="vp-btn" id="vpFwdBtn" title="+10 сек (→)">
          <svg viewBox="0 0 24 24"><path d="M18.4 10.6C16.55 9 14.15 8 11.5 8c-4.65 0-8.58 3.03-9.96 7.22L3.9 16c1.05-3.19 4.05-5.5 7.6-5.5 1.95 0 3.73.72 5.12 1.88L13 16h9V7l-3.6 3.6z"/><text x="12" y="20" font-size="6" text-anchor="middle" fill="currentColor" font-family="sans-serif">10</text></svg>
        </button>
        <div class="vp-vol-wrap">
          <button class="vp-btn" id="vpMuteBtn" title="Звук (M)">
            <svg id="vpVolIco" viewBox="0 0 24 24"><path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02z"/></svg>
          </button>
          <input class="vp-vol-slider" id="vpVol" type="range" min="0" max="1" step="0.05" value="1">
        </div>
        <span class="vp-time" id="vpTime">00:00 / 00:00</span>
        <div class="vp-spacer"></div>
        <div class="vp-speed-wrap">
          <button class="vp-speed-btn" id="vpSpeedBtn">1x</button>
          <div class="vp-speed-menu" id="vpSpeedMenu">
            ${VP_SPEEDS.map(s=>`<div class="vp-speed-item${s===1?' selected':''}" data-spd="${s}">${s}x</div>`).join('')}
          </div>
        </div>
        <button class="vp-btn" id="vpFsBtn" title="Полный экран (F)">
          <svg viewBox="0 0 24 24"><path d="M7 14H5v5h5v-2H7v-3zm-2-4h2V7h3V5H5v5zm12 7h-3v2h5v-5h-2v3zM14 5v2h3v3h2V5h-5z"/></svg>
        </button>
      </div>
    </div>`;

  const vid = wrapper.querySelector('video');
  if (origSources.length) {
    origSources.forEach(s => { const el = document.createElement('source'); el.src = s.src; el.type = s.type; vid.appendChild(el); });
  } else if (origSrc) { vid.src = origSrc; }

  const $ = id => wrapper.querySelector('#'+id);
  const ripple=$('vpRipple'),toast=$('vpToast'),prog=$('vpProg'),buf=$('vpBuf'),
        fill=$('vpFill'),thumb=$('vpThumb'),tip=$('vpTip'),playBtn=$('vpPlayBtn'),
        playIco=$('vpPlayIco'),rewBtn=$('vpRewBtn'),fwdBtn=$('vpFwdBtn'),
        muteBtn=$('vpMuteBtn'),volIco=$('vpVolIco'),volSlider=$('vpVol'),
        timeEl=$('vpTime'),speedBtn=$('vpSpeedBtn'),speedMenu=$('vpSpeedMenu'),
        fsBtn=$('vpFsBtn'),prevStepBtn=$('vpPrevStepBtn'),nextStepBtn=$('vpNextStepBtn');

  let toastTimer, uiTimer, scrubbing=false, curSpeed=1;

  function showToast(msg) { toast.textContent=msg; toast.classList.add('show'); clearTimeout(toastTimer); toastTimer=setTimeout(()=>toast.classList.remove('show'),1100); }
  function showUI() { wrapper.classList.add('show-ui'); clearTimeout(uiTimer); uiTimer=setTimeout(()=>wrapper.classList.remove('show-ui'),2400); }
  function showRipple(icon) { ripple.querySelector('svg').innerHTML=icon; ripple.classList.remove('show'); void ripple.offsetWidth; ripple.classList.add('show'); setTimeout(()=>ripple.classList.remove('show'),380); }

  function updatePlayIcon() {
    playIco.innerHTML = vid.paused ? '<path d="M8 5v14l11-7z"/>' : '<path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/>';
    wrapper.classList.toggle('paused', vid.paused);
  }
  function updateVolIcon() {
    const m = vid.muted||vid.volume===0;
    volIco.innerHTML = m
      ? '<path d="M16.5 12c0-1.77-1.02-3.29-2.5-4.03v2.21l2.45 2.45c.03-.2.05-.41.05-.63zm2.5 0c0 .94-.2 1.82-.54 2.64l1.51 1.51C20.63 14.91 21 13.5 21 12c0-4.28-2.99-7.86-7-8.77v2.06c2.89.86 5 3.54 5 6.71zM4.27 3L3 4.27 7.73 9H3v6h4l5 5v-6.73l4.25 4.25c-.67.52-1.42.93-2.25 1.18v2.06c1.38-.31 2.63-.95 3.69-1.81L19.73 21 21 19.73l-9-9L4.27 3zM12 4L9.91 6.09 12 8.18V4z"/>'
      : '<path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM14 3.23v2.06c2.89.86 5 3.54 5 6.71s-2.11 5.85-5 6.71v2.06c4.01-.91 7-4.49 7-8.77s-2.99-7.86-7-8.77z"/>';
  }
  function updateProgress() {
    if (!vid.duration) return;
    const pct = vid.currentTime/vid.duration*100;
    fill.style.width=pct+'%'; thumb.style.left=pct+'%';
    timeEl.textContent=vpFmt(vid.currentTime)+' / '+vpFmt(vid.duration);
    vpSavePos(vpKey,vid.currentTime);
  }
  function updateBuffer() {
    if (!vid.duration||!vid.buffered.length) return;
    buf.style.width=(vid.buffered.end(vid.buffered.length-1)/vid.duration*100)+'%';
  }

  vid.addEventListener('loadedmetadata',()=>{
    const saved=vpLoadPos(vpKey);
    if(saved>2&&saved<vid.duration-3){vid.currentTime=saved;showToast('Продолжаем с '+vpFmt(saved));}
    updateProgress();
  });
  vid.addEventListener('timeupdate',updateProgress);
  vid.addEventListener('progress',updateBuffer);
  vid.addEventListener('play',updatePlayIcon);
  vid.addEventListener('pause',updatePlayIcon);
  vid.addEventListener('volumechange',updateVolIcon);

  vid.addEventListener('click',()=>{
    if(endOverlay?.classList.contains('show')) return; // ignore clicks when end overlay is visible
    if(vid.paused){vid.play();showRipple('<path d="M8 5v14l11-7z"/>');}
    else{vid.pause();showRipple('<path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/>');}
    showUI();
  });

  playBtn.addEventListener('click',e=>{e.stopPropagation();vid.paused?vid.play():vid.pause();});
  rewBtn.addEventListener('click',e=>{e.stopPropagation();vid.currentTime=Math.max(0,vid.currentTime-10);showToast('⏪ −10 сек');showUI();});
  fwdBtn.addEventListener('click',e=>{e.stopPropagation();vid.currentTime=Math.min(vid.duration||0,vid.currentTime+10);showToast('⏩ +10 сек');showUI();});

  muteBtn.addEventListener('click',e=>{
    e.stopPropagation();
    vid.muted=!vid.muted;
    volSlider.value=vid.muted?0:(vid.volume||1);
    if(!vid.muted&&!vid.volume)vid.volume=1;
    showToast(vid.muted?'🔇 Звук выключен':'🔊 Звук включён');
  });
  volSlider.addEventListener('input',e=>{e.stopPropagation();vid.volume=parseFloat(volSlider.value);vid.muted=vid.volume===0;});

  speedBtn.addEventListener('click',e=>{e.stopPropagation();speedMenu.classList.toggle('open');});
  speedMenu.querySelectorAll('.vp-speed-item').forEach(item=>{
    item.addEventListener('click',e=>{
      e.stopPropagation();
      curSpeed=parseFloat(item.dataset.spd);
      vid.playbackRate=curSpeed;
      speedBtn.textContent=curSpeed+'x';
      speedMenu.querySelectorAll('.vp-speed-item').forEach(i=>i.classList.remove('selected'));
      item.classList.add('selected');
      speedMenu.classList.remove('open');
      showToast('Скорость: '+curSpeed+'x');
    });
  });
  document.addEventListener('click',()=>speedMenu.classList.remove('open'));

  function seekTo(e) {
    const r=prog.getBoundingClientRect();
    vid.currentTime=Math.max(0,Math.min(1,(e.clientX-r.left)/r.width))*(vid.duration||0);
  }
  prog.addEventListener('mousedown',e=>{scrubbing=true;seekTo(e);});
  prog.addEventListener('mousemove',e=>{
    const r=prog.getBoundingClientRect();
    const pct=Math.max(0,Math.min(1,(e.clientX-r.left)/r.width));
    tip.textContent=vpFmt(pct*(vid.duration||0));
    tip.style.left=(pct*100)+'%';
    if(scrubbing)seekTo(e);
  });
  document.addEventListener('mouseup',()=>{scrubbing=false;});

  fsBtn.addEventListener('click',e=>{e.stopPropagation();document.fullscreenElement?document.exitFullscreen?.():wrapper.requestFullscreen?.();});

  function doStepNav(dir) {
    if(typeof vpNavigateStep==='function') vpNavigateStep(dir, wrapper);
  }

  // ── End-of-video overlay ──────────────────────────────────────
  const endOverlay = $('vpEndOverlay');
  const endNextBtn = $('vpEndNextBtn');
  const endLabel   = $('vpEndLabel');
  const endCountdown = $('vpEndCountdown');
  const autoToggle = $('vpAutoNext');
  const autoLabel  = $('vpAutoLabel');
  const ringFill   = $('vpEndRingFill');

  const AUTO_KEY = 'vp_autonext';
  let countdownTimer = null;
  const COUNTDOWN = 3;
  // circumference of circle r=19
  const CIRC = 2 * Math.PI * 19;
  const RING_DURATION_MS = COUNTDOWN * 1000;

  function loadAutoNext() {
    try { return localStorage.getItem(AUTO_KEY) !== 'off'; } catch(_){ return true; }
  }
  function saveAutoNext(v) { try { localStorage.setItem(AUTO_KEY, v ? 'on' : 'off'); } catch(_){} }

  function hasNextStep() {
    if(typeof _flatList === 'undefined') return false;
    const f = typeof currentFlatIdx === 'function' ? currentFlatIdx() : -1;
    return f >= 0 && f < _flatList.length - 1;
  }

  function showEndOverlay() {
    if(!endOverlay || !hasNextStep()) return;
    endOverlay.classList.add('show');
    const auto = loadAutoNext();
    autoToggle.checked = auto;
    autoLabel.textContent = auto ? 'Автопереход включён' : 'Автопереход выключен';
    // Reset ring to empty without transition, then animate via countdown
    if(ringFill) {
      ringFill.style.transition = 'none';
      ringFill.style.strokeDasharray = CIRC;
      ringFill.style.strokeDashoffset = CIRC;
      // Force reflow so transition:none takes effect before we re-enable it
      void ringFill.getBoundingClientRect();
      ringFill.style.transition = 'stroke-dashoffset 1s linear';
    }
    if(auto) startCountdown();
  }

  function hideEndOverlay() {
    endOverlay?.classList.remove('show');
    clearCountdown();
  }

  function startCountdown() {
    // Animate ring via CSS: set target offset=0 with transition=COUNTDOWN s
    if(ringFill) {
      ringFill.style.transition = `stroke-dashoffset ${RING_DURATION_MS}ms linear`;
      ringFill.style.strokeDashoffset = '0';
    }
    // Update text countdown with setInterval every second
    let remaining = COUNTDOWN;
    updateCountdownText(remaining);
    const textTimer = setInterval(() => {
      remaining--;
      updateCountdownText(remaining);
    }, 1000);
    // Navigate exactly when animation ends
    countdownTimer = setTimeout(() => {
      clearInterval(textTimer);
      doStepNav(1);
    }, RING_DURATION_MS);
  }

  function updateCountdownText(rem) {
    if(endCountdown) endCountdown.textContent = rem > 0 ? `откроется через ${rem}...` : '';
  }

  function updateCountdownUI(rem) { updateCountdownText(rem); }

  function clearCountdown() {
    clearTimeout(countdownTimer);
    countdownTimer = null;
    // Also stop ring animation mid-way
    if(ringFill) {
      const computed = window.getComputedStyle(ringFill).strokeDashoffset;
      ringFill.style.transition = 'none';
      ringFill.style.strokeDashoffset = computed;
    }
  }

  if(endNextBtn) endNextBtn.addEventListener('click', e => { e.stopPropagation(); hideEndOverlay(); doStepNav(1); });
  if(endLabel)   endLabel.addEventListener('click',   e => { e.stopPropagation(); hideEndOverlay(); doStepNav(1); });
  // Click on overlay background → cancel autoplay, hide overlay
  if(endOverlay) endOverlay.addEventListener('click', e => {
    if(e.target === endOverlay) { clearCountdown(); hideEndOverlay(); }
  });

  if(autoToggle) autoToggle.addEventListener('change', () => {
    const on = autoToggle.checked;
    saveAutoNext(on);
    autoLabel.textContent = on ? 'Автопереход включён' : 'Автопереход выключен';
    if(on) startCountdown(); else clearCountdown();
  });

  vid.addEventListener('ended', () => { showEndOverlay(); });
  vid.addEventListener('play',  () => { hideEndOverlay(); });

  wrapper.addEventListener('vp-step-changed', e => {
    showToast('▶ ' + (e.detail.title||''));
    showUI();
    updatePlayIcon();
  });
  if(prevStepBtn) prevStepBtn.addEventListener('click',e=>{e.stopPropagation();doStepNav(-1);});
  if(nextStepBtn) nextStepBtn.addEventListener('click',e=>{e.stopPropagation();doStepNav(1);});
  wrapper.addEventListener('mousemove',showUI);

  document.addEventListener('keydown',e=>{
    const tag=document.activeElement?.tagName;
    if(tag==='INPUT'||tag==='TEXTAREA'||tag==='SELECT')return;
    if(wrapper.getBoundingClientRect().width===0)return;
    switch(e.key){
      case ' ':case 'k':case 'K':
        e.preventDefault();
        if(vid.paused){vid.play();showToast('▶ Воспроизведение');}else{vid.pause();showToast('⏸ Пауза');}
        showUI();break;
      case 'ArrowLeft':if(e.shiftKey){e.preventDefault();doStepNav(-1);}else{e.preventDefault();vid.currentTime=Math.max(0,vid.currentTime-10);showToast('⏪ −10 сек');}showUI();break;
      case 'ArrowRight':if(e.shiftKey){e.preventDefault();doStepNav(1);}else{e.preventDefault();vid.currentTime=Math.min(vid.duration||0,vid.currentTime+10);showToast('⏩ +10 сек');}showUI();break;
      case 'ArrowUp':e.preventDefault();vid.volume=Math.min(1,vid.volume+0.1);vid.muted=false;volSlider.value=vid.volume;showToast('🔊 '+Math.round(vid.volume*100)+'%');showUI();break;
      case 'ArrowDown':e.preventDefault();vid.volume=Math.max(0,vid.volume-0.1);volSlider.value=vid.volume;showToast('🔉 '+Math.round(vid.volume*100)+'%');showUI();break;
      case 'm':case 'M':vid.muted=!vid.muted;if(!vid.muted&&!vid.volume){vid.volume=1;volSlider.value=1;}showToast(vid.muted?'🔇 Звук выключен':'🔊 Звук включён');showUI();break;
      case 'f':case 'F':document.fullscreenElement?document.exitFullscreen?.():wrapper.requestFullscreen?.();break;
      case '>':{const si=VP_SPEEDS.indexOf(curSpeed);if(si<VP_SPEEDS.length-1){curSpeed=VP_SPEEDS[si+1];vid.playbackRate=curSpeed;speedBtn.textContent=curSpeed+'x';showToast('Скорость: '+curSpeed+'x');}}break;
      case '<':{const si=VP_SPEEDS.indexOf(curSpeed);if(si>0){curSpeed=VP_SPEEDS[si-1];vid.playbackRate=curSpeed;speedBtn.textContent=curSpeed+'x';showToast('Скорость: '+curSpeed+'x');}}break;
    }
  });
}

function initAllPlayers() {
  document.querySelectorAll('.video-wrapper').forEach(w=>{
    if(!w.dataset.vpInit){w.dataset.vpInit='1';initVideoPlayer(w);}
  });
}
