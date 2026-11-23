<script>
  import { onMount } from 'svelte';
  import { apiRequest, getToken, setToken, uploadImage } from '../lib/api.js';

  let materials = [];
  let loading = true;
  let error = '';
  let activeTab = 'materials'; // 'materials' or 'lessons'

  // Modal / Form state for Material
  let showMaterialModal = false;
  let materialForm = {
    id: null,
    title: '',
    slug: '',
    description: '',
    order_index: 1
  };

  // Modal / Form state for Lesson
  let showLessonModal = false;
  let lessonForm = {
    id: null,
    material_id: null,
    title: '',
    slug: '',
    order_index: 1,
    image_url: '',
    content: '',
    notes: '',
    starter_code: '',
    blocks: '[]',
    enable_dartpad: true
  };

  // Active blocks array in Lesson Editor
  let blocks = [];
  let uploadingBlockId = null;

  onMount(async () => {
    if (!getToken()) {
      window.location.hash = '/admin/login';
      return;
    }
    await loadData();
  });

  async function loadData() {
    loading = true;
    error = '';
    try {
      materials = await apiRequest('/materials');
    } catch (err) {
      error = err.message;
      if (err.message.includes('401') || err.message.includes('Token')) {
        setToken(null);
        window.location.hash = '/admin/login';
      }
    } finally {
      loading = false;
    }
  }

  function handleLogout() {
    setToken(null);
    window.location.hash = '/admin/login';
  }

  // --- Material CRUD ---
  function openCreateMaterial() {
    materialForm = { id: null, title: '', slug: '', description: '', order_index: (materials.length + 1) };
    showMaterialModal = true;
  }

  function openEditMaterial(item) {
    materialForm = { ...item };
    showMaterialModal = true;
  }

  async function saveMaterial() {
    try {
      if (materialForm.id) {
        await apiRequest(`/admin/materials/${materialForm.id}`, {
          method: 'PUT',
          body: JSON.stringify(materialForm)
        });
      } else {
        await apiRequest('/admin/materials', {
          method: 'POST',
          body: JSON.stringify(materialForm)
        });
      }
      showMaterialModal = false;
      await loadData();
    } catch (err) {
      alert(`Gagal menyimpan materi: ${err.message}`);
    }
  }

  async function deleteMaterial(id) {
    if (!confirm('Apakah Anda yakin ingin menghapus materi ini beserta seluruh bab di dalamnya?')) return;
    try {
      await apiRequest(`/admin/materials/${id}`, { method: 'DELETE' });
      await loadData();
    } catch (err) {
      alert(`Gagal menghapus materi: ${err.message}`);
    }
  }

  // --- Dynamic Block Builder Functions ---
  function generateBlockId() {
    return 'blk_' + Date.now() + '_' + Math.random().toString(36).substr(2, 5);
  }

  function addBlock(type) {
    let newBlock;
    if (type === 'text') {
      newBlock = { id: generateBlockId(), type: 'text', content: '' };
    } else if (type === 'image') {
      newBlock = { id: generateBlockId(), type: 'image', url: '', caption: '' };
    } else if (type === 'note') {
      newBlock = { id: generateBlockId(), type: 'note', content: '' };
    } else if (type === 'code') {
      newBlock = { id: generateBlockId(), type: 'code', title: 'Contoh Kode Dart', code: 'void main() {\n  print("Hello Dart!");\n}' };
    }
    if (newBlock) {
      blocks = [...blocks, newBlock];
    }
  }

  function removeBlock(index) {
    blocks = blocks.filter((_, i) => i !== index);
  }

  function moveBlockUp(index) {
    if (index === 0) return;
    const temp = blocks[index];
    blocks[index] = blocks[index - 1];
    blocks[index - 1] = temp;
    blocks = [...blocks];
  }

  function moveBlockDown(index) {
    if (index === blocks.length - 1) return;
    const temp = blocks[index];
    blocks[index] = blocks[index + 1];
    blocks[index + 1] = temp;
    blocks = [...blocks];
  }

  async function handleBlockImageUpload(e, blockId) {
    const file = e.target.files[0];
    if (!file) return;

    uploadingBlockId = blockId;
    try {
      const res = await uploadImage(file);
      const bIndex = blocks.findIndex(b => b.id === blockId);
      if (bIndex !== -1) {
        blocks[bIndex].url = res.url;
        blocks = [...blocks];
      }
    } catch (err) {
      alert(err.message || 'Gagal mengupload gambar blok');
    } finally {
      uploadingBlockId = null;
    }
  }

  // --- Lesson CRUD ---
  function openCreateLesson(defaultMaterialId = null) {
    blocks = [
      { id: generateBlockId(), type: 'text', content: '### Judul Pengantar\n\nTuliskan penjelasan materi dalam format Markdown...' }
    ];
    lessonForm = {
      id: null,
      material_id: defaultMaterialId || (materials[0] ? materials[0].id : null),
      title: '',
      slug: '',
      order_index: 1,
      image_url: '',
      content: '',
      notes: '',
      starter_code: 'void main() {\n  print("Hello Dart!");\n}',
      blocks: JSON.stringify(blocks),
      enable_dartpad: true
    };
    showLessonModal = true;
  }

  function openEditLesson(lesson) {
    lessonForm = { ...lesson };
    
    // Parse blocks or fallback from legacy content
    if (lesson.blocks && lesson.blocks.trim() !== '' && lesson.blocks.trim() !== '[]') {
      try {
        blocks = JSON.parse(lesson.blocks);
      } catch (e) {
        blocks = [];
      }
    } else {
      // Create blocks from legacy fields if any
      blocks = [];
      if (lesson.image_url) {
        blocks.push({ id: generateBlockId(), type: 'image', url: lesson.image_url, caption: '' });
      }
      if (lesson.content) {
        blocks.push({ id: generateBlockId(), type: 'text', content: lesson.content });
      }
      if (lesson.notes) {
        blocks.push({ id: generateBlockId(), type: 'note', content: lesson.notes });
      }
      if (lesson.starter_code) {
        blocks.push({ id: generateBlockId(), type: 'code', title: 'Starter Code', code: lesson.starter_code });
      }
    }

    if (blocks.length === 0) {
      blocks = [{ id: generateBlockId(), type: 'text', content: '' }];
    }

    showLessonModal = true;
  }

  async function saveLesson() {
    try {
      if (!lessonForm.material_id) {
        alert('Pilih modul materi induk terlebih dahulu');
        return;
      }
      lessonForm.material_id = parseInt(lessonForm.material_id);
      lessonForm.order_index = parseInt(lessonForm.order_index) || 1;
      lessonForm.blocks = JSON.stringify(blocks);

      // Sinkronisasi legacy fields untuk backward compatibility
      const firstText = blocks.find(b => b.type === 'text');
      const firstImg = blocks.find(b => b.type === 'image');
      const firstNote = blocks.find(b => b.type === 'note');
      const firstCode = blocks.find(b => b.type === 'code');
      if (firstText) lessonForm.content = firstText.content;
      if (firstImg) lessonForm.image_url = firstImg.url;
      if (firstNote) lessonForm.notes = firstNote.content;
      if (firstCode && !lessonForm.starter_code) lessonForm.starter_code = firstCode.code;

      if (lessonForm.id) {
        await apiRequest(`/admin/lessons/${lessonForm.id}`, {
          method: 'PUT',
          body: JSON.stringify(lessonForm)
        });
      } else {
        await apiRequest('/admin/lessons', {
          method: 'POST',
          body: JSON.stringify(lessonForm)
        });
      }
      showLessonModal = false;
      await loadData();
    } catch (err) {
      alert(`Gagal menyimpan bab: ${err.message}`);
    }
  }

  async function deleteLesson(id) {
    if (!confirm('Apakah Anda yakin ingin menghapus bab ini?')) return;
    try {
      await apiRequest(`/admin/lessons/${id}`, { method: 'DELETE' });
      await loadData();
    } catch (err) {
      alert(`Gagal menghapus bab: ${err.message}`);
    }
  }

  function autoSlug(title) {
    return title.toLowerCase().trim()
      .replace(/[^\w\s-]/g, '')
      .replace(/[\s_-]+/g, '-')
      .replace(/^-+|-+$/g, '');
  }
