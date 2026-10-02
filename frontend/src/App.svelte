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

  let uploadFile;
  let uploadPreview = '';
  let tagsInput = '';
  let description = '';
  let dragActive = false;
  let uploading = false;
  let uploadMessage = '';

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
      memes = [meme, ...memes.filter((item) => item.id !== meme.id)];
      total += 1;
      resetUpload();
      if (meme.description_status === 'pending') pollDescription(meme.id);
    } catch (err) {
      uploadMessage = err.message;
    } finally {
      uploading = false;
    }
  }

  async function pollDescription(id) {
    for (let attempt = 0; attempt < 30; attempt += 1) {
      await new Promise((resolve) => setTimeout(resolve, 2500));
      try {
        const pending = memes.some((item) => item.id === id && item.description_status === 'pending');
        if (!pending) return;
        const params = new URLSearchParams({ limit: String(PAGE_SIZE) });
        if (selectedTag) params.set('tag', selectedTag);
        const refreshed = await api(`/api/memes?${params}`);
        const updated = refreshed.memes.find((item) => item.id === id);
        if (updated) memes = memes.map((item) => (item.id === id ? updated : item));
        if (updated && updated.description_status !== 'pending') return;
      } catch {
        return;
      }
    }
  }

  async function retryDescription(meme) {
    try {
      const pending = await api(`/api/memes/${meme.id}/describe`, { method: 'POST' });
      memes = memes.map((item) => (item.id === meme.id ? pending : item));
      pollDescription(meme.id);
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
    observer = new IntersectionObserver(([entry]) => {
      if (entry.isIntersecting && nextCursor && !loadingMore) load(false);
    }, { rootMargin: '500px' });
    if (sentinel) observer.observe(sentinel);
    return () => {
      observer?.disconnect();
      window.removeEventListener('dragover', preventWindowDrop);
      window.removeEventListener('drop', preventWindowDrop);
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
    <span class="anonymous"><span class="dot"></span> anonymous archive</span>
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
    <span class="drop-hint">PNG, JPG, GIF or WebP · up to 20 MB</span>
    <input bind:this={fileInput} class="visually-hidden" type="file" accept="image/jpeg,image/png,image/gif,image/webp" on:change={onFileInput} />
  </section>

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
        <article class="meme-card">
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
