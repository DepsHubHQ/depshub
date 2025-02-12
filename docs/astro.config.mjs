// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";

// https://astro.build/config
export default defineConfig({
  redirects: {
    "/": "/what-is-depshub",
  },
  integrations: [
    starlight({
      editLink: {
        baseUrl: "https://github.com/depshubhq/depshub/edit/main/docs/",
      },
      title: "DepsHub",
      social: {
        github: "https://github.com/depshubhq/depshub",
      },
      logo: {
        src: "/public/logo.svg",
        alt: "DepsHub logo",
        replacesTitle: true,
      },
      favicon: "/favicon.ico",
      sidebar: [
        {
          label: "Getting started",
          items: [
            { label: "What is DepsHub?", slug: "what-is-depshub" },
            { label: "Why?", slug: "why" },
            { label: "Installation", slug: "installation" },
          ],
        },
        {
          label: "Usage",
          items: [
            { label: "Examples", slug: "guides/examples" },
            { label: "CI/CD integrations", slug: "guides/integrations" },
            //{ label: "Creating custom rules", slug: "guides/custom" },
          ],
        },
        {
          label: "Reference",
          autogenerate: { directory: "reference" },
        },
        {
          label: "Misc",
          items: [
            { label: "Supported package managers", slug: "misc/supported" },
            { label: "Technical details", slug: "misc/technical-details" },
            { label: "Contributions", slug: "misc/contributions" },
          ],
        },
      ],
      customCss: ["./src/styles/custom.css"],
      pagination: false,
      components: {
        ThemeSelect: "./src/components/ThemeSelect.astro",
        ThemeProvider: "./src/components/ThemeProvider.astro",
      },
    }),
  ],
});
