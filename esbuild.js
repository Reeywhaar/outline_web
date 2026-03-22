const { sassPlugin, postcssModules } = require("esbuild-sass-plugin");
const esbuild = require("esbuild");

async function main() {
  const options = {
    entryPoints: ["front/index.tsx"],
    bundle: true,
    outfile: "static/dist/main.js",
    sourcemap: true,
    minify: !process.env.WATCH,
    plugins: [
      sassPlugin({
        transform: postcssModules({
          basedir: "./front",
        }),
      }),
    ],
  };
  if (process.env.WATCH) {
    const ctx = await esbuild.context(options);

    await ctx.watch();
  } else {
    await esbuild.build(options)
  };
}

main().catch((e) => {
  console.error(e);
});
