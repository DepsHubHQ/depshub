// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	integrations: [
		starlight({
			title: 'DepsHub Docs',
			social: {
				github: 'https://github.com/depshubhq/depshub',
			},
			sidebar: [
        {
          label: 'Getting started',
          items: [
            { label: 'What is DepsHub?', slug: 'what-is-depshub' },
            { label: 'Installation', slug: 'installation' },
          ],
        },
        {
          label: 'Guides',
          items: [
            { label: 'Linter', slug: 'guides/linter' },
            { label: 'Updater', slug: 'guides/updater' },
            { label: 'CI/CD integrations', slug: 'guides/integrations' },
          ],
        },
        {
          label: 'Configuration',
          autogenerate: { directory: 'reference' },
        },
        {
          label: 'Misc',
          items: [
            { label: 'Supported languages', slug: 'misc/supported' },
            { label: 'Technical details', slug: 'misc/technical-details' },
            { label: 'Contributions', slug: 'misc/contributions' },
          ],
        },
			],
      customCss: [
        './src/styles/custom.css',
      ],
      pagination: false,
      expressiveCode: {
        themes: [],
        useStarlightDarkModeSwitch: false,
      },
		}),
	],
});
