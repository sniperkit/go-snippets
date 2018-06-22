import babel from "rollup-plugin-babel";
import nodeResolve from "rollup-plugin-node-resolve";
import commonJS from 'rollup-plugin-commonjs';
import multientry from 'rollup-plugin-multi-entry';

export default {
  entry: "src/*.js",
  dest: "lib/index.js",
  banner: '#!/usr/bin/env node',
  exports: 'named',
  plugins: [
    multientry(),
    babel(),
    nodeResolve({
      jsnext: true,
      main: true,
    }),
    commonJS({
    }),
  ],
  format: "cjs",
  treeshake: false,
  external: ['readline'],
};
