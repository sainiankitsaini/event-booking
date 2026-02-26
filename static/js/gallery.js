/* Image Gallery Component */
function initGallery(containerId, images) {
  const container = document.getElementById(containerId);
  if (!container || !images || images.length === 0) {
    if (container) container.innerHTML = '<div style="padding:2rem;text-align:center;color:#8888aa;">No images available</div>';
    return;
  }

  let currentIndex = 0;

  const html = `
    <div class="gallery-wrapper">
      <div class="gallery-main">
        <img id="gallery-featured" src="/uploads/${images[0]}" alt="Event image" class="gallery-featured-img" />
        <div class="gallery-counter" id="gallery-counter">1 / ${images.length}</div>
        <button class="gallery-nav gallery-prev" id="gallery-prev" aria-label="Previous">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none"><path d="M15 18l-6-6 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </button>
        <button class="gallery-nav gallery-next" id="gallery-next" aria-label="Next">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none"><path d="M9 18l6-6-6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </button>
      </div>
      <div class="gallery-thumbs" id="gallery-thumbs">
        ${images.map(function(img, i) {
          return '<img src="/uploads/' + img + '" class="gallery-thumb' + (i === 0 ? ' active' : '') + '" data-index="' + i + '" alt="Thumbnail ' + (i+1) + '" />';
        }).join('')}
      </div>
    </div>
  `;

  container.innerHTML = html;

  const featured = document.getElementById('gallery-featured');
  const counter = document.getElementById('gallery-counter');
  const thumbs = container.querySelectorAll('.gallery-thumb');

  function showImage(index) {
    if (index < 0) index = images.length - 1;
    if (index >= images.length) index = 0;
    currentIndex = index;

    featured.style.opacity = '0';
    setTimeout(function() {
      featured.src = '/uploads/' + images[currentIndex];
      featured.style.opacity = '1';
    }, 200);

    counter.textContent = (currentIndex + 1) + ' / ' + images.length;
    thumbs.forEach(function(t, i) {
      t.classList.toggle('active', i === currentIndex);
    });
  }

  document.getElementById('gallery-prev').addEventListener('click', function() {
    showImage(currentIndex - 1);
  });

  document.getElementById('gallery-next').addEventListener('click', function() {
    showImage(currentIndex + 1);
  });

  thumbs.forEach(function(thumb) {
    thumb.addEventListener('click', function() {
      showImage(parseInt(this.dataset.index));
    });
  });
}
