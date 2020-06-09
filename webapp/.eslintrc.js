module.exports = {
  root: true,
  env: {
    node: true,
  },
  extends: [
    "plugin:vue/essential",
  ],
  parserOptions: {
    parser: "babel-eslint",
    ecmaVersion: 2017,
  },
  rules: {
    "quotes": ["error", "double"],
    "semi": ["error", "always"]
  },
};
