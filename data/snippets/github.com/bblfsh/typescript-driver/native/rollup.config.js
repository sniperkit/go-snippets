import nodeResolve from "rollup-plugin-node-resolve";
import commonJS from 'rollup-plugin-commonjs';

export default {
  input: 'src/index.js',
  output: {
    file: 'lib/index.js',
    format: 'cjs'
  },
  banner: '#!/usr/bin/env node',
  exports: 'named',
  plugins: [
    nodeResolve({
      jsnext: true,
      main: true
    }),
    commonJS({})
  ],
  treeshake: false,
  external: ['readline']
};