export default {
    env: {
        browser: true,
        node: true,
        es2021: true,
    },

    // Extend from ESLint's recommended rules.
    extends: ['eslint:recommended'],

    // Parser options allow you to specify the ECMAScript version and module system.
    parserOptions: {
        ecmaVersion: 12,
        sourceType: 'module',
    },

    // Custom rules to override the defaults from the extended configurations.
    rules: {
        'no-unused-vars': 'warn',
      // You can add more custom rules here:
      // 'indent': ['error', 2],
      // 'quotes': ['error', 'single'],
    },
};
