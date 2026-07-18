<script lang="ts">
  import { Upload, X, Loader2 } from '@lucide/svelte';
  import { uploadMedia } from '$lib/api/media';
  import { Alert } from '$lib/components/ui';

  interface Props {
    photoUrl?: string;
    disabled?: boolean;
    size?: 'sm' | 'md' | 'lg';
    entityType?: string;
    entityId?: number;
    onPhotoChange?: (detail: { file: File | null; previewUrl: string; uploadedUrl?: string }) => void;
  }

  let { photoUrl = '', disabled = false, size = 'md', entityType = 'product', entityId, onPhotoChange }: Props = $props();

  let fileInput: HTMLInputElement;
  let previewUrl: string = $state(photoUrl);
  let isDragging: boolean = $state(false);
  let uploading: boolean = $state(false);
  let error: string = $state('');

  $effect(() => {
    previewUrl = photoUrl;
  });

  function handleFileSelect(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (file) {
      processFile(file);
    }
  }

  function handleDragOver(event: DragEvent) {
    event.preventDefault();
    isDragging = true;
  }

  function handleDragLeave(event: DragEvent) {
    event.preventDefault();
    isDragging = false;
  }

  function handleDrop(event: DragEvent) {
    event.preventDefault();
    isDragging = false;
    const file = event.dataTransfer?.files?.[0];
    if (file) {
      processFile(file);
    }
  }

  async function processFile(file: File) {
    error = '';
    if (!file.type.startsWith('image/')) {
      error = 'Por favor, selecione apenas arquivos de imagem (PNG, JPG, WEBP)';
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      error = 'A imagem deve ter no máximo 5MB';
      return;
    }

    // Preview local imediato
    const reader = new FileReader();
    reader.onload = (e) => {
      previewUrl = e.target?.result as string;
    };
    reader.readAsDataURL(file);

    // Upload real
    uploading = true;
    try {
      const uploaded = await uploadMedia(file, entityType, entityId);
      const uploadedUrl = `/${uploaded.FilePath}`;
      previewUrl = uploadedUrl;
      onPhotoChange?.({ file, previewUrl: uploadedUrl, uploadedUrl });
    } catch (error: any) {
      error = 'Erro ao fazer upload da imagem: ' + (error?.message || 'Tente novamente');
      previewUrl = '';
      if (fileInput) {
        fileInput.value = '';
      }
    } finally {
      uploading = false;
    }
  }

  function handleRemove() {
    previewUrl = '';
    if (fileInput) {
      fileInput.value = '';
    }
    onPhotoChange?.({ file: null, previewUrl: '' });
  }

  function triggerFileInput() {
    fileInput.click();
  }
</script>

<div class="photo-upload" class:size class:dragging={isDragging} class:has-photo={previewUrl}>
  {#if error}
    <Alert variant="error" dismissible onDismiss={() => error = ''} class="upload-alert">
      {error}
    </Alert>
  {/if}

  <input
    type="file"
    bind:this={fileInput}
    accept="image/png,image/jpeg,image/webp"
    onchange={handleFileSelect}
    disabled={disabled}
    class="file-input"
  />

  {#if uploading}
    <div class="upload-area">
      <div class="upload-icon uploading">
        <Loader2 size={32} class="spin" />
      </div>
      <div class="upload-text">
        <span class="upload-title">Enviando...</span>
      </div>
    </div>
  {:else if previewUrl}
    <div class="photo-preview">
      <img src={previewUrl} alt="Preview do produto" class="preview-image" />
      <button
        class="remove-button"
        onclick={handleRemove}
        disabled={disabled}
        title="Remover foto"
      >
        <X size={16} />
      </button>
    </div>
  {:else}
    <div
      class="upload-area"
      ondragover={handleDragOver}
      ondragleave={handleDragLeave}
      ondrop={handleDrop}
      onclick={triggerFileInput}
      role="button"
      tabindex="0"
      onkeydown={(e) => e.key === 'Enter' && triggerFileInput()}
    >
      <div class="upload-icon">
        <Upload size={32} />
      </div>
      <div class="upload-text">
        <span class="upload-title">Adicionar Foto</span>
        <span class="upload-subtitle">PNG • JPG • WEBP</span>
      </div>
    </div>
  {/if}
</div>

<style>
  .photo-upload {
    position: relative;
    width: 100%;
    border-radius: 12px;
    overflow: hidden;
    background: #f8fafc;
    border: 2px dashed #e2e8f0;
    transition: all 0.15s ease;
  }

  .photo-upload.sm {
    height: 120px;
  }

  .photo-upload.md {
    height: 200px;
  }

  .photo-upload.lg {
    height: 280px;
  }

  .photo-upload.dragging {
    border-color: #6366f1;
    background: #eef2ff;
  }

  .photo-upload.has-photo {
    border-style: solid;
    border-color: #f1f5f9;
    background: #ffffff;
  }

  .file-input {
    display: none;
  }

  .upload-alert {
    position: absolute;
    top: 0.5rem;
    left: 0.5rem;
    right: 0.5rem;
    z-index: 10;
  }

  .upload-area {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    gap: 1rem;
    cursor: pointer;
    padding: 1.5rem;
  }

  .upload-icon {
    color: #94a3b8;
    transition: color 0.15s ease;
  }

  .upload-icon.uploading {
    color: #6366f1;
  }

  .spin {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  .photo-upload:hover .upload-icon {
    color: #6366f1;
  }

  .upload-text {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.25rem;
    text-align: center;
  }

  .upload-title {
    font-size: 0.875rem;
    font-weight: 500;
    color: #0f172a;
  }

  .upload-subtitle {
    font-size: 0.75rem;
    color: #64748b;
  }

  .photo-preview {
    position: relative;
    width: 100%;
    height: 100%;
  }

  .preview-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .remove-button {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: rgba(0, 0, 0, 0.6);
    border: none;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: background 0.15s ease;
  }

  .remove-button:hover {
    background: rgba(0, 0, 0, 0.8);
  }

  .remove-button:focus {
    outline: none;
    box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.5);
  }

  .upload-area:focus {
    outline: none;
    box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.3);
  }
</style>
