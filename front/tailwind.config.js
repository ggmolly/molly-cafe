/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["index.html", "tw.css"],
  theme: {
    extend: {
      fontFamily: {
        'mono': ['"Berkeley Mono"', 'monospace']
      }
    },
  },
  plugins: [],
}

