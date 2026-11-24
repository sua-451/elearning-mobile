<script>
  import { onMount } from 'svelte';
  import { getToken } from './lib/api.js';
  import { isDark, toggleTheme, initTheme } from './lib/theme.js';
  import Logo from './components/Logo.svelte';
  import Home from './pages/Home.svelte';
  import MaterialDetail from './pages/MaterialDetail.svelte';
  import AdminLogin from './pages/AdminLogin.svelte';
  import AdminDashboard from './pages/AdminDashboard.svelte';
  import { Sun, Moon, BookOpen, Shield } from 'lucide-svelte';

  let currentHash = window.location.hash || '#/';
  let routeParams = {};
  let isLoggedIn = !!getToken();

  // State to gracefully track if background & corner mascot images exist
  let bgWallpaperExists = true;
  let cornerMascotExists = true;

  function parseRoute() {
    currentHash = window.location.hash || '#/';
    isLoggedIn = !!getToken();

    const cleanPath = currentHash.replace(/^#/, '') || '/';

    // Route matching
    if (cleanPath === '/' || cleanPath === '') {
      routeParams = { page: 'home' };
    } else if (cleanPath.startsWith('/materials/')) {
      const slug = cleanPath.replace('/materials/', '').split('/')[0];
      routeParams = { page: 'material-detail', slug };
    } else if (cleanPath === '/admin/login') {
      routeParams = { page: 'admin-login' };
    } else if (cleanPath === '/admin/dashboard' || cleanPath.startsWith('/admin')) {
      routeParams = { page: 'admin-dashboard' };
    } else {
      routeParams = { page: 'home' };
    }
  }

  onMount(() => {
    initTheme();
    parseRoute();
    window.addEventListener('hashchange', parseRoute);
    return () => {
      window.removeEventListener('hashchange', parseRoute);
    };
  });
</script>

<div class="min-h-screen flex flex-col bg-slate-50 dark:bg-slate-900 text-slate-800 dark:text-slate-100 transition-colors duration-200 relative overflow-x-hidden">
  
  <!-- ========================================================= -->
  <!-- 1. SLOT BACKGROUND DEKORASI (Wallpaper & Subtle Overlay)  -->
  <!-- ========================================================= -->
  <div class="fixed inset-0 pointer-events-none z-0 overflow-hidden select-none" aria-hidden="true">
    <!-- Optional Custom Anime Wallpaper (fallback gracefully if not placed yet) -->
    {#if bgWallpaperExists}
      <img
        src="/decorations/bg.webp"
        alt=""
        class="w-full h-full object-cover object-center opacity-15 dark:opacity-20 transition-opacity duration-500"
        on:error={() => bgWallpaperExists = false}
      />
    {/if}

    <!-- Dynamic Soft Gradient & Dot Pattern Overlays to preserve text readability -->
    <div class="absolute inset-0 anime-bg-pattern"></div>
    <div class="absolute inset-0 bg-gradient-to-b from-slate-50/70 via-transparent to-slate-50/80 dark:from-slate-900/70 dark:via-transparent dark:to-slate-900/80"></div>
  </div>

  <!-- ========================================================= -->
  <!-- 2. GLOBAL NAVIGATION HEADER                               -->
  <!-- ========================================================= -->
  <header class="h-16 bg-white/85 dark:bg-slate-950/85 backdrop-blur-md border-b border-slate-200/80 dark:border-slate-800/80 sticky top-0 z-30 flex items-center justify-between px-3 sm:px-8 transition-colors">
    <a href="#/" class="flex items-center gap-2 sm:gap-3 font-bold text-slate-900 dark:text-white tracking-tight hover:opacity-90 transition min-h-[44px]">
      <Logo />
      <div class="flex flex-col">
        <span class="leading-none text-slate-900 dark:text-white font-extrabold text-sm sm:text-base">MobileIndex</span>
        <span class="text-[10px] text-slate-500 dark:text-slate-400 font-medium leading-tight">Dart & Flutter Hub</span>
      </div>
    </a>

    <nav class="flex items-center gap-1.5 sm:gap-4 text-sm font-medium">
      <a 
        href="#/" 
        class="flex items-center gap-1.5 px-3 py-2 rounded-lg transition text-xs sm:text-sm min-h-[44px] min-w-[44px] justify-center {routeParams.page === 'home' ? 'text-[#0175C2] font-semibold bg-sky-50 dark:bg-sky-950/50' : 'text-slate-600 dark:text-slate-300 hover:text-[#0175C2] dark:hover:text-sky-400'}"
      >
        <BookOpen class="w-4 h-4 flex-shrink-0" />
        <span class="hidden xs:inline sm:inline">Materi</span>
      </a>
      
      {#if isLoggedIn}
        <a 
          href="#/admin/dashboard" 
          class="flex items-center gap-1.5 px-3 py-2 rounded-lg bg-[#E1F5FE] dark:bg-sky-950 text-[#0175C2] dark:text-sky-400 font-semibold hover:bg-sky-100 dark:hover:bg-sky-900 transition text-xs min-h-[44px] justify-center"
        >
          <Shield class="w-3.5 h-3.5 flex-shrink-0" />
          <span class="hidden xs:inline">Admin</span>
        </a>
      {:else}
        <a 
          href="#/admin/login" 
          class="flex items-center gap-1.5 px-2.5 py-2 text-xs text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 transition min-h-[44px] justify-center"
        >
          <Shield class="w-3.5 h-3.5 flex-shrink-0" />
          <span class="hidden sm:inline">Admin</span> Login
        </a>
      {/if}

      <!-- Dark Mode Toggle Button with >=44px touch area -->
      <button 
        type="button" 
        on:click={toggleTheme}
        class="p-2.5 rounded-lg text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition min-h-[44px] min-w-[44px] flex items-center justify-center cursor-pointer"
        title={$isDark ? 'Beralih ke Mode Terang' : 'Beralih ke Mode Gelap'}
        aria-label="Toggle tema gelap/terang"
      >
        {#if $isDark}
          <Sun class="w-4 h-4 text-amber-400" />
        {:else}
          <Moon class="w-4 h-4 text-slate-600" />
        {/if}
      </button>
    </nav>
  </header>

  <!-- ========================================================= -->
  <!-- 3. DYNAMIC PAGE CONTENT (Layered above background)        -->
  <!-- ========================================================= -->
  <main class="flex-1 flex flex-col min-w-0 relative z-10">
    {#if routeParams.page === 'home'}
      <Home />
    {:else if routeParams.page === 'material-detail'}
      <MaterialDetail slug={routeParams.slug} />
    {:else if routeParams.page === 'admin-login'}
      <AdminLogin />
    {:else if routeParams.page === 'admin-dashboard'}
      <AdminDashboard />
    {/if}
  </main>

  <!-- ========================================================= -->
  <!-- 4. FLOATING CORNER MASCOT (Anime Chibi in Bottom Right)   -->
  <!-- ========================================================= -->
  <!-- Hidden on mobile when reading detail to avoid covering code/buttons -->
  {#if cornerMascotExists && routeParams.page !== 'material-detail'}
    <div 
      class="fixed bottom-0 right-1 sm:right-2 pointer-events-none z-20 hidden sm:block select-none"
      aria-hidden="true"
    >
      <img
        src="/decorations/ellen-removebg.png"
        alt="Ellen Mascot"
        class="w-14 sm:w-18 md:w-22 lg:w-24 max-h-28 object-contain object-bottom drop-shadow-md hover:scale-110 transition-transform duration-300 opacity-90 dark:opacity-85"
        on:error={() => cornerMascotExists = false}
      />
    </div>
  {/if}

  <!-- ========================================================= -->
  <!-- 5. FOOTER                                                 -->
  <!-- ========================================================= -->
  {#if routeParams.page !== 'material-detail'}
    <footer class="py-6 border-t border-slate-200/80 dark:border-slate-800/80 text-center text-xs text-slate-500 dark:text-slate-400 bg-white/70 dark:bg-slate-950/70 backdrop-blur-md transition-colors relative z-10">
      <div class="max-w-6xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-center gap-2">
        <div class="flex items-center gap-2 flex-wrap justify-center">
          <Logo />
          <span>© 2026 <strong>MobileIndex</strong>. Platform Pembelajaran Mobile Development.</span>
        </div>
      </div>
    </footer>
  {/if}

</div>
