<script>
  import { onMount } from 'svelte';
  import { listMemes, pollChangedMemes } from './api';

  const PAGE_SIZE = 36;
  const POLL_INTERVAL_MS = 8000;
  // A reorder touches every meme's updated_at, so one tick may have to drain
  // several delta pages before the cursor is caught up.
  const MAX_DELTA_PAGES = 8;
  const MAX_UPLOAD_BYTES = 20 * 1024 * 1024;
  const API_ROOT = '/api/custom';

  let memes = [];
  let nextCursor = '';
  let total = 0;
  let loading = true;
  let loadingMore = false;
  let error = '';
  let toastMessage = '';
  let pollError = '';
  let selectedTag = '';
  let tagQuery = '';
  let sentinel;
  let memeGrid;
  let fileInput;
  let observer;
  let pollTimer;
  let filterTimer;
  let pollInFlight = false;
  let pollFailures = 0;
  let pollWatermark = '';
  let pollCursor = '';
  let viewerId = '';
  let viewerTouchStartX = 0;
  let viewerTouchStartY = 0;
  let viewerTouchActive = false;
  let viewerAnimationKey = 0;
  let viewerDirection = '';
  let viewerImageScroll;
  let viewerScale = 1;
  let viewerFit = true;
  let viewerPinching = false;
  let viewerPinchStartDistance = 0;
  let viewerPinchStartScale = 1;
  $: viewerIndex = viewerId ? memes.findIndex((item) => item.id === viewerId) : -1;
  $: viewerMeme = viewerIndex >= 0 ? memes[viewerIndex] : null;

  let uploadFile;
  let uploadPreview = '';
  let tagsInput = '';
  let tagInput = '';
  let addingTags = false;
  let tagDeleteMemeId = '';
  let tagDeleteName = '';
  let tagDeleteTimer;
  let tagDeleteActive = false;
  let tagDeleteOver = false;
  let tagDeletePointerId = null;
  let suppressTagClick = false;
  let dragActive = false;
  let deleteActive = false;
  let draggedMemeId = '';
  let dragStartMemes = null;
  let dragDropHandled = false;
  let dragCancelled = false;
  let dragCommitTargetId = '';
  let dragCommitAfter = false;
  let longPressTimer;
  let touchPointerId = null;
  let touchStartX = 0;
  let touchStartY = 0;
  let touchDragActive = false;
  let touchOverDelete = false;
  let suppressClick = false;
  let reorderTargetId = '';
  let reorderTargetAfter = false;
  let deletingId = '';
  let uploading = false;
  let uploadMessage = '';
  let deleteError = '';

  async function api(path, options = {}) {
    const response = await fetch(path, options);
    let body = null;
    try {
      body = await response.json();
    } catch {
      // A useful error is still shown below when the response is not JSON.
    }
    if (!response.ok) {
      throw new Error(body?.error || `Request failed (${response.status})`);
    }
    return body;
  }

  async function load(reset = false) {
    if (reset) {
      loading = true;
      nextCursor = '';
      // The watermark is about to be re-seeded for a new page/tag window, so a
      // cursor left over from a truncated delta under the old filter must not
      // skip rows that changed since the new watermark.
      pollCursor = '';
    } else if (loadingMore || !nextCursor) {
      return;
    }
    if (!reset) loadingMore = true;
    error = '';
    try {
      const result = await listMemes({
        limit: PAGE_SIZE,
        offset: reset ? 0 : Number(nextCursor),
        tag: selectedTag,
      });
      if (reset) {
        memes = result.memes;
        total = result.total || 0;
        seedPollWatermark(result.memes);
      } else {
        // New memes shift the offset window, so merge by id instead of
        // appending and risking a duplicate keyed element.
        applyMemes(result.memes);
        total = result.total || total;
      }
      nextCursor = result.next_cursor || '';
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
      loadingMore = false;
    }
  }

  // Poll deltas compare against meme.updated_at, so start the watermark just
  // before the newest row on screen to catch writes that raced the page load.
  function seedPollWatermark(items) {
    const newest = (items || []).reduce(
      (latest, item) =>
        Date.parse(item.updated_at || '') > latest ? Date.parse(item.updated_at) : latest,
      0,
    );
    pollWatermark = new Date(Math.max(0, (newest || Date.now()) - 60_000)).toISOString();
  }

  /**
   * Apply a poll delta (or a meme this tab just wrote) to the timeline. The
   * server is authoritative, but a delta can be a few milliseconds behind an
   * in-flight write, so the freshest `updated_at` for an id wins and the
   * visible order is the server's own `sort_order`, not the order rows arrived
   * in. Rows the delta reports as `deleted` (or that vanished from it) are
   * dropped, which is how a delete in another tab reaches this one.
   */
  function applyMemes(incoming) {
    const byID = new Map(memes.map((item) => [item.id, item]));
    let dirty = false;
    for (const item of incoming || []) {
      const existing = byID.get(item.id);
      if (item.deleted) {
        if (existing) {
          byID.delete(item.id);
          dirty = true;
          if (viewerId === item.id) closeViewer();
        }
        continue;
      }
      if (
        existing &&
        Date.parse(item.updated_at || '') <= Date.parse(existing.updated_at || '')
      ) {
        // A stale delta for a row this tab already knows better about.
        continue;
      }
      byID.set(item.id, item);
      dirty = true;
    }
    if (!dirty) return;
    memes = [...byID.values()].sort(compareTimeline);
  }

  // Mirrors the server's timeline ordering: sort_order DESC, created_at DESC,
  // id DESC.
  function compareTimeline(a, b) {
    return (
      (b.sort_order || 0) - (a.sort_order || 0) ||
      Date.parse(b.created_at || '') - Date.parse(a.created_at || '') ||
      (b.id || '').localeCompare(a.id || '')
    );
  }

  function delay(ms) {
    return new Promise((resolve) => window.setTimeout(resolve, ms));
  }

  /**
   * Re-check a row this tab wrote, and keep following it while a vision job is
   * still landing on it. Tag edits are final the moment the API answers, so
   * they only need the one re-read; uploads and description retries are queued.
   */
  async function syncAfterWrite(meme) {
    if (meme.description_status === 'pending') {
      void syncMeme(meme.id);
      return;
    }
    void syncNow(meme.id);
  }

  /** Re-read one meme from the database and apply it if it is newer. */
  async function syncNow(id) {
    try {
      const fresh = await fetchMeme(id);
      applyMemes([fresh]);
      pollError = '';
      return memes.find((item) => item.id === id) || null;
    } catch (err) {
      pollError = err.message;
      return null;
    }
  }

  /**
   * Follow a write this tab made until the row stops changing. Uploads and
   * description retries are answered immediately by the API while the vision
   * worker is still queueing, so a single refetch is not enough: the modal
   * would settle on the pre-write snapshot and never move again.
   */
  async function syncMeme(id) {
    const settled = await syncNow(id);
    let stamp = settled?.updated_at || '';
    // Keep following the row while it is still moving (a vision job rewrites
    // description/tags a few seconds later), then stop once it holds still.
    for (let attempt = 0; attempt < 8; attempt += 1) {
      await delay(2000);
      const current = await syncNow(id);
      if (!current) break;
      if (current.updated_at === stamp && current.description_status !== 'pending') break;
      stamp = current.updated_at;
    }
  }

  /**
   * Timer tick. Pulls only the rows whose metadata changed since the last tick
   * (plus a short overlap), so the cost is a handful of rows instead of the
   * whole first page, and deletions in another tab are picked up here. A
   * reorder transaction stamps every meme, so one tick drains several delta
   * pages before the cursor catches up with the database.
   */
  async function refreshLatest() {
    if (pollInFlight || loading || document.visibilityState === 'hidden') return;
    pollInFlight = true;
    try {
      let pages = 0;
      let cursor = pollCursor;
      do {
        const tick = await pollChangedMemes({
          after: pollWatermark || new Date(Date.now() - 60_000).toISOString(),
          cursor,
          tag: selectedTag,
        });
        if (tick.memes.length) applyMemes(tick.memes);
        if (!tick.cursor) pollWatermark = tick.watermark;
        cursor = tick.cursor;
        pages += 1;
      } while (cursor && pages < MAX_DELTA_PAGES);
      pollCursor = cursor;
      pollError = '';
      pollFailures = 0;
    } catch (err) {
      // Loud, but non-blocking: the timeline stays on screen and the next tick
      // retries from the last good watermark. Back off while failures stack up
      // so a backend that is down is not hammered every few seconds.
      pollFailures += 1;
      pollError = err.message;
      startPolling(Math.min(POLL_INTERVAL_MS * 2 ** pollFailures, 60_000));
    } finally {
      pollInFlight = false;
    }
  }

  function startPolling(intervalMs = POLL_INTERVAL_MS) {
    window.clearInterval(pollTimer);
    pollTimer = window.setInterval(refreshLatest, intervalMs);
  }

  function onVisibilityChange() {
    if (document.visibilityState !== 'visible') return;
    // Returning to a tab that was hidden through several ticks catches up
    // immediately, then restores the normal cadence after a failed-run backoff.
    startPolling();
    void refreshLatest();
  }

  function onPaste(event) {
    const imageItem = [...(event.clipboardData?.items || [])].find(
      (item) => item.kind === 'file' && item.type.startsWith('image/')
    );
    const file = imageItem?.getAsFile();
    if (!file) return;
    event.preventDefault();
    setUploadFile(file);
  }

  function chooseTag(tag) {
    if (suppressTagClick) {
      suppressTagClick = false;
      return;
    }
    selectedTag = selectedTag === tag ? '' : tag;
    tagQuery = selectedTag ? `#${selectedTag}` : '';
    load(true);
  }

  function applyTag(event) {
    event?.preventDefault();
    const value = tagQuery.trim().replace(/^#/, '').split(/[\s,]+/)[0].toLowerCase();
    selectedTag = value;
    tagQuery = value ? `#${value}` : '';
    load(true);
  }

  function onTagInput() {
    if (filterTimer) window.clearTimeout(filterTimer);
    filterTimer = window.setTimeout(() => applyTag(), 300);
  }

  function clearTag() {
    if (filterTimer) window.clearTimeout(filterTimer);
    selectedTag = '';
    tagQuery = '';
    load(true);
  }

  function setUploadFile(file) {
    if (!file) return;
    if (!file.type.startsWith('image/')) {
      uploadMessage = 'Please choose an image file.';
      return;
    }
    if (file.size > MAX_UPLOAD_BYTES) {
      uploadMessage = 'That image is over the 20 MB limit.';
      return;
    }
    if (uploadPreview) URL.revokeObjectURL(uploadPreview);
    uploadFile = file;
    uploadPreview = URL.createObjectURL(file);
    uploadMessage = '';
  }

  function onFileInput(event) {
    setUploadFile(event.target.files?.[0]);
    event.target.value = '';
  }

  function onDrop(event) {
    event.preventDefault();
    dragActive = false;
    setUploadFile(event.dataTransfer.files?.[0]);
  }

  function preventWindowDrop(event) {
    event.preventDefault();
  }

  function onMemeContextMenu(event) {
    if (event.pointerType === 'touch' || event.sourceCapabilities?.firesTouchEvents) event.preventDefault();
  }

  function routeState() {
    const path = window.location.pathname.split('/').filter(Boolean);
    if (path[0] === 'meme' && path[1]) return { mode: 'modal', id: path[1] };
    if (path[0] === 'story' && path[1]) return { mode: 'story', id: path[1] };
    const params = new URLSearchParams(window.location.search);
    if (params.get('meme')) return { mode: 'modal', id: params.get('meme') };
    if (params.get('at')) return { mode: 'story', id: params.get('at') };
    return { mode: '', id: '' };
  }

  function setRoute(path, replace = false) {
    window.history[replace ? 'replaceState' : 'pushState']({}, '', path);
  }

  function showViewer(id, replace = false, direction = '') {
    viewerId = id;
    viewerDirection = direction;
    viewerScale = 1;
    viewerFit = true;
    viewerPinching = false;
    viewerTouchActive = false;
    viewerAnimationKey += 1;
    document.body.style.overflow = 'hidden';
    requestAnimationFrame(() => viewerImageScroll?.scrollTo(0, 0));
    setRoute(`/meme/${id}`, replace);
  }

  function openViewer(event, meme) {
    if (suppressClick) {
      suppressClick = false;
      event.preventDefault();
      return;
    }
    if (event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
    event.preventDefault();
    showViewer(meme.id);
  }

  function closeViewer(updateRoute = true) {
    viewerId = '';
    viewerDirection = '';
    document.body.style.overflow = '';
    if (updateRoute && routeState().mode === 'modal') setRoute('/');
  }

  function viewerPrevious() {
    if (viewerIndex > 0) showViewer(memes[viewerIndex - 1].id, true, 'prev');
  }

  function touchDistance(touches) {
    const [first, second] = touches;
    return Math.hypot(second.clientX - first.clientX, second.clientY - first.clientY);
  }

  function touchCenter(touches) {
    const [first, second] = touches;
    return {
      x: (first.clientX + second.clientX) / 2,
      y: (first.clientY + second.clientY) / 2,
    };
  }

  function setViewerScale(value, clientX, clientY, minimum = 0.25) {
    const next = Math.min(4, Math.max(minimum, value));
    const element = viewerImageScroll;
    if (!element || next === viewerScale) {
      viewerScale = next;
      return;
    }
    const rect = element.getBoundingClientRect();
    const localX = clientX === undefined ? rect.width / 2 : clientX - rect.left;
    const localY = clientY === undefined ? rect.height / 2 : clientY - rect.top;
    const contentX = (element.scrollLeft + localX) / viewerScale;
    const contentY = (element.scrollTop + localY) / viewerScale;
    viewerScale = next;
    requestAnimationFrame(() => {
      element.scrollLeft = Math.max(0, contentX * next - localX);
      element.scrollTop = Math.max(0, contentY * next - localY);
    });
  }

  function toggleViewerZoom() {
    viewerFit = !viewerFit;
    viewerScale = 1;
  }

  function onViewerWheel(event) {
    // Preserve normal scrolling and browser zoom. Shift-wheel is the explicit
    // desktop gesture for image zoom.
    if (!event.shiftKey || event.ctrlKey) return;
    event.preventDefault();
    event.stopPropagation();
    viewerFit = false;
    setViewerScale(viewerScale * Math.pow(1.002, -event.deltaY), event.clientX, event.clientY, 0.25);
  }

  function onViewerTouchStart(event) {
    if (event.touches.length === 2) {
      viewerPinching = true;
      viewerTouchActive = false;
      viewerPinchStartDistance = touchDistance(event.touches);
      viewerPinchStartScale = viewerScale;
      event.preventDefault();
      return;
    }
    if (event.touches.length !== 1) return;
    viewerTouchActive = true;
    viewerTouchStartX = event.touches[0].clientX;
    viewerTouchStartY = event.touches[0].clientY;
  }

  function onViewerTouchMove(event) {
    if (!viewerPinching || event.touches.length < 2) return;
    event.preventDefault();
    const center = touchCenter(event.touches);
    viewerFit = false;
    setViewerScale(
      viewerPinchStartScale * (touchDistance(event.touches) / viewerPinchStartDistance),
      center.x,
      center.y,
      1,
    );
  }

  function onViewerTouchEnd(event) {
    if (viewerPinching) {
      if (event.touches.length < 2) {
        viewerPinching = false;
        viewerTouchActive = false;
      }
      return;
    }
    if (!viewerTouchActive || !viewerMeme || viewerScale > 1.01 || event.changedTouches.length !== 1) {
      viewerTouchActive = false;
      return;
    }
    viewerTouchActive = false;
    const touch = event.changedTouches[0];
    const dx = touch.clientX - viewerTouchStartX;
    const dy = touch.clientY - viewerTouchStartY;
    if (Math.abs(dx) < 45 || Math.abs(dx) < Math.abs(dy) * 1.2) return;
    if (dx < 0) void viewerNext();
    else viewerPrevious();
  }

  async function viewerNext() {
    if (viewerIndex < memes.length - 1) {
      showViewer(memes[viewerIndex + 1].id, true, 'next');
      return;
    }
    if (!nextCursor || loadingMore) return;
    await load(false);
    if (viewerIndex < memes.length - 1) showViewer(memes[viewerIndex + 1].id, true);
  }

  async function applyRoute() {
    const route = routeState();
    if (!route.id) {
      if (viewerId) closeViewer(false);
      return;
    }
    let attempts = 0;
    while (!memes.some((item) => item.id === route.id) && nextCursor && attempts < 100) {
      await load(false);
      attempts += 1;
    }
    if (route.mode === 'modal' && memes.some((item) => item.id === route.id)) {
      showViewer(route.id, true);
    } else if (route.mode === 'story') {
      requestAnimationFrame(() => {
        const card = [...document.querySelectorAll('.meme-card')].find((item) => item.dataset.memeId === route.id);
        card?.scrollIntoView({ behavior: 'smooth', block: 'center' });
      });
    }
  }

  function onKeyDown(event) {
    cancelMemeDrag(event);
    if (!viewerMeme) return;
    if (event.key === 'Escape') {
      event.preventDefault();
      closeViewer();
    } else if (event.key === 'ArrowLeft') {
      event.preventDefault();
      viewerPrevious();
    } else if (event.key === 'ArrowRight') {
      event.preventDefault();
      void viewerNext();
    }
  }

  function startMemeDrag(event, meme) {
    draggedMemeId = meme.id;
    dragStartMemes = [...memes];
    dragDropHandled = false;
    dragCancelled = false;
    dragCommitTargetId = '';
    dragCommitAfter = false;
    event.dataTransfer.effectAllowed = 'move';
    event.dataTransfer.setData('application/x-memebrary-id', meme.id);
    event.dataTransfer.setData('text/plain', meme.id);
  }

  function endMemeDrag() {
    const sourceID = draggedMemeId;
    const targetID = dragCommitTargetId || reorderTargetId;
    const after = dragCommitTargetId ? dragCommitAfter : reorderTargetAfter;
    const original = dragStartMemes;
    // dragend is the reliable release signal. Do not gate this on dropEffect:
    // Chromium can report "none" when the DOM reflows beneath a native drag.
    if (original && sourceID && targetID && sourceID !== targetID && !dragDropHandled && !dragCancelled) {
      dragDropHandled = true;
      void commitOrder(sourceID, targetID, after, original);
    } else if (original && !dragDropHandled) {
      memes = original;
    }
    draggedMemeId = '';
    dragStartMemes = null;
    dragDropHandled = false;
    dragCancelled = false;
    dragCommitTargetId = '';
    dragCommitAfter = false;
    reorderTargetId = '';
    reorderTargetAfter = false;
    deleteActive = false;
  }

  function cancelMemeDrag(event) {
    if (event.key !== 'Escape' || !draggedMemeId) return;
    dragCancelled = true;
    if (dragStartMemes) memes = dragStartMemes;
    dragCommitTargetId = '';
    dragCommitAfter = false;
    reorderTargetId = '';
    reorderTargetAfter = false;
  }

  function isMemeDrag(event) {
    return Array.from(event.dataTransfer.types).includes('application/x-memebrary-id');
  }

  function dropIsAfter(event, card) {
    const rect = card.getBoundingClientRect();
    const dx = event.clientX - (rect.left + rect.width / 2);
    const dy = event.clientY - (rect.top + rect.height / 2);
    return Math.abs(dx) > Math.abs(dy) ? dx > 0 : dy > 0;
  }

  function applyReorderPreview(meme, after) {
    if (!meme || draggedMemeId === meme.id) return false;
    const sourceIndex = memes.findIndex((item) => item.id === draggedMemeId);
    const targetIndex = memes.findIndex((item) => item.id === meme.id);
    if (sourceIndex >= 0 && targetIndex >= 0) {
      const next = [...memes];
      const [source] = next.splice(sourceIndex, 1);
      const adjustedTargetIndex = next.findIndex((item) => item.id === meme.id);
      next.splice(Math.max(0, adjustedTargetIndex + (after ? 1 : 0)), 0, source);
      memes = next;
    }
    dragCommitTargetId = meme.id;
    dragCommitAfter = after;
    reorderTargetId = meme.id;
    reorderTargetAfter = after;
    return true;
  }

  function previewReorder(event, card, meme, after = dropIsAfter(event, card)) {
    if (!card || !meme || draggedMemeId === meme.id || !isMemeDrag(event)) return false;
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
    return applyReorderPreview(meme, after);
  }

  function onCardDragOver(event, meme) {
    previewReorder(event, event.currentTarget, meme);
  }

  function clearLongPress() {
    if (longPressTimer) window.clearTimeout(longPressTimer);
    longPressTimer = undefined;
  }

  function onMemePointerDown(event, meme) {
    if (event.pointerType === 'mouse' || event.button !== 0) return;
    touchPointerId = event.pointerId;
    touchStartX = event.clientX;
    touchStartY = event.clientY;
    touchDragActive = false;
    touchOverDelete = false;
    clearLongPress();
    longPressTimer = window.setTimeout(() => {
      touchDragActive = true;
      draggedMemeId = meme.id;
      dragStartMemes = [...memes];
      dragDropHandled = false;
      dragCancelled = false;
      dragCommitTargetId = '';
      dragCommitAfter = false;
      suppressClick = true;
      event.currentTarget.setPointerCapture?.(event.pointerId);
      updateMemePointer(event);
    }, 450);
  }

  function updateMemePointer(event) {
    if (!touchDragActive || event.pointerId !== touchPointerId) return;
    event.preventDefault();
    const element = document.elementFromPoint(event.clientX, event.clientY);
    const deleteZone = element?.closest?.('.delete-zone');
    if (deleteZone) {
      if (dragStartMemes) memes = dragStartMemes;
      dragCommitTargetId = '';
      dragCommitAfter = false;
      reorderTargetId = '';
      reorderTargetAfter = false;
      touchOverDelete = true;
      deleteActive = true;
      return;
    }
    touchOverDelete = false;
    deleteActive = false;
    const target = findDropTarget({ clientX: event.clientX, clientY: event.clientY });
    const meme = memes.find((item) => item.id === target?.card?.dataset.memeId);
    if (target && meme) applyReorderPreview(meme, target.after);
  }

  function onMemePointerMove(event) {
    if (event.pointerId !== touchPointerId) return;
    if (!touchDragActive) {
      if (Math.hypot(event.clientX - touchStartX, event.clientY - touchStartY) > 10) clearLongPress();
      return;
    }
    updateMemePointer(event);
  }

  async function finishMemePointer(event, cancelled = false) {
    if (event.pointerId !== touchPointerId) return;
    clearLongPress();
    if (!touchDragActive) {
      touchPointerId = null;
      return;
    }
    event.preventDefault();
    if (cancelled) dragCancelled = true;
    if (touchOverDelete) {
      const id = draggedMemeId;
      const original = dragStartMemes || memes;
      dragDropHandled = true;
      memes = original;
      await deleteMemeByID(id);
    } else if (dragCommitTargetId && dragCommitTargetId !== draggedMemeId && !dragCancelled) {
      const original = dragStartMemes || memes;
      dragDropHandled = true;
      await commitOrder(draggedMemeId, dragCommitTargetId, dragCommitAfter, original);
    } else if (dragStartMemes) {
      memes = dragStartMemes;
    }
    draggedMemeId = '';
    dragStartMemes = null;
    dragDropHandled = false;
    touchPointerId = null;
    touchDragActive = false;
    touchOverDelete = false;
    deleteActive = false;
    reorderTargetId = '';
    reorderTargetAfter = false;
    dragCommitTargetId = '';
    dragCommitAfter = false;
    window.setTimeout(() => (suppressClick = false), 0);
  }

  function onMemePointerUp(event) {
    void finishMemePointer(event);
  }

  function onMemePointerCancel(event) {
    void finishMemePointer(event, true);
  }

  function onCardDragLeave(event) {
    if (!event.relatedTarget || !event.currentTarget.contains(event.relatedTarget)) {
      reorderTargetId = '';
      reorderTargetAfter = false;
    }
  }

  function findDropTarget(event) {
    if (!memeGrid) return null;
    const cards = [...memeGrid.querySelectorAll('.meme-card')].map((card) => ({
      card,
      rect: card.getBoundingClientRect()
    }));
    if (cards.length === 0) return null;

    const rows = [];
    for (const item of cards) {
      let row = rows.find((candidate) => Math.abs(candidate.top - item.rect.top) < 4);
      if (!row) {
        row = { top: item.rect.top, bottom: item.rect.bottom, cards: [] };
        rows.push(row);
      }
      row.bottom = Math.max(row.bottom, item.rect.bottom);
      row.cards.push(item);
    }
    rows.sort((a, b) => a.top - b.top);
    for (const row of rows) row.cards.sort((a, b) => a.rect.left - b.rect.left);

    const row = rows.reduce((closest, candidate) => {
      const distance = event.clientY < candidate.top
        ? candidate.top - event.clientY
        : event.clientY > candidate.bottom
          ? event.clientY - candidate.bottom
          : 0;
      return !closest || distance < closest.distance ? { row: candidate, distance } : closest;
    }, null)?.row;
    if (!row) return null;

    // The horizontal space beyond a row is an explicit first/last insertion
    // zone, rather than requiring a precise drop on the edge card.
    if (event.clientX < row.cards[0].rect.left) {
      return { card: row.cards[0].card, after: false };
    }
    const last = row.cards[row.cards.length - 1];
    if (event.clientX > last.rect.right) {
      return { card: last.card, after: true };
    }

    let nearest = row.cards[0];
    let distance = Infinity;
    for (const item of row.cards) {
      const dx = event.clientX - (item.rect.left + item.rect.width / 2);
      const dy = event.clientY - (item.rect.top + item.rect.height / 2);
      const nextDistance = dx * dx + dy * dy;
      if (nextDistance < distance) {
        distance = nextDistance;
        nearest = item;
      }
    }
    return { card: nearest.card, after: dropIsAfter(event, nearest.card) };
  }

  function onTimelineDragOver(event) {
    if (!isMemeDrag(event)) return;
    const target = findDropTarget(event);
    const meme = memes.find((item) => item.id === target?.card?.dataset.memeId);
    if (!target || !meme) return;
    previewReorder(event, target.card, meme, target.after);
  }

  function onTimelineDragLeave(event) {
    if (!event.relatedTarget || !event.currentTarget.contains(event.relatedTarget)) {
      reorderTargetId = '';
      reorderTargetAfter = false;
    }
  }

  async function commitOrder(sourceID, targetID, after, oldMemes) {
    try {
      const response = await fetch(`${API_ROOT}/memes/${sourceID}/order`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(after ? { after_id: targetID } : { before_id: targetID })
      });
      if (!response.ok) {
        let body = null;
        try { body = await response.json(); } catch {}
        throw new Error(body?.error || `Reorder failed (${response.status})`);
      }
      // Reconcile the optimistic preview with the server's canonical order.
      await refreshLatest();
    } catch (err) {
      memes = oldMemes;
      error = err.message;
    }
  }

  async function reorderMeme(event, target, after = reorderTargetAfter) {
    event.preventDefault();
    const sourceID = event.dataTransfer.getData('application/x-memebrary-id') || event.dataTransfer.getData('text/plain');
    reorderTargetId = '';
    reorderTargetAfter = false;
    if (!sourceID || sourceID === target.id) return;
    const oldMemes = dragStartMemes || memes;
    dragDropHandled = true;
    await commitOrder(sourceID, target.id, after, oldMemes);
  }

  async function onCardDrop(event, target) {
    await reorderMeme(event, target, dropIsAfter(event, event.currentTarget));
  }

  async function onTimelineDrop(event) {
    if (!reorderTargetId) return;
    const target = memes.find((item) => item.id === reorderTargetId);
    if (target) await reorderMeme(event, target, reorderTargetAfter);
  }

  function onDeleteDragOver(event) {
    if (!Array.from(event.dataTransfer.types).includes('application/x-memebrary-id')) return;
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
    if (dragStartMemes) memes = dragStartMemes;
    dragCommitTargetId = '';
    dragCommitAfter = false;
    reorderTargetId = '';
    reorderTargetAfter = false;
    deleteActive = true;
  }

  function onDeleteDragLeave(event) {
    if (!event.currentTarget.contains(event.relatedTarget)) deleteActive = false;
  }

  async function deleteMemeByID(id) {
    const meme = memes.find((item) => item.id === id);
    if (!meme || deletingId || !window.confirm('Delete this meme permanently?')) return;
    const oldMemes = dragStartMemes || memes;
    deletingId = meme.id;
    deleteError = '';
    memes = oldMemes;
    try {
      const response = await fetch(`${API_ROOT}/memes/${meme.id}`, { method: 'DELETE' });
      if (!response.ok) {
        let body = null;
        try { body = await response.json(); } catch {}
        throw new Error(body?.error || `Delete failed (${response.status})`);
      }
      memes = memes.filter((item) => item.id !== meme.id);
      total = Math.max(0, total - 1);
      if (viewerId === meme.id) closeViewer();
    } catch (err) {
      memes = oldMemes;
      deleteError = err.message;
    } finally {
      deletingId = '';
      draggedMemeId = '';
    }
  }

  async function onDeleteDrop(event) {
    event.preventDefault();
    deleteActive = false;
    const id = event.dataTransfer.getData('application/x-memebrary-id') || event.dataTransfer.getData('text/plain');
    if (!memes.some((item) => item.id === id)) return;
    dragDropHandled = true;
    if (dragStartMemes) memes = dragStartMemes;
    await deleteMemeByID(id);
  }

  function resetUpload() {
    if (uploadPreview) URL.revokeObjectURL(uploadPreview);
    uploadFile = null;
    uploadPreview = '';
    tagsInput = '';
    uploadMessage = '';
  }

  async function submitUpload() {
    if (!uploadFile || uploading) return;
    uploading = true;
    uploadMessage = '';
    const form = new FormData();
    form.set('file', uploadFile);
    form.set('tags', tagsInput);
    try {
      const meme = await api(`${API_ROOT}/memes`, { method: 'POST', body: form });
      if (!selectedTag || meme.tags?.includes(selectedTag)) {
        applyMemes([meme]);
      }
      resetUpload();
      if (!uploadMessage) {
        toastMessage = uploadFile.name;
        window.setTimeout(() => { toastMessage = ''; }, 5000);
      }
      // The vision worker is still queueing when the upload is answered, so
      // keep following this row until its description/tags land.
      void syncMeme(meme.id);
    } catch (err) {
      uploadMessage = err.message;
    } finally {
      uploading = false;
    }
  }

  function clearTagDelete() {
    if (tagDeleteTimer) window.clearTimeout(tagDeleteTimer);
    tagDeleteTimer = undefined;
    tagDeleteMemeId = '';
    tagDeleteName = '';
    tagDeleteActive = false;
    tagDeleteOver = false;
    tagDeletePointerId = null;
  }

  function onTagPointerDown(event, meme, tag) {
    if (event.pointerType === 'mouse' && event.button !== 0) return;
    clearTagDelete();
    tagDeleteMemeId = meme.id;
    tagDeleteName = tag;
    tagDeletePointerId = event.pointerId;
    tagDeleteTimer = window.setTimeout(() => {
      tagDeleteActive = true;
      suppressTagClick = true;
      event.currentTarget.setPointerCapture?.(event.pointerId);
    }, 550);
  }

  function onTagPointerMove(event) {
    if (!tagDeleteMemeId || event.pointerId !== tagDeletePointerId) return;
    if (!tagDeleteActive) return;
    const button = document.elementFromPoint(event.clientX, event.clientY);
    tagDeleteOver = button === event.currentTarget;
  }

  async function onTagPointerUp(event) {
    if (!tagDeleteMemeId) return;
    if (tagDeleteTimer) window.clearTimeout(tagDeleteTimer);
    const shouldDelete = tagDeleteActive && tagDeleteOver;
    const memeId = tagDeleteMemeId;
    const tag = tagDeleteName;
    if (tagDeleteActive) suppressTagClick = true;
    clearTagDelete();
    if (shouldDelete) await removeTag(memeId, tag);
  }

  function onTagPointerCancel() {
    clearTagDelete();
  }

  async function removeTag(memeId, tag) {
    try {
      const updated = await api(`${API_ROOT}/memes/${memeId}/tags/${encodeURIComponent(tag)}`, { method: 'DELETE' });
      applyMemes([updated]);
      // A vision job may still be landing on this row; keep the modal in step.
      void syncAfterWrite(updated);
    } catch (err) {
      error = err.message;
    }
  }

  async function addTags() {
    if (!viewerMeme || !tagInput.trim() || addingTags) return;
    addingTags = true;
    try {
      const updated = await api(`${API_ROOT}/memes/${viewerMeme.id}/tags`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ tags: tagInput.split(/[\\s,]+/).filter(Boolean) })
      });
      applyMemes([updated]);
      tagInput = '';
      // A vision job may still be landing on this row; keep the modal in step.
      void syncAfterWrite(updated);
    } catch (err) {
      error = err.message;
    } finally {
      addingTags = false;
    }
  }

  async function retryDescription(meme) {
    if (!meme) return;
    try {
      const pending = await api(`${API_ROOT}/memes/${meme.id}/describe`, { method: 'POST' });
      applyMemes([pending]);
      // The worker answers immediately and rewrites the row seconds later.
      void syncMeme(pending.id);
    } catch (err) {
      error = err.message;
    }
  }

  onMount(() => {
    load(true).then(applyRoute);
    window.addEventListener('popstate', applyRoute);
    window.addEventListener('dragover', preventWindowDrop);
    window.addEventListener('drop', preventWindowDrop);
    window.addEventListener('paste', onPaste);
    window.addEventListener('keydown', onKeyDown)
    window.addEventListener('visibilitychange', onVisibilityChange);
    startPolling();
    observer = new IntersectionObserver(([entry]) => {
      if (entry.isIntersecting && nextCursor && !loadingMore) load(false);
    }, { rootMargin: '500px' });
    if (sentinel) observer.observe(sentinel);
    return () => {
      observer?.disconnect();
      window.removeEventListener('popstate', applyRoute);
      window.removeEventListener('dragover', preventWindowDrop);
      window.removeEventListener('drop', preventWindowDrop);
      window.removeEventListener('paste', onPaste);
      window.removeEventListener('keydown', onKeyDown)
      window.removeEventListener('visibilitychange', onVisibilityChange);
      window.clearInterval(pollTimer);
      if (filterTimer) window.clearTimeout(filterTimer);
      if (uploadPreview) URL.revokeObjectURL(uploadPreview);
      document.body.style.overflow = '';
    };
  });
