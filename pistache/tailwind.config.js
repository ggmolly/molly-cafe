/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["*.html", "pistache-tw.css"],
  theme: {
    extend: {
      fontFamily: {
        'mono': ['"Berkeley Mono"', 'monospace']
      }
    },
  },
  plugins: [],
  corePlugins: {
    preflight: false,
  }
}

