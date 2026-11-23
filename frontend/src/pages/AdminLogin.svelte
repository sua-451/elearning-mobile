<script>
  import { apiRequest, setToken, setUser } from '../lib/api.js';

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

<div class="min-h-[calc(100vh-100px)] flex items-center justify-center px-4 py-12">
  <div class="max-w-md w-full bg-white rounded-2xl border border-slate-200 p-8 shadow-sm">
    <div class="text-center mb-6">
      <div class="w-12 h-12 bg-blue-50 text-[#0175C2] rounded-xl flex items-center justify-center mx-auto mb-3 font-bold text-xl">
        📱
      </div>
      <h2 class="text-2xl font-extrabold text-slate-900">Admin Login</h2>
      <p class="text-xs text-slate-500 mt-1">Kelola materi dan bab pembelajaran Dart & Flutter</p>
    </div>

    {#if error}
      <div class="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 text-xs rounded-lg">
        {error}
      </div>
    {/if}

    <form on:submit|preventDefault={handleLogin} class="space-y-4">
      <div>
        <label for="username" class="block text-xs font-semibold text-slate-700 mb-1">Username</label>
        <input 
          id="username"
          type="text" 
          bind:value={username}
          placeholder="admin"
          required
          class="w-full px-3.5 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#0175C2] focus:border-transparent"
        />
      </div>

      <div>
        <label for="password" class="block text-xs font-semibold text-slate-700 mb-1">Password</label>
        <input 
          id="password"
          type="password" 
          bind:value={password}
          placeholder="••••••••"
          required
          class="w-full px-3.5 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#0175C2] focus:border-transparent"
        />
      </div>

      <button
        type="submit"
        disabled={loading}
        class="w-full py-2.5 px-4 bg-[#0175C2] hover:bg-[#02569B] text-white text-sm font-semibold rounded-lg transition disabled:opacity-50"
      >
        {loading ? 'Memproses...' : 'Masuk ke Dashboard'}
      </button>
    </form>

    <div class="mt-6 pt-6 border-t border-slate-100 text-center">
      <p class="text-xs text-slate-500">
        Default Kredensial: <span class="font-mono font-semibold text-slate-700">admin</span> / <span class="font-mono font-semibold text-slate-700">admin123</span>
      </p>
    </div>
  </div>
</div>
