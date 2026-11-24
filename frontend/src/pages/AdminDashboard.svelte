<script>
  import { onMount } from 'svelte';
  import { apiRequest, getToken, setToken, uploadImage } from '../lib/api.js';
  import { 
    Plus, 
    Edit, 
    Trash2, 
    ArrowUp, 
    ArrowDown, 
    X, 
    Upload, 
    ExternalLink, 
    LogOut, 
    BookOpen, 
    Layers, 
    Loader2, 
    AlertCircle,
    Type,
    Image,
    Lightbulb,
    Code,
    Check,
    Users,
    KeyRound,
    ShieldCheck
  } from 'lucide-svelte';

  let materials = [];
  let users = [];
  let currentUser = null;
  let loading = true;
  let error = '';
  let activeTab = 'materials'; // 'materials', 'lessons', or 'admins'

  // Admin Management state
  let adminMessage = '';
  let adminError = '';
  let adminActionLoading = false;
  let newAdminUsername = '';
  let newAdminPassword = '';
  let currentPassword = '';
  let newPassword = '';
  let confirmNewPassword = '';

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
      const [mats, meData, usersData] = await Promise.all([
        apiRequest('/materials'),
        apiRequest('/admin/me').catch(() => null),
        apiRequest('/admin/users').catch(() => [])
      ]);
      materials = mats;
      currentUser = meData;
      users = usersData;
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

  async function loadUsers() {
    try {
      users = await apiRequest('/admin/users');
    } catch (err) {
      adminError = err.message;
    }
  }

  async function handleCreateAdmin() {
    adminMessage = '';
    adminError = '';
    if (!newAdminUsername.trim() || !newAdminPassword) {
      adminError = 'Username dan password wajib diisi';
      return;
    }
    adminActionLoading = true;
    try {
      await apiRequest('/admin/users', {
        method: 'POST',
        body: JSON.stringify({
          username: newAdminUsername.trim(),
          password: newAdminPassword
        })
      });
      adminMessage = `Admin "${newAdminUsername}" berhasil ditambahkan!`;
      newAdminUsername = '';
      newAdminPassword = '';
      await loadUsers();
    } catch (err) {
      adminError = err.message;
    } finally {
      adminActionLoading = false;
    }
  }

  async function handleDeleteAdmin(userId, username) {
    adminMessage = '';
    adminError = '';
    if (!confirm(`Yakin ingin menghapus akun admin "${username}"?`)) return;
    adminActionLoading = true;
    try {
      await apiRequest(`/admin/users/${userId}`, {
        method: 'DELETE'
      });
      adminMessage = `Admin "${username}" berhasil dihapus.`;
      await loadUsers();
    } catch (err) {
      adminError = err.message;
    } finally {
      adminActionLoading = false;
    }
  }

  async function handleChangePassword() {
    adminMessage = '';
    adminError = '';
    if (!currentPassword || !newPassword) {
      adminError = 'Password lama dan baru wajib diisi';
      return;
    }
    if (newPassword !== confirmNewPassword) {
      adminError = 'Konfirmasi password baru tidak cocok';
      return;
    }
    if (newPassword.length < 5) {
      adminError = 'Password baru minimal 5 karakter';
      return;
    }
    adminActionLoading = true;
    try {
      await apiRequest('/admin/users/password', {
        method: 'PUT',
        body: JSON.stringify({
          current_password: currentPassword,
          new_password: newPassword
        })
      });
      adminMessage = 'Kata sandi berhasil diperbarui!';
      currentPassword = '';
      newPassword = '';
      confirmNewPassword = '';
    } catch (err) {
      adminError = err.message;
    } finally {
      adminActionLoading = false;
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

<div class="max-w-6xl mx-auto px-4 sm:px-6 py-6 sm:py-8">
  <!-- Top Bar -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between pb-6 mb-6 border-b border-slate-200 dark:border-slate-800 gap-4">
    <div>
      <span class="text-xs font-bold text-[#0175C2] dark:text-sky-400 uppercase tracking-wider">MobileIndex Admin</span>
      <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white">Manajemen Konten E-Learning</h1>
    </div>

    <div class="flex items-center gap-2 sm:gap-3 flex-wrap">
      <a 
        href="#/" 
        class="inline-flex items-center gap-1.5 text-xs font-semibold px-3 py-2 text-slate-600 dark:text-slate-300 bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl hover:bg-slate-50 dark:hover:bg-slate-900 transition min-h-[40px]"
      >
        <ExternalLink class="w-3.5 h-3.5" />
        <span>Web Publik</span>
      </a>
      <button 
        on:click={handleLogout}
        class="inline-flex items-center gap-1.5 text-xs font-semibold px-3 py-2 text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950/40 border border-red-100 dark:border-red-900 rounded-xl hover:bg-red-100 dark:hover:bg-red-900/60 transition cursor-pointer min-h-[40px]"
      >
        <LogOut class="w-3.5 h-3.5" />
        <span>Logout</span>
      </button>
    </div>
  </div>

  <!-- Tabs Navigation with Horizontal Scroll for Mobile -->
  <div class="flex border-b border-slate-200 dark:border-slate-800 mb-6 gap-3 sm:gap-6 overflow-x-auto no-scrollbar">
    <button 
      on:click={() => activeTab = 'materials'}
      class="inline-flex items-center gap-2 pb-3 text-xs sm:text-sm font-semibold transition border-b-2 whitespace-nowrap min-h-[40px] cursor-pointer {activeTab === 'materials' ? 'border-[#0175C2] text-[#0175C2] dark:text-sky-400' : 'border-transparent text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200'}"
    >
      <BookOpen class="w-4 h-4 flex-shrink-0" />
      <span>Materi ({materials.length})</span>
    </button>
    <button 
      on:click={() => activeTab = 'lessons'}
      class="inline-flex items-center gap-2 pb-3 text-xs sm:text-sm font-semibold transition border-b-2 whitespace-nowrap min-h-[40px] cursor-pointer {activeTab === 'lessons' ? 'border-[#0175C2] text-[#0175C2] dark:text-sky-400' : 'border-transparent text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200'}"
    >
      <Layers class="w-4 h-4 flex-shrink-0" />
      <span>Bab Pembelajaran</span>
    </button>
    <button 
      on:click={() => activeTab = 'admins'}
      class="inline-flex items-center gap-2 pb-3 text-xs sm:text-sm font-semibold transition border-b-2 whitespace-nowrap min-h-[40px] cursor-pointer {activeTab === 'admins' ? 'border-[#0175C2] text-[#0175C2] dark:text-sky-400' : 'border-transparent text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200'}"
    >
      <Users class="w-4 h-4 flex-shrink-0" />
      <span>Admin ({users.length})</span>
    </button>
  </div>

  {#if loading}
    <div class="py-16 text-center text-slate-400">
      <Loader2 class="w-8 h-8 animate-spin text-[#0175C2] mx-auto" />
      <p class="mt-2 text-sm">Memuat data admin...</p>
    </div>
  {:else if error}
    <div class="p-4 bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-900 text-red-700 dark:text-red-300 text-sm rounded-xl mb-6 flex items-center gap-2">
      <AlertCircle class="w-4 h-4 flex-shrink-0" />
      <span>{error}</span>
    </div>
  {:else}
    <!-- Tab 1: Materials Management -->
    {#if activeTab === 'materials'}
      <div class="space-y-4">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4">
          <p class="text-xs text-slate-500 dark:text-slate-400">Kelola modul materi MobileIndex tingkat utama</p>
          <button 
            on:click={openCreateMaterial}
            class="inline-flex items-center justify-center gap-1.5 px-4 py-2.5 bg-[#0175C2] hover:bg-[#02569B] text-white text-xs font-semibold rounded-xl transition shadow-xs min-h-[40px] cursor-pointer"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>Tambah Materi Baru</span>
          </button>
        </div>

        <!-- Desktop Table View -->
        <div class="hidden md:block bg-white dark:bg-slate-950 rounded-xl border border-slate-200 dark:border-slate-800 overflow-x-auto shadow-xs">
          <table class="w-full text-left text-sm min-w-[600px]">
            <thead class="bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800 text-xs font-semibold text-slate-600 dark:text-slate-400">
              <tr>
                <th class="p-3 w-16 text-center">Urutan</th>
                <th class="p-3">Judul Materi</th>
                <th class="p-3">Slug URL</th>
                <th class="p-3 w-28 text-center">Total Bab</th>
                <th class="p-3 w-36 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
              {#each materials as item}
                <tr class="hover:bg-slate-50/50 dark:hover:bg-slate-900/50">
                  <td class="p-3 text-center font-mono text-xs text-slate-500 dark:text-slate-400">{item.order_index}</td>
                  <td class="p-3 font-semibold text-slate-900 dark:text-white">{item.title}</td>
                  <td class="p-3 font-mono text-xs text-slate-500 dark:text-slate-400">/materials/{item.slug}</td>
                  <td class="p-3 text-center text-xs">
                    <span class="px-2 py-0.5 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 font-medium">
                      {item.lessons ? item.lessons.length : 0}
                    </span>
                  </td>
                  <td class="p-3 text-right space-x-2">
                    <button 
                      on:click={() => openCreateLesson(item.id)}
                      class="text-xs text-[#0175C2] dark:text-sky-400 hover:underline font-medium p-1 cursor-pointer"
                    >
                      + Bab
                    </button>
                    <button 
                      on:click={() => openEditMaterial(item)}
                      class="text-xs text-slate-600 dark:text-slate-300 hover:text-slate-900 dark:hover:text-white font-medium p-1 cursor-pointer"
                    >
                      Edit
                    </button>
                    <button 
                      on:click={() => deleteMaterial(item.id)}
                      class="text-xs text-red-600 dark:text-red-400 hover:text-red-800 font-medium p-1 cursor-pointer"
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

        <!-- Mobile Card List View (<768px) -->
        <div class="md:hidden space-y-3">
          {#each materials as item}
            <div class="bg-white dark:bg-slate-950 rounded-xl border border-slate-200 dark:border-slate-800 p-4 space-y-3 shadow-xs">
              <div class="flex items-start justify-between gap-2">
                <div>
                  <span class="text-[10px] font-mono uppercase bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 px-2 py-0.5 rounded">
                    Urutan {item.order_index}
                  </span>
                  <h3 class="text-base font-bold text-slate-900 dark:text-white mt-1">{item.title}</h3>
                  <p class="text-xs font-mono text-slate-400 mt-0.5">/materials/{item.slug}</p>
                </div>
                <span class="text-xs px-2.5 py-1 rounded-full bg-sky-50 dark:bg-sky-950 text-[#0175C2] dark:text-sky-300 font-semibold flex-shrink-0">
                  {item.lessons ? item.lessons.length : 0} Bab
                </span>
              </div>

              {#if item.description}
                <p class="text-xs text-slate-600 dark:text-slate-400 line-clamp-2">{item.description}</p>
              {/if}

              <div class="pt-2 border-t border-slate-100 dark:border-slate-800 flex items-center justify-end gap-2">
                <button 
                  on:click={() => openCreateLesson(item.id)}
                  class="px-3 py-2 text-xs text-[#0175C2] dark:text-sky-400 bg-sky-50 dark:bg-sky-950/60 font-semibold rounded-lg min-h-[40px] cursor-pointer"
                >
                  + Bab
                </button>
                <button 
                  on:click={() => openEditMaterial(item)}
                  class="px-3 py-2 text-xs text-slate-700 dark:text-slate-300 bg-slate-100 dark:bg-slate-800 font-medium rounded-lg min-h-[40px] cursor-pointer"
                >
                  Edit
                </button>
                <button 
                  on:click={() => deleteMaterial(item.id)}
                  class="px-3 py-2 text-xs text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950/40 font-medium rounded-lg min-h-[40px] cursor-pointer"
                >
                  Hapus
                </button>
              </div>
            </div>
          {/each}
          {#if materials.length === 0}
            <div class="p-6 text-center text-slate-400 text-xs bg-white dark:bg-slate-950 rounded-xl border border-slate-200 dark:border-slate-800">
              Belum ada materi. Klik tombol tambah di atas.
            </div>
          {/if}
        </div>
      </div>
    {/if}

    <!-- Tab 2: Lessons Management -->
    {#if activeTab === 'lessons'}
      <div class="space-y-4">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4">
          <p class="text-xs text-slate-500 dark:text-slate-400">Kelola bab pelajaran di dalam setiap materi</p>
          <button 
            on:click={() => openCreateLesson()}
            class="inline-flex items-center justify-center gap-1.5 px-4 py-2.5 bg-[#0175C2] hover:bg-[#02569B] text-white text-xs font-semibold rounded-xl transition shadow-xs min-h-[40px] cursor-pointer"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>Tambah Bab Baru</span>
          </button>
        </div>

        {#each materials as mat}
          <div class="bg-white dark:bg-slate-950 rounded-xl border border-slate-200 dark:border-slate-800 overflow-hidden shadow-xs mb-6">
            <div class="p-3.5 sm:p-4 bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800 flex flex-col sm:flex-row sm:items-center justify-between gap-2">
              <span class="font-bold text-sm text-slate-900 dark:text-white">{mat.title}</span>
              <button 
                on:click={() => openCreateLesson(mat.id)}
                class="text-xs font-semibold text-[#0175C2] dark:text-sky-400 hover:underline self-start sm:self-auto py-1"
              >
                + Tambah Bab ke Materi Ini
              </button>
            </div>

            <!-- Desktop Lessons Table -->
            <div class="hidden md:block overflow-x-auto">
              <table class="w-full text-left text-sm min-w-[550px]">
                <thead class="text-xs font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-100 dark:border-slate-800">
                  <tr>
                    <th class="p-3 w-16 text-center">Urutan</th>
                    <th class="p-3">Judul Bab</th>
                    <th class="p-3">Slug</th>
                    <th class="p-3 w-28 text-center">DartPad</th>
                    <th class="p-3 w-32 text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
                  {#if mat.lessons && mat.lessons.length > 0}
                    {#each mat.lessons as lesson}
                      <tr class="hover:bg-slate-50/50 dark:hover:bg-slate-900/50">
                        <td class="p-3 text-center font-mono text-xs text-slate-400">{lesson.order_index}</td>
                        <td class="p-3 font-medium text-slate-800 dark:text-slate-200">{lesson.title}</td>
                        <td class="p-3 font-mono text-xs text-slate-500 dark:text-slate-400">{lesson.slug}</td>
                        <td class="p-3 text-center">
                          {#if lesson.enable_dartpad}
                            <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-emerald-50 dark:bg-emerald-950/60 text-emerald-700 dark:text-emerald-300">Aktif</span>
                          {:else}
                            <span class="text-[10px] font-medium px-2 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400">Nonaktif</span>
                          {/if}
                        </td>
                        <td class="p-3 text-right space-x-2">
                          <button 
                            on:click={() => openEditLesson(lesson)}
                            class="text-xs text-slate-600 dark:text-slate-300 hover:text-slate-900 dark:hover:text-white font-medium p-1 cursor-pointer"
                          >
                            Edit
                          </button>
                          <button 
                            on:click={() => deleteLesson(lesson.id)}
                            class="text-xs text-red-600 dark:text-red-400 hover:text-red-800 font-medium p-1 cursor-pointer"
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

            <!-- Mobile Lessons Card List (<768px) -->
            <div class="md:hidden divide-y divide-slate-100 dark:divide-slate-800">
              {#if mat.lessons && mat.lessons.length > 0}
                {#each mat.lessons as lesson}
                  <div class="p-3.5 space-y-2">
                    <div class="flex items-start justify-between gap-2">
                      <div class="flex items-center gap-2">
                        <span class="text-xs font-mono px-2 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300">
                          #{lesson.order_index}
                        </span>
                        <h4 class="font-semibold text-slate-900 dark:text-white text-sm">{lesson.title}</h4>
                      </div>
                      {#if lesson.enable_dartpad}
                        <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-emerald-50 dark:bg-emerald-950/60 text-emerald-700 dark:text-emerald-300 flex-shrink-0">
                          DartPad
                        </span>
                      {/if}
                    </div>

                    <div class="flex items-center justify-between text-xs text-slate-400 pt-1">
                      <span class="font-mono text-[11px] truncate">{lesson.slug}</span>
                      <div class="flex items-center gap-2">
                        <button 
                          on:click={() => openEditLesson(lesson)}
                          class="px-2.5 py-1.5 text-xs text-slate-700 dark:text-slate-300 bg-slate-100 dark:bg-slate-800 font-medium rounded-lg min-h-[36px] cursor-pointer"
                        >
                          Edit
                        </button>
                        <button 
                          on:click={() => deleteLesson(lesson.id)}
                          class="px-2.5 py-1.5 text-xs text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950/40 font-medium rounded-lg min-h-[36px] cursor-pointer"
                        >
                          Hapus
                        </button>
                      </div>
                    </div>
                  </div>
                {/each}
              {:else}
                <div class="p-4 text-center text-slate-400 text-xs">Belum ada bab untuk materi ini.</div>
              {/if}
            </div>

          </div>
        {/each}
      </div>
    {/if}

    <!-- Tab 3: Admin Management -->
    {#if activeTab === 'admins'}
      <div class="space-y-6">
        {#if adminMessage}
          <div class="p-3.5 bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-800 text-emerald-800 dark:text-emerald-300 text-xs rounded-xl flex items-center justify-between">
            <span class="flex items-center gap-2">
              <ShieldCheck class="w-4 h-4 text-emerald-600 flex-shrink-0" />
              <span>{adminMessage}</span>
            </span>
            <button on:click={() => adminMessage = ''} class="text-emerald-600 hover:text-emerald-800 p-1"><X class="w-3.5 h-3.5" /></button>
          </div>
        {/if}

        {#if adminError}
          <div class="p-3.5 bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-900 text-red-700 dark:text-red-300 text-xs rounded-xl flex items-center justify-between">
            <span class="flex items-center gap-2">
              <AlertCircle class="w-4 h-4 text-red-600 flex-shrink-0" />
              <span>{adminError}</span>
            </span>
            <button on:click={() => adminError = ''} class="text-red-600 hover:text-red-800 p-1"><X class="w-3.5 h-3.5" /></button>
          </div>
        {/if}

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Card: Tambah Admin Baru -->
          <div class="bg-white dark:bg-slate-950 rounded-2xl border border-slate-200 dark:border-slate-800 p-5 shadow-xs">
            <div class="flex items-center gap-2 mb-4 pb-3 border-b border-slate-100 dark:border-slate-800">
              <div class="p-2 bg-sky-50 dark:bg-sky-950/50 rounded-lg text-[#0175C2]">
                <Users class="w-4 h-4" />
              </div>
              <div>
                <h3 class="text-sm font-bold text-slate-900 dark:text-white">Tambah Akun Admin Baru</h3>
                <p class="text-[11px] text-slate-500 dark:text-slate-400">Buat akun untuk pengelola materi lainnya</p>
              </div>
            </div>

            <form on:submit|preventDefault={handleCreateAdmin} class="space-y-3">
              <div>
                <label for="new-admin-user" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Username</label>
                <input 
                  id="new-admin-user"
                  type="text" 
                  bind:value={newAdminUsername}
                  placeholder="Contoh: amiya"
                  required
                  class="w-full px-3 py-2.5 border border-slate-300 dark:border-slate-700 rounded-xl text-xs bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-hidden focus:border-[#0175C2]"
                />
              </div>

              <div>
                <label for="new-admin-pass" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Password</label>
                <input 
                  id="new-admin-pass"
                  type="password" 
                  bind:value={newAdminPassword}
                  placeholder="Minimal 5 karakter"
                  required
                  class="w-full px-3 py-2.5 border border-slate-300 dark:border-slate-700 rounded-xl text-xs bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-hidden focus:border-[#0175C2]"
                />
              </div>

              <button 
                type="submit"
                disabled={adminActionLoading}
                class="w-full py-2.5 bg-[#0175C2] hover:bg-[#02569B] text-white text-xs font-semibold rounded-xl transition disabled:opacity-50 cursor-pointer min-h-[44px]"
              >
                {adminActionLoading ? 'Menyimpan...' : 'Tambah Admin'}
              </button>
            </form>
          </div>

          <!-- Card: Ganti Password Saya -->
          <div class="bg-white dark:bg-slate-950 rounded-2xl border border-slate-200 dark:border-slate-800 p-5 shadow-xs">
            <div class="flex items-center gap-2 mb-4 pb-3 border-b border-slate-100 dark:border-slate-800">
              <div class="p-2 bg-amber-50 dark:bg-amber-950/50 rounded-lg text-amber-600">
                <KeyRound class="w-4 h-4" />
              </div>
              <div>
                <h3 class="text-sm font-bold text-slate-900 dark:text-white">Ganti Password Saya</h3>
                <p class="text-[11px] text-slate-500 dark:text-slate-400">
                  Akun aktif: <strong class="text-slate-800 dark:text-slate-200">{currentUser?.username || 'Saya'}</strong>
                </p>
              </div>
            </div>

            <form on:submit|preventDefault={handleChangePassword} class="space-y-3">
              <div>
                <label for="curr-pass" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Password Saat Ini</label>
                <input 
                  id="curr-pass"
                  type="password" 
                  bind:value={currentPassword}
                  placeholder="Masukkan kata sandi lama"
                  required
                  class="w-full px-3 py-2.5 border border-slate-300 dark:border-slate-700 rounded-xl text-xs bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-hidden focus:border-[#0175C2]"
                />
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div>
                  <label for="new-pass" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Password Baru</label>
                  <input 
                    id="new-pass"
                    type="password" 
                    bind:value={newPassword}
                    placeholder="Minimal 5 karakter"
                    required
                    class="w-full px-3 py-2.5 border border-slate-300 dark:border-slate-700 rounded-xl text-xs bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-hidden focus:border-[#0175C2]"
                  />
                </div>
                <div>
                  <label for="confirm-pass" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Ulangi Password Baru</label>
                  <input 
                    id="confirm-pass"
                    type="password" 
                    bind:value={confirmNewPassword}
                    placeholder="Ulangi kata sandi baru"
                    required
                    class="w-full px-3 py-2.5 border border-slate-300 dark:border-slate-700 rounded-xl text-xs bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-hidden focus:border-[#0175C2]"
                  />
                </div>
              </div>

              <button 
                type="submit"
                disabled={adminActionLoading}
                class="w-full py-2.5 bg-slate-900 hover:bg-black dark:bg-slate-800 dark:hover:bg-slate-700 text-white text-xs font-semibold rounded-xl transition disabled:opacity-50 cursor-pointer min-h-[44px]"
              >
                {adminActionLoading ? 'Memproses...' : 'Simpan Password Baru'}
              </button>
            </form>
          </div>
        </div>

        <!-- Tabel Daftar Semua Akun Admin (Responsive wrapper) -->
        <div class="bg-white dark:bg-slate-950 rounded-2xl border border-slate-200 dark:border-slate-800 overflow-hidden shadow-xs">
          <div class="p-4 bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800 flex justify-between items-center">
            <span class="font-bold text-xs text-slate-700 dark:text-slate-300">Daftar Admin Terdaftar</span>
            <span class="text-xs text-slate-400">Total: {users.length} akun</span>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full text-left text-sm min-w-[400px]">
              <thead class="text-xs font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-100 dark:border-slate-800">
                <tr>
                  <th class="p-3 w-16 text-center">ID</th>
                  <th class="p-3">Username</th>
                  <th class="p-3">Status</th>
                  <th class="p-3 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
                {#each users as u}
                  <tr class="hover:bg-slate-50/50 dark:hover:bg-slate-900/50">
                    <td class="p-3 text-center font-mono text-xs text-slate-400">{u.id}</td>
                    <td class="p-3 font-semibold text-slate-800 dark:text-slate-200 flex items-center gap-2 flex-wrap">
                      <span>{u.username}</span>
                      {#if currentUser && currentUser.id === u.id}
                        <span class="text-[10px] bg-sky-100 dark:bg-sky-950/70 text-sky-700 dark:text-sky-300 font-semibold px-2 py-0.5 rounded-full">
                          Anda (Sedang Login)
                        </span>
                      {/if}
                    </td>
                    <td class="p-3 text-xs text-slate-500 dark:text-slate-400">
                      <span class="inline-flex items-center gap-1 text-emerald-600 dark:text-emerald-400">
                        <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                        <span>Aktif</span>
                      </span>
                    </td>
                    <td class="p-3 text-right">
                      {#if currentUser && currentUser.id === u.id}
                        <span class="text-xs text-slate-400 italic">Akun aktif</span>
                      {:else}
                        <button 
                          on:click={() => handleDeleteAdmin(u.id, u.username)}
                          class="text-xs text-red-600 dark:text-red-400 hover:text-red-800 font-medium cursor-pointer p-1.5"
                        >
                          Hapus Admin
                        </button>
                      {/if}
                    </td>
                  </tr>
                {/each}
                {#if users.length === 0}
                  <tr>
                    <td colspan="4" class="p-4 text-center text-slate-400 text-xs">Belum ada akun admin lain.</td>
                  </tr>
                {/if}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    {/if}
  {/if}
</div>

<!-- MODAL: Form Material (Responsive for mobile) -->
{#if showMaterialModal}
  <div class="fixed inset-0 bg-slate-900/50 dark:bg-black/70 backdrop-blur-xs flex items-center justify-center p-3 sm:p-4 z-50 overflow-y-auto">
    <div class="bg-white dark:bg-slate-950 rounded-2xl max-w-lg w-full p-5 sm:p-6 shadow-xl border border-slate-200 dark:border-slate-800 max-h-[90vh] overflow-y-auto">
      <h3 class="text-lg font-bold text-slate-900 dark:text-white mb-4">
        {materialForm.id ? 'Edit Materi' : 'Tambah Materi Baru'}
      </h3>

      <div class="space-y-4">
        <div>
          <label for="mat-title" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Judul Materi</label>
          <input 
            id="mat-title"
            type="text" 
            bind:value={materialForm.title} 
            on:input={() => { if (!materialForm.id) materialForm.slug = autoSlug(materialForm.title); }}
            placeholder="Contoh: Pengenalan Dasar Dart" 
            class="w-full px-3 py-2.5 border border-slate-300 dark:border-slate-700 rounded-xl text-sm bg-white dark:bg-slate-900 text-slate-900 dark:text-white"
          />
        </div>

        <div>
          <label for="mat-slug" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Slug URL</label>
          <input 
            id="mat-slug"
            type="text" 
            bind:value={materialForm.slug} 
            placeholder="pengenalan-dasar-dart" 
            class="w-full px-3 py-2.5 border border-slate-300 dark:border-slate-700 rounded-xl text-sm font-mono text-xs bg-white dark:bg-slate-900 text-slate-900 dark:text-white"
          />
        </div>

        <div>
          <label for="mat-desc" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Deskripsi Singkat</label>
          <textarea 
            id="mat-desc"
            bind:value={materialForm.description} 
            rows="3" 
            placeholder="Ringkasan tentang materi ini..." 
            class="w-full px-3 py-2.5 border border-slate-300 dark:border-slate-700 rounded-xl text-sm bg-white dark:bg-slate-900 text-slate-900 dark:text-white"
          ></textarea>
        </div>

        <div>
          <label for="mat-order" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Nomor Urutan Tampil</label>
          <input 
            id="mat-order"
            type="number" 
            bind:value={materialForm.order_index} 
            class="w-full px-3 py-2.5 border border-slate-300 dark:border-slate-700 rounded-xl text-sm bg-white dark:bg-slate-900 text-slate-900 dark:text-white"
          />
        </div>
      </div>

      <div class="mt-6 flex justify-end gap-2">
        <button 
          on:click={() => showMaterialModal = false}
          class="px-4 py-2.5 text-xs font-semibold text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-xl transition min-h-[40px] cursor-pointer"
        >
          Batal
        </button>
        <button 
          on:click={saveMaterial}
          class="px-4 py-2.5 text-xs font-semibold text-white bg-[#0175C2] hover:bg-[#02569B] rounded-xl transition min-h-[40px] cursor-pointer"
        >
          Simpan Materi
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- MODAL: Form Lesson (Dynamic Vertical Block Builder, fully mobile scrollable & responsive) -->
{#if showLessonModal}
  <div class="fixed inset-0 bg-slate-900/50 dark:bg-black/70 backdrop-blur-xs flex items-center justify-center p-2 sm:p-4 z-50 overflow-y-auto">
    <div class="bg-white dark:bg-slate-950 rounded-2xl max-w-4xl w-full p-4 sm:p-6 shadow-xl border border-slate-200 dark:border-slate-800 my-4 sm:my-8 flex flex-col max-h-[92vh]">
      
      <!-- Modal Header -->
      <div class="flex items-center justify-between pb-3 sm:pb-4 border-b border-slate-200 dark:border-slate-800 flex-shrink-0">
        <div>
          <h3 class="text-base sm:text-lg font-bold text-slate-900 dark:text-white">
            {lessonForm.id ? 'Edit Bab Pelajaran' : 'Tambah Bab Baru'}
          </h3>
          <p class="text-[11px] sm:text-xs text-slate-500 dark:text-slate-400">Susun konten bab secara vertikal & modular</p>
        </div>
        <button on:click={() => showLessonModal = false} class="p-2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 cursor-pointer">
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Modal Body (Scrollable) -->
      <div class="space-y-4 sm:space-y-6 overflow-y-auto py-3 sm:py-4 pr-1 flex-1">
        
        <!-- Metadata Header -->
        <div class="bg-slate-50/80 dark:bg-slate-900/60 p-3.5 sm:p-4 rounded-xl border border-slate-200 dark:border-slate-800 space-y-3 sm:space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
            <div>
              <label for="les-mat" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Pilih Materi Induk</label>
              <select 
                id="les-mat"
                bind:value={lessonForm.material_id} 
                class="w-full px-3 py-2 border border-slate-300 dark:border-slate-700 rounded-lg text-sm bg-white dark:bg-slate-900 text-slate-900 dark:text-white"
              >
                {#each materials as m}
                  <option value={m.id}>{m.title}</option>
                {/each}
              </select>
            </div>

            <div>
              <label for="les-order" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Nomor Urutan Bab</label>
              <input 
                id="les-order"
                type="number" 
                bind:value={lessonForm.order_index} 
                class="w-full px-3 py-2 border border-slate-300 dark:border-slate-700 rounded-lg text-sm bg-white dark:bg-slate-900 text-slate-900 dark:text-white"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
            <div>
              <label for="les-title" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Judul Bab</label>
              <input 
                id="les-title"
                type="text" 
                bind:value={lessonForm.title} 
                on:input={() => { if (!lessonForm.id) lessonForm.slug = autoSlug(lessonForm.title); }}
                placeholder="Contoh: 1. Pengantar Dart & Struktur main()" 
                class="w-full px-3 py-2 border border-slate-300 dark:border-slate-700 rounded-lg text-sm bg-white dark:bg-slate-900 text-slate-900 dark:text-white"
              />
            </div>

            <div>
              <label for="les-slug" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Slug URL Bab</label>
              <input 
                id="les-slug"
                type="text" 
                bind:value={lessonForm.slug} 
                placeholder="pengantar-dart-dan-main" 
                class="w-full px-3 py-2 border border-slate-300 dark:border-slate-700 rounded-lg text-sm font-mono text-xs bg-white dark:bg-slate-900 text-slate-900 dark:text-white"
              />
            </div>
          </div>
        </div>

        <!-- Dynamic Vertical Blocks Builder -->
        <div class="space-y-4">
          <div class="flex items-center justify-between flex-wrap gap-2">
            <div>
              <h4 class="text-sm font-bold text-slate-900 dark:text-white">Blok Konten Bab ({blocks.length})</h4>
              <p class="text-[11px] sm:text-xs text-slate-500 dark:text-slate-400">Atur urutan blok vertikal untuk menyusun materi</p>
            </div>

            <!-- Quick Add Bar -->
            <div class="flex items-center gap-1.5 flex-wrap">
              <span class="text-xs text-slate-400 font-medium mr-1">+ Tambah:</span>
              <button 
                type="button" 
                on:click={() => addBlock('text')}
                class="inline-flex items-center gap-1 px-2.5 py-1 text-xs font-semibold rounded-lg bg-slate-100 dark:bg-slate-800 hover:bg-blue-50 dark:hover:bg-sky-950/60 hover:text-[#0175C2] dark:hover:text-sky-300 text-slate-700 dark:text-slate-300 transition border border-slate-200 dark:border-slate-700 min-h-[36px]"
              >
                <Type class="w-3.5 h-3.5" />
                <span>Text</span>
              </button>
              <button 
                type="button" 
                on:click={() => addBlock('image')}
                class="inline-flex items-center gap-1 px-2.5 py-1 text-xs font-semibold rounded-lg bg-slate-100 dark:bg-slate-800 hover:bg-emerald-50 dark:hover:bg-emerald-950/60 hover:text-emerald-700 dark:hover:text-emerald-300 text-slate-700 dark:text-slate-300 transition border border-slate-200 dark:border-slate-700 min-h-[36px]"
              >
                <Image class="w-3.5 h-3.5" />
                <span>Image</span>
              </button>
              <button 
                type="button" 
                on:click={() => addBlock('note')}
                class="inline-flex items-center gap-1 px-2.5 py-1 text-xs font-semibold rounded-lg bg-slate-100 dark:bg-slate-800 hover:bg-amber-50 dark:hover:bg-amber-950/60 hover:text-amber-700 dark:hover:text-amber-300 text-slate-700 dark:text-slate-300 transition border border-slate-200 dark:border-slate-700 min-h-[36px]"
              >
                <Lightbulb class="w-3.5 h-3.5" />
                <span>Note</span>
              </button>
              <button 
                type="button" 
                on:click={() => addBlock('code')}
                class="inline-flex items-center gap-1 px-2.5 py-1 text-xs font-semibold rounded-lg bg-slate-100 dark:bg-slate-800 hover:bg-purple-50 dark:hover:bg-purple-950/60 hover:text-purple-700 dark:hover:text-purple-300 text-slate-700 dark:text-slate-300 transition border border-slate-200 dark:border-slate-700 min-h-[36px]"
              >
                <Code class="w-3.5 h-3.5" />
                <span>Code</span>
              </button>
            </div>
          </div>

          <!-- Blocks List -->
          {#if blocks.length === 0}
            <div class="text-center py-8 bg-slate-50 dark:bg-slate-900 border-2 border-dashed border-slate-200 dark:border-slate-800 rounded-xl">
              <p class="text-xs text-slate-400">Belum ada blok konten. Gunakan tombol di atas untuk menambahkan blok pertama.</p>
            </div>
          {:else}
            <div class="space-y-4">
              {#each blocks as block, index (block.id || index)}
                <div class="p-3.5 sm:p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 shadow-xs hover:border-slate-300 dark:hover:border-slate-700 transition-all">
                  
                  <!-- Block Card Top Header -->
                  <div class="flex items-center justify-between pb-2.5 mb-3 border-b border-slate-100 dark:border-slate-800">
                    <div class="flex items-center gap-2">
                      <span class="text-xs font-mono font-bold px-2 py-0.5 rounded 
                        {block.type === 'text' ? 'bg-blue-100 dark:bg-sky-950 text-[#0175C2] dark:text-sky-300' : ''}
                        {block.type === 'image' ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-800 dark:text-emerald-300' : ''}
                        {block.type === 'note' ? 'bg-amber-100 dark:bg-amber-950 text-amber-800 dark:text-amber-300' : ''}
                        {block.type === 'code' ? 'bg-purple-100 dark:bg-purple-950 text-purple-800 dark:text-purple-300' : ''}"
                      >
                        {index + 1}. {block.type.toUpperCase()}
                      </span>
                    </div>

                    <!-- Reorder & Delete Actions -->
                    <div class="flex items-center gap-1">
                      <button
                        type="button"
                        on:click={() => moveBlockUp(index)}
                        disabled={index === 0}
                        class="p-1.5 text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-100 disabled:opacity-30 rounded hover:bg-slate-100 dark:hover:bg-slate-800 transition cursor-pointer min-h-[32px] min-w-[32px] flex items-center justify-center"
                        title="Geser ke Atas"
                      >
                        <ArrowUp class="w-4 h-4" />
                      </button>
                      <button
                        type="button"
                        on:click={() => moveBlockDown(index)}
                        disabled={index === blocks.length - 1}
                        class="p-1.5 text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-100 disabled:opacity-30 rounded hover:bg-slate-100 dark:hover:bg-slate-800 transition cursor-pointer min-h-[32px] min-w-[32px] flex items-center justify-center"
                        title="Geser ke Bawah"
                      >
                        <ArrowDown class="w-4 h-4" />
                      </button>
                      <button
                        type="button"
                        on:click={() => removeBlock(index)}
                        class="p-1.5 ml-1 text-red-500 hover:text-red-700 rounded hover:bg-red-50 dark:hover:bg-red-950/40 transition cursor-pointer min-h-[32px] min-w-[32px] flex items-center justify-center"
                        title="Hapus Blok Ini"
                      >
                        <Trash2 class="w-4 h-4" />
                      </button>
                    </div>
                  </div>

                  <!-- Block Content Based on Type -->
                  {#if block.type === 'text'}
                    <div>
                      <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Teks / Markdown</label>
                      <textarea
                        bind:value={block.content}
                        rows="4"
                        placeholder="Tuliskan penjelasan bab, poin-poin materi, atau teks Markdown di sini..."
                        class="w-full px-3 py-2 border border-slate-300 dark:border-slate-700 rounded-lg text-xs font-mono bg-white dark:bg-slate-950 text-slate-900 dark:text-white"
                      ></textarea>
                    </div>

                  {:else if block.type === 'image'}
                    <div class="space-y-3">
                      <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">URL Gambar atau Upload File</label>
                        <div class="flex flex-col sm:flex-row gap-2">
                          <input
                            type="text"
                            bind:value={block.url}
                            placeholder="https://... atau upload file lokal"
                            class="flex-1 px-3 py-2 border border-slate-300 dark:border-slate-700 rounded-lg text-xs bg-white dark:bg-slate-950 text-slate-900 dark:text-white"
                          />
                          <label class="cursor-pointer px-3 py-2 bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 border border-slate-300 dark:border-slate-700 rounded-lg text-xs font-medium text-slate-700 dark:text-slate-300 flex items-center justify-center gap-1.5 flex-shrink-0 min-h-[36px]">
                            {#if uploadingBlockId === block.id}
                              <Loader2 class="w-3.5 h-3.5 animate-spin" />
                              <span>Mengunggah...</span>
                            {:else}
                              <Upload class="w-3.5 h-3.5" />
                              <span>Upload File</span>
                            {/if}
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
                        <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Keterangan / Caption Gambar (Opsional)</label>
                        <input
                          type="text"
                          bind:value={block.caption}
                          placeholder="Contoh: Diagram Alur Program Dart"
                          class="w-full px-3 py-2 border border-slate-300 dark:border-slate-700 rounded-lg text-xs bg-white dark:bg-slate-950 text-slate-900 dark:text-white"
                        />
                      </div>

                      {#if block.url}
                        <div class="p-2 bg-slate-50 dark:bg-slate-950 border dark:border-slate-800 rounded-lg flex items-center gap-3">
                          <img src={block.url} alt={block.caption || 'Preview'} class="h-14 w-20 object-cover rounded border dark:border-slate-700 flex-shrink-0" />
                          <span class="text-xs text-slate-500 dark:text-slate-400 truncate">{block.caption || block.url}</span>
                        </div>
                      {/if}
                    </div>

                  {:else if block.type === 'note'}
                    <div>
                      <label class="block text-xs font-semibold text-amber-800 dark:text-amber-400 mb-1">Catatan Penting / Tip Callout</label>
                      <textarea
                        bind:value={block.content}
                        rows="2"
                        placeholder="💡 Contoh: Selalu gunakan widget Const jika properti tidak pernah berubah..."
                        class="w-full px-3 py-2 border border-amber-300 dark:border-amber-800 bg-amber-50/40 dark:bg-amber-950/20 rounded-lg text-xs text-amber-950 dark:text-amber-200"
                      ></textarea>
                    </div>

                  {:else if block.type === 'code'}
                    <div class="space-y-2">
                      <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Judul / Label Cuplikan Kode</label>
                        <input
                          type="text"
                          bind:value={block.title}
                          placeholder="Contoh: Contoh 1: Fungsi Named Parameter"
                          class="w-full px-3 py-2 border border-slate-300 dark:border-slate-700 rounded-lg text-xs bg-white dark:bg-slate-950 text-slate-900 dark:text-white font-medium"
                        />
                      </div>
                      <div>
                        <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Kode Dart</label>
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
        <div class="p-3.5 sm:p-4 bg-sky-50/60 dark:bg-sky-950/40 border border-sky-100 dark:border-sky-900/40 rounded-xl space-y-3">
          <div class="flex items-center justify-between gap-2">
            <div>
              <span class="text-xs font-bold text-[#0175C2] dark:text-sky-300">Area Simulasi Dart (Interactive DartPad)</span>
              <p class="text-[11px] text-slate-500 dark:text-slate-400">Aktifkan tombol toggle panel DartPad interaktif</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer flex-shrink-0">
              <input type="checkbox" bind:checked={lessonForm.enable_dartpad} class="sr-only peer" />
              <div class="w-11 h-6 bg-slate-200 dark:bg-slate-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[#0175C2]"></div>
            </label>
          </div>

          {#if lessonForm.enable_dartpad}
            <div>
              <label for="starter-code" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Default Starter Code (Jika Ingin Disediakan)</label>
              <textarea 
                id="starter-code"
                bind:value={lessonForm.starter_code} 
                rows="3" 
                placeholder="void main() &#123; print('Hello Dart'); &#125;" 
                class="w-full px-3 py-2 border border-slate-300 dark:border-slate-700 rounded-lg text-xs font-mono bg-slate-900 text-slate-100"
              ></textarea>
            </div>
          {/if}
        </div>

      </div>

      <!-- Modal Footer -->
      <div class="flex justify-end gap-2 pt-3 sm:pt-4 border-t border-slate-200 dark:border-slate-800 flex-shrink-0">
        <button 
          type="button"
          on:click={() => showLessonModal = false}
          class="px-4 py-2 text-xs font-semibold text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-xl transition min-h-[40px] cursor-pointer"
        >
          Batal
        </button>
        <button 
          type="button"
          on:click={saveLesson}
          class="px-4 py-2 text-xs font-semibold text-white bg-[#0175C2] hover:bg-[#02569B] rounded-xl transition min-h-[40px] cursor-pointer"
        >
          Simpan Seluruh Bab & Blok
        </button>
      </div>

    </div>
  </div>
{/if}
