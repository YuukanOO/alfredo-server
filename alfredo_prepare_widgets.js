console.log(require("babel-core").transform(process.argv[2], {
  plugins: [
    "transform-react-jsx",
    "transform-es2015-arrow-functions",
  ],
}).code);