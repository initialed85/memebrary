<script>
  import { onMount } from 'svelte';

  const PAGE_SIZE = 36;
  const MAX_UPLOAD_BYTES = 20 * 1024 * 1024;

  let memes = [];
  let nextCursor = '';
  let total = 0;
  let loading = true;
  let loadingMore = false;
  let error = '';
  let selectedTag = '';
  let tagQuery = '';
  let sentinel;
  let fileInput;
  let observer;
  let pollTimer;
  let pollInFlight = false;

  let uploadFile;
  let uploadPreview = '';
  let tagsInput = '';
  let description = '';
  let dragActive = false;
  let deleteActive = false;
  let draggedMemeId = '';
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
    } else if (loadingMore || !nextCursor) {
      return;
    }
    if (!reset) loadingMore = true;
    error = '';
    const params = new URLSearchParams({ limit: String(PAGE_SIZE) });
    if (nextCursor && !reset) params.set('cursor', nextCursor);
    if (selectedTag) params.set('tag', selectedTag);
    try {
      const result = await api(`/api/memes?${params}`);
      memes = reset ? result.memes : [...memes, ...result.memes];
      nextCursor = result.next_cursor || '';
      total = result.total || 0;
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
      loadingMore = false;
    }
  }

  async function refreshLatest() {
    if (pollInFlight || loading || loadingMore || document.visibilityState === 'hidden') return;
    pollInFlight = true;
    const params = new URLSearchParams({ limit: String(PAGE_SIZE) });
    if (selectedTag) params.set('tag', selectedTag);
    try {
      const result = await api(`/api/memes?${params}`);
      const remote = result.memes || [];
      const remoteIDs = new Set(remote.map((item) => item.id));
      // Keep already-loaded older pages while replacing the latest window. This
      // lets another browser's upload/AI result appear without jumping scroll.
      memes = [...remote, ...memes.filter((item) => !remoteIDs.has(item.id))];
      total = result.total || 0;
      if (!nextCursor || memes.length <= remote.length) nextCursor = result.next_cursor || '';
    } catch (err) {
      // A transient poll failure should not disrupt a timeline that is already visible.
      if (memes.length === 0) error = err.message;
    } finally {
      pollInFlight = false;
    }
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

  function clearTag() {
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

  function startMemeDrag(event, meme) {
    draggedMemeId = meme.id;
    event.dataTransfer.effectAllowed = 'move';
    event.dataTransfer.setData('application/x-memebrary-id', meme.id);
    event.dataTransfer.setData('text/plain', meme.id);
  }

  function endMemeDrag() {
    draggedMemeId = '';
    deleteActive = false;
  }

  function onDeleteDragOver(event) {
    if (!Array.from(event.dataTransfer.types).includes('application/x-memebrary-id')) return;
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
    deleteActive = true;
  }

  function onDeleteDragLeave() {
    deleteActive = false;
  }

  async function onDeleteDrop(event) {
    event.preventDefault();
    deleteActive = false;
    const id = event.dataTransfer.getData('application/x-memebrary-id') || event.dataTransfer.getData('text/plain');
    const meme = memes.find((item) => item.id === id);
    if (!meme || deletingId) return;
    if (!window.confirm('Delete this meme permanently?')) return;
    deletingId = meme.id;
    deleteError = '';
    try {
      const response = await fetch(`/api/memes/${meme.id}`, { method: 'DELETE' });
      if (!response.ok) {
        let body = null;
        try { body = await response.json(); } catch {}
        throw new Error(body?.error || `Delete failed (${response.status})`);
      }
      memes = memes.filter((item) => item.id !== meme.id);
      total = Math.max(0, total - 1);
    } catch (err) {
      deleteError = err.message;
    } finally {
      deletingId = '';
      draggedMemeId = '';
    }
  }

  function resetUpload() {
    if (uploadPreview) URL.revokeObjectURL(uploadPreview);
    uploadFile = null;
    uploadPreview = '';
    tagsInput = '';
    description = '';
    uploadMessage = '';
  }

  async function submitUpload() {
    if (!uploadFile || uploading) return;
    uploading = true;
    uploadMessage = '';
    const form = new FormData();
    form.set('file', uploadFile);
    form.set('tags', tagsInput);
    form.set('description', description.trim());
    try {
      const meme = await api('/api/memes', { method: 'POST', body: form });
      if (!selectedTag || meme.tags?.includes(selectedTag)) {
        memes = [meme, ...memes.filter((item) => item.id !== meme.id)];
        total += 1;
      }
      resetUpload();
    } catch (err) {
      uploadMessage = err.message;
    } finally {
      uploading = false;
    }
  }

  async function retryDescription(meme) {
    try {
      const pending = await api(`/api/memes/${meme.id}/describe`, { method: 'POST' });
      memes = memes.map((item) => (item.id === meme.id ? pending : item));
    } catch (err) {
      error = err.message;
    }
  }

  function formatDate(value) {
    const date = new Date(value);
    const now = new Date();
    if (date.toDateString() === now.toDateString()) {
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }
    return date.toLocaleDateString([], { month: 'short', day: 'numeric', year: date.getFullYear() === now.getFullYear() ? undefined : 'numeric' });
  }

  function shortName(name) {
    if (!name || name.length <= 32) return name;
    return `${name.slice(0, 29)}…`;
  }

  onMount(() => {
    load(true);
    window.addEventListener('dragover', preventWindowDrop);
    window.addEventListener('drop', preventWindowDrop);
    window.addEventListener('paste', onPaste);
    pollTimer = window.setInterval(refreshLatest, 8000);
    observer = new IntersectionObserver(([entry]) => {
      if (entry.isIntersecting && nextCursor && !loadingMore) load(false);
    }, { rootMargin: '500px' });
    if (sentinel) observer.observe(sentinel);
    return () => {
      observer?.disconnect();
      window.removeEventListener('dragover', preventWindowDrop);
      window.removeEventListener('drop', preventWindowDrop);
      window.removeEventListener('paste', onPaste);
      window.clearInterval(pollTimer);
      if (uploadPreview) URL.revokeObjectURL(uploadPreview);
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
    <span class="anonymous" title="The timeline refreshes every 8 seconds"><span class="dot"></span> anonymous · live</span>
    <form class="filter" on:submit={applyTag}>
      <label for="tag-filter">filter</label>
      <input id="tag-filter" bind:value={tagQuery} placeholder="#cats" autocomplete="off" />
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
        on:dragover={onDeleteDragOver}
        on:dragleave={onDeleteDragLeave}
        on:drop={onDeleteDrop}
      >
        <span class="trash-icon" aria-hidden="true">⌫</span>
        <div>
          <strong>{deletingId ? 'Deleting…' : 'Drag here to delete'}</strong>
          <span>{deleteError || 'release to remove forever'}</span>
        </div>
      </section>
    </div>

    {#if uploadFile}
      <section class="upload-editor" aria-label="New meme details">
        <div class="editor-preview"><img src={uploadPreview} alt="Selected meme preview" /></div>
        <div class="editor-fields">
          <label>
            <span>hashtags <small>(optional)</small></span>
            <input bind:value={tagsInput} placeholder="#reaction #work #animals" maxlength="500" />
          </label>
          <label>
            <span>description <small>(optional — AI fills this in)</small></span>
            <input bind:value={description} placeholder="What is happening in this meme?" maxlength="500" />
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

  <div class="timeline-heading">
    <div>
      <h1>{selectedTag ? `#${selectedTag}` : 'Latest memes'}</h1>
      {#if total}<span class="result-count">{total.toLocaleString()} {total === 1 ? 'meme' : 'memes'}</span>{/if}
    </div>
    {#if selectedTag}<button class="quiet-button" on:click={clearTag}>show everything</button>{/if}
  </div>

  {#if error}
    <div class="notice error" role="alert"><strong>Couldn’t load the library.</strong> {error} <button on:click={() => load(true)}>Try again</button></div>
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
    <section class="meme-grid" aria-label="Meme timeline">
      {#each memes as meme (meme.id)}
        <article
          class:dragging={draggedMemeId === meme.id}
          class="meme-card"
          draggable="true"
          on:dragstart={(event) => startMemeDrag(event, meme)}
          on:dragend={endMemeDrag}
        >
          <a class="image-frame" href={`/media/${meme.id}`} target="_blank" rel="noreferrer">
            <img loading="lazy" src={`/media/${meme.id}`} alt={meme.description || 'Meme image'} />
          </a>
          <div class="card-details">
            {#if meme.tags?.length}
              <div class="tags">
                {#each meme.tags as tag}<button on:click={() => chooseTag(tag)}>#{tag}</button>{/each}
              </div>
            {/if}
            {#if meme.description}
              <p class="description">{meme.description}</p>
            {:else if meme.description_status === 'pending'}
              <p class="description pending"><span class="mini-spinner"></span> writing a description…</p>
            {:else if meme.description_status === 'failed'}
              <p class="description failed">Description unavailable <button on:click={() => retryDescription(meme)}>retry</button></p>
            {:else}
              <p class="description muted">No description</p>
            {/if}
            <div class="card-footer"><time datetime={meme.created_at}>{formatDate(meme.created_at)}</time><span title={meme.original_name}>{shortName(meme.original_name)}</span></div>
          </div>
        </article>
      {/each}
    </section>
  {/if}

  <div bind:this={sentinel} class="load-sentinel" aria-hidden="true"></div>
  {#if loadingMore}<p class="loading-more">loading more…</p>{/if}
  {#if !loading && !nextCursor && memes.length > 0}<p class="end-note">— end of the archive —</p>{/if}
</main>

<footer><span>meme/brary</span><span>no accounts · no tracking · just memes</span></footer>
