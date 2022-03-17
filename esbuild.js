const { sassPlugin, postcssModules } = require("esbuild-sass-plugin");

require("esbuild")
  .build({
    entryPoints: ["front/index.tsx"],
    bundle: true,
    outfile: "static/dist/main.js",
    watch: !!process.env.WATCH,
    sourcemap: true,
    plugins: [
      sassPlugin({
        transform: postcssModules({
          basedir: "./front",
        }),
      }),
    ],
  })
  .catch(() => process.exit(1));
