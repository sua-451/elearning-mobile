<script>
  import { onMount } from 'svelte';
  import { apiRequest } from '../lib/api.js';
  import { marked } from 'marked';

  export let slug = '';

  let material = null;
  let activeLesson = null;
  let activeBlocks = [];
  let loading = true;
  let error = '';
  let sidebarOpen = true;

  // DartPad panel toggle state: user can toggle DartPad panel on/off
  // Kept alive without re-mounting iframe
  let dartpadOpen = true;
  let activeMobileTab = 'material'; // 'material' or 'dartpad'
  let copiedBlockId = null;

  onMount(async () => {
    await fetchMaterial();
  });

  async function fetchMaterial() {
    loading = true;
    error = '';
    try {
      material = await apiRequest(`/materials/${slug}`);
      if (material && material.lessons && material.lessons.length > 0) {
        selectLesson(material.lessons[0]);
      }
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  function selectLesson(lesson) {
    activeLesson = lesson;
    copiedBlockId = null;
    activeMobileTab = 'material';

    // Parse blocks or generate fallback blocks from legacy content
    if (lesson.blocks && lesson.blocks.trim() !== '' && lesson.blocks.trim() !== '[]') {
      try {
        activeBlocks = JSON.parse(lesson.blocks);
      } catch (e) {
        activeBlocks = [];
      }
    } else {
      activeBlocks = [];
      if (lesson.image_url) {
        activeBlocks.push({ id: 'legacy_img', type: 'image', url: lesson.image_url, caption: '' });
      }
      if (lesson.content) {
        activeBlocks.push({ id: 'legacy_content', type: 'text', content: lesson.content });
      }
      if (lesson.notes) {
        activeBlocks.push({ id: 'legacy_notes', type: 'note', content: lesson.notes });
      }
      if (lesson.starter_code) {
        activeBlocks.push({ id: 'legacy_code', type: 'code', title: 'Starter Code', code: lesson.starter_code });
      }
    }

    // Set dartpad default state based on lesson configuration
    if (lesson.enable_dartpad) {
      dartpadOpen = true;
    } else {
      dartpadOpen = false;
    }
  }

  function copyCode(codeText, blockId) {
    if (!codeText) return;
    navigator.clipboard.writeText(codeText);
    copiedBlockId = blockId;
    setTimeout(() => {
      copiedBlockId = null;
    }, 2000);
  }

  function toggleSidebar() {
    sidebarOpen = !sidebarOpen;
  }

  function toggleDartpad() {
    dartpadOpen = !dartpadOpen;
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

        <!-- DartPad Action Controls -->
        <div class="flex items-center gap-2">
          {#if activeLesson && activeLesson.enable_dartpad}
            <!-- Desktop Toggle Button for DartPad Keep-Alive Panel -->
            <button
              on:click={toggleDartpad}
              class="hidden md:flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition {dartpadOpen ? 'bg-sky-100 text-[#0175C2] hover:bg-sky-200' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}"
              title={dartpadOpen ? 'Tutup Panel DartPad' : 'Buka Panel DartPad'}
            >
              <span>💻 DartPad</span>
              <span class="text-[10px] px-1 py-0.2 rounded {dartpadOpen ? 'bg-[#0175C2] text-white' : 'bg-slate-300 text-slate-700'}">
                {dartpadOpen ? 'Aktif' : 'Off'}
              </span>
            </button>

            <!-- Mobile Toggle Tabs (if DartPad enabled) -->
            <div class="flex md:hidden bg-slate-100 p-0.5 rounded-lg border border-slate-200 text-xs">
              <button 
                on:click={() => activeMobileTab = 'material'} 
                class="px-2.5 py-1 rounded-md font-medium transition {activeMobileTab === 'material' ? 'bg-white text-[#0175C2] shadow-xs' : 'text-slate-600'}"
              >
                Materi
              </button>
              <button 
                on:click={() => { activeMobileTab = 'dartpad'; dartpadOpen = true; }} 
                class="px-2.5 py-1 rounded-md font-medium transition {activeMobileTab === 'dartpad' ? 'bg-[#0175C2] text-white shadow-xs' : 'text-slate-600'}"
              >
                DartPad
              </button>
            </div>
          {/if}
        </div>
      </div>

      <!-- Reading and Practice Area -->
      <div class="flex-1 overflow-hidden relative">
        {#if !activeLesson}
          <div class="flex items-center justify-center h-full text-slate-400 text-sm">
            Pilih salah satu bab dari panel kiri untuk mulai membaca.
          </div>
        {:else}
          <!-- Adaptive Layout: Split-screen desktop if Dartpad is enabled & open; Keep-alive via hidden class -->
          <div class="h-full flex flex-col md:flex-row overflow-hidden">
            
            <!-- Column 1: Reading Content (Dynamic Vertical Blocks) -->
            <div class="flex-1 overflow-y-auto p-4 sm:p-8 space-y-6 {activeLesson.enable_dartpad && dartpadOpen ? 'md:border-r border-slate-200' : ''} {activeMobileTab === 'dartpad' ? 'hidden md:block' : 'block'}">
              
              <div class="{activeLesson.enable_dartpad && dartpadOpen ? 'max-w-none' : 'max-w-3xl mx-auto'} space-y-6">
                <!-- Lesson Header -->
                <div>
                  <span class="text-xs font-semibold text-[#0175C2] tracking-wider uppercase">Bab Pelajaran</span>
                  <h1 class="text-2xl sm:text-3xl font-extrabold text-slate-900 tracking-tight mt-1">
                    {activeLesson.title}
                  </h1>
                </div>

                <!-- Dynamic Vertical Content Blocks Renderer -->
                <div class="space-y-6">
                  {#each activeBlocks as block (block.id)}
                    
                    <!-- 1. Text / Markdown Block -->
                    {#if block.type === 'text'}
                      <div class="prose prose-slate max-w-none text-slate-800 leading-relaxed text-sm sm:text-base">
                        {@html marked.parse(block.content || '')}
                      </div>

                    <!-- 2. Image Block -->
                    {:else if block.type === 'image' && block.url}
                      <div class="rounded-xl overflow-hidden border border-slate-200 bg-slate-100 shadow-2xs">
                        <img src={block.url} alt={block.caption || activeLesson.title} class="w-full max-h-96 object-cover" />
                        {#if block.caption}
                          <p class="p-2.5 bg-white text-xs text-center text-slate-500 border-t border-slate-100">
                            {block.caption}
                          </p>
                        {/if}
                      </div>

                    <!-- 3. Note / Callout Block -->
                    {:else if block.type === 'note' && block.content}
                      <div class="p-4 rounded-xl bg-sky-50/80 border border-sky-200 text-sky-950 text-xs sm:text-sm flex gap-3 shadow-2xs">
                        <div class="text-[#0175C2] font-bold text-lg leading-none">💡</div>
                        <div class="leading-relaxed whitespace-pre-line">{block.content}</div>
                      </div>

                    <!-- 4. Code Block with One-Click Copy -->
                    {:else if block.type === 'code' && block.code}
                      <div class="rounded-xl border border-slate-800 bg-slate-900 text-slate-100 overflow-hidden shadow-xs">
                        <div class="flex items-center justify-between px-4 py-2 bg-slate-950 border-b border-slate-800 text-xs text-slate-400">
                          <span class="font-mono font-medium text-slate-300">{block.title || 'code.dart'}</span>
                          <button
                            type="button"
                            on:click={() => copyCode(block.code, block.id)}
                            class="flex items-center gap-1.5 px-2.5 py-1 rounded bg-slate-800 hover:bg-[#0175C2] text-white transition text-xs font-semibold"
                          >
                            {#if copiedBlockId === block.id}
                              <svg class="w-3.5 h-3.5 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
                              <span>Tersalin!</span>
                            {:else}
                              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg>
                              <span>Salin Kode</span>
                            {/if}
                          </button>
                        </div>
                        <pre class="p-4 text-xs font-mono overflow-x-auto leading-relaxed">{block.code}</pre>
                      </div>
                    {/if}

                  {/each}

                  {#if activeBlocks.length === 0}
                    <div class="p-8 text-center text-slate-400 text-sm bg-slate-50 rounded-xl border border-dashed border-slate-200">
                      Bab ini belum memiliki konten blok.
                    </div>
                  {/if}
                </div>

              </div>

            </div>

            <!-- Column 2: Sticky Embedded DartPad with Keep-Alive (CSS hidden/flex, preserved iframe state) -->
            {#if activeLesson.enable_dartpad}
              <div 
                class="flex-1 flex-col h-full bg-[#1e1e1e] overflow-hidden 
                {dartpadOpen ? (activeMobileTab === 'material' ? 'hidden md:flex' : 'flex') : 'hidden'}"
              >
                <!-- DartPad Bar -->
                <div class="px-4 py-2.5 bg-slate-900 border-b border-slate-800 flex items-center justify-between text-xs text-slate-300 flex-shrink-0">
                  <div class="flex items-center gap-2">
                    <div class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></div>
                    <span class="font-bold tracking-wide text-white">DartPad Interactive Runner</span>
                  </div>
                  <div class="flex items-center gap-3">
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
                    <button
                      type="button"
                      on:click={toggleDartpad}
                      class="text-slate-400 hover:text-slate-200 text-xs hidden md:inline-block ml-1"
                      title="Sembunyikan Panel"
                    >
                      ✕
                    </button>
                  </div>
                </div>

                <!-- Iframe DartPad (Never destroyed when toggled, preventing code loss or reload) -->
                <div class="flex-1 w-full h-full relative">
                  <iframe
                    title="DartPad Interactive Playground"
                    src="https://dartpad.dev/embed-dart.html?theme=dark"
                    class="w-full h-full border-0 absolute inset-0"
                    loading="lazy"
                  ></iframe>
                </div>

              </div>
            {/if}

          </div>
        {/if}
      </div>

    </main>

  </div>
{/if}
