<script>
  import { onMount } from 'svelte';
  import { apiRequest } from '../lib/api.js';
  import { marked } from 'marked';

  export let slug = '';

  let material = null;
  let activeLesson = null;
  let loading = true;
  let error = '';
  let sidebarOpen = true;
  let activeMobileTab = 'material'; // 'material' or 'dartpad'
  let copied = false;

  onMount(async () => {
    await fetchMaterial();
  });

  async function fetchMaterial() {
    loading = true;
    error = '';
    try {
      material = await apiRequest(`/materials/${slug}`);
      if (material && material.lessons && material.lessons.length > 0) {
        // default select first lesson
        activeLesson = material.lessons[0];
      }
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  function selectLesson(lesson) {
    activeLesson = lesson;
    copied = false;
    activeMobileTab = 'material';
  }

  function copyCode() {
    if (!activeLesson || !activeLesson.starter_code) return;
    navigator.clipboard.writeText(activeLesson.starter_code);
    copied = true;
    setTimeout(() => {
      copied = false;
    }, 2000);
  }

  function toggleSidebar() {
    sidebarOpen = !sidebarOpen;
  }
</script>

{#if loading}
  <div class="flex justify-center items-center h-[calc(100vh-64px)] text-slate-500">
    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-[#0175C2]"></div>
    <span class="ml-3 font-medium">Memuat materi bab...</span>
  </div>
{:else if error}
  <div class="max-w-3xl mx-auto px-4 py-16 text-center">
    <div class="p-6 bg-red-50 border border-red-200 rounded-xl text-red-700">
      <h3 class="font-bold text-lg mb-2">Gagal Memuat Materi</h3>
      <p class="text-sm">{error}</p>
      <a href="#/" class="inline-block mt-4 text-xs font-semibold px-4 py-2 bg-[#0175C2] text-white rounded-md hover:bg-[#02569B] transition">
        &larr; Kembali ke Beranda
      </a>
    </div>
  </div>
{:else if material}
  <div class="flex h-[calc(100vh-64px)] overflow-hidden bg-slate-50">
    
    <!-- Collapsible Sidebar -->
    <aside 
      class="{sidebarOpen ? 'w-72 sm:w-80' : 'w-0'} flex-shrink-0 bg-white border-r border-slate-200 transition-all duration-300 flex flex-col overflow-hidden z-20 shadow-sm"
    >
      <div class="p-4 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
        <div class="overflow-hidden pr-2">
          <span class="text-xs font-bold uppercase tracking-wider text-[#0175C2]">Daftar Bab</span>
          <h2 class="font-bold text-slate-900 text-sm truncate" title={material.title}>{material.title}</h2>
        </div>
        <button 
          on:click={toggleSidebar} 
          class="p-1 rounded text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition" 
          title="Tutup Sidebar"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7"></path></svg>
        </button>
      </div>

      <div class="flex-1 overflow-y-auto p-3 space-y-1.5">
        {#if material.lessons && material.lessons.length > 0}
          {#each material.lessons as lesson, idx}
            <button
              on:click={() => selectLesson(lesson)}
              class="w-full text-left px-3 py-2.5 rounded-lg text-sm flex items-start gap-2.5 transition-all {activeLesson && activeLesson.id === lesson.id ? 'bg-[#E1F5FE] text-[#0175C2] font-semibold shadow-xs' : 'text-slate-600 hover:bg-slate-100'}"
            >
              <span class="text-xs px-1.5 py-0.5 rounded {activeLesson && activeLesson.id === lesson.id ? 'bg-[#0175C2] text-white' : 'bg-slate-200 text-slate-600'} font-mono mt-0.5">
                {idx + 1}
              </span>
              <span class="flex-1 line-clamp-2 leading-snug">{lesson.title}</span>
              {#if lesson.enable_dartpad}
                <span class="text-[10px] bg-sky-100 text-[#0175C2] font-bold px-1.5 py-0.5 rounded uppercase self-center">
                  DartPad
                </span>
              {/if}
            </button>
          {/each}
        {:else}
          <p class="text-xs text-slate-400 p-4 text-center">Belum ada bab yang ditambahkan.</p>
        {/if}
      </div>

      <div class="p-3 border-t border-slate-100 bg-slate-50/50">
        <a href="#/" class="text-xs text-slate-500 hover:text-[#0175C2] flex items-center gap-1.5 py-1">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
          Kembali ke Daftar Materi
        </a>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-1 flex flex-col min-w-0 overflow-hidden bg-white">
      
      <!-- Top Action Bar -->
      <div class="h-12 border-b border-slate-100 px-4 flex items-center justify-between bg-white flex-shrink-0">
        <div class="flex items-center gap-2">
          {#if !sidebarOpen}
            <button 
              on:click={toggleSidebar} 
              class="p-1.5 rounded-md text-slate-600 hover:bg-slate-100 border border-slate-200 transition text-xs flex items-center gap-1.5"
              title="Buka Daftar Bab"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7"></path></svg>
              <span>Daftar Bab</span>
            </button>
          {/if}
          <span class="text-xs text-slate-400 font-medium truncate">
            {material.title} &nbsp;/&nbsp; <span class="text-slate-700">{activeLesson ? activeLesson.title : 'Pilih Bab'}</span>
          </span>
        </div>

        <!-- Mobile Toggle Tabs (if DartPad enabled) -->
        {#if activeLesson && activeLesson.enable_dartpad}
          <div class="flex md:hidden bg-slate-100 p-0.5 rounded-lg border border-slate-200 text-xs">
            <button 
              on:click={() => activeMobileTab = 'material'} 
              class="px-2.5 py-1 rounded-md font-medium transition {activeMobileTab === 'material' ? 'bg-white text-[#0175C2] shadow-xs' : 'text-slate-600'}"
            >
              Materi
            </button>
            <button 
              on:click={() => activeMobileTab = 'dartpad'} 
              class="px-2.5 py-1 rounded-md font-medium transition {activeMobileTab === 'dartpad' ? 'bg-[#0175C2] text-white shadow-xs' : 'text-slate-600'}"
            >
              DartPad
            </button>
          </div>
        {/if}
      </div>

      <!-- Reading and Practice Area -->
      <div class="flex-1 overflow-hidden relative">
        {#if !activeLesson}
          <div class="flex items-center justify-center h-full text-slate-400 text-sm">
            Pilih salah satu bab dari panel kiri untuk mulai membaca.
          </div>
        {:else if !activeLesson.enable_dartpad}
          <!-- Centered Optimal Reading Width for non-Dartpad lessons -->
          <div class="h-full overflow-y-auto px-4 py-8">
            <div class="max-w-3xl mx-auto space-y-6">
              
              <h1 class="text-2xl sm:text-3xl font-extrabold text-slate-900 tracking-tight">
                {activeLesson.title}
              </h1>

              <!-- Optional Image -->
              {#if activeLesson.image_url}
                <div class="rounded-xl overflow-hidden border border-slate-200 bg-slate-100">
                  <img src={activeLesson.image_url} alt={activeLesson.title} class="w-full max-h-96 object-cover" />
                </div>
              {/if}

              <!-- Content Area (Markdown) -->
              <div class="prose prose-slate max-w-none text-slate-800 leading-relaxed text-base space-y-4">
                {@html marked.parse(activeLesson.content || '')}
              </div>

              <!-- Notes / Callout -->
              {#if activeLesson.notes}
                <div class="p-4 rounded-xl bg-blue-50/70 border border-blue-200 text-blue-900 text-sm flex gap-3">
                  <div class="text-[#0175C2] font-bold text-lg leading-none">ℹ</div>
                  <div class="leading-relaxed whitespace-pre-line">{activeLesson.notes}</div>
                </div>
              {/if}

            </div>
          </div>
        {:else}
          <!-- Split-Screen 2 Kolom Layout for Dartpad lessons -->
          <div class="h-full flex flex-col md:flex-row overflow-hidden">
            
            <!-- Left Column: Reading Content -->
            <div class="flex-1 overflow-y-auto p-4 sm:p-8 space-y-6 md:border-r border-slate-200 {activeMobileTab === 'dartpad' ? 'hidden md:block' : 'block'}">
              <h1 class="text-2xl sm:text-3xl font-extrabold text-slate-900 tracking-tight">
                {activeLesson.title}
              </h1>

              <!-- Optional Image -->
              {#if activeLesson.image_url}
                <div class="rounded-xl overflow-hidden border border-slate-200 bg-slate-100">
                  <img src={activeLesson.image_url} alt={activeLesson.title} class="w-full max-h-80 object-cover" />
                </div>
              {/if}

              <!-- Content Area (Markdown) -->
              <div class="prose prose-slate max-w-none text-slate-800 leading-relaxed text-sm sm:text-base space-y-4">
                {@html marked.parse(activeLesson.content || '')}
              </div>

              <!-- Notes / Callout -->
              {#if activeLesson.notes}
                <div class="p-4 rounded-xl bg-sky-50 border border-sky-200 text-sky-950 text-xs sm:text-sm flex gap-3">
                  <div class="text-[#0175C2] font-bold text-lg leading-none">💡</div>
                  <div class="leading-relaxed whitespace-pre-line">{activeLesson.notes}</div>
                </div>
              {/if}

              <!-- Starter Code Block with Copy Button -->
              {#if activeLesson.starter_code}
                <div class="rounded-xl border border-slate-800 bg-slate-900 text-slate-100 overflow-hidden shadow-xs">
                  <div class="flex items-center justify-between px-4 py-2 bg-slate-950 border-b border-slate-800 text-xs text-slate-400">
                    <span class="font-mono font-medium text-slate-300">starter_code.dart</span>
                    <button
                      on:click={copyCode}
                      class="flex items-center gap-1.5 px-2.5 py-1 rounded bg-slate-800 hover:bg-[#0175C2] text-white transition text-xs font-semibold"
                    >
                      {#if copied}
                        <svg class="w-3.5 h-3.5 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
                        <span>Tersalin!</span>
                      {:else}
                        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg>
                        <span>Salin Kode</span>
                      {/if}
                    </button>
                  </div>
                  <pre class="p-4 text-xs font-mono overflow-x-auto leading-relaxed">{activeLesson.starter_code}</pre>
                </div>
              {/if}

            </div>

            <!-- Right Column: Sticky Embedded DartPad -->
            <div class="flex-1 flex flex-col h-full bg-[#1e1e1e] overflow-hidden {activeMobileTab === 'material' ? 'hidden md:flex' : 'flex'}">
              
              <!-- DartPad Bar -->
              <div class="px-4 py-2.5 bg-slate-900 border-b border-slate-800 flex items-center justify-between text-xs text-slate-300 flex-shrink-0">
                <div class="flex items-center gap-2">
                  <div class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></div>
                  <span class="font-bold tracking-wide text-white">DartPad Interactive Runner</span>
                </div>
                <a
                  href="https://dartpad.dev" 
                  target="_blank" 
                  rel="noopener noreferrer"
                  class="flex items-center gap-1 text-sky-400 hover:text-sky-300 font-medium transition"
                  title="Buka di tab baru"
                >
                  <span>Buka di DartPad.dev</span>
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path></svg>
                </a>
              </div>

              <!-- Iframe DartPad -->
              <div class="flex-1 w-full h-full relative">
                <iframe
                  title="DartPad Interactive Playground"
                  src="https://dartpad.dev/embed-dart.html?theme=dark"
                  class="w-full h-full border-0 absolute inset-0"
                  loading="lazy"
                ></iframe>
              </div>

            </div>

          </div>
        {/if}
      </div>

    </main>

  </div>
{/if}
