module.exports = {
  content: ['./dist/**/*.{html,js}'],
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
