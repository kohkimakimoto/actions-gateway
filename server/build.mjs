import * as esbuild from 'esbuild';

await esbuild.build({
  entryPoints: ['resources/js/app.js'],
  bundle: true,
  minify: true,
  outdir: 'public/assets',
})