</script>

<div class="max-w-6xl mx-auto px-4 py-8">
  <!-- Top Bar -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between pb-6 mb-6 border-b border-slate-200 gap-4">
    <div>
      <span class="text-xs font-bold text-[#0175C2] uppercase tracking-wider">Admin Workspace</span>
      <h1 class="text-2xl font-extrabold text-slate-900">Manajemen Konten E-Learning</h1>
    </div>

    <div class="flex items-center gap-3">
      <a href="#/" class="text-xs font-semibold px-3 py-2 text-slate-600 bg-white border border-slate-200 rounded-lg hover:bg-slate-50 transition">
        Lihat Web Publik
      </a>
      <button 
        on:click={handleLogout}
        class="text-xs font-semibold px-3 py-2 text-red-600 bg-red-50 border border-red-100 rounded-lg hover:bg-red-100 transition"
      >
        Logout
      </button>
    </div>
  </div>

  <!-- Tabs Navigation -->
  <div class="flex border-b border-slate-200 mb-6 gap-6">
    <button 
      on:click={() => activeTab = 'materials'}
      class="pb-3 text-sm font-semibold transition border-b-2 {activeTab === 'materials' ? 'border-[#0175C2] text-[#0175C2]' : 'border-transparent text-slate-500 hover:text-slate-800'}"
    >
      Daftar Materi ({materials.length})
    </button>
    <button 
      on:click={() => activeTab = 'lessons'}
      class="pb-3 text-sm font-semibold transition border-b-2 {activeTab === 'lessons' ? 'border-[#0175C2] text-[#0175C2]' : 'border-transparent text-slate-500 hover:text-slate-800'}"
    >
      Semua Bab Pembelajaran
    </button>
  </div>

  {#if loading}
    <div class="py-16 text-center text-slate-400">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-[#0175C2]"></div>
      <p class="mt-2 text-sm">Memuat data admin...</p>
    </div>
  {:else if error}
    <div class="p-4 bg-red-50 border border-red-200 text-red-700 text-sm rounded-lg mb-6">
      {error}
    </div>
  {:else}
    <!-- Tab 1: Materials Management -->
    {#if activeTab === 'materials'}
      <div class="space-y-4">
        <div class="flex justify-between items-center mb-4">
          <p class="text-xs text-slate-500">Kelola modul kursus tingkat utama</p>
          <button 
            on:click={openCreateMaterial}
            class="px-4 py-2 bg-[#0175C2] hover:bg-[#02569B] text-white text-xs font-semibold rounded-lg transition"
          >
            + Tambah Materi Baru
          </button>
        </div>

        <div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
          <table class="w-full text-left text-sm">
            <thead class="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-600">
              <tr>
                <th class="p-3 w-16 text-center">Urutan</th>
                <th class="p-3">Judul Materi</th>
                <th class="p-3">Slug URL</th>
                <th class="p-3 w-28 text-center">Total Bab</th>
                <th class="p-3 w-36 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              {#each materials as item}
                <tr class="hover:bg-slate-50/50">
                  <td class="p-3 text-center font-mono text-xs text-slate-500">{item.order_index}</td>
                  <td class="p-3 font-semibold text-slate-900">{item.title}</td>
                  <td class="p-3 font-mono text-xs text-slate-500">/materials/{item.slug}</td>
                  <td class="p-3 text-center text-xs">
                    <span class="px-2 py-0.5 rounded-full bg-slate-100 text-slate-700 font-medium">
                      {item.lessons ? item.lessons.length : 0}
                    </span>
                  </td>
                  <td class="p-3 text-right space-x-2">
                    <button 
                      on:click={() => openCreateLesson(item.id)}
                      class="text-xs text-[#0175C2] hover:underline font-medium"
                    >
                      + Bab
                    </button>
                    <button 
                      on:click={() => openEditMaterial(item)}
                      class="text-xs text-slate-600 hover:text-slate-900 font-medium"
                    >
                      Edit
                    </button>
                    <button 
                      on:click={() => deleteMaterial(item.id)}
                      class="text-xs text-red-600 hover:text-red-800 font-medium"
                    >
                      Hapus
                    </button>
                  </td>
                </tr>
              {/each}
              {#if materials.length === 0}
                <tr>
                  <td colspan="5" class="p-6 text-center text-slate-400 text-xs">Belum ada materi. Klik tombol tambah di atas.</td>
                </tr>
              {/if}
            </tbody>
          </table>
        </div>
      </div>
    {/if}

    <!-- Tab 2: Lessons Management -->
    {#if activeTab === 'lessons'}
      <div class="space-y-4">
        <div class="flex justify-between items-center mb-4">
          <p class="text-xs text-slate-500">Kelola bab pelajaran di dalam setiap materi</p>
          <button 
            on:click={() => openCreateLesson()}
            class="px-4 py-2 bg-[#0175C2] hover:bg-[#02569B] text-white text-xs font-semibold rounded-lg transition"
          >
            + Tambah Bab Baru
          </button>
        </div>

        {#each materials as mat}
          <div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs mb-6">
            <div class="p-3.5 bg-slate-50 border-b border-slate-200 flex justify-between items-center">
              <span class="font-bold text-sm text-slate-900">{mat.title}</span>
              <button 
                on:click={() => openCreateLesson(mat.id)}
                class="text-xs font-semibold text-[#0175C2] hover:underline"
              >
                + Tambah Bab ke Materi Ini
              </button>
            </div>

            <table class="w-full text-left text-sm">
              <thead class="text-xs font-semibold text-slate-500 border-b border-slate-100">
                <tr>
                  <th class="p-3 w-16 text-center">Urutan</th>
                  <th class="p-3">Judul Bab</th>
                  <th class="p-3">Slug</th>
                  <th class="p-3 w-28 text-center">DartPad</th>
                  <th class="p-3 w-32 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                {#if mat.lessons && mat.lessons.length > 0}
                  {#each mat.lessons as lesson}
                    <tr class="hover:bg-slate-50/50">
                      <td class="p-3 text-center font-mono text-xs text-slate-400">{lesson.order_index}</td>
                      <td class="p-3 font-medium text-slate-800">{lesson.title}</td>
                      <td class="p-3 font-mono text-xs text-slate-500">{lesson.slug}</td>
                      <td class="p-3 text-center">
                        {#if lesson.enable_dartpad}
                          <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-emerald-50 text-emerald-700">Aktif</span>
                        {:else}
                          <span class="text-[10px] font-medium px-2 py-0.5 rounded bg-slate-100 text-slate-500">Nonaktif</span>
                        {/if}
                      </td>
                      <td class="p-3 text-right space-x-2">
                        <button 
                          on:click={() => openEditLesson(lesson)}
                          class="text-xs text-slate-600 hover:text-slate-900 font-medium"
                        >
                          Edit
                        </button>
                        <button 
                          on:click={() => deleteLesson(lesson.id)}
                          class="text-xs text-red-600 hover:text-red-800 font-medium"
                        >
                          Hapus
                        </button>
                      </td>
                    </tr>
                  {/each}
                {:else}
                  <tr>
                    <td colspan="5" class="p-4 text-center text-slate-400 text-xs">Belum ada bab untuk materi ini.</td>
                  </tr>
                {/if}
              </tbody>
            </table>
          </div>
        {/each}
      </div>
    {/if}
  {/if}
</div>

<!-- MODAL: Form Material -->
{#if showMaterialModal}
  <div class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
    <div class="bg-white rounded-2xl max-w-lg w-full p-6 shadow-xl border border-slate-200">
      <h3 class="text-lg font-bold text-slate-900 mb-4">
        {materialForm.id ? 'Edit Materi' : 'Tambah Materi Baru'}
      </h3>

      <div class="space-y-4">
        <div>
          <label for="mat-title" class="block text-xs font-semibold text-slate-700 mb-1">Judul Materi</label>
          <input 
            id="mat-title"
            type="text" 
            bind:value={materialForm.title} 
            on:input={() => { if (!materialForm.id) materialForm.slug = autoSlug(materialForm.title); }}
            placeholder="Contoh: Pengenalan Dasar Dart" 
            class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm"
          />
        </div>

        <div>
          <label for="mat-slug" class="block text-xs font-semibold text-slate-700 mb-1">Slug URL</label>
          <input 
            id="mat-slug"
            type="text" 
            bind:value={materialForm.slug} 
            placeholder="pengenalan-dasar-dart" 
            class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm font-mono text-xs"
          />
        </div>

        <div>
          <label for="mat-desc" class="block text-xs font-semibold text-slate-700 mb-1">Deskripsi Singkat</label>
          <textarea 
            id="mat-desc"
            bind:value={materialForm.description} 
            rows="3" 
            placeholder="Ringkasan tentang materi ini..." 
            class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm"
          ></textarea>
        </div>

        <div>
          <label for="mat-order" class="block text-xs font-semibold text-slate-700 mb-1">Nomor Urutan Tampil</label>
          <input 
            id="mat-order"
            type="number" 
            bind:value={materialForm.order_index} 
            class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm"
          />
        </div>
      </div>

      <div class="mt-6 flex justify-end gap-2">
        <button 
          on:click={() => showMaterialModal = false}
          class="px-4 py-2 text-xs font-semibold text-slate-600 hover:bg-slate-100 rounded-lg transition"
        >
          Batal
        </button>
        <button 
          on:click={saveMaterial}
          class="px-4 py-2 text-xs font-semibold text-white bg-[#0175C2] hover:bg-[#02569B] rounded-lg transition"
        >
          Simpan Materi
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- MODAL: Form Lesson (Dynamic Vertical Block Builder) -->
{#if showLessonModal}
  <div class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50 overflow-y-auto">
    <div class="bg-white rounded-2xl max-w-4xl w-full p-6 shadow-xl border border-slate-200 my-8 flex flex-col max-h-[90vh]">
      
      <!-- Modal Header -->
      <div class="flex items-center justify-between pb-4 border-b border-slate-200 flex-shrink-0">
        <div>
          <h3 class="text-lg font-bold text-slate-900">
            {lessonForm.id ? 'Edit Bab Pelajaran' : 'Tambah Bab Baru'}
          </h3>
          <p class="text-xs text-slate-500">Susun konten bab secara vertikal & modular seperti Google Forms</p>
        </div>
        <button on:click={() => showLessonModal = false} class="text-slate-400 hover:text-slate-600 text-lg font-bold">
          ✕
        </button>
      </div>

      <!-- Modal Body (Scrollable) -->
      <div class="space-y-6 overflow-y-auto py-4 pr-1 flex-1">
        
        <!-- Metadata Header -->
        <div class="bg-slate-50/80 p-4 rounded-xl border border-slate-200 space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label for="les-mat" class="block text-xs font-semibold text-slate-700 mb-1">Pilih Materi Induk</label>
              <select 
                id="les-mat"
                bind:value={lessonForm.material_id} 
                class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm bg-white"
              >
                {#each materials as m}
                  <option value={m.id}>{m.title}</option>
                {/each}
              </select>
            </div>

            <div>
              <label for="les-order" class="block text-xs font-semibold text-slate-700 mb-1">Nomor Urutan Bab</label>
              <input 
                id="les-order"
                type="number" 
                bind:value={lessonForm.order_index} 
                class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm bg-white"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label for="les-title" class="block text-xs font-semibold text-slate-700 mb-1">Judul Bab</label>
              <input 
                id="les-title"
                type="text" 
                bind:value={lessonForm.title} 
                on:input={() => { if (!lessonForm.id) lessonForm.slug = autoSlug(lessonForm.title); }}
                placeholder="Contoh: 1. Pengantar Dart & Struktur main()" 
                class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm bg-white"
              />
            </div>

            <div>
              <label for="les-slug" class="block text-xs font-semibold text-slate-700 mb-1">Slug URL Bab</label>
              <input 
                id="les-slug"
                type="text" 
                bind:value={lessonForm.slug} 
                placeholder="pengantar-dart-dan-main" 
                class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm font-mono text-xs bg-white"
              />
            </div>
          </div>
        </div>

        <!-- Dynamic Vertical Blocks Builder -->
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-bold text-slate-900">Blok Konten Bab ({blocks.length})</h4>
              <p class="text-xs text-slate-500">Tambahkan dan atur urutan blok vertikal untuk menyusun materi</p>
            </div>

            <!-- Quick Add Bar -->
            <div class="flex items-center gap-1.5 flex-wrap">
              <span class="text-xs text-slate-400 font-medium mr-1">+ Tambah:</span>
              <button 
                type="button" 
                on:click={() => addBlock('text')}
                class="px-2.5 py-1 text-xs font-semibold rounded bg-slate-100 hover:bg-blue-50 hover:text-[#0175C2] text-slate-700 transition border border-slate-200"
              >
                📝 Text
              </button>
              <button 
                type="button" 
                on:click={() => addBlock('image')}
                class="px-2.5 py-1 text-xs font-semibold rounded bg-slate-100 hover:bg-emerald-50 hover:text-emerald-700 text-slate-700 transition border border-slate-200"
              >
                🖼️ Image
              </button>
              <button 
                type="button" 
                on:click={() => addBlock('note')}
                class="px-2.5 py-1 text-xs font-semibold rounded bg-slate-100 hover:bg-amber-50 hover:text-amber-700 text-slate-700 transition border border-slate-200"
              >
                💡 Note
              </button>
              <button 
                type="button" 
                on:click={() => addBlock('code')}
                class="px-2.5 py-1 text-xs font-semibold rounded bg-slate-100 hover:bg-purple-50 hover:text-purple-700 text-slate-700 transition border border-slate-200"
              >
                💻 Code
              </button>
            </div>
          </div>

          <!-- Blocks List -->
          {#if blocks.length === 0}
            <div class="text-center py-8 bg-slate-50 border-2 border-dashed border-slate-200 rounded-xl">
              <p class="text-xs text-slate-400">Belum ada blok konten. Gunakan tombol di atas untuk menambahkan blok pertama.</p>
            </div>
          {:else}
            <div class="space-y-4">
              {#each blocks as block, index (block.id || index)}
                <div class="p-4 rounded-xl border border-slate-200 bg-white shadow-xs hover:border-slate-300 transition-all">
                  
                  <!-- Block Card Top Header -->
                  <div class="flex items-center justify-between pb-3 mb-3 border-b border-slate-100">
                    <div class="flex items-center gap-2">
                      <span class="text-xs font-mono font-bold px-2 py-0.5 rounded 
                        {block.type === 'text' ? 'bg-blue-100 text-[#0175C2]' : ''}
                        {block.type === 'image' ? 'bg-emerald-100 text-emerald-800' : ''}
                        {block.type === 'note' ? 'bg-amber-100 text-amber-800' : ''}
                        {block.type === 'code' ? 'bg-purple-100 text-purple-800' : ''}"
                      >
                        {index + 1}. Blok {block.type.toUpperCase()}
                      </span>
                    </div>

                    <!-- Reorder & Delete Actions -->
                    <div class="flex items-center gap-1">
                      <button
                        type="button"
                        on:click={() => moveBlockUp(index)}
                        disabled={index === 0}
                        class="p-1 text-slate-500 hover:text-slate-800 disabled:opacity-30 rounded hover:bg-slate-100 transition"
                        title="Geser ke Atas"
                      >
                        ▲
                      </button>
                      <button
                        type="button"
                        on:click={() => moveBlockDown(index)}
                        disabled={index === blocks.length - 1}
                        class="p-1 text-slate-500 hover:text-slate-800 disabled:opacity-30 rounded hover:bg-slate-100 transition"
                        title="Geser ke Bawah"
                      >
                        ▼
                      </button>
                      <button
                        type="button"
                        on:click={() => removeBlock(index)}
                        class="p-1 ml-2 text-red-500 hover:text-red-700 rounded hover:bg-red-50 transition"
                        title="Hapus Blok Ini"
                      >
                        ✕
                      </button>
                    </div>
                  </div>

                  <!-- Block Content Based on Type -->
                  {#if block.type === 'text'}
                    <div>
                      <label class="block text-xs font-semibold text-slate-700 mb-1">Teks / Markdown</label>
                      <textarea
                        bind:value={block.content}
                        rows="4"
                        placeholder="Tuliskan penjelasan bab, poin-poin materi, atau teks Markdown di sini..."
                        class="w-full px-3 py-2 border border-slate-300 rounded-lg text-xs font-mono bg-white"
                      ></textarea>
                    </div>

                  {:else if block.type === 'image'}
                    <div class="space-y-3">
                      <div>
                        <label class="block text-xs font-semibold text-slate-700 mb-1">URL Gambar atau Upload File</label>
                        <div class="flex gap-2">
                          <input
                            type="text"
                            bind:value={block.url}
                            placeholder="https://... atau upload file lokal"
                            class="flex-1 px-3 py-2 border border-slate-300 rounded-lg text-xs bg-white"
                          />
                          <label class="cursor-pointer px-3 py-2 bg-slate-100 hover:bg-slate-200 border border-slate-300 rounded-lg text-xs font-medium text-slate-700 flex-shrink-0">
                            {uploadingBlockId === block.id ? 'Mengunggah...' : 'Upload Gambar'}
                            <input
                              type="file"
                              accept="image/*"
                              on:change={(e) => handleBlockImageUpload(e, block.id)}
                              class="hidden"
                              disabled={uploadingBlockId === block.id}
                            />
                          </label>
                        </div>
                      </div>

                      <div>
                        <label class="block text-xs font-semibold text-slate-700 mb-1">Keterangan / Caption Gambar (Opsional)</label>
                        <input
                          type="text"
                          bind:value={block.caption}
                          placeholder="Contoh: Diagram Alur Program Dart"
                          class="w-full px-3 py-2 border border-slate-300 rounded-lg text-xs bg-white"
                        />
                      </div>

                      {#if block.url}
                        <div class="p-2 bg-slate-50 border rounded-lg flex items-center gap-3">
                          <img src={block.url} alt={block.caption || 'Preview'} class="h-16 w-24 object-cover rounded border" />
                          <span class="text-xs text-slate-500 truncate">{block.caption || block.url}</span>
                        </div>
                      {/if}
                    </div>

                  {:else if block.type === 'note'}
                    <div>
                      <label class="block text-xs font-semibold text-amber-800 mb-1">Catatan Penting / Tip Callout</label>
                      <textarea
                        bind:value={block.content}
                        rows="2"
                        placeholder="💡 Contoh: Selalu gunakan widget Const jika properti tidak pernah berubah..."
                        class="w-full px-3 py-2 border border-amber-300 bg-amber-50/40 rounded-lg text-xs text-amber-950"
                      ></textarea>
                    </div>

                  {:else if block.type === 'code'}
                    <div class="space-y-2">
                      <div>
                        <label class="block text-xs font-semibold text-slate-700 mb-1">Judul / Label Cuplikan Kode</label>
                        <input
                          type="text"
                          bind:value={block.title}
                          placeholder="Contoh: Contoh 1: Fungsi Named Parameter"
                          class="w-full px-3 py-2 border border-slate-300 rounded-lg text-xs bg-white font-medium"
                        />
                      </div>
                      <div>
                        <label class="block text-xs font-semibold text-slate-700 mb-1">Kode Dart</label>
                        <textarea
                          bind:value={block.code}
                          rows="4"
                          placeholder="void main() &#123; ... &#125;"
                          class="w-full px-3 py-2 border border-slate-700 rounded-lg text-xs font-mono bg-slate-900 text-slate-100"
                        ></textarea>
                      </div>
                    </div>
                  {/if}

                </div>
              {/each}
            </div>
          {/if}
        </div>

        <!-- Dart Simulation Settings -->
        <div class="p-4 bg-sky-50/60 border border-sky-100 rounded-xl space-y-3">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-xs font-bold text-[#0175C2]">Area Simulasi Dart (Interactive DartPad)</span>
              <p class="text-[11px] text-slate-500">Aktifkan tombol toggle panel DartPad interaktif di sisi kanan</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" bind:checked={lessonForm.enable_dartpad} class="sr-only peer" />
              <div class="w-11 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[#0175C2]"></div>
            </label>
          </div>

          {#if lessonForm.enable_dartpad}
            <div>
              <label for="starter-code" class="block text-xs font-semibold text-slate-700 mb-1">Default Starter Code (Jika Ingin Disediakan)</label>
              <textarea 
                id="starter-code"
                bind:value={lessonForm.starter_code} 
                rows="3" 
                placeholder="void main() &#123; print('Hello Dart'); &#125;" 
                class="w-full px-3 py-2 border border-slate-300 rounded-lg text-xs font-mono bg-slate-900 text-slate-100"
              ></textarea>
            </div>
          {/if}
        </div>

      </div>

      <!-- Modal Footer -->
      <div class="flex justify-end gap-2 pt-4 border-t border-slate-200 flex-shrink-0">
        <button 
          type="button"
          on:click={() => showLessonModal = false}
          class="px-4 py-2 text-xs font-semibold text-slate-600 hover:bg-slate-100 rounded-lg transition"
        >
          Batal
        </button>
        <button 
          type="button"
          on:click={saveLesson}
          class="px-4 py-2 text-xs font-semibold text-white bg-[#0175C2] hover:bg-[#02569B] rounded-lg transition"
        >
          Simpan Seluruh Bab & Blok
        </button>
      </div>

    </div>
  </div>
{/if}
