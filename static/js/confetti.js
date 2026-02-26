/* Confetti wrapper — loads canvas-confetti from CDN and provides a helper */
(function() {
  var script = document.createElement('script');
  script.src = 'https://cdn.jsdelivr.net/npm/canvas-confetti@1.9.2/dist/confetti.browser.min.js';
  script.async = true;
  document.head.appendChild(script);

  window.fireConfetti = function() {
    if (typeof confetti === 'function') {
      confetti({ particleCount: 150, spread: 80, origin: { y: 0.6 } });
      setTimeout(function() {
        confetti({ particleCount: 80, spread: 100, origin: { y: 0.7 } });
      }, 300);
    }
  };
})();
