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
    enable_dartpad: true
  };

  let imageUploading = false;

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

  // --- Lesson CRUD ---
  function openCreateLesson(defaultMaterialId = null) {
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
      enable_dartpad: true
    };
    showLessonModal = true;
  }

  function openEditLesson(lesson) {
    lessonForm = { ...lesson };
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

  // Image Upload handler
  async function handleImageUpload(e) {
    const file = e.target.files[0];
    if (!file) return;

    imageUploading = true;
    try {
      const res = await uploadImage(file);
      lessonForm.image_url = res.url;
    } catch (err) {
      alert(err.message);
    } finally {
      imageUploading = false;
    }
  }

  function autoSlug(title, targetObj) {
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
          <label class="block text-xs font-semibold text-slate-700 mb-1">Judul Materi</label>
          <input 
            type="text" 
            bind:value={materialForm.title} 
            on:input={() => { if (!materialForm.id) materialForm.slug = autoSlug(materialForm.title); }}
            placeholder="Contoh: Pengenalan Dasar Dart" 
            class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Slug URL</label>
          <input 
            type="text" 
            bind:value={materialForm.slug} 
            placeholder="pengenalan-dasar-dart" 
            class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm font-mono text-xs"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Deskripsi Singkat</label>
          <textarea 
            bind:value={materialForm.description} 
            rows="3" 
            placeholder="Ringkasan tentang materi ini..." 
            class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm"
          ></textarea>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Nomor Urutan Tampil</label>
          <input 
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

<!-- MODAL: Form Lesson -->
{#if showLessonModal}
  <div class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50 overflow-y-auto">
    <div class="bg-white rounded-2xl max-w-2xl w-full p-6 shadow-xl border border-slate-200 my-8">
      <h3 class="text-lg font-bold text-slate-900 mb-4">
        {lessonForm.id ? 'Edit Bab Pelajaran' : 'Tambah Bab Baru'}
      </h3>

      <div class="space-y-4 max-h-[70vh] overflow-y-auto pr-2">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">Pilih Materi Induk</label>
            <select 
              bind:value={lessonForm.material_id} 
              class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm bg-white"
            >
              {#each materials as m}
                <option value={m.id}>{m.title}</option>
              {/each}
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">Nomor Urutan Bab</label>
            <input 
              type="number" 
              bind:value={lessonForm.order_index} 
              class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">Judul Bab</label>
            <input 
              type="text" 
              bind:value={lessonForm.title} 
              on:input={() => { if (!lessonForm.id) lessonForm.slug = autoSlug(lessonForm.title); }}
              placeholder="Contoh: 1. Dasar Variabel & Tipe Data" 
              class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm"
            />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1">Slug URL Bab</label>
            <input 
              type="text" 
              bind:value={lessonForm.slug} 
              placeholder="dasar-variabel" 
              class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm font-mono text-xs"
            />
          </div>
        </div>

        <!-- Image URL / Upload -->
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Gambar Materi (Opsional)</label>
          <div class="flex gap-2 items-center">
            <input 
              type="text" 
              bind:value={lessonForm.image_url} 
              placeholder="Input URL gambar atau upload langsung ->" 
              class="flex-1 px-3 py-2 border border-slate-300 rounded-lg text-xs"
            />
            <label class="cursor-pointer px-3 py-2 bg-slate-100 hover:bg-slate-200 border border-slate-300 rounded-lg text-xs font-medium text-slate-700 flex-shrink-0">
              {imageUploading ? 'Uploading...' : 'Pilih File'}
              <input type="file" accept="image/*" on:change={handleImageUpload} class="hidden" disabled={imageUploading} />
            </label>
          </div>
          {#if lessonForm.image_url}
            <div class="mt-2 text-xs text-slate-500 flex items-center gap-2">
              <span class="truncate">Preview:</span>
              <img src={lessonForm.image_url} alt="Preview" class="h-10 w-16 object-cover rounded border" />
            </div>
          {/if}
        </div>

        <!-- Markdown Content -->
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Konten Utama (Markdown Didukung)</label>
          <textarea 
            bind:value={lessonForm.content} 
            rows="6" 
            placeholder="Tuliskan materi pelajaran di sini menggunakan format Markdown..." 
            class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm font-mono"
          ></textarea>
        </div>

        <!-- Notes / Callout -->
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Area Notes / Tips Penting (Opsional)</label>
          <textarea 
            bind:value={lessonForm.notes} 
            rows="2" 
            placeholder="Catatan atau tips khusus untuk siswa..." 
            class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm"
          ></textarea>
        </div>

        <!-- Dart Simulation Settings -->
        <div class="p-4 bg-sky-50/50 border border-sky-100 rounded-xl space-y-3">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-xs font-bold text-[#0175C2]">Area Simulasi Dart (Interactive DartPad)</span>
              <p class="text-[11px] text-slate-500">Aktifkan layout split-screen dengan embed DartPad di sisi kanan</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" bind:checked={lessonForm.enable_dartpad} class="sr-only peer" />
              <div class="w-11 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[#0175C2]"></div>
            </label>
          </div>

          {#if lessonForm.enable_dartpad}
            <div>
              <label class="block text-xs font-semibold text-slate-700 mb-1">Starter Code (Dart)</label>
              <textarea 
                bind:value={lessonForm.starter_code} 
                rows="4" 
                placeholder="void main() &#123; print('Hello Dart'); &#125;" 
                class="w-full px-3 py-2 border border-slate-300 rounded-lg text-xs font-mono bg-slate-900 text-slate-100"
              ></textarea>
            </div>
          {/if}
        </div>

      </div>

      <div class="mt-6 flex justify-end gap-2 pt-4 border-t border-slate-100">
        <button 
          on:click={() => showLessonModal = false}
          class="px-4 py-2 text-xs font-semibold text-slate-600 hover:bg-slate-100 rounded-lg transition"
        >
          Batal
        </button>
        <button 
          on:click={saveLesson}
          class="px-4 py-2 text-xs font-semibold text-white bg-[#0175C2] hover:bg-[#02569B] rounded-lg transition"
        >
          Simpan Bab
        </button>
      </div>
    </div>
  </div>
{/if}
