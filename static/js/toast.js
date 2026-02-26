/* Toast Notification System */
(function() {
  let container = null;

  function getContainer() {
    if (!container) {
      container = document.createElement('div');
      container.id = 'toast-container';
      container.style.cssText = 'position:fixed;top:1.5rem;right:1.5rem;z-index:10000;display:flex;flex-direction:column;gap:0.75rem;pointer-events:none;';
      document.body.appendChild(container);
    }
    return container;
  }

  const icons = {
    success: '<svg width="20" height="20" viewBox="0 0 20 20" fill="none"><circle cx="10" cy="10" r="10" fill="#2ecc71"/><path d="M6 10l3 3 5-6" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>',
    error: '<svg width="20" height="20" viewBox="0 0 20 20" fill="none"><circle cx="10" cy="10" r="10" fill="#e74c3c"/><path d="M7 7l6 6M13 7l-6 6" stroke="#fff" stroke-width="2" stroke-linecap="round"/></svg>',
    info: '<svg width="20" height="20" viewBox="0 0 20 20" fill="none"><circle cx="10" cy="10" r="10" fill="#4a90d9"/><path d="M10 9v5M10 6.5v.01" stroke="#fff" stroke-width="2" stroke-linecap="round"/></svg>'
  };

  window.showToast = function(message, type) {
    type = type || 'info';
    const c = getContainer();

    const toast = document.createElement('div');
    toast.style.cssText = 'pointer-events:auto;display:flex;align-items:center;gap:0.75rem;padding:0.85rem 1.25rem;border-radius:10px;background:rgba(26,26,36,0.95);backdrop-filter:blur(12px);border:1px solid rgba(255,255,255,0.1);color:#f0f0f5;font-size:0.92rem;font-family:Inter,sans-serif;box-shadow:0 8px 32px rgba(0,0,0,0.3);min-width:280px;max-width:420px;transform:translateX(120%);transition:transform 0.35s cubic-bezier(0.4,0,0.2,1),opacity 0.35s ease;opacity:0;';

    const icon = document.createElement('span');
    icon.innerHTML = icons[type] || icons.info;
    icon.style.cssText = 'flex-shrink:0;display:flex;';

    const text = document.createElement('span');
    text.textContent = message;
    text.style.cssText = 'flex:1;line-height:1.4;';

    const close = document.createElement('button');
    close.innerHTML = '&times;';
    close.style.cssText = 'background:none;border:none;color:#8888aa;font-size:1.3rem;cursor:pointer;padding:0 0 0 0.5rem;line-height:1;flex-shrink:0;';
    close.onclick = function() { dismiss(toast); };

    toast.appendChild(icon);
    toast.appendChild(text);
    toast.appendChild(close);
    c.appendChild(toast);

    requestAnimationFrame(function() {
      requestAnimationFrame(function() {
        toast.style.transform = 'translateX(0)';
        toast.style.opacity = '1';
      });
    });

    const timer = setTimeout(function() { dismiss(toast); }, 4000);
    toast._timer = timer;
  };

  function dismiss(toast) {
    if (toast._dismissed) return;
    toast._dismissed = true;
    clearTimeout(toast._timer);
    toast.style.transform = 'translateX(120%)';
    toast.style.opacity = '0';
    setTimeout(function() {
      if (toast.parentNode) toast.parentNode.removeChild(toast);
    }, 400);
  }
})();
