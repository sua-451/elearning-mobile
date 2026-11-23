<script>
  import { onMount } from 'svelte';
  import { apiRequest } from '../lib/api.js';
  import { Smartphone, BookOpen, ArrowRight, Layers, AlertCircle, Loader2 } from 'lucide-svelte';

  let materials = [];
  let loading = true;
  let error = '';

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

<div class="max-w-6xl mx-auto px-4 py-10">
  <!-- Minimalist Header -->
  <header class="mb-10 text-center sm:text-left">
    <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-semibold bg-[#E1F5FE] dark:bg-sky-950/60 text-[#0175C2] dark:text-sky-300 mb-3 border border-sky-100 dark:border-sky-900/40">
      <Smartphone class="w-3.5 h-3.5" />
      <span>Mobile Development</span>
      <span>•</span>
      <span>Dart & Flutter</span>
    </div>
    <h1 class="text-3xl sm:text-4xl font-extrabold text-slate-900 dark:text-white tracking-tight">
      MobileIndex Katalog Pembelajaran
    </h1>
    <p class="text-slate-600 dark:text-slate-400 mt-2 text-base sm:text-lg">
      Pelajari sintaks Dart dan konsep Flutter secara ringkas, cepat, dan langsung praktik interaktif di DartPad.
    </p>
  </header>

  <!-- Materials Grid -->
  {#if loading}
    <div class="flex justify-center items-center py-20 text-slate-400">
      <Loader2 class="w-8 h-8 animate-spin text-[#0175C2]" />
      <span class="ml-3 font-medium">Memuat materi...</span>
    </div>
  {:else if error}
    <div class="p-4 bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-900 text-red-700 dark:text-red-300 rounded-lg text-sm flex items-center gap-2">
      <AlertCircle class="w-4 h-4 flex-shrink-0" />
      <span>{error}. Pastikan server backend sudah berjalan di port 8080.</span>
    </div>
  {:else if materials.length === 0}
    <div class="text-center py-16 bg-white dark:bg-slate-950 rounded-xl border border-dashed border-slate-200 dark:border-slate-800 p-8">
      <BookOpen class="w-8 h-8 text-slate-400 mx-auto mb-2" />
      <p class="text-slate-500 dark:text-slate-400 font-medium">Belum ada materi pembelajaran yang tersedia.</p>
    </div>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each materials as item}
        <a 
          href={`#/materials/${item.slug}`}
          class="group bg-white dark:bg-slate-950 rounded-xl border border-slate-200 dark:border-slate-800 p-6 flex flex-col justify-between hover:border-[#0175C2] dark:hover:border-sky-500 hover:shadow-md transition-all cursor-pointer block no-underline"
        >
          <div>
            <div class="flex items-center justify-between mb-3">
              <span class="text-xs font-bold uppercase tracking-wider text-[#0175C2] dark:text-sky-400 bg-blue-50 dark:bg-sky-950/60 px-2.5 py-1 rounded border border-blue-100 dark:border-sky-900/30">
                Modul {item.order_index || 1}
              </span>
              <span class="text-xs text-slate-500 dark:text-slate-400 flex items-center gap-1">
                <Layers class="w-3.5 h-3.5" />
                {item.lessons ? item.lessons.length : 0} Bab
              </span>
            </div>
            
            <h2 class="text-xl font-bold text-slate-900 dark:text-white group-hover:text-[#0175C2] dark:group-hover:text-sky-400 transition-colors mb-2">
              {item.title}
            </h2>
            
            <p class="text-slate-600 dark:text-slate-400 text-sm line-clamp-3 leading-relaxed mb-4">
              {item.description || 'Tidak ada deskripsi singkat.'}
            </p>
          </div>

          <div class="pt-4 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between text-xs font-medium text-slate-600 dark:text-slate-400">
            <span class="text-[#0175C2] dark:text-sky-400 font-semibold flex items-center gap-1 group-hover:translate-x-1 transition-transform">
              <span>Mulai Belajar</span>
              <ArrowRight class="w-3.5 h-3.5" />
            </span>
          </div>
        </a>
      {/each}
    </div>
  {/if}
</div>
