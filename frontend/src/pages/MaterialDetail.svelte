<script>
  import { onMount, onDestroy } from 'svelte';
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
    FileText,
    Maximize2,
    Minimize2
  } from 'lucide-svelte';

  export let slug = '';

  let material = null;
  let activeLesson = null;
  let activeBlocks = [];
  let loading = true;
  let error = '';
  let sidebarOpen = false; // default false on mobile for clean start, initialized based on window width

  // DartPad panel toggle state (Keep-Alive via CSS)
  let dartpadOpen = false; // Toggle state for DartPad
  let copiedBlockId = null;

  function handleResize() {
    if (typeof window !== 'undefined') {
      if (window.innerWidth >= 1024) {
        sidebarOpen = true;
      }
    }
  }

  // Manage body scroll lock when DartPad is open on mobile (<768px)
  $: {
    if (typeof document !== 'undefined') {
      const isMobile = window.innerWidth < 768;
      if (dartpadOpen && isMobile && activeLesson && activeLesson.enable_dartpad) {
        document.body.classList.add('overflow-hidden-mobile');
      } else {
        document.body.classList.remove('overflow-hidden-mobile');
      }
    }
  }

  onMount(async () => {
    if (typeof window !== 'undefined') {
      sidebarOpen = window.innerWidth >= 1024;
      window.addEventListener('resize', handleResize);
    }
    await fetchMaterial();
  });

  onDestroy(() => {
    if (typeof document !== 'undefined') {
      document.body.classList.remove('overflow-hidden-mobile');
    }
    if (typeof window !== 'undefined') {
      window.removeEventListener('resize', handleResize);
    }
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

    // Automatically close mobile sidebar when selecting a lesson
    if (typeof window !== 'undefined' && window.innerWidth < 768) {
      sidebarOpen = false;
    }

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

    // Default desktop opens if enabled; mobile stays false until user taps Buka DartPad
    if (typeof window !== 'undefined' && window.innerWidth >= 768) {
      dartpadOpen = !!lesson.enable_dartpad;
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

  function closeDartpad() {
    dartpadOpen = false;
  }
</script>

{#if loading}
  <div class="flex justify-center items-center h-[calc(100vh-64px)] text-slate-500 dark:text-slate-400">
    <Loader2 class="w-8 h-8 animate-spin text-[#0175C2]" />
    <span class="ml-3 font-medium text-sm sm:text-base">Memuat materi bab...</span>
  </div>
{:else if error}
  <div class="max-w-3xl mx-auto px-4 py-16 text-center">
    <div class="p-6 bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-900 rounded-xl text-red-700 dark:text-red-300">
      <h3 class="font-bold text-lg mb-2">Gagal Memuat Materi</h3>
      <p class="text-sm">{error}</p>
      <a href="#/" class="inline-flex items-center gap-1.5 mt-4 text-xs font-semibold px-4 py-2 bg-[#0175C2] text-white rounded-lg hover:bg-[#02569B] transition min-h-[44px]">
        <ArrowLeft class="w-3.5 h-3.5" />
        <span>Kembali ke Beranda</span>
      </a>
    </div>
  </div>
{:else if material}
  <div class="flex h-[calc(100vh-64px)] overflow-hidden bg-slate-50 dark:bg-slate-900 relative">
    
    <!-- Mobile Backdrop for Sidebar -->
    {#if sidebarOpen}
      <button 
        type="button"
        on:click={() => sidebarOpen = false} 
        class="fixed inset-0 bg-slate-900/60 z-30 md:hidden backdrop-blur-xs cursor-pointer border-0 w-full h-full"
        aria-label="Tutup menu samping"
      ></button>
    {/if}

    <!-- Collapsible Sidebar (Drawer on mobile, inline flex column on desktop) -->
    <aside 
      class="{sidebarOpen ? 'w-72 sm:w-80 translate-x-0' : 'w-0 -translate-x-full md:translate-x-0'} fixed md:static inset-y-0 left-0 top-16 md:top-0 z-40 md:z-20 flex-shrink-0 bg-white dark:bg-slate-950 border-r border-slate-200 dark:border-slate-800 transition-all duration-300 flex flex-col overflow-hidden shadow-lg md:shadow-none"
    >
      <div class="p-4 border-b border-slate-100 dark:border-slate-800 flex items-center justify-between bg-slate-50/50 dark:bg-slate-900/50 flex-shrink-0">
        <div class="overflow-hidden pr-2">
          <span class="text-xs font-bold uppercase tracking-wider text-[#0175C2] dark:text-sky-400">Daftar Bab</span>
          <h2 class="font-bold text-slate-900 dark:text-white text-sm truncate" title={material.title}>{material.title}</h2>
        </div>
        <button 
          on:click={toggleSidebar} 
          class="p-2 rounded-lg text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition min-h-[44px] min-w-[44px] flex items-center justify-center cursor-pointer" 
          title="Tutup Sidebar"
          aria-label="Tutup Sidebar"
        >
          <ChevronLeft class="w-5 h-5" />
        </button>
      </div>

      <div class="flex-1 overflow-y-auto p-3 space-y-1.5 overscroll-contain">
        {#if material.lessons && material.lessons.length > 0}
          {#each material.lessons as lesson, idx}
            <button
              on:click={() => selectLesson(lesson)}
              class="w-full text-left px-3 py-2.5 rounded-xl text-sm flex items-start gap-2.5 transition-all min-h-[44px] cursor-pointer {activeLesson && activeLesson.id === lesson.id ? 'bg-[#E1F5FE] dark:bg-sky-950 text-[#0175C2] dark:text-sky-300 font-semibold shadow-xs' : 'text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-900'}"
            >
              <span class="text-xs px-2 py-0.5 rounded-md {activeLesson && activeLesson.id === lesson.id ? 'bg-[#0175C2] text-white' : 'bg-slate-200 dark:bg-slate-800 text-slate-600 dark:text-slate-300'} font-mono mt-0.5 flex-shrink-0">
                {idx + 1}
              </span>
              <span class="flex-1 line-clamp-2 leading-snug">{lesson.title}</span>
              {#if lesson.enable_dartpad}
                <span class="text-[10px] bg-sky-100 dark:bg-sky-900/60 text-[#0175C2] dark:text-sky-300 font-bold px-1.5 py-0.5 rounded uppercase self-center flex-shrink-0">
                  DartPad
                </span>
              {/if}
            </button>
          {/each}
        {:else}
          <p class="text-xs text-slate-400 p-4 text-center">Belum ada bab yang ditambahkan.</p>
        {/if}
      </div>

      <div class="p-3 border-t border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50 flex-shrink-0">
        <a href="#/" class="text-xs text-slate-500 dark:text-slate-400 hover:text-[#0175C2] dark:hover:text-sky-400 flex items-center gap-2 py-2 min-h-[44px]">
          <ArrowLeft class="w-4 h-4 flex-shrink-0" />
          <span>Kembali ke Daftar Materi</span>
        </a>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-1 flex flex-col min-w-0 overflow-hidden bg-white dark:bg-slate-900">
      
      <!-- Top Action Bar -->
      <div class="h-14 border-b border-slate-100 dark:border-slate-800 px-3 sm:px-6 flex items-center justify-between bg-white dark:bg-slate-950 flex-shrink-0 gap-2">
        <div class="flex items-center gap-2 min-w-0">
          <button 
            on:click={toggleSidebar} 
            class="p-2 rounded-lg text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 border border-slate-200 dark:border-slate-800 transition text-xs flex items-center gap-1.5 min-h-[40px] flex-shrink-0 cursor-pointer"
            title="Daftar Bab"
            aria-label="Toggle Bab"
          >
            <Menu class="w-4 h-4" />
            <span class="hidden sm:inline font-medium">Daftar Bab</span>
          </button>
          
          <div class="text-xs text-slate-400 dark:text-slate-500 font-medium truncate">
            <span class="hidden md:inline">{material.title} / </span>
            <span class="text-slate-800 dark:text-slate-200 font-semibold">{activeLesson ? activeLesson.title : 'Pilih Bab'}</span>
          </div>
        </div>

        <!-- DartPad Action Controls (Mobile & Desktop) -->
        <div class="flex items-center gap-2 flex-shrink-0">
          {#if activeLesson && activeLesson.enable_dartpad}
            <!-- Toggle DartPad Button (Visible clearly on both mobile & desktop) -->
            <button
              on:click={toggleDartpad}
              class="flex items-center gap-1.5 px-3 py-2 rounded-xl text-xs font-semibold transition min-h-[40px] cursor-pointer shadow-2xs {dartpadOpen ? 'bg-sky-100 dark:bg-sky-950 text-[#0175C2] dark:text-sky-300 border border-sky-200 dark:border-sky-800' : 'bg-[#0175C2] hover:bg-[#02569B] text-white shadow-xs'}"
              title={dartpadOpen ? 'Sembunyikan DartPad' : 'Buka DartPad'}
            >
              <Code2 class="w-4 h-4 flex-shrink-0" />
              <span>{dartpadOpen ? 'Tutup DartPad' : 'Buka DartPad'}</span>
              <span class="text-[10px] px-1.5 py-0.5 rounded-full {dartpadOpen ? 'bg-[#0175C2] text-white' : 'bg-white/20 text-white font-mono'}">
                {dartpadOpen ? 'Aktif' : 'Coba'}
              </span>
            </button>
          {/if}
        </div>
      </div>

      <!-- Reading and Practice Area -->
      <div class="flex-1 overflow-hidden relative">
        {#if !activeLesson}
          <div class="flex items-center justify-center h-full text-slate-400 text-sm p-4 text-center">
            Pilih salah satu bab dari panel samping untuk mulai membaca.
          </div>
        {:else}
          <!-- Adaptive Layout:
               - Desktop: Split-screen horizontal row (Desktop Keep-Alive via CSS)
               - Mobile (<768px): Single-column reading view + Full-Screen Fixed Overlay for DartPad
          -->
          <div class="h-full flex flex-col md:flex-row overflow-hidden relative">
            
            <!-- Column 1: Reading Content (Single-column stack on mobile, responsive split on desktop) -->
            <div class="flex-1 overflow-y-auto px-4 py-6 sm:p-8 space-y-6 {activeLesson.enable_dartpad && dartpadOpen ? 'md:border-r border-slate-200 dark:border-slate-800' : ''} overscroll-contain">
              
              <div class="{activeLesson.enable_dartpad && dartpadOpen ? 'max-w-none' : 'max-w-3xl mx-auto'} space-y-6">
                <!-- Lesson Header -->
                <div>
                  <span class="text-xs font-semibold text-[#0175C2] dark:text-sky-400 tracking-wider uppercase">Bab Pelajaran</span>
                  <h1 class="text-xl sm:text-3xl font-extrabold text-slate-900 dark:text-white tracking-tight mt-1 leading-snug">
                    {activeLesson.title}
                  </h1>
                </div>

                <!-- Floating DartPad CTA banner for mobile readers if DartPad is enabled but currently closed -->
                {#if activeLesson.enable_dartpad && !dartpadOpen}
                  <div class="md:hidden p-3.5 bg-gradient-to-r from-sky-50 to-blue-50 dark:from-sky-950/40 dark:to-slate-900 border border-sky-200 dark:border-sky-900/60 rounded-xl flex items-center justify-between gap-3">
                    <div class="flex items-center gap-2.5">
                      <div class="p-2 bg-[#0175C2] text-white rounded-lg flex-shrink-0">
                        <Code2 class="w-4 h-4" />
                      </div>
                      <div class="text-xs">
                        <strong class="text-slate-900 dark:text-white block font-semibold">Tersedia Simulasi DartPad</strong>
                        <span class="text-slate-500 dark:text-slate-400">Jalankan dan coba kode langsung di HP</span>
                      </div>
                    </div>
                    <button
                      type="button"
                      on:click={toggleDartpad}
                      class="px-3 py-2 bg-[#0175C2] text-white text-xs font-bold rounded-lg hover:bg-[#02569B] transition flex-shrink-0 min-h-[40px] cursor-pointer"
                    >
                      Buka
                    </button>
                  </div>
                {/if}

                <!-- Dynamic Vertical Content Blocks Renderer -->
                <div class="space-y-6 max-w-full overflow-hidden">
                  {#each activeBlocks as block (block.id)}
                    
                    <!-- 1. Text / Markdown Block -->
                    {#if block.type === 'text'}
                      <div class="prose prose-slate dark:prose-invert max-w-none text-slate-800 dark:text-slate-200 leading-relaxed text-sm sm:text-base break-words">
                        {@html marked.parse(block.content || '')}
                      </div>

                    <!-- 2. Image Block (Responsive SVG / Image) -->
                    {:else if block.type === 'image' && block.url}
                      <div class="rounded-xl overflow-hidden border border-slate-200 dark:border-slate-800 bg-slate-100 dark:bg-slate-950 shadow-2xs max-w-full">
                        <img 
                          src={block.url} 
                          alt={block.caption || activeLesson.title} 
                          class="w-full h-auto max-h-96 object-contain sm:object-cover mx-auto" 
                          loading="lazy"
                        />
                        {#if block.caption}
                          <p class="p-2.5 bg-white dark:bg-slate-950 text-xs text-center text-slate-500 dark:text-slate-400 border-t border-slate-100 dark:border-slate-800">
                            {block.caption}
                          </p>
                        {/if}
                      </div>

                    <!-- 3. Note / Callout Block -->
                    {:else if block.type === 'note' && block.content}
                      <div class="p-3.5 sm:p-4 rounded-xl bg-sky-50/80 dark:bg-sky-950/40 border border-sky-200 dark:border-sky-900 text-sky-950 dark:text-sky-200 text-xs sm:text-sm flex gap-3 shadow-2xs">
                        <Lightbulb class="w-5 h-5 text-[#0175C2] dark:text-sky-400 flex-shrink-0 mt-0.5" />
                        <div class="leading-relaxed whitespace-pre-line break-words flex-1">{block.content}</div>
                      </div>

                    <!-- 4. Code Block with One-Click Copy & Responsive Horizontal Scroll -->
                    {:else if block.type === 'code' && block.code}
                      <div class="rounded-xl border border-slate-800 dark:border-slate-700 bg-slate-900 dark:bg-slate-950 text-slate-100 overflow-hidden shadow-xs max-w-full">
                        <div class="flex items-center justify-between px-3 sm:px-4 py-2 bg-slate-950 dark:bg-black/40 border-b border-slate-800 text-xs text-slate-400 gap-2">
                          <span class="font-mono font-medium text-slate-300 truncate">{block.title || 'code.dart'}</span>
                          <button
                            type="button"
                            on:click={() => copyCode(block.code, block.id)}
                            class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-[#0175C2] text-white transition text-xs font-semibold min-h-[36px] flex-shrink-0 cursor-pointer"
                            aria-label="Salin Kode"
                          >
                            {#if copiedBlockId === block.id}
                              <Check class="w-3.5 h-3.5 text-emerald-400" />
                              <span>Tersalin!</span>
                            {:else}
                              <Copy class="w-3.5 h-3.5" />
                              <span>Salin</span>
                            {/if}
                          </button>
                        </div>
                        <pre class="p-3 sm:p-4 text-xs font-mono overflow-x-auto leading-relaxed max-w-full text-slate-200">{block.code}</pre>
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

            <!-- Column 2 / Mobile Overlay: Embedded DartPad (Keep-Alive via CSS)
                 - On Mobile (<768px): fixed inset-0 z-50 h-[100dvh] w-full flex-col
                 - On Desktop (>=768px): flex-1 h-full flex-col relative
            -->
            {#if activeLesson.enable_dartpad}
              <div 
                class="bg-[#1e1e1e] overflow-hidden transition-all
                  {dartpadOpen 
                    ? 'fixed inset-0 z-50 h-[100dvh] w-full flex flex-col md:static md:inset-auto md:z-auto md:flex-1 md:h-full' 
                    : 'hidden'}"
              >
                <!-- Compact Header Bar on DartPad (optimized for mobile & desktop) -->
                <div class="px-3 sm:px-4 py-2.5 bg-slate-900 border-b border-slate-800 flex items-center justify-between text-xs text-slate-300 flex-shrink-0 gap-2 min-h-[48px]">
                  <div class="flex items-center gap-2 min-w-0">
                    <div class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse flex-shrink-0"></div>
                    <span class="font-bold tracking-wide text-white truncate text-xs sm:text-sm">DartPad Runner</span>
                  </div>

                  <div class="flex items-center gap-2 sm:gap-3 flex-shrink-0">
                    <a
                      href="https://dartpad.dev" 
                      target="_blank" 
                      rel="noopener noreferrer"
                      class="flex items-center gap-1 text-sky-400 hover:text-sky-300 font-medium transition py-1 px-1.5 text-[11px] sm:text-xs min-h-[36px]"
                      title="Buka di tab baru"
                    >
                      <span>DartPad.dev</span>
                      <ExternalLink class="w-3.5 h-3.5" />
                    </a>

                    <!-- Close / Exit Fullscreen Button (especially important on Mobile) -->
                    <button
                      type="button"
                      on:click={closeDartpad}
                      class="flex items-center gap-1 px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-white font-semibold text-xs transition min-h-[36px] cursor-pointer"
                      title="Keluar dari DartPad / Kembali ke Materi"
                      aria-label="Tutup DartPad"
                    >
                      <X class="w-4 h-4 text-slate-300" />
                      <span class="sm:inline">Tutup</span>
                    </button>
                  </div>
                </div>

                <!-- Iframe DartPad Runner with 100% dimensions inside container -->
                <div class="flex-1 w-full h-full relative bg-[#1e1e1e]">
                  <iframe
                    title="DartPad Interactive Playground"
                    src="https://dartpad.dev/embed-dart.html?theme=dark"
                    class="w-full h-full border-0 absolute inset-0"
                    loading="lazy"
                    allow="clipboard-write"
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
