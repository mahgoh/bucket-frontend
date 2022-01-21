module.exports = {
  content: ['./dist/**/*.{html,js}'],
  safelist: [
    {
      pattern:
        /text-(red|orange|amber|yellow|lime|green|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose|gray)-(500|700)/,
    },
    {
      pattern:
        /bg-(red|orange|amber|yellow|lime|green|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose|gray)-(100)/,
    },
  ],
  theme: {
    extend: {
      gridTemplateColumns: {
        13: 'repeat(13, minmax(0, 1fr))',
        15: 'repeat(15, minmax(0, 1fr))',
      },
    },
  },
  plugins: [require('@tailwindcss/forms')],
};
