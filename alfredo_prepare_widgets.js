require("babel-plugin-transform-react-jsx");
console.log(require("babel-core").transform(process.argv[2], {
  plugins: ["transform-react-jsx"]
}).code);