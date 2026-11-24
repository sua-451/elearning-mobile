<script>
  import { apiRequest, setToken, setUser } from '../lib/api.js';
  import Logo from '../components/Logo.svelte';
  import { Shield, Lock, User, AlertCircle, Loader2 } from 'lucide-svelte';

  let username = '';
  let password = '';
  let loading = false;
  let error = '';

  async function handleLogin() {
    if (!username || !password) {
      error = 'Username dan password wajib diisi';
      return;
    }

    loading = true;
    error = '';
    try {
      const data = await apiRequest('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ username, password })
      });

      setToken(data.token);
      setUser(data.user);
      window.location.hash = '/admin/dashboard';
    } catch (err) {
      error = err.message || 'Login gagal';
    } finally {
      loading = false;
    }
  }
</script>

<div class="min-h-[calc(100vh-100px)] flex items-center justify-center px-4 py-8 sm:py-12">
  <div class="max-w-md w-full bg-white dark:bg-slate-950 rounded-2xl border border-slate-200 dark:border-slate-800 p-6 sm:p-8 shadow-sm transition-colors">
    <div class="text-center mb-6">
      <div class="flex justify-center mb-3">
        <Logo />
      </div>
      <h2 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white leading-tight">MobileIndex Admin Login</h2>
      <p class="text-xs text-slate-500 dark:text-slate-400 mt-1.5">Kelola materi dan bab pembelajaran Dart & Flutter</p>
    </div>

    {#if error}
      <div class="mb-4 p-3.5 bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-900 text-red-700 dark:text-red-300 text-xs rounded-xl flex items-center gap-2">
        <AlertCircle class="w-4 h-4 flex-shrink-0" />
        <span class="break-words">{error}</span>
      </div>
    {/if}

    <form on:submit|preventDefault={handleLogin} class="space-y-4">
      <div>
        <label for="username" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Username</label>
        <div class="relative">
          <input 
            id="username"
            type="text" 
            bind:value={username}
            placeholder="Masukkan username"
            required
            class="w-full pl-9 pr-3.5 py-2.5 border border-slate-300 dark:border-slate-700 rounded-xl text-sm bg-white dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-[#0175C2]"
          />
          <User class="w-4 h-4 text-slate-400 absolute left-3 top-3.5" />
        </div>
      </div>

      <div>
        <label for="password" class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Password</label>
        <div class="relative">
          <input 
            id="password"
            type="password" 
            bind:value={password}
            placeholder="Masukkan password"
            required
            class="w-full pl-9 pr-3.5 py-2.5 border border-slate-300 dark:border-slate-700 rounded-xl text-sm bg-white dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-[#0175C2]"
          />
          <Lock class="w-4 h-4 text-slate-400 absolute left-3 top-3.5" />
        </div>
      </div>

      <button
        type="submit"
        disabled={loading}
        class="w-full py-3 px-4 bg-[#0175C2] hover:bg-[#02569B] text-white text-sm font-semibold rounded-xl transition disabled:opacity-50 flex items-center justify-center gap-2 shadow-xs cursor-pointer min-h-[44px]"
      >
        {#if loading}
          <Loader2 class="w-4 h-4 animate-spin" />
          <span>Memproses...</span>
        {:else}
          <Shield class="w-4 h-4" />
          <span>Masuk ke Dashboard</span>
        {/if}
      </button>
    </form>
  </div>
</div>
