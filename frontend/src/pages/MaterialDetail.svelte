<script>
  import { onMount } from 'svelte';
  import { apiRequest } from '../lib/api.js';
  import { marked } from 'marked';
  import { 
    ChevronLeft, 
    Menu, 
    ExternalLink, 
    Copy, 
    Check, 
    Code2, 
    Lightbulb, 
    ArrowLeft, 
    Loader2, 
    X,
    FileText
  } from 'lucide-svelte';

  export let slug = '';

  let material = null;
  let activeLesson = null;
  let activeBlocks = [];
  let loading = true;
  let error = '';
  let sidebarOpen = true;

  // DartPad panel toggle state (Keep-Alive via CSS)
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
  <div class="flex justify-center items-center h-[calc(100vh-64px)] text-slate-500 dark:text-slate-400">
    <Loader2 class="w-8 h-8 animate-spin text-[#0175C2]" />
    <span class="ml-3 font-medium">Memuat materi bab...</span>
  </div>
{:else if error}
  <div class="max-w-3xl mx-auto px-4 py-16 text-center">
    <div class="p-6 bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-900 rounded-xl text-red-700 dark:text-red-300">
      <h3 class="font-bold text-lg mb-2">Gagal Memuat Materi</h3>
      <p class="text-sm">{error}</p>
      <a href="#/" class="inline-flex items-center gap-1.5 mt-4 text-xs font-semibold px-4 py-2 bg-[#0175C2] text-white rounded-md hover:bg-[#02569B] transition">
        <ArrowLeft class="w-3.5 h-3.5" />
        <span>Kembali ke Beranda</span>
      </a>
    </div>
  </div>
{:else if material}
  <div class="flex h-[calc(100vh-64px)] overflow-hidden bg-slate-50 dark:bg-slate-900">
    
    <!-- Collapsible Sidebar -->
    <aside 
      class="{sidebarOpen ? 'w-72 sm:w-80' : 'w-0'} flex-shrink-0 bg-white dark:bg-slate-950 border-r border-slate-200 dark:border-slate-800 transition-all duration-300 flex flex-col overflow-hidden z-20 shadow-sm"
    >
      <div class="p-4 border-b border-slate-100 dark:border-slate-800 flex items-center justify-between bg-slate-50/50 dark:bg-slate-900/50">
        <div class="overflow-hidden pr-2">
          <span class="text-xs font-bold uppercase tracking-wider text-[#0175C2] dark:text-sky-400">Daftar Bab</span>
          <h2 class="font-bold text-slate-900 dark:text-white text-sm truncate" title={material.title}>{material.title}</h2>
        </div>
        <button 
          on:click={toggleSidebar} 
          class="p-1 rounded text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition" 
          title="Tutup Sidebar"
        >
          <ChevronLeft class="w-5 h-5" />
        </button>
      </div>

      <div class="flex-1 overflow-y-auto p-3 space-y-1.5">
        {#if material.lessons && material.lessons.length > 0}
          {#each material.lessons as lesson, idx}
            <button
              on:click={() => selectLesson(lesson)}
              class="w-full text-left px-3 py-2.5 rounded-lg text-sm flex items-start gap-2.5 transition-all {activeLesson && activeLesson.id === lesson.id ? 'bg-[#E1F5FE] dark:bg-sky-950 text-[#0175C2] dark:text-sky-300 font-semibold shadow-xs' : 'text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-900'}"
            >
              <span class="text-xs px-1.5 py-0.5 rounded {activeLesson && activeLesson.id === lesson.id ? 'bg-[#0175C2] text-white' : 'bg-slate-200 dark:bg-slate-800 text-slate-600 dark:text-slate-300'} font-mono mt-0.5">
                {idx + 1}
              </span>
              <span class="flex-1 line-clamp-2 leading-snug">{lesson.title}</span>
              {#if lesson.enable_dartpad}
                <span class="text-[10px] bg-sky-100 dark:bg-sky-900/60 text-[#0175C2] dark:text-sky-300 font-bold px-1.5 py-0.5 rounded uppercase self-center">
                  DartPad
                </span>
              {/if}
            </button>
          {/each}
        {:else}
          <p class="text-xs text-slate-400 p-4 text-center">Belum ada bab yang ditambahkan.</p>
        {/if}
      </div>

      <div class="p-3 border-t border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50">
        <a href="#/" class="text-xs text-slate-500 dark:text-slate-400 hover:text-[#0175C2] dark:hover:text-sky-400 flex items-center gap-1.5 py-1">
          <ArrowLeft class="w-3.5 h-3.5" />
          <span>Kembali ke Daftar Materi</span>
        </a>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-1 flex flex-col min-w-0 overflow-hidden bg-white dark:bg-slate-900">
      
      <!-- Top Action Bar -->
      <div class="h-12 border-b border-slate-100 dark:border-slate-800 px-4 flex items-center justify-between bg-white dark:bg-slate-950 flex-shrink-0">
        <div class="flex items-center gap-2">
          {#if !sidebarOpen}
            <button 
              on:click={toggleSidebar} 
              class="p-1.5 rounded-md text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 border border-slate-200 dark:border-slate-800 transition text-xs flex items-center gap-1.5"
              title="Buka Daftar Bab"
            >
              <Menu class="w-4 h-4" />
              <span>Daftar Bab</span>
            </button>
          {/if}
          <span class="text-xs text-slate-400 dark:text-slate-500 font-medium truncate">
            {material.title} &nbsp;/&nbsp; <span class="text-slate-700 dark:text-slate-200">{activeLesson ? activeLesson.title : 'Pilih Bab'}</span>
          </span>
        </div>

        <!-- DartPad Action Controls -->
        <div class="flex items-center gap-2">
          {#if activeLesson && activeLesson.enable_dartpad}
            <!-- Desktop Toggle Button for DartPad Keep-Alive Panel -->
            <button
              on:click={toggleDartpad}
              class="hidden md:flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition {dartpadOpen ? 'bg-sky-100 dark:bg-sky-950 text-[#0175C2] dark:text-sky-300 hover:bg-sky-200 dark:hover:bg-sky-900' : 'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700'}"
              title={dartpadOpen ? 'Tutup Panel DartPad' : 'Buka Panel DartPad'}
            >
              <Code2 class="w-3.5 h-3.5" />
              <span>DartPad</span>
              <span class="text-[10px] px-1 py-0.2 rounded {dartpadOpen ? 'bg-[#0175C2] text-white' : 'bg-slate-300 dark:bg-slate-700 text-slate-700 dark:text-slate-300'}">
                {dartpadOpen ? 'Aktif' : 'Off'}
              </span>
            </button>

            <!-- Mobile Toggle Tabs -->
            <div class="flex md:hidden bg-slate-100 dark:bg-slate-800 p-0.5 rounded-lg border border-slate-200 dark:border-slate-700 text-xs">
              <button 
                on:click={() => activeMobileTab = 'material'} 
                class="px-2.5 py-1 rounded-md font-medium transition {activeMobileTab === 'material' ? 'bg-white dark:bg-slate-900 text-[#0175C2] dark:text-sky-400 shadow-xs' : 'text-slate-600 dark:text-slate-400'}"
              >
                Materi
              </button>
              <button 
                on:click={() => { activeMobileTab = 'dartpad'; dartpadOpen = true; }} 
                class="px-2.5 py-1 rounded-md font-medium transition {activeMobileTab === 'dartpad' ? 'bg-[#0175C2] text-white shadow-xs' : 'text-slate-600 dark:text-slate-400'}"
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
            <div class="flex-1 overflow-y-auto p-4 sm:p-8 space-y-6 {activeLesson.enable_dartpad && dartpadOpen ? 'md:border-r border-slate-200 dark:border-slate-800' : ''} {activeMobileTab === 'dartpad' ? 'hidden md:block' : 'block'}">
              
              <div class="{activeLesson.enable_dartpad && dartpadOpen ? 'max-w-none' : 'max-w-3xl mx-auto'} space-y-6">
                <!-- Lesson Header -->
                <div>
                  <span class="text-xs font-semibold text-[#0175C2] dark:text-sky-400 tracking-wider uppercase">Bab Pelajaran</span>
                  <h1 class="text-2xl sm:text-3xl font-extrabold text-slate-900 dark:text-white tracking-tight mt-1">
                    {activeLesson.title}
                  </h1>
                </div>

                <!-- Dynamic Vertical Content Blocks Renderer -->
                <div class="space-y-6">
                  {#each activeBlocks as block (block.id)}
                    
                    <!-- 1. Text / Markdown Block -->
                    {#if block.type === 'text'}
                      <div class="prose prose-slate dark:prose-invert max-w-none text-slate-800 dark:text-slate-200 leading-relaxed text-sm sm:text-base">
                        {@html marked.parse(block.content || '')}
                      </div>

                    <!-- 2. Image Block -->
                    {:else if block.type === 'image' && block.url}
                      <div class="rounded-xl overflow-hidden border border-slate-200 dark:border-slate-800 bg-slate-100 dark:bg-slate-950 shadow-2xs">
                        <img src={block.url} alt={block.caption || activeLesson.title} class="w-full max-h-96 object-cover" />
                        {#if block.caption}
                          <p class="p-2.5 bg-white dark:bg-slate-950 text-xs text-center text-slate-500 dark:text-slate-400 border-t border-slate-100 dark:border-slate-800">
                            {block.caption}
                          </p>
                        {/if}
                      </div>

                    <!-- 3. Note / Callout Block -->
                    {:else if block.type === 'note' && block.content}
                      <div class="p-4 rounded-xl bg-sky-50/80 dark:bg-sky-950/40 border border-sky-200 dark:border-sky-900 text-sky-950 dark:text-sky-200 text-xs sm:text-sm flex gap-3 shadow-2xs">
                        <Lightbulb class="w-5 h-5 text-[#0175C2] dark:text-sky-400 flex-shrink-0 mt-0.5" />
                        <div class="leading-relaxed whitespace-pre-line">{block.content}</div>
                      </div>

                    <!-- 4. Code Block with One-Click Copy -->
                    {:else if block.type === 'code' && block.code}
                      <div class="rounded-xl border border-slate-800 dark:border-slate-700 bg-slate-900 dark:bg-slate-950 text-slate-100 overflow-hidden shadow-xs">
                        <div class="flex items-center justify-between px-4 py-2 bg-slate-950 dark:bg-black/40 border-b border-slate-800 dark:border-slate-800 text-xs text-slate-400">
                          <span class="font-mono font-medium text-slate-300">{block.title || 'code.dart'}</span>
                          <button
                            type="button"
                            on:click={() => copyCode(block.code, block.id)}
                            class="flex items-center gap-1.5 px-2.5 py-1 rounded bg-slate-800 hover:bg-[#0175C2] text-white transition text-xs font-semibold"
                          >
                            {#if copiedBlockId === block.id}
                              <Check class="w-3.5 h-3.5 text-emerald-400" />
                              <span>Tersalin!</span>
                            {:else}
                              <Copy class="w-3.5 h-3.5" />
                              <span>Salin Kode</span>
                            {/if}
                          </button>
                        </div>
                        <pre class="p-4 text-xs font-mono overflow-x-auto leading-relaxed">{block.code}</pre>
                      </div>
                    {/if}

                  {/each}

                  {#if activeBlocks.length === 0}
                    <div class="p-8 text-center text-slate-400 text-sm bg-slate-50 dark:bg-slate-950 rounded-xl border border-dashed border-slate-200 dark:border-slate-800">
                      <FileText class="w-8 h-8 mx-auto mb-2 opacity-50" />
                      Bab ini belum memiliki konten blok.
                    </div>
                  {/if}
                </div>

              </div>

            </div>

            <!-- Column 2: Sticky Embedded DartPad with Keep-Alive -->
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
                      <ExternalLink class="w-3.5 h-3.5" />
                    </a>
                    <button
                      type="button"
                      on:click={toggleDartpad}
                      class="text-slate-400 hover:text-slate-200 text-xs hidden md:inline-block ml-1"
                      title="Sembunyikan Panel"
                    >
                      <X class="w-4 h-4" />
                    </button>
                  </div>
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
            {/if}

          </div>
        {/if}
      </div>

    </main>

  </div>
{/if}
