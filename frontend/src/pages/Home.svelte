<script>
  import { onMount } from 'svelte';
  import { apiRequest } from '../lib/api.js';

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

  function navigate(path) {
    window.location.hash = path;
  }
</script>

<div class="max-w-6xl mx-auto px-4 py-10">
  <!-- Minimalist Header -->
  <header class="mb-10 text-center sm:text-left">
    <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-semibold bg-[#E1F5FE] text-[#0175C2] mb-3">
      <span>📱 Mobile Development</span>
      <span>•</span>
      <span>Dart & Flutter</span>
    </div>
    <h1 class="text-3xl sm:text-4xl font-extrabold text-slate-900 tracking-tight">
      Katalog Materi Pembelajaran Mobile Dev
    </h1>
    <p class="text-slate-600 mt-2 text-base sm:text-lg">
      Pelajari sintaks Dart dan arsitektur Flutter secara ringkas, cepat, dan langsung praktik interaktif di DartPad.
    </p>
  </header>

  <!-- Materials Grid -->
  {#if loading}
    <div class="flex justify-center items-center py-20 text-slate-400">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-[#0175C2]"></div>
      <span class="ml-3 font-medium">Memuat materi...</span>
    </div>
  {:else if error}
    <div class="p-4 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
      ⚠️ {error}. Pastikan server backend sudah berjalan di port 8080.
    </div>
  {:else if materials.length === 0}
    <div class="text-center py-16 bg-white rounded-xl border border-dashed border-slate-200 p-8">
      <p class="text-slate-500 font-medium">Belum ada materi pembelajaran yang tersedia.</p>
    </div>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each materials as item}
        <a 
          href={`#/materials/${item.slug}`}
          class="group bg-white rounded-xl border border-slate-200 p-6 flex flex-col justify-between hover:border-[#0175C2] hover:shadow-md transition-all cursor-pointer block no-underline"
        >
          <div>
            <div class="flex items-center justify-between mb-3">
              <span class="text-xs font-bold uppercase tracking-wider text-[#0175C2] bg-blue-50 px-2.5 py-1 rounded">
                Modul {item.order_index || 1}
              </span>
              <span class="text-xs text-slate-500">
                {item.lessons ? item.lessons.length : 0} Bab Pelajaran
              </span>
            </div>
            
            <h2 class="text-xl font-bold text-slate-900 group-hover:text-[#0175C2] transition-colors mb-2">
              {item.title}
            </h2>
            
            <p class="text-slate-600 text-sm line-clamp-3 leading-relaxed mb-4">
              {item.description || 'Tidak ada deskripsi singkat.'}
            </p>
          </div>

          <div class="pt-4 border-t border-slate-100 flex items-center justify-between text-xs font-medium text-slate-600">
            <span class="text-[#0175C2] font-semibold flex items-center gap-1 group-hover:translate-x-0.5 transition-transform">
              Mulai Belajar &rarr;
            </span>
          </div>
        </a>
      {/each}
    </div>
  {/if}
</div>
