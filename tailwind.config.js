module.exports = {
  content: ['./dist/**/*.{html,js}'],
  theme: {
    extend: {
      gridTemplateColumns: {
        15: 'repeat(15, minmax(0, 1fr))',
      },
    },
  },
  plugins: [require('@tailwindcss/forms')],
};
