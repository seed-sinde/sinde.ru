export default defineNitroPlugin((nitroApp) => {
  nitroApp.hooks.hook("render:html", (html) => {
    html.head.unshift(`
      <script>
        (function() {
          const s = localStorage.getItem('theme') || 'system';
          const isDark = s === 'dark' || (s === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches);
          if (isDark) document.documentElement.classList.add('dark');
        })();
      </script>
    `);
  });
});
