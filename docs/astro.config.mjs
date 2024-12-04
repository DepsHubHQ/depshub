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
            { label: 'What is DepsHub?', slug: 'guides/example' },
            { label: 'Installation', slug: 'getting-started' },
          ],
        },
        {
          label: 'Usage',
          items: [
            { label: 'Linter', slug: 'guides/example' },
            { label: 'Updater', slug: 'getting-started' },
          ],
        },
        {
          label: 'Configuration',
          autogenerate: { directory: 'reference' },
        },
				{
          label: 'Automation',
					items: [
						{ label: 'GitHub Actions', slug: 'guides/example' },
            { label: 'GitLab CI', slug: 'guides/example' },
            { label: 'Bitbucket Pipelines', slug: 'guides/example' },
            { label: 'Jenkins', slug: 'guides/example' },
            { label: 'Azure DevOps', slug: 'guides/example' },
            { label: 'Travis CI', slug: 'guides/example' },
					],
				},
        {
          label: 'Misc',
          items: [
            { label: 'Technical details', slug: 'guides/example' },
            { label: 'Contributions', slug: 'guides/example' },
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
