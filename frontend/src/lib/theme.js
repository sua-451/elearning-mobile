import { writable } from 'svelte/store';

function getInitialDark() {
  if (typeof window === 'undefined') return false;
  const saved = localStorage.getItem('theme');
  if (saved === 'dark') return true;
  if (saved === 'light') return false;
  return window.matchMedia('(prefers-color-scheme: dark)').matches;
}

export const isDark = writable(getInitialDark());

function applyTheme(dark) {
  if (typeof document === 'undefined') return;
  const root = document.documentElement;
  if (dark) {
    root.classList.add('dark');
  } else {
    root.classList.remove('dark');
  }
}

export function toggleTheme() {
  isDark.update(val => {
    const nextVal = !val;
    if (typeof window !== 'undefined') {
      localStorage.setItem('theme', nextVal ? 'dark' : 'light');
    }
    applyTheme(nextVal);
    return nextVal;
  });
}

export function initTheme() {
  const dark = getInitialDark();
  isDark.set(dark);
  applyTheme(dark);
}
