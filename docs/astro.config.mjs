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
        { label: 'Getting started', slug: 'getting-started' },
				{
					label: 'Guides',
					items: [
						{ label: 'Example Guide', slug: 'guides/example' },
					],
				},
				{
					label: 'Reference',
					autogenerate: { directory: 'reference' },
				},
			],
      customCss: [
        './src/styles/custom.css',
      ],
		}),
	],
});
