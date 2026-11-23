<script>
  import { onMount } from 'svelte';
  import { getToken, setToken } from './lib/api.js';
  import Home from './pages/Home.svelte';
  import MaterialDetail from './pages/MaterialDetail.svelte';
  import AdminLogin from './pages/AdminLogin.svelte';
  import AdminDashboard from './pages/AdminDashboard.svelte';

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
    parseRoute();
    window.addEventListener('hashchange', parseRoute);
    return () => {
      window.removeEventListener('hashchange', parseRoute);
    };
  });
</script>

<div class="min-h-screen flex flex-col bg-slate-50 text-slate-800">
  
  <!-- Global Navigation Header -->
  <header class="h-16 bg-white border-b border-slate-200 sticky top-0 z-30 flex items-center justify-between px-4 sm:px-8">
    <a href="#/" class="flex items-center gap-2.5 font-bold text-slate-900 text-lg tracking-tight hover:opacity-90">
      <div class="w-8 h-8 rounded-lg bg-[#0175C2] flex items-center justify-center text-white font-extrabold text-sm shadow-xs">
        D
      </div>
      <div class="flex flex-col">
        <span class="leading-none text-slate-900 font-extrabold">Dart & Flutter Hub</span>
        <span class="text-[10px] text-slate-500 font-medium leading-tight">Mobile Dev Learning</span>
      </div>
    </a>

    <nav class="flex items-center gap-4 text-sm font-medium">
      <a href="#/" class="text-slate-600 hover:text-[#0175C2] transition {routeParams.page === 'home' ? 'text-[#0175C2] font-semibold' : ''}">
        Katalog Materi
      </a>
      
      {#if isLoggedIn}
        <a href="#/admin/dashboard" class="px-3 py-1.5 rounded-lg bg-[#E1F5FE] text-[#0175C2] font-semibold hover:bg-sky-100 transition text-xs">
          Dashboard Admin
        </a>
      {:else}
        <a href="#/admin/login" class="text-xs text-slate-500 hover:text-slate-800 transition">
          Admin Login
        </a>
      {/if}
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
    <footer class="py-6 border-t border-slate-200 text-center text-xs text-slate-500 bg-white">
      <div class="max-w-6xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-2">
        <p>© 2026 Dart & Flutter Mobile Learning Platform. Dibangun dengan Golang, Svelte, dan MySQL.</p>
        <p class="text-slate-400">Ringan • Minimalis • Cepat</p>
      </div>
    </footer>
  {/if}

</div>
