/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['internal/templates/*.templ'],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
  daisyui: {
    themes: [
      "light",
      "dark",
      "cupcake",
      "bumblebee",
      "emerald",
      "corporate",
      "synthwave",
      "retro",
      "cyberpunk",
      "valentine",
      "halloween",
      "garden",
      "forest",
      "aqua",
      "lofi",
      "pastel",
      "fantasy",
      "wireframe",
      "black",
      "luxury",
      "dracula",
      "cmyk",
      "autumn",
      "business",
      "acid",
      "lemonade",
      "night",
      "coffee",
      "winter",
      "dim",
      "nord",
      "sunset",
    ],
  },
  safelist: [
    "from-slate-400", "from-slate-500", "from-gray-400", "from-gray-500", "from-zinc-400", "from-zinc-500", "from-neutral-400", "from-neutral-500", "from-stone-400", "from-stone-500", "from-red-400", "from-red-500", "from-orange-400", "from-orange-500", "from-amber-400", "from-amber-500", "from-yellow-400", "from-yellow-500", "from-lime-400", "from-lime-500", "from-green-400", "from-green-500", "from-emerald-400", "from-emerald-500", "from-teal-400", "from-teal-500", "from-cyan-400", "from-cyan-500", "from-sky-400", "from-sky-500", "from-blue-400", "from-blue-500", "from-indigo-400", "from-indigo-500", "from-violet-400", "from-violet-500", "from-purple-400", "from-purple-500", "from-fuchsia-400", "from-fuchsia-500", "from-pink-400", "from-pink-500", "from-rose-400", "from-rose-500"
    ,"to-slate-400", "to-slate-500", "to-gray-400", "to-gray-500", "to-zinc-400", "to-zinc-500", "to-neutral-400", "to-neutral-500", "to-stone-400", "to-stone-500", "to-red-400", "to-red-500", "to-orange-400", "to-orange-500", "to-amber-400", "to-amber-500", "to-yellow-400", "to-yellow-500", "to-lime-400", "to-lime-500", "to-green-400", "to-green-500", "to-emerald-400", "to-emerald-500", "to-teal-400", "to-teal-500", "to-cyan-400", "to-cyan-500", "to-sky-400", "to-sky-500", "to-blue-400", "to-blue-500", "to-indigo-400", "to-indigo-500", "to-violet-400", "to-violet-500", "to-purple-400", "to-purple-500", "to-fuchsia-400", "to-fuchsia-500", "to-pink-400", "to-pink-500", "to-rose-400", "to-rose-500"
  ]
}

