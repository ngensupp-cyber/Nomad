/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        desert: {
          900: '#1A120B', // Deep Desert
          800: '#2D2013', // Roasted Sand
          700: '#3E2C1C', // Dune Shadow
          primary: '#E2725B', // Terracotta
          accent: '#F4A460', // Sandy Brown
          text: '#EADBC8', // Parchment
        }
      },
      fontFamily: {
        desert: ['Outfit', 'sans-serif'],
      },
      backgroundImage: {
        'desert-gradient': 'linear-gradient(to bottom right, #1A120B, #2D2013)',
      }
    },
  },
  plugins: [],
}
