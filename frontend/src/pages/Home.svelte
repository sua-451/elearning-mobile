<script>
  import { onMount } from 'svelte';
  import { apiRequest } from '../lib/api.js';
  import { Smartphone, BookOpen, ArrowRight, Layers, AlertCircle, Loader2, Sparkles } from 'lucide-svelte';

  let materials = [];
  let loading = true;
  let error = '';

  // State to track if hero mascot image exists in public/decorations/
  let heroMascotExists = true;

  onMount(async () => {
    try {
      materials = await apiRequest('/materials');
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  });
</script>

<div class="max-w-6xl mx-auto px-4 sm:px-6 py-6 sm:py-10">
  
  <!-- ========================================================= -->
  <!-- HERO BANNER DENGAN SLOT MASKOT ANIME                      -->
  <!-- ========================================================= -->
  <header class="mb-10 sm:mb-12 relative overflow-hidden rounded-3xl bg-gradient-to-br from-white/90 via-sky-50/50 to-blue-50/60 dark:from-slate-950/90 dark:via-slate-900/80 dark:to-sky-950/40 border border-slate-200/80 dark:border-slate-800/80 p-6 sm:p-10 shadow-sm backdrop-blur-md">
    
    <!-- Background glow ornament -->
    <div class="absolute -top-16 -right-16 w-64 h-64 bg-sky-400/10 dark:bg-sky-500/10 rounded-full blur-3xl pointer-events-none" aria-hidden="true"></div>
    <div class="absolute -bottom-16 -left-16 w-64 h-64 bg-blue-500/10 dark:bg-blue-600/10 rounded-full blur-3xl pointer-events-none" aria-hidden="true"></div>

    <div class="flex flex-col-reverse md:flex-row items-center justify-between gap-6 relative z-10">
      
      <!-- Left Column: Text & Badges -->
      <div class="flex-1 text-center md:text-left">
        <div class="inline-flex items-center gap-1.5 sm:gap-2 px-3.5 py-1.5 rounded-full text-xs font-semibold bg-[#E1F5FE] dark:bg-sky-950/80 text-[#0175C2] dark:text-sky-300 mb-3 border border-sky-200/60 dark:border-sky-900/60 shadow-2xs">
          <Smartphone class="w-3.5 h-3.5 flex-shrink-0" />
          <span>Mobile Development</span>
          <span>•</span>
          <span>Dart & Flutter</span>
          <Sparkles class="w-3 h-3 text-amber-500" />
        </div>

        <h1 class="text-2xl sm:text-4xl lg:text-5xl font-extrabold text-slate-900 dark:text-white tracking-tight leading-tight">
          MobileIndex Katalog Pembelajaran
        </h1>

        <p class="text-slate-600 dark:text-slate-300 mt-3 text-sm sm:text-base lg:text-lg leading-relaxed max-w-xl">
          Pelajari sintaks Dart dan konsep Flutter secara ringkas, cepat, dan langsung praktik interaktif di DartPad bersama panduan terstruktur.
        </p>

        <div class="mt-5 flex items-center justify-center md:justify-start gap-4 text-xs font-medium text-slate-500 dark:text-slate-400">
          <span class="flex items-center gap-1.5">
            <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
            Siap Praktik Interaktif
          </span>
          <span>•</span>
          <span>DartPad Keep-Alive</span>
        </div>
      </div>

      <!-- Right Column: SLOT KHUSUS ILUSTRASI / MASKOT ANIME -->
      <div class="w-full md:w-auto flex justify-center flex-shrink-0 relative">
        <div class="relative w-44 h-44 sm:w-56 sm:h-56 md:w-64 md:h-64 flex items-center justify-center">
          
          {#if heroMascotExists}
            <img
              src="/decorations/miyabi-removebg.png"
              alt="Miyabi - Maskot MobileIndex"
              class="w-full h-full object-contain pointer-events-none drop-shadow-2xl hover:scale-105 transition-transform duration-500 select-none animate-in fade-in"
              on:error={() => heroMascotExists = false}
            />
          {:else}
            <!-- Elegant Fallback Graphic if hero-mascot.png is not yet added -->
            <div class="w-32 h-32 sm:w-40 sm:h-40 rounded-2xl bg-gradient-to-tr from-[#0175C2] to-sky-400 p-0.5 shadow-lg flex items-center justify-center text-white relative group">
              <div class="w-full h-full bg-white dark:bg-slate-900 rounded-[14px] flex flex-col items-center justify-center p-4 text-center">
                <Smartphone class="w-10 h-10 text-[#0175C2] dark:text-sky-400 mb-1" />
                <span class="text-xs font-bold text-slate-800 dark:text-slate-200">MobileIndex</span>
                <span class="text-[10px] text-slate-400">Dart & Flutter</span>
              </div>
            </div>
          {/if}

        </div>
      </div>

    </div>
  </header>

  <!-- ========================================================= -->
  <!-- MATERIALS GRID (Kartu materi semi-transparan estetis)    -->
  <!-- ========================================================= -->
  {#if loading}
    <div class="flex justify-center items-center py-20 text-slate-400">
      <Loader2 class="w-8 h-8 animate-spin text-[#0175C2]" />
      <span class="ml-3 font-medium">Memuat materi...</span>
    </div>
  {:else if error}
    <div class="p-4 bg-red-50/90 dark:bg-red-950/50 backdrop-blur-sm border border-red-200 dark:border-red-900 text-red-700 dark:text-red-300 rounded-2xl text-sm flex items-center gap-2">
      <AlertCircle class="w-4 h-4 flex-shrink-0" />
      <span>{error}. Pastikan server backend sudah berjalan di port 8080.</span>
    </div>
  {:else if materials.length === 0}
    <div class="text-center py-16 bg-white/80 dark:bg-slate-950/80 backdrop-blur-md rounded-2xl border border-dashed border-slate-200 dark:border-slate-800 p-8">
      <BookOpen class="w-8 h-8 text-slate-400 mx-auto mb-2" />
      <p class="text-slate-500 dark:text-slate-400 font-medium">Belum ada materi pembelajaran yang tersedia.</p>
    </div>
  {:else}
    <!-- Single-column on mobile (<768px), 2 cols on md, 3 cols on lg -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
      {#each materials as item}
        <a 
          href={`#/materials/${item.slug}`}
          class="group bg-white/85 dark:bg-slate-950/85 backdrop-blur-sm rounded-2xl border border-slate-200/90 dark:border-slate-800/90 p-5 sm:p-6 flex flex-col justify-between hover:border-[#0175C2] dark:hover:border-sky-500 hover:shadow-lg transition-all duration-200 cursor-pointer block no-underline min-h-[160px] relative overflow-hidden"
        >
          <!-- Corner accent hover light -->
          <div class="absolute -top-12 -right-12 w-24 h-24 bg-sky-500/0 group-hover:bg-sky-500/10 dark:group-hover:bg-sky-400/15 rounded-full blur-xl transition-all duration-300 pointer-events-none"></div>

          <div>
            <div class="flex items-center justify-between mb-3 gap-2">
              <span class="text-xs font-bold uppercase tracking-wider text-[#0175C2] dark:text-sky-400 bg-blue-50/90 dark:bg-sky-950/80 px-2.5 py-1 rounded-md border border-blue-100 dark:border-sky-900/40">
                Modul {item.order_index || 1}
              </span>
              <span class="text-xs text-slate-500 dark:text-slate-400 flex items-center gap-1 flex-shrink-0">
                <Layers class="w-3.5 h-3.5" />
                {item.lessons ? item.lessons.length : 0} Bab
              </span>
            </div>
            
            <h2 class="text-lg sm:text-xl font-bold text-slate-900 dark:text-white group-hover:text-[#0175C2] dark:group-hover:text-sky-400 transition-colors mb-2 leading-snug">
              {item.title}
            </h2>
            
            <p class="text-slate-600 dark:text-slate-400 text-xs sm:text-sm line-clamp-3 leading-relaxed mb-4">
              {item.description || 'Tidak ada deskripsi singkat.'}
            </p>
          </div>

          <div class="pt-3 border-t border-slate-100/90 dark:border-slate-800/90 flex items-center justify-between text-xs font-medium text-slate-600 dark:text-slate-400 min-h-[44px]">
            <span class="text-[#0175C2] dark:text-sky-400 font-semibold flex items-center gap-1 group-hover:translate-x-1 transition-transform">
              <span>Mulai Belajar</span>
              <ArrowRight class="w-3.5 h-3.5" />
            </span>
            <span class="text-[11px] text-slate-400 font-normal">Interaktif DartPad</span>
          </div>
        </a>
      {/each}
    </div>
  {/if}
</div>