</script>

<svelte:head>
  <title>meme/brary</title>
</svelte:head>

<header class="topbar">
  <div class="topbar-inner">
    <a class="wordmark" href="/" aria-label="meme/brary home">
      <span class="mark">m</span><span>meme<span class="slash">/</span>brary</span>
    </a>
    <span
      class="anonymous"
      class:stale={pollError}
      title={pollError || 'The timeline refreshes every 8 seconds'}
    ><span class="dot"></span> anonymous · {pollError ? 'stale' : 'live'}</span>
    <form class="filter" on:submit={applyTag}>
      <input id="tag-filter" bind:value={tagQuery} on:input={onTagInput} placeholder="#search tags" aria-label="Search hashtags" autocomplete="off" />
      {#if selectedTag}
        <button class="clear-filter" type="button" on:click={clearTag} aria-label="Clear tag filter">×</button>
      {/if}
    </form>
  </div>
</header>

<main class="page">
  <section class="control-dock" aria-label="Meme controls">
    <div class="upload-row">
      <section
        class:drag-active={dragActive}
        class="dropzone"
        aria-label="Image upload drop zone"
        on:dragenter|preventDefault={() => (dragActive = true)}
        on:dragover|preventDefault={() => (dragActive = true)}
        on:dragleave|preventDefault={() => (dragActive = false)}
        on:drop={onDrop}
      >
        <div class="drop-copy">
          <span class="upload-icon">＋</span>
          <div>
            <strong>Drop a meme here</strong>
            <span>or <button type="button" class="link-button" on:click={() => fileInput?.click()}>choose an image</button></span>
          </div>
        </div>
        <span class="drop-hint">PNG, JPG, GIF or WebP · up to 20 MB · Ctrl+V works too</span>
        <input bind:this={fileInput} class="visually-hidden" type="file" accept="image/jpeg,image/png,image/gif,image/webp" on:change={onFileInput} />
      </section>

      <section
        class:delete-active={deleteActive}
        class="delete-zone"
        aria-label="Delete meme drop zone"
        on:dragenter={onDeleteDragOver}
        on:dragover={onDeleteDragOver}
        on:dragleave={onDeleteDragLeave}
        on:drop={onDeleteDrop}
      >
        <span class="trash-icon" aria-hidden="true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
            <path d="M4 7h16" />
            <path d="M9 7V4h6v3" />
            <path d="M7 7l1 13h8l1-13" />
            <path d="M10 11v5M14 11v5" />
          </svg>
        </span>
        <div>
          <strong>{deletingId ? 'Deleting…' : 'Drag here to delete'}</strong>
          <span>{deleteError || 'release to remove forever'}</span>
        </div>
      </section>
    </div>

    {#if uploadFile}
      {#if toastMessage}
        <div class="notice success" role="status" aria-label="Upload confirmation">
          <strong>uploaded</strong> {toastMessage}
          <button class="link-button" on:click={() => toastMessage = ''}>×</button>
        </div>
      {/if}
      <section class="upload-editor" aria-label="New meme details">
        <div class="editor-preview"><img src={uploadPreview} alt="Selected meme preview" /></div>
        <div class="editor-fields">
          <label>
            <span>hashtags <small>(optional)</small></span>
            <input bind:value={tagsInput} placeholder="#reaction #work #animals" maxlength="500" />
          </label>
          {#if uploadMessage}<p class="form-error">{uploadMessage}</p>{/if}
          <div class="editor-actions">
            <button class="button primary" disabled={uploading} on:click={submitUpload}>{uploading ? 'Uploading…' : 'Add to library'}</button>
            <button class="button" disabled={uploading} on:click={resetUpload}>Cancel</button>
          </div>
        </div>
      </section>
    {/if}
  </section>

  {#if error}
    <div class="notice error" role="alert"><strong>Couldn’t load the library.</strong> {error} <button on:click={() => load(true)}>Try again</button></div>
  {/if}

  {#if pollError}
    <div class="notice error" role="alert">
      <strong>Recent changes aren’t landing.</strong> The timeline may be out of date. {pollError}
      <button on:click={() => refreshLatest()}>Retry now</button>
    </div>
  {/if}

  {#if loading}
    <div class="loading-grid" aria-label="Loading memes">{#each Array(8) as _}<div class="skeleton"></div>{/each}</div>
  {:else if memes.length === 0}
    <div class="empty-state">
      <div class="empty-mark">¯\_(ツ)_/¯</div>
      <h2>{selectedTag ? 'Nothing with that tag yet.' : 'The library is empty.'}</h2>
      <p>{selectedTag ? 'Try another hashtag or clear the filter.' : 'Be the first to drop a meme into the archive.'}</p>
    </div>
  {:else}
    <section
      bind:this={memeGrid}
      class="meme-grid"
      aria-label="Meme timeline"
      on:dragover={onTimelineDragOver}
      on:dragleave={onTimelineDragLeave}
      on:drop={onTimelineDrop}
    >
      {#each memes as meme (meme.id)}
        <article
          class:dragging={draggedMemeId === meme.id}
          class:reorder-target={reorderTargetId === meme.id}
          class:reorder-target-after={reorderTargetId === meme.id && reorderTargetAfter}
          class="meme-card"
          data-meme-id={meme.id}
          title="Drag to rearrange · drag to the red area to delete"
          draggable="true"
          on:dragstart={(event) => startMemeDrag(event, meme)}
          on:dragend={endMemeDrag}
          on:pointerdown={(event) => onMemePointerDown(event, meme)}
          on:pointermove={onMemePointerMove}
          on:pointerup={onMemePointerUp}
          on:pointercancel={onMemePointerCancel}
          on:contextmenu={onMemeContextMenu}
          on:dragover|stopPropagation={(event) => onCardDragOver(event, meme)}
          on:dragleave={onCardDragLeave}
          on:drop|stopPropagation={(event) => onCardDrop(event, meme)}
        >
          <a class="image-frame" href={`${API_ROOT}/media/${meme.id}`} target="_blank" rel="noreferrer" on:click={(event) => openViewer(event, meme)}>
            <img loading="lazy" src={`${API_ROOT}/media/${meme.id}`} alt={meme.description || 'Meme image'} />
          </a>
          {#if meme.tags?.length}
            <div class="card-details">
              <div class="tags">
                {#each meme.tags as tag}<button on:click={() => chooseTag(tag)}>#{tag}</button>{/each}
              </div>
            </div>
          {/if}
        </article>
      {/each}
    </section>
  {/if}

  <div bind:this={sentinel} class="load-sentinel" aria-hidden="true"></div>
  {#if loadingMore}<p class="loading-more">loading more…</p>{/if}
</main>

{#if viewerMeme}
  <div class="viewer-backdrop" role="presentation" on:click={closeViewer}>
    <dialog open class="viewer" aria-label="Meme viewer" on:click|stopPropagation>
      <button class="viewer-close" aria-label="Close image viewer" on:click={closeViewer}>×</button>
      <button class="viewer-fit-toggle" aria-label={viewerFit ? 'Use current image zoom' : 'Fit image to modal'} title={viewerFit ? 'Use current image zoom' : 'Fit image to modal'} on:click={toggleViewerZoom}>
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path d="M8 3H3v5M16 3h5v5M8 21H3v-5M16 21h5v-5" />
          <path d="M3 3l6 6M21 3l-6 6M3 21l6-6M21 21l-6-6" />
        </svg>
      </button>
      {#key viewerAnimationKey}
      <div class="viewer-image-stage">
        <div
          bind:this={viewerImageScroll}
          class="viewer-image-scroll"
          class:viewer-fit={viewerFit}
          role="application"
          aria-label="Zoomable image"
          class:viewer-slide-next={viewerDirection === 'next'}
          class:viewer-slide-prev={viewerDirection === 'prev'}
          on:wheel={onViewerWheel}
          on:touchstart={onViewerTouchStart}
          on:touchmove={onViewerTouchMove}
          on:touchend={onViewerTouchEnd}
        >
          <img
            class="viewer-image"
            style={viewerFit ? '' : `width: ${viewerScale * 100}%; max-width: ${viewerScale > 1 ? 'none' : '100%'};`}
            src={`${API_ROOT}/media/${viewerMeme.id}`}
            alt={viewerMeme.description || 'Meme image'}
            draggable="false"
          />
        </div>
        <button class="viewer-arrow viewer-prev" aria-label="Previous meme" disabled={viewerIndex <= 0} on:click={viewerPrevious}>‹</button>
        <button class="viewer-arrow viewer-next" aria-label="Next meme" disabled={viewerIndex >= memes.length - 1 && !nextCursor} on:click={() => void viewerNext()}>›</button>
      </div>
      {/key}
      <div class="viewer-details">
        {#if viewerMeme.description_status === 'pending'}
          <div class="viewer-description viewer-pending">writing a description…</div>
        {:else if viewerMeme.description}
          <div class="viewer-description">{viewerMeme.description}</div>
        {:else if viewerMeme.description_status === 'failed'}
          <div class="viewer-description viewer-failed">Description unavailable · <button on:click={() => retryDescription(viewerMeme)}>retry</button></div>
        {/if}
        <div class="viewer-tags">
          <button
            type="button"
            class="viewer-regenerate"
            title="Regenerate description and AI tags"
            aria-label="Regenerate description and AI tags"
            disabled={viewerMeme.description_status === 'pending'}
            on:click={() => retryDescription(viewerMeme)}
          >↻</button>
          {#each viewerMeme.tags || [] as tag}
            <button
              type="button"
              class:tag-delete-ready={tagDeleteActive && tagDeleteMemeId === viewerMeme.id && tagDeleteName === tag && tagDeleteOver}
              on:click={() => chooseTag(tag)}
              on:pointerdown={(event) => onTagPointerDown(event, viewerMeme, tag)}
              on:pointermove={onTagPointerMove}
              on:pointerup={onTagPointerUp}
              on:pointercancel={onTagPointerCancel}
            >#{tag}</button>
          {/each}
          <form on:submit|preventDefault={addTags}>
            <input bind:value={tagInput} placeholder="add tags" aria-label="Add tags" />
            <button type="submit" disabled={addingTags || !tagInput.trim()}>+</button>
          </form>
        </div>
      </div>
    </dialog>
  </div>
{/if}

<footer><span>meme/brary</span><span>no accounts · no tracking · just memes</span></footer>
