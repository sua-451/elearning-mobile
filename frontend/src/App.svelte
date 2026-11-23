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

<div class="min-h-screen flex flex-col bg-slate-50 dark:bg-slate-900 text-slate-800 dark:text-slate-100 transition-colors duration-200">
  
  <!-- Global Navigation Header -->
  <header class="h-16 bg-white dark:bg-slate-950 border-b border-slate-200 dark:border-slate-800 sticky top-0 z-30 flex items-center justify-between px-4 sm:px-8 transition-colors">
    <a href="#/" class="flex items-center gap-3 font-bold text-slate-900 dark:text-white tracking-tight hover:opacity-90 transition">
      <Logo />
      <div class="flex flex-col">
        <span class="leading-none text-slate-900 dark:text-white font-extrabold text-base">MobileIndex</span>
        <span class="text-[10px] text-slate-500 dark:text-slate-400 font-medium leading-tight">Dart & Flutter Hub</span>
      </div>
    </a>

    <nav class="flex items-center gap-3 sm:gap-4 text-sm font-medium">
      <a 
        href="#/" 
        class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg transition text-xs sm:text-sm {routeParams.page === 'home' ? 'text-[#0175C2] font-semibold bg-sky-50 dark:bg-sky-950/50' : 'text-slate-600 dark:text-slate-300 hover:text-[#0175C2] dark:hover:text-sky-400'}"
      >
        <BookOpen class="w-4 h-4" />
        <span class="hidden sm:inline">Katalog</span> Materi
      </a>
      
      {#if isLoggedIn}
        <a 
          href="#/admin/dashboard" 
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-[#E1F5FE] dark:bg-sky-950 text-[#0175C2] dark:text-sky-400 font-semibold hover:bg-sky-100 dark:hover:bg-sky-900 transition text-xs"
        >
          <Shield class="w-3.5 h-3.5" />
          <span>Dashboard</span>
        </a>
      {:else}
        <a 
          href="#/admin/login" 
          class="flex items-center gap-1 text-xs text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 transition"
        >
          <Shield class="w-3.5 h-3.5" />
          <span class="hidden sm:inline">Admin</span> Login
        </a>
      {/if}

      <!-- Dark Mode Toggle Button -->
      <button 
        type="button" 
        on:click={toggleTheme}
        class="p-2 rounded-lg text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition"
        title={$isDark ? 'Beralih ke Mode Terang' : 'Beralih ke Mode Gelap'}
      >
        {#if $isDark}
          <Sun class="w-4 h-4 text-amber-400" />
        {:else}
          <Moon class="w-4 h-4 text-slate-600" />
        {/if}
      </button>
    </nav>
  </header>

  <!-- Dynamic Page Content -->
  <main class="flex-1">
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

  <!-- Footer (only on Home or non-lesson full view) -->
  {#if routeParams.page !== 'material-detail'}
    <footer class="py-6 border-t border-slate-200 dark:border-slate-800 text-center text-xs text-slate-500 dark:text-slate-400 bg-white dark:bg-slate-950 transition-colors">
      <div class="max-w-6xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-center gap-2">
        <div class="flex items-center gap-2">
          <Logo />
          <span>© 2026 <strong>MobileIndex</strong>. Platform Pembelajaran Mobile Development.</span>
        </div>
      </div>
    </footer>
  {/if}

</div>
