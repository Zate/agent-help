/* copy.js — auto-wires clipboard copy buttons to every <pre> block */
(function () {
  'use strict';

  var SVG_COPY = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>';
  var SVG_CHECK = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="20 6 9 17 4 12"/></svg>';

  function wrapAndBind(pre) {
    // Don't double-wrap
    if (pre.parentNode && pre.parentNode.classList.contains('pre-wrap')) return;
    // Skip prompt pre blocks that already have a dedicated copy button sibling
    if (pre.nextElementSibling && pre.nextElementSibling.classList.contains('prompt-copy-btn')) return;

    var wrap = document.createElement('div');
    wrap.className = 'pre-wrap';
    pre.parentNode.insertBefore(wrap, pre);
    wrap.appendChild(pre);

    var btn = document.createElement('button');
    btn.className = 'copy-btn';
    btn.setAttribute('aria-label', 'Copy to clipboard');
    btn.innerHTML = SVG_COPY;
    wrap.appendChild(btn);

    btn.addEventListener('click', function () {
      var text = pre.textContent || '';
      if (!navigator.clipboard) return;
      navigator.clipboard.writeText(text.trim()).then(function () {
        btn.innerHTML = SVG_CHECK;
        btn.classList.add('copied');
        setTimeout(function () {
          btn.innerHTML = SVG_COPY;
          btn.classList.remove('copied');
        }, 1600);
      });
    });
  }

  function init() {
    document.querySelectorAll('pre').forEach(wrapAndBind);
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
